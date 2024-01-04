package fsl

import (
	"context"
	"encoding/json"
)

// PassState by default passes its input to its output, performing no work.
type PassState struct {
	// Type is the type's name of PassState
	// +Required
	// MUST be `Pass`
	Type StateType `json:"Type"`

	// Comment provided for human-readable description of the state.
	// +Optional
	Comment string `json:"Comment,omitempty"`

	// InputPath is an Path, which is applied to a State's raw input to select some
	// or all of it, that selection is used by the state.
	// +Optional
	// Defaults to '$'
	InputPath Path `json:"InputPath,omitempty"`

	// OutputPath is an Path, which is applied to the state's output after the application of `ResultPath`,
	// producing the effective output which serves as the raw input for the next state.
	// +Optional
	// Defaults to '$'
	OutputPath Path `json:"OutputPath,omitempty"`

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
	ResultPath ReferencePath `json:"ResultPath,omitempty"`

	// Parameters is a Payload Template which is a JSON object, whose input is the result of
	// applying the `InputPath` to the raw input.
	// If the `Parameters` is provided, its payload, after the extraction and embedding,
	// becomes the effective result.
	// +Optional
	Parameters *PayloadTemplate `json:"Parameters,omitempty"`

	// ------ State specified field

	// Result is the output of the virtual task, if it not provided,
	// the output is the input.
	// +Optional
	Result *PayloadTemplate `json:"Result,omitempty"`
}

// UnmarshalJSON implement json.Unmarshaler and will set defaults
// value for PassState.
func (s *PassState) UnmarshalJSON(data []byte) error {
	s.InputPath = "$"
	s.OutputPath = "$"
	s.ResultPath = "$"

	type PassStateForUnmarshal PassState
	return json.Unmarshal(data, (*PassStateForUnmarshal)(s))
}

func (s *PassState) ApplyInput(ctx context.Context, sc StateContext) ([]byte, error) {
	inputObj, err := s.InputPath.Apply(ctx, pathContext{
		Input:       sc.Input(),
		ContextData: sc.ContextData(),
	})
	if err != nil {
		return nil, err
	}

	if s.Parameters == nil {
		return inputObj, nil
	}

	inputObj, err = s.Parameters.Apply(ctx, payloadTemplateContext{
		Input:       inputObj,
		ContextData: sc.ContextData(),
	})
	if err != nil {
		return nil, err
	}

	return inputObj, nil
}

func (s *PassState) ApplyOutput(ctx context.Context, sc StateContext) ([]byte, error) {
	outputObj := sc.Output()
	var err error
	if s.Result != nil {
		outputObj, err = s.Result.Apply(ctx, payloadTemplateContext{
			Input:       outputObj,
			ContextData: sc.ContextData(),
		})
		if err != nil {
			return nil, err
		}
	}

	outputObj, err = s.ResultPath.Apply(ctx, referencePathContext{
		Input:  sc.Input(),
		Output: outputObj,
	})
	if err != nil {
		return nil, err
	}

	outputObj, err = s.OutputPath.Apply(ctx, pathContext{
		Input:       outputObj,
		ContextData: sc.ContextData(),
	})
	if err != nil {
		return nil, err
	}

	return outputObj, nil
}

// Validate will validate the PassState configuration
func (s *PassState) Validate(_ context.Context) error {
	return nil
}
