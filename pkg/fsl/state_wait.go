package fsl

import (
	"context"
	"encoding/json"
	"time"

	"github.com/lsytj0413/ena/xerrors"
)

// WaitState causes the interpreter to delay the machine from continuing
// for a specified time.
type WaitState struct {
	// Type is the type's name of WaitState
	// +Required
	// MUST be `Wait`
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

	// ------ State specified field

	// Seconds is the second duration to delay.
	// +Optional
	Seconds *int `json:"Seconds,omitempty"`

	// SecondsPath is a Reference Path to the seconds data.
	// +Optional
	SecondsPath *ReferencePath `json:"SecondsPath,omitempty"`

	// Timestamp is absolute expiry time specified as an ISO-8601 extended offset data-time format string.
	// +Optional
	Timestamp *time.Time `json:"Timestamp,omitempty"`

	// Timestamp is an Refercence Path to the Timestamp data.
	// +Optional
	TimestampPath *ReferencePath `json:"TimestampPath,omitempty"`
}

// UnmarshalJSON implement json.Unmarshaler and will set defaults
// value for WaitState.
func (s *WaitState) UnmarshalJSON(data []byte) error {
	s.InputPath = "$"
	s.OutputPath = "$"

	type WaitStateForUnmarshal WaitState
	return json.Unmarshal(data, (*WaitStateForUnmarshal)(s))
}

func (s *WaitState) ApplyInput(ctx context.Context, sc StateContext) ([]byte, error) {
	inputObj, err := s.InputPath.Apply(ctx, pathContext{
		Input:       sc.Input(),
		ContextData: sc.ContextData(),
	})
	if err != nil {
		return nil, err
	}

	return inputObj, nil
}

func (s *WaitState) ApplyOutput(ctx context.Context, sc StateContext) ([]byte, error) {
	outputObj, err := s.OutputPath.Apply(ctx, pathContext{
		Input:       sc.Output(),
		ContextData: sc.ContextData(),
	})
	if err != nil {
		return nil, err
	}

	return outputObj, nil
}

// Validate will validate the WaitState configuration
func (s *WaitState) Validate(_ context.Context) error {
	count := 0
	if s.Seconds != nil {
		count++

		if *s.Seconds <= 0 {
			return xerrors.New("Seconds MUST be positive")
		}
	}

	if s.SecondsPath != nil {
		count++
	}

	if s.Timestamp != nil {
		count++
	}

	if s.TimestampPath != nil {
		count++
	}

	if count != 1 {
		return xerrors.New("MUST contain exactly one of `Seconds`、`SecondsPath`、`Timestamp` or `TimestampPath`")
	}
	return nil
}
