// Package executor ...
package executor

import (
	"fmt"
	"time"

	"github.com/fornext-io/fornext/pkg/asl"
	"github.com/google/uuid"
)

// Executor ...
type Executor struct {
	sm  *asl.StateMachine
	ctx *ExecutionContext

	activityContextes  map[string]*ActivityContext
	branchContextes    map[string]*BranchContext
	iterationContextes map[string]*IterationContext

	done chan interface{}
	ev   chan interface{}
}

// NewExecutor ...
func NewExecutor(sm *asl.StateMachine) *Executor {
	v := &Executor{
		sm:                 sm,
		done:               make(chan interface{}),
		ev:                 make(chan interface{}, 100),
		activityContextes:  map[string]*ActivityContext{},
		branchContextes:    map[string]*BranchContext{},
		iterationContextes: map[string]*IterationContext{},
	}
	return v
}

// Run ...
func (e *Executor) Run() *ExecutionContext {
	e.ev <- &StartExecutionCommand{
		ID: uuid.NewString(),
	}
	go e.run()

	return e.ctx
}

// WaitExecutionDone ...
func (e *Executor) WaitExecutionDone() *ExecutionContext {
	<-e.done
	return e.ctx
}

func findState(sm *asl.StateMachine, stateName string) asl.State {
	vv, ok := sm.States[stateName]
	if ok {
		return vv
	}

	for _, vvv := range sm.States {
		vv = findStateInState(vvv, stateName)
		if vv != nil {
			return vv
		}
	}

	return nil
}

func findStateInState(ss asl.State, stateName string) asl.State {
	switch sss := ss.(type) {
	case *asl.ParallelState:
		for _, branch := range sss.Branches {
			vv, ok := branch.States[stateName]
			if ok {
				return vv
			}

			for _, vvv := range branch.States {
				vv = findStateInState(vvv, stateName)
				if vv != nil {
					return vv
				}
			}
		}
	case *asl.MapState:
		vv, ok := sss.ItemProcessor.States[stateName]
		if ok {
			return vv
		}

		for _, vvv := range sss.ItemProcessor.States {
			vv = findStateInState(vvv, stateName)
			if vv != nil {
				return vv
			}
		}
	}

	return nil
}

func (e *Executor) run() {
	for ev := range e.ev {
		switch evv := ev.(type) {
		case *StartExecutionCommand:
			e.ev <- &ExecutionStartedEvent{
				ID: evv.ID,
			}
			e.ctx = &ExecutionContext{
				ID: evv.ID,
			}
			e.ev <- &StartStateCommand{
				ActivityID: uuid.NewString(),
				StateName:  e.sm.StartAt,
			}
		case *StartStateCommand:
			e.processStartStateCommand(evv)
		case *CreateTaskCommand:
			e.processCreateTaskCommand(evv)
		case *ActivateTaskCommand:
			e.processActivateTaskCommand(evv)
		case *CompleteTaskCommand:
			e.processCompleteTaskCommand(evv)
		case *CompleteStateCommand:
			e.processCompleteStateCommand(evv)
		case *StartBranchCommand:
			e.processStartBranchCommand(evv)
		case *CompleteBranchCommand:
			e.processCompleteBranchCommand(evv)
		case *StartIterationCommand:
			e.processStartIterationCommand(evv)
		case *CompleteIterationCommand:
			e.processCompleteIterationCommand(evv)
		case *CompleteExecutionCommand:
			e.processCompleteExecutionCommand(evv)
		default:
			fmt.Printf("receieve event: %#v\n", ev)
		}
	}
}

