package fsl

import (
	"context"
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

type payloadTemplateContext struct {
	Input       []byte
	ContextData []byte
}

// Apply will apply this payload template on provided context
func (p PayloadTemplate) Apply(_ context.Context, pc payloadTemplateContext) ([]byte, error) {
	var payloadObj interface{}
	err := json.Unmarshal([]byte(string(p)), &payloadObj)
	if err != nil {
		return nil, err
	}

	var inputObj interface{}
	err = json.Unmarshal(pc.Input, &inputObj)
	if err != nil {
		return nil, err
	}

	var contextObj interface{}
	err = json.Unmarshal(pc.ContextData, &contextObj)
	if err != nil {
		return nil, err
	}

	resultObj, err := renderPayloadTemplate(payloadObj, payloadTemplateRenderContext{
		inputObject:   inputObj,
		contextObject: contextObj,
		Input:         pc.Input,
		ContextData:   pc.ContextData,
	})
	if err != nil {
		return nil, err
	}

	result, err := json.Marshal(resultObj)
	if err != nil {
		return nil, err
	}
	return result, nil
}

type payloadTemplateRenderContext struct {
	inputObject   interface{}
	contextObject interface{}

	Input       []byte
	ContextData []byte
}

func renderPayloadTemplate(templateObj interface{}, pc payloadTemplateRenderContext) (interface{}, error) {
	switch templateObj2 := templateObj.(type) {
	case map[string]interface{}:
		return renderMapPayloadTemplate(templateObj2, pc)
	case []interface{}:
		return renderSlicePayloadTemplate(templateObj2, pc)
	default:
		return templateObj, nil
	}
}

func renderSlicePayloadTemplate(templateObject []interface{}, pc payloadTemplateRenderContext) ([]interface{}, error) {
	targetObj := make([]interface{}, 0, len(templateObject))
	for _, obj := range templateObject {
		pobj, err := renderPayloadTemplate(obj, pc)
		if err != nil {
			return nil, err
		}

		targetObj = append(targetObj, pobj)
	}
	return targetObj, nil
}

func renderMapPayloadTemplate(
	templateObject map[string]interface{},
	pc payloadTemplateRenderContext,
) (map[string]interface{}, error) {
	targetObj := map[string]interface{}{}
	for k, v := range templateObject {
		if !strings.HasSuffix(k, ".$") {
			obj, err := renderPayloadTemplate(v, pc)
			if err != nil {
				return nil, err
			}

			targetObj[k] = obj
			continue
		}

		k = strings.TrimSuffix(k, ".$")
		vv, ok := v.(string)
		if !ok {
			return nil, xerrors.Errorf("key with .$ suffix, it value must be string")
		}

		switch {
		case strings.HasPrefix(vv, "$"):
			p := Path(vv)
			result, err := p.Apply(context.Background(), pathContext{
				Input:       pc.Input,
				ContextData: pc.ContextData,
			})
			if err != nil {
				return nil, err
			}

			var resultObj interface{}
			err = json.Unmarshal(result, &resultObj)
			if err != nil {
				return nil, err
			}
			targetObj[k] = resultObj
		case strings.HasPrefix(vv, "States."):
			return nil, xerrors.Errorf("not support Instrinsic now")
		default:
			return nil, xerrors.Errorf("key with .$ suffix, it value must be Path or Instrinsic")
		}
	}
	return targetObj, nil
}
