package asl

import (
	"context"
)

// SucceedState either terminates a state machine successfully, ends a branch
// of a Parallel Satte, or ends an iteration of a Map State.
// The output of a Succeed State is the same as its input, possibly modified by
// `InputPath` and/or `OutputPath`.
type SucceedState struct {
	// Type is the type's name of SucceedState
	// +Required
	// MUST be `Succeed`
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

	// ------ State specified field
}

// Validate will validate the SucceedState configuration
func (s *SucceedState) Validate(_ context.Context) error {
	return nil
}
