package fsl

import (
	"context"
	"encoding/json"

	"github.com/lsytj0413/ena"
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
}

// UnmarshalJSON implement json.Unmarshaler and will set defaults
// value for TaskState.
func (s *TaskState) UnmarshalJSON(data []byte) error {
	s.InputPath = "$"
	s.OutputPath = "$"
	s.ResultPath = "$"

	type TaskStateForUnmarshal TaskState
	err := json.Unmarshal(data, (*TaskStateForUnmarshal)(s))
	if err != nil {
		return err
	}

	if s.TimeoutSecondsPath == nil && s.TimeoutSeconds == nil {
		s.TimeoutSeconds = ena.PointerTo(60)
	}
	return nil
}

// Validate will validate the TaskState configuration
func (s *TaskState) Validate(_ context.Context) error {
	return nil
}

func (s *TaskState) ApplyInput(ctx context.Context, sc StateContext) ([]byte, error) {
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

func (s *TaskState) ApplyOutput(ctx context.Context, sc StateContext) ([]byte, error) {
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