func (e *Executor) processStartStateCommand(cmd *StartStateCommand) {
	fmt.Printf("start state: %+v\n", cmd)
	e.ev <- &StateStartedEvent{
		ActivityID: cmd.ActivityID,
	}

	e.activityContextes[cmd.ActivityID] = &ActivityContext{
		ID:                cmd.ActivityID,
		StateName:         cmd.StateName,
		ParentBranchID:    cmd.ParentBranchID,
		ParentIterationID: cmd.ParentIterationID,
	}
	state := findState(e.sm, cmd.StateName)
	switch sss := state.(type) {
	case *asl.TaskState:
		e.ev <- &CreateTaskCommand{
			ID:         uuid.NewString(),
			ActivityID: cmd.ActivityID,
			Resource:   sss.Resource,
		}
	case *asl.WaitState:
		e.ev <- &CreateTaskCommand{
			ID:         uuid.NewString(),
			ActivityID: cmd.ActivityID,
			Resource:   "__sleep",
		}
	case *asl.PassState:
		e.ev <- &CompleteStateCommand{
			ActivityID: cmd.ActivityID,
		}
	case *asl.ChoiceState:
		e.ev <- &CompleteStateCommand{
			ActivityID: cmd.ActivityID,
		}
	case *asl.SucceedState:
		e.ev <- &CompleteStateCommand{
			ActivityID: cmd.ActivityID,
		}
	case *asl.FailState:
		e.ev <- &CompleteStateCommand{
			ActivityID: cmd.ActivityID,
		}
	case *asl.ParallelState:
		e.activityContextes[cmd.ActivityID].BranchStatus = &ActivityBranchStatus{
			Max:  len(sss.Branches),
			Done: 0,
		}
		for i := range sss.Branches {
			e.ev <- &StartBranchCommand{
				BranchID:   uuid.NewString(),
				ActivityID: cmd.ActivityID,
				Index:      i,
			}
		}
	case *asl.MapState:
		e.activityContextes[cmd.ActivityID].IterationStatus = &ActivityIterationStatus{
			Max:  2,
			Done: 0,
		}
		e.ev <- &StartIterationCommand{
			IterationID: uuid.NewString(),
			ActivityID:  cmd.ActivityID,
			Index:       0,
		}
		e.ev <- &StartIterationCommand{
			IterationID: uuid.NewString(),
			ActivityID:  cmd.ActivityID,
			Index:       1,
		}
	default:
		panic(fmt.Errorf("unsupport state %T", state))
	}
}

func (e *Executor) processCreateTaskCommand(cmd *CreateTaskCommand) {
	fmt.Printf("create task %+v\n", cmd)
	e.ev <- &TaskCreatedEvent{
		ID: cmd.ID,
	}

	go func() {
		if cmd.Resource == "__sleep" {
			time.Sleep(5 * time.Second)
			e.ev <- &CompleteTaskCommand{
				TaskID:     cmd.ID,
				ActivityID: cmd.ActivityID,
			}
		} else {
			e.ev <- &CompleteTaskCommand{
				TaskID:     cmd.ID,
				ActivityID: cmd.ActivityID,
			}
		}
	}()
}

func (e *Executor) processActivateTaskCommand(_ *ActivateTaskCommand) {

}

func (e *Executor) processCompleteTaskCommand(cmd *CompleteTaskCommand) {
	fmt.Printf("complete task command: %+v\n", cmd)

	e.ev <- &TaskCompletedEvent{
		TaskID: cmd.TaskID,
	}
	e.ev <- &CompleteStateCommand{
		ActivityID: cmd.ActivityID,
	}
}

