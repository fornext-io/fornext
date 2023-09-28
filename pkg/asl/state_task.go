package asl

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/lsytj0413/ena/xerrors"
	"github.com/ohler55/ojg/jp"
)

// TaskState causes the interpreter to execute the work identified by the state's
// `Resource` field.
type TaskState struct {
	// Type is the type's name of TaskState
	// +Required
	// MUST be `Task`
	Type StateType `json:"Type"`

	// Comment provided for human-readable description of the state.
	// +Optional
	Comment string `json:"Comment,omitempty"`

	// InputPath is an Path, which is applied to a State's raw input to select some
	// or all of it, that selection is used by the state.
	// +Optional
	// Defaults to '$'
	InputPath *Path `json:"InputPath,omitempty"`

	// OutputPath is an Path, which is applied to the state's output after the application of `ResultPath`,
	// producing the effective output which serves as the raw input for the next state.
	// +Optional
	// Defaults to '$'
	OutputPath *Path `json:"OutputPath,omitempty"`

	// Next is the name of state the interpreter follows a transition to.
	// It MUST exactly and case-sensitively match the name of the another state.
	// +Optional
	Next string `json:"Next,omitempty"`

	// End causes the interpreter to terminate the machine.
	// +Optional
	End bool `json:"End,omitempty"`

	// ResultPath is a Reference Path, which specifies the raw input's combination with or
	// replacement by the state's result.
	// The value of `ResultPath` MUST NOT begin with '$$'
	// +Optional
	// Defaults to '$'
	ResultPath *ReferencePath `json:"ResultPath,omitempty"`

	// Parameters is a Payload Template which is a JSON object, whose input is the result of
	// applying the `InputPath` to the raw input.
	// If the `Parameters` is provided, its payload, after the extraction and embedding,
	// becomes the effective result.
	// +Optional
	Parameters *PayloadTemplate `json:"Parameters,omitempty"`

	// ResultSelector is a Payload Template, whose input is the result, and whose
	// payload replaces and becomes the effective result.
	// +Optional
	ResultSelector *PayloadTemplate `json:"ResultSelector,omitempty"`

	// Retry is an array of Retrier, when a state reports an error, the interpreter scans through the
	// Retriers and, when the Error Name appears in the value of
	// a Retrier's `ErrorEquals` field, implements the retry policy described in that Retrier.
	// +Optional
	Retry Retriers `json:"Retry,omitempty"`

	// Catch is an array of Catcher, when a state reports an error and either there is no Retrier,
	// or retries have failed to resolve the error, the interpreter scans through the Catchers in
	// array order, and when the Error Name appears in the value of a Catcher's `ErrorEquals` field,
	// transitions the machine to the state named in the value of the `Next` field.
	// +Optional
	Catch Catchers `json:"Catch,omitempty"`

	// ------ State specified field

	// Resource is an URI that uniquely identifies the specific task to execute.
	// +Required
	Resource string `json:"Resource"`

	// If the state runs longer than the specified timeout, or if more time than the specified
	// heartbeat elapses between heartbeats from the task, then the interpreter fails the state with
	// a `States.Timeout` Error Name.

	// TimeoutSeconds is an specify timeouts.
	// +Optional
	// MUST be positive integers.
	// Defaults to 60.
	TimeoutSeconds *int `json:"TimeoutSeconds,omitempty"`

	// HeartbeatSeconds is an specify timeouts.
	// +Optional
	// MUST be positive integers.
	// MUST be smaller than `TimeoutSeconds` value.
	HeartbeatSeconds *int `json:"HeartbeatSeconds,omitempty"`

	// TimeoutSecondsPath is a Reference Path
	// +Optional
	// MUST select fields whose values are positive integers.
	TimeoutSecondsPath *ReferencePath `json:"TimeoutSecondsPath,omitempty"`

	// HeartbeatSecondsPath is a Reference Path
	// +Optional
	// MUST select fields whose values are positive integers.
	HeartbeatSecondsPath *ReferencePath `json:"HeartbeatSecondsPath,omitempty"`

	// Credentials MUST be a JSON object whose value is defined by the interpreter.
	// +Optional
	// DIDN'T Support except Defaults value.
	Credentials *Credential `json:"Credentials,omitempty"`
}

// Credential specifies a target role the state machine's execution role
// must assume before invoking the specified Resource
type Credential struct {
}

// Validate will validate the TaskState configuration
func (s *TaskState) Validate(_ context.Context) error {
	return nil
}

