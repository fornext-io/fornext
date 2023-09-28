package asl

import (
	"encoding/json"
	"strings"

	"github.com/lsytj0413/ena/xerrors"
)

// PayloadTemplate ia a JSON Object to reshape input data and output data to meet the format expectations of users.
// It MUST be a JSON object.
//
//  1. If any field within the Payload Template has a name ending with '.$',
//     its value is transformed according to rules below and the the field is renamed to strip the '.$' suffix.
//     1.1. If the field value begins with only one '$', the value MUST be a Path
//     1.2. If the field value begins with '$$', the first dollar sign is stripped and the remainder
//     MUST be a Path. In this case, the Path is applied to the Context Object and is the new field value.
//     1.3. If the field value doesn't begin with '$', it MUST be an Instrinsic Function.
//     1.4. If the path is legal but cannot be applied successfully, the interpreter fails the machine execution
//     with an Error Name of "States.ParameterPathFailure". If the Intrinsic Function fails during evaluation,
//     the interpreter fails the machine execution with an Error Name of "States.IntrinsicFailure".
//
// A JSON object MUST NOT have duplicate field names after fields ending with '$' are renamed to strip the suffix.
type PayloadTemplate string

// UnmarshalJSON will return the object of json
func (p *PayloadTemplate) UnmarshalJSON(b []byte) error {
	*p = PayloadTemplate(string(b))
	return nil
}

// MarshalJSON will return the json bytes of object
func (p *PayloadTemplate) MarshalJSON() ([]byte, error) {
	return []byte(string(*p)), nil
}

// ApplyPayloadTemplate ...
func ApplyPayloadTemplate(payload *string, obj interface{}, contextObj interface{}) (interface{}, error) {
	if payload == nil {
		return obj, nil
	}

	p := *payload
	var payloadObj interface{}
	err := json.Unmarshal([]byte(p), &payloadObj)
	if err != nil {
		return nil, xerrors.Errorf("payload MUST be a JSON object, %w", err)
	}

	return ApplyPayloadTemplateOnObject(payloadObj, obj, contextObj)
}

// ApplyPayloadTemplateOnObject ...
func ApplyPayloadTemplateOnObject(
	payloadObject interface{},
	obj interface{},
	contextObj interface{},
) (interface{}, error) {
	switch pObj := payloadObject.(type) {
	case map[string]interface{}:
		return ApplyPayloadTemplateMapOnObject(pObj, obj, contextObj)
	case []interface{}:
		return ApplyPayloadTemplateSliceOnObject(pObj, obj, contextObj)
	default:
	}

	return payloadObject, nil
}

// ApplyPayloadTemplateSliceOnObject ...
func ApplyPayloadTemplateSliceOnObject(
	payloadObject []interface{},
	obj interface{},
	contextObj interface{},
) (interface{}, error) {
	var targetObj []interface{}
	for _, v := range payloadObject {
		vv, err := ApplyPayloadTemplateOnObject(v, obj, contextObj)
		if err != nil {
			return nil, err
		}

		targetObj = append(targetObj, vv)
	}

	return targetObj, nil
}

// ApplyPayloadTemplateMapOnObject ...
func ApplyPayloadTemplateMapOnObject(
	payloadObject map[string]interface{},
	obj interface{},
	contextObj interface{},
) (interface{}, error) {
	targetObj := make(map[string]interface{})

	for k, v := range payloadObject {
		if !strings.HasSuffix(k, ".$") {
			vv, err := ApplyPayloadTemplateOnObject(v, obj, contextObj)
			if err != nil {
				return nil, err
			}

			targetObj[k] = vv
			continue
		}

		// If the field within the Payload Template has a name ending with '.$', it's value MUST be an string
		vv, ok := v.(string)
		if !ok {
			return nil, xerrors.Errorf("key with .$ suffix, it value must be string")
		}

		k = strings.TrimSuffix(k, ".$")
		switch {
		case strings.HasPrefix(vv, "$$"):
			vvv, err := ApplyPathOnObject(vv[1:], contextObj)
			if err != nil {
				return nil, err
			}

			targetObj[k] = vvv
		case strings.HasPrefix(vv, "$"):
			vvv, err := ApplyPathOnObject(vv, obj)
			if err != nil {
				return nil, err
			}

			targetObj[k] = vvv
		case strings.HasPrefix(vv, "States."):
			return nil, xerrors.Errorf("not support Instrinsic now")
		default:
			return nil, xerrors.Errorf("key with .$ suffix, it value must be Path or Instrinsic")
		}
	}

	return targetObj, nil
}