func (e *Executor) processCompleteStateCommand(cmd *CompleteStateCommand) {
	fmt.Printf("complete state command: %+v\n", cmd)

	e.ev <- &StateCompletedEvent{
		ActivityID: cmd.ActivityID,
	}
	state := findState(e.sm, e.activityContextes[cmd.ActivityID].StateName)
	at := e.activityContextes[cmd.ActivityID]

	switch sss := state.(type) {
	case *asl.TaskState:
		if sss.End {
			if at.ParentBranchID != nil {
				// This is under an Parallel State
				e.ev <- &CompleteBranchCommand{
					BranchID: *at.ParentBranchID,
				}
				return
			} else if at.ParentIterationID != nil {
				e.ev <- &CompleteIterationCommand{
					IterationID: *at.ParentIterationID,
				}
				return
			}

			e.ev <- &CompleteExecutionCommand{
				ID: "",
			}
			return
		}

		e.ev <- &StartStateCommand{
			ActivityID:        uuid.NewString(),
			StateName:         sss.Next,
			ParentBranchID:    at.ParentBranchID,
			ParentIterationID: at.ParentIterationID,
		}
		return
	case *asl.WaitState:
		if sss.End {
			if at.ParentBranchID != nil {
				// This is under an Parallel State
				e.ev <- &CompleteBranchCommand{
					BranchID: *at.ParentBranchID,
				}
				return
			} else if at.ParentIterationID != nil {
				e.ev <- &CompleteIterationCommand{
					IterationID: *at.ParentIterationID,
				}
				return
			}

			e.ev <- &CompleteExecutionCommand{
				ID: "",
			}
			return
		}

		e.ev <- &StartStateCommand{
			ActivityID:        uuid.NewString(),
			StateName:         sss.Next,
			ParentBranchID:    at.ParentBranchID,
			ParentIterationID: at.ParentIterationID,
		}
		return
	case *asl.PassState:
		if sss.End {
			if at.ParentBranchID != nil {
				// This is under an Parallel State
				e.ev <- &CompleteBranchCommand{
					BranchID: *at.ParentBranchID,
				}
				return
			} else if at.ParentIterationID != nil {
				e.ev <- &CompleteIterationCommand{
					IterationID: *at.ParentIterationID,
				}
				return
			}

			e.ev <- &CompleteExecutionCommand{
				ID: "",
			}
			return
		}

		e.ev <- &StartStateCommand{
			ActivityID:        uuid.NewString(),
			StateName:         sss.Next,
			ParentBranchID:    at.ParentBranchID,
			ParentIterationID: at.ParentIterationID,
		}
		return
	case *asl.ChoiceState:
		e.ev <- &StartStateCommand{
			ActivityID:        uuid.NewString(),
			StateName:         *sss.Default,
			ParentBranchID:    at.ParentBranchID,
			ParentIterationID: at.ParentIterationID,
		}
	case *asl.SucceedState:
		if at.ParentBranchID != nil {
			// This is under an Parallel State
			e.ev <- &CompleteBranchCommand{
				BranchID: *at.ParentBranchID,
			}
			return
		} else if at.ParentIterationID != nil {
			e.ev <- &CompleteIterationCommand{
				IterationID: *at.ParentIterationID,
			}
			return
		}

		e.ev <- &CompleteExecutionCommand{
			ID: "",
		}
		return
	case *asl.FailState:
		if at.ParentBranchID != nil {
			// This is under an Parallel State
			e.ev <- &CompleteBranchCommand{
				BranchID: *at.ParentBranchID,
			}
			return
		} else if at.ParentIterationID != nil {
			e.ev <- &CompleteIterationCommand{
				IterationID: *at.ParentIterationID,
			}
			return
		}

		e.ev <- &CompleteExecutionCommand{
			ID: "",
		}
		return
	case *asl.ParallelState:
		if sss.End {
			if at.ParentBranchID != nil {
				// This is under an Parallel State
				e.ev <- &CompleteBranchCommand{
					BranchID: *at.ParentBranchID,
				}
				return
			} else if at.ParentIterationID != nil {
				e.ev <- &CompleteIterationCommand{
					IterationID: *at.ParentIterationID,
				}
				return
			}

			e.ev <- &CompleteExecutionCommand{
				ID: "",
			}
			return
		}

		e.ev <- &StartStateCommand{
			ActivityID:        uuid.NewString(),
			StateName:         sss.Next,
			ParentBranchID:    at.ParentBranchID,
			ParentIterationID: at.ParentIterationID,
		}
		return
	case *asl.MapState:
		if sss.End {
			if at.ParentBranchID != nil {
				// This is under an Parallel State
				e.ev <- &CompleteBranchCommand{
					BranchID: *at.ParentBranchID,
				}
				return
			} else if at.ParentIterationID != nil {
				e.ev <- &CompleteIterationCommand{
					IterationID: *at.ParentIterationID,
				}
				return
			}

			e.ev <- &CompleteExecutionCommand{
				ID: "",
			}
			return
		}

		e.ev <- &StartStateCommand{
			ActivityID:        uuid.NewString(),
			StateName:         sss.Next,
			ParentBranchID:    at.ParentBranchID,
			ParentIterationID: at.ParentIterationID,
		}
		return
	default:
		panic(fmt.Errorf("unsupport state %T", state))
	}
}