// GetJobInput ...
func (s *TaskState) GetJobInput(obj interface{}, contextObject interface{}) (interface{}, error) {
	targetObj := obj

	// 1. first, input path
	if s.InputPath != nil {
		// If the value of InputPath is null, that means that the raw input is discarded, and the effective input for the state is an empty JSON object, {}. Note that having a value of null is different from the "InputPath" field being absent.

		if *s.InputPath == "" {
			return nil, xerrors.Errorf("invalid empty InputPath")
		}

		inputObj := obj
		path := string(*s.InputPath)
		if strings.HasPrefix(string(*s.InputPath), "$$") {
			// context object
			inputObj = contextObject
			path = path[1:]
		}

		expr, err := jp.ParseString(path)
		if err != nil {
			return nil, err
		}

		// If the result is multiple values, we must change to an JSON array containing all of them
		// If the result is null (no element? one element with nil?), the input change to {} (empty object)
		result := expr.Get(inputObj)
		if len(result) != 1 {
			return nil, xerrors.Errorf("must only one element InputPath")
		}

		targetObj = result[0]
	}

	// 2. use parameter to change input to Task arguments
	if s.Parameters != nil {
		var pObj interface{}
		err := json.Unmarshal([]byte(*s.Parameters), &pObj)
		if err != nil {
			return nil, err
		}

		targetObj, err = renderPayloadTemplate(pObj, targetObj, contextObject)
		if err != nil {
			return nil, err
		}
	}

	// 3. ResultSelector
	if s.ResultSelector != nil {
		var pObj interface{}
		err := json.Unmarshal([]byte(*s.ResultSelector), &pObj)
		if err != nil {
			return nil, err
		}

		targetObj, err = renderPayloadTemplate(pObj, targetObj, contextObject)
		if err != nil {
			return nil, err
		}
	}

	// 4. ResultPath
	if s.ResultPath != nil {
		// If the value of ResultPath is null, that means that the state’s result is discarded and its raw input becomes its result.

		expr, err := jp.ParseString((string)(*s.ResultPath))
		if err != nil {
			return nil, err
		}

		err = expr.Set(obj, targetObj)
		if err != nil {
			return nil, err
		}
	}

	// 5. OutputPath
	if s.OutputPath != nil {
		// If the value of OutputPath is null, that means the input and result are discarded, and the effective output from the state is an empty JSON object, {}.

		inputObj := targetObj
		path := (string)(*s.OutputPath)
		if strings.HasPrefix(path, "$$") {
			path = path[1:]
			inputObj = contextObject
		}

		expr, err := jp.ParseString(path)
		if err != nil {
			return nil, err
		}

		// If the result is multiple values, we must change to an JSON array containing all of them
		result := expr.Get(inputObj)
		if len(result) != 1 {
			return nil, xerrors.Errorf("must only one element InputPath")
		}

		targetObj = result[0]
	}

	return targetObj, nil
}

func renderPayloadTemplate(templateObj interface{}, inputObj interface{}, contextObject interface{}) (interface{}, error) {
	switch templateObj2 := templateObj.(type) {
	case map[string]interface{}:
		var targetObj map[string]interface{}
		for k, v := range templateObj2 {
			if !strings.HasSuffix(k, ".$") {
				obj, err := renderPayloadTemplate(v, inputObj, contextObject)
				if err != nil {
					return nil, err
				}

				targetObj[k] = obj
				continue
			}

			vv, ok := v.(string)
			if !ok {
				return nil, xerrors.Errorf("key with .$ suffix, it value must be string")
			}

			if strings.HasPrefix(vv, "$$") {
				// It was from contextObject
				path := vv[1:]
				expr, err := jp.ParseString(path)
				if err != nil {
					return nil, err
				}

				result := expr.Get(contextObject)
				if len(result) != 1 {
					return nil, xerrors.Errorf("must only one element InputPath")
				}

				targetObj[k] = result[0]
			} else if strings.HasPrefix(vv, "$") {
				expr, err := jp.ParseString(vv)
				if err != nil {
					return nil, err
				}

				result := expr.Get(contextObject)
				if len(result) != 1 {
					return nil, xerrors.Errorf("must only one element InputPath")
				}

				targetObj[k] = result[0]
			} else if strings.HasPrefix(vv, "States.") {
				return nil, xerrors.Errorf("not support Instrinsic now")
			} else {
				return nil, xerrors.Errorf("key with .$ suffix, it value must be Path or Instrinsic")
			}
		}
		return targetObj, nil
	case []interface{}:
		var targetObj []interface{}
		for _, obj := range templateObj2 {
			obj2, err := renderPayloadTemplate(obj, inputObj, contextObject)
			if err != nil {
				return nil, err
			}

			targetObj = append(targetObj, obj2)
		}
		return targetObj, nil
	default:
		return templateObj, nil
	}
}
