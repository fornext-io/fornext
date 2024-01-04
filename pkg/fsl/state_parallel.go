package fsl

import (
	"context"
	"encoding/json"
)

// ParallelState causes parallel executino of branches.
type ParallelState struct {
	// Type is the type's name of ParallelState
	// +Required
	// MUST be `Parallel`
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

	// Branches is an array of `States` are executed in parallel.
	// +Required
	Branches []BranchProcessor `json:"Branches"`
}

// BranchProcessor is the `States` executed in parallel.
type BranchProcessor struct {
	// StartAt is the name of state which inerpreter will start running.
	// The value MUST exactly match one of names of the `States` field.
	// +Required
	StartAt string `json:"StartAt"`

	// States represent the states.
	// +Required
	States States `json:"States"`
}

// UnmarshalJSON implement json.Unmarshaler and will set defaults
// value for ParallelState.
func (s *ParallelState) UnmarshalJSON(data []byte) error {
	s.InputPath = "$"
	s.OutputPath = "$"
	s.ResultPath = "$"

	type ParallelStateForUnmarshal ParallelState
	return json.Unmarshal(data, (*ParallelStateForUnmarshal)(s))
}

// Validate will validate the ParallelState configuration
func (s *ParallelState) Validate(_ context.Context) error {
	return nil
}

func (s *ParallelState) ApplyInput(ctx context.Context, sc StateContext) ([]byte, error) {
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

func (s *ParallelState) ApplyOutput(ctx context.Context, sc StateContext) ([]byte, error) {
	outputObj := sc.Output()
	var err error
	if s.ResultSelector != nil {
		outputObj, err = s.ResultSelector.Apply(ctx, payloadTemplateContext{
			Input:       sc.Output(),
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
