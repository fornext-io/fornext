// Package asl provide the description and model of ASL
package asl

import (
	"context"
	"encoding/json"

	"github.com/lsytj0413/ena/xerrors"
)

// StateMachine is the state machine represent of asl
type StateMachine struct {
	// Comment provided for human-readable description of the machine.
	// +Optional
	Comment string `json:"Comment,omitempty"`

	// StartAt is the name of state which inerpreter will start running.
	// The value MUST exactly match one of names of the `States` field.
	// +Required
	StartAt string `json:"StartAt"`

	// Version gives the version of ASL used in the machine.
	// +Optional
	// Defaults to 1.0
	Version string `json:"Version,omitempty"`

	// TimeoutSeconds is the maximum number of seconds the machine is allowed to run.
	// If the machine runs longer than the specified time, then the interpreter fails
	// the machine with a `States.Timeout` Error Name.
	// +Optional
	TimeoutSeconds *int `json:"TimeoutSeconds,omitempty"`

	// States represent the states.
	// +Required
	States States `json:"States"`
}

type stateMachineForUnmarshal StateMachine

// UnmarshalJSON set the default value when unmarshal
func (s *StateMachine) UnmarshalJSON(b []byte) error {
	v := &stateMachineForUnmarshal{
		Version: "1.0",
	}

	err := json.Unmarshal(b, v)
	if err != nil {
		return err
	}

	(*s) = StateMachine(*v)
	return nil
}

// Validate will validate the StateMachine configuration
func (s *StateMachine) Validate(_ context.Context) error {
	if s.StartAt == "" {
		return xerrors.New("StartAt can not be empty")
	}

	if s.States == nil || len(s.States) == 0 {
		return xerrors.New("States can not be empty")
	}

	if _, ok := s.States[s.StartAt]; !ok {
		return xerrors.New("StartAt can not found in states")
	}

	if s.TimeoutSeconds != nil && *s.TimeoutSeconds <= 0 {
		return xerrors.New("TimeoutSeconds must be positive")
	}

	return nil
}
