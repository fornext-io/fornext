package asl

import (
	"context"
	"encoding/json"

	"github.com/lsytj0413/ena/xerrors"
)

// State represent a step of state machine
type State interface {
	// Validate will validate the state configuration
	Validate(context.Context) error
}

// StateType present the type of State object
type StateType = string

const (
	// StateTypeWait is the type of Wait State
	StateTypeWait StateType = "Wait"

	// StateTypeTask is the type of Task State
	StateTypeTask StateType = "Task"

	// StateTypeSucceed is the type of Succeed State
	StateTypeSucceed StateType = "Succeed"

	// StateTypePass is the type of Pass State
	StateTypePass StateType = "Pass"

	// StateTypeMap is the type of Map State
	StateTypeMap StateType = "Map"

	// StateTypeFail is the type of Fail State
	StateTypeFail StateType = "Fail"

	// StateTypeChoice is the type of Choice State
	StateTypeChoice StateType = "Choice"

	// StateTypeParallel is the type of Parallel State
	StateTypeParallel StateType = "Parallel"
)

// StateFactory to construct a type reference State object
type StateFactory func() State

var (
	stateFactories = map[StateType]StateFactory{
		StateTypeWait: func() State {
			return &WaitState{}
		},
		StateTypeTask: func() State {
			return &TaskState{}
		},
		StateTypeSucceed: func() State {
			return &SucceedState{}
		},
		StateTypePass: func() State {
			return &PassState{}
		},
		StateTypeMap: func() State {
			return &MapState{}
		},
		StateTypeFail: func() State {
			return &FailState{}
		},
		StateTypeChoice: func() State {
			return &ChoiceState{}
		},
		StateTypeParallel: func() State {
			return &ParallelState{}
		},
	}
)

// States present a collection of state
type States map[string]State

// UnmarshalJSON return the States map from json bytes
func (s *States) UnmarshalJSON(b []byte) error {
	rawStatesMap := map[string]json.RawMessage{}
	err := json.Unmarshal(b, &rawStatesMap)
	if err != nil {
		return err
	}

	v := make(map[string]State, len(rawStatesMap))
	for stateName, stateData := range rawStatesMap {
		o, err := UnmarshalStateFromJSON(stateData)
		if err != nil {
			return err
		}

		v[stateName] = o
	}

	(*s) = v
	return nil
}

// Validate will validate the States configuration
func (s *States) Validate(ctx context.Context) error {
	for _, state := range *s {
		err := state.Validate(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

// BaseState is the struct for unmarshal State's Type field
type BaseState struct {
	Type StateType `json:"Type"`
}

// UnmarshalStateFromJSON will unmarshal the json's byte slice to State object
func UnmarshalStateFromJSON(b []byte) (State, error) {
	bs := &BaseState{}
	err := json.Unmarshal(b, bs)
	if err != nil {
		return nil, err
	}

	objectFactory, ok := stateFactories[bs.Type]
	if !ok {
		return nil, xerrors.Errorf("unsupport state type: %s", bs.Type)
	}

	ss := objectFactory()
	err = json.Unmarshal(b, ss)
	if err != nil {
		return nil, err
	}
	return ss, nil
}
