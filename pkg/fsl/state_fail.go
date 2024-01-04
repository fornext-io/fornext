package fsl

import (
	"context"
)

// FailState terminates the machine and marks it as a failure.
type FailState struct {
	// Type is the type's name of FailState
	// +Required
	// MUST be `Fail`
	Type StateType `json:"Type"`

	// Comment provided for human-readable description of the state.
	// +Optional
	Comment string `json:"Comment,omitempty"`

	// ------ State specified field

	// Error is the an error name that can be used for error handling, operational
	// or diagnostic perposes.
	// +Required
	Error string `json:"Error"`

	// Cause used to provide a human-readable message.
	// +Required
	Cause string `json:"Cause"`
}

// Validate will validate the FailState configuration
func (s *FailState) Validate(_ context.Context) error {
	return nil
}