func (e *Executor) processStartBranchCommand(cmd *StartBranchCommand) {
	fmt.Printf("start branch: %v\n", cmd)
	e.ev <- &BranchStartedEvent{
		BranchID: cmd.BranchID,
	}
	e.branchContextes[cmd.BranchID] = &BranchContext{
		BranchID:   cmd.BranchID,
		Index:      cmd.Index,
		ActivityID: cmd.ActivityID,
	}
	at := e.activityContextes[cmd.ActivityID]

	state := findState(e.sm, at.StateName).(*asl.ParallelState)
	e.ev <- &StartStateCommand{
		StateName:         state.Branches[cmd.Index].StartAt,
		ActivityID:        uuid.NewString(),
		ParentBranchID:    &cmd.BranchID,
		ParentIterationID: nil,
	}
}

func (e *Executor) processCompleteBranchCommand(cmd *CompleteBranchCommand) {
	fmt.Printf("complete branch: %v\n", cmd)

	e.ev <- &BranchCompletedEvent{
		BranchID: cmd.BranchID,
	}
	at := e.activityContextes[e.branchContextes[cmd.BranchID].ActivityID]
	at.BranchStatus.Done++
	if at.BranchStatus.Done == at.BranchStatus.Max {
		e.ev <- &CompleteStateCommand{
			ActivityID: e.branchContextes[cmd.BranchID].ActivityID,
		}
	}
}

func (e *Executor) processStartIterationCommand(cmd *StartIterationCommand) {
	fmt.Printf("start iteration: %v\n", cmd)
	e.ev <- &IterationStartedEvent{
		IterationID: cmd.IterationID,
	}
	e.iterationContextes[cmd.IterationID] = &IterationContext{
		IterationID: cmd.IterationID,
		Index:       cmd.Index,
		ActivityID:  cmd.ActivityID,
	}
	at := e.activityContextes[cmd.ActivityID]

	state := findState(e.sm, at.StateName).(*asl.MapState)
	e.ev <- &StartStateCommand{
		StateName:         state.ItemProcessor.StartAt,
		ActivityID:        uuid.NewString(),
		ParentBranchID:    nil,
		ParentIterationID: &cmd.IterationID,
	}
}

func (e *Executor) processCompleteIterationCommand(cmd *CompleteIterationCommand) {
	fmt.Printf("complete iteration: %v\n", cmd)
	e.ev <- &IterationCompletedEvent{
		IterationID: cmd.IterationID,
	}
	at := e.activityContextes[e.iterationContextes[cmd.IterationID].ActivityID]
	at.IterationStatus.Done++
	if at.IterationStatus.Done == at.IterationStatus.Max {
		e.ev <- &CompleteStateCommand{
			ActivityID: e.iterationContextes[cmd.IterationID].ActivityID,
		}
	}
}

func (e *Executor) processCompleteExecutionCommand(cmd *CompleteExecutionCommand) {
	fmt.Printf("complete execution: %v\n", cmd)

	close(e.done)
}
