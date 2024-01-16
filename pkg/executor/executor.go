// Package executor ...
package executor

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/fornext-io/fornext/pkg/fsl"
	"github.com/fornext-io/fornext/pkg/utils"
	"github.com/google/uuid"
)

type stateContextHolder struct {
	input           []byte
	contextData     []byte
	effectiveInput  []byte
	output          []byte
	effectiveOutput []byte
}

var (
	_ fsl.StateContext = (*stateContextHolder)(nil)
)

func (h *stateContextHolder) Input() []byte {
	return h.input
}

func (h *stateContextHolder) ContextData() []byte {
	return h.contextData
}

func (h *stateContextHolder) EffectiveInput() []byte {
	return h.effectiveInput
}

func (h *stateContextHolder) Output() []byte {
	return h.output
}

func (h *stateContextHolder) EffectiveOutput() []byte {
	return h.effectiveOutput
}

// Executor ...
type Executor struct {
	sm  *fsl.StateMachine
	ctx *ExecutionContext

	// iterationContextes map[string]*IterationContext
	store *Storage

	done chan interface{}
	ev   chan interface{}

	c            *hybridLogicalClock
	taskHandlers map[string]func(cmd *CreateTaskCommand) []byte
}

// NewExecutor ...
func NewExecutor(sm *fsl.StateMachine, handlers map[string]func(*CreateTaskCommand) []byte) *Executor {
	store, err := NewStorage("./fornext")
	if err != nil {
		panic(err)
	}

	v := &Executor{
		sm:   sm,
		done: make(chan interface{}),
		ev:   make(chan interface{}, 100),
		// iterationContextes: map[string]*IterationContext{},
		store:        store,
		c:            newHybridLogicalClock(),
		taskHandlers: map[string]func(*CreateTaskCommand) []byte{},
	}
	v.taskHandlers["__sleep"] = v.executeSleepTask

	for k, vv := range handlers {
		v.taskHandlers[k] = vv
	}

	return v
}

// Run ...
func (e *Executor) Run(input []byte) *ExecutionContext {
	t := e.c.Next()
	e.ev <- &StartExecutionCommand{
		ID:         "/tenant/namespace/e/" + t.AsString(),
		WorkflowID: "example",
		Input:      input,
		Timestamp:  t.AsUint64(),
	}
	go e.run()

	return e.ctx
}

// WaitExecutionDone ...
func (e *Executor) WaitExecutionDone() *ExecutionContext {
	<-e.done
	return e.ctx
}

func findState(sm *fsl.StateMachine, stateName string) fsl.State {
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

func findStateInState(ss fsl.State, stateName string) fsl.State {
	switch sss := ss.(type) {
	case *fsl.ParallelState:
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
	case *fsl.MapState:
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
				ID:         evv.ID,
				Input:      evv.Input,
				Status:     "Running",
				Timestamp:  evv.Timestamp,
				WorkflowID: evv.WorkflowID,
			}
			t := e.c.Next()
			e.ev <- &StartStateCommand{
				ActivityID: evv.ID + "/" + t.AsString(),
				StateName:  e.sm.StartAt,
				Input:      evv.Input,
				Timestamp:  t.AsUint64(),
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
			slog.InfoContext(context.Background(), "receieve event", slog.Any("Event", ev))
		}
	}
}

func (e *Executor) processStartStateCommand(cmd *StartStateCommand) {
	slog.InfoContext(context.Background(), "start state", slog.Any("Command", cmd), slog.String("Input", string(cmd.Input)))

	e.ev <- &StateStartedEvent{
		ActivityID: cmd.ActivityID,
	}

	err := Set(context.Background(), e.store, cmd.ActivityID, &ActivityContext{
		ID:                cmd.ActivityID,
		StateName:         cmd.StateName,
		ParentBranchID:    cmd.ParentBranchID,
		ParentIterationID: cmd.ParentIterationID,
		Input:             cmd.Input,
	})
	if err != nil {
		panic(err)
	}

	state := findState(e.sm, cmd.StateName)
	switch sss := state.(type) {
	case *fsl.TaskState:
		err := (&taskStateProcessor{}).StartState(context.Background(), cmd, e, sss)
		if err != nil {
			panic(err)
		}
	case *fsl.WaitState:
		err := (&waitStateProcessor{}).StartState(context.Background(), cmd, e, sss)
		if err != nil {
			panic(err)
		}
	case *fsl.PassState:
		err := (&paasStateProcessor{}).StartState(context.Background(), cmd, e, sss)
		if err != nil {
			panic(err)
		}

	case *fsl.ChoiceState:
		err := (&choiceStateProcessor{}).StartState(context.Background(), cmd, e, sss)
		if err != nil {
			panic(err)
		}
	case *fsl.SucceedState:
		err := (&succeedStateProcessor{}).StartState(context.Background(), cmd, e, sss)
		if err != nil {
			panic(err)
		}
	case *fsl.FailState:
		err := (&failStateProcessor{}).StartState(context.Background(), cmd, e, sss)
		if err != nil {
			panic(err)
		}
	case *fsl.ParallelState:
		err := (&parallelStateProcessor{}).StartState(context.Background(), cmd, e, sss)
		if err != nil {
			panic(err)
		}
	case *fsl.MapState:
		err := (&mapStateProcessor{}).StartState(context.Background(), cmd, e, sss)
		if err != nil {
			panic(err)
		}
	default:
		panic(fmt.Errorf("unsupport state %T", state))
	}
}

func (e *Executor) executeTask(cmd *CreateTaskCommand) {
	go func() {
		h := e.taskHandlers[cmd.Resource]
		if h == nil {
			panic(fmt.Errorf("cannot find handler for task %v", cmd.Resource))
		}

		output := h(cmd)
		e.ev <- &CompleteTaskCommand{
			TaskID: cmd.ID,
			Output: output,
		}
	}()
}

func (e *Executor) executeSleepTask(cmd *CreateTaskCommand) []byte {
	time.Sleep(5 * time.Second)
	return cmd.Input
}

func (e *Executor) processCreateTaskCommand(cmd *CreateTaskCommand) {
	e.ev <- &TaskCreatedEvent{
		ID: cmd.ID,
	}
	err := Set(context.Background(), e.store, cmd.ID, &TaskContext{
		ActivityID: cmd.ActivityID,
	})
	if err != nil {
		panic(err)
	}

	e.executeTask(cmd)
}

func (e *Executor) processActivateTaskCommand(_ *ActivateTaskCommand) {

}

func (e *Executor) processCompleteTaskCommand(cmd *CompleteTaskCommand) {
	fmt.Printf("complete task command: %+v\n", cmd)

	taskCtx, err := Get[TaskContext](context.Background(), e.store, cmd.TaskID)
	if err != nil {
		panic(err)
	}

	e.ev <- &TaskCompletedEvent{
		TaskID: cmd.TaskID,
	}
	e.ev <- &CompleteStateCommand{
		ActivityID: taskCtx.ActivityID,
		Output:     cmd.Output,
	}
}

func (e *Executor) processCompleteStateCommand(cmd *CompleteStateCommand) {
	slog.InfoContext(context.Background(), "complete state command", slog.Any("Command", cmd), slog.String("Output", string(cmd.Output)))

	e.ev <- &StateCompletedEvent{
		ActivityID: cmd.ActivityID,
	}
	act, err := Get[ActivityContext](context.Background(), e.store, cmd.ActivityID)
	if err != nil {
		panic(err)
	}

	state := findState(e.sm, act.StateName)
	switch sss := state.(type) {
	case *fsl.TaskState:
		err := (&taskStateProcessor{}).CompleteState(context.Background(), cmd, e, sss)
		if err != nil {
			panic(err)
		}
	case *fsl.WaitState:
		err := (&waitStateProcessor{}).CompleteState(context.Background(), cmd, e, sss)
		if err != nil {
			panic(err)
		}
	case *fsl.PassState:
		err := (&paasStateProcessor{}).CompleteState(context.Background(), cmd, e, sss)
		if err != nil {
			panic(err)
		}
	case *fsl.ChoiceState:
		err := (&choiceStateProcessor{}).CompleteState(context.Background(), cmd, e, sss)
		if err != nil {
			panic(err)
		}
	case *fsl.SucceedState:
		err := (&succeedStateProcessor{}).CompleteState(context.Background(), cmd, e, sss)
		if err != nil {
			panic(err)
		}
	case *fsl.FailState:
		err := (&failStateProcessor{}).CompleteState(context.Background(), cmd, e, sss)
		if err != nil {
			panic(err)
		}
	case *fsl.ParallelState:
		err := (&parallelStateProcessor{}).CompleteState(context.Background(), cmd, e, sss)
		if err != nil {
			panic(err)
		}
	case *fsl.MapState:
		err := (&mapStateProcessor{}).CompleteState(context.Background(), cmd, e, sss)
		if err != nil {
			panic(err)
		}
	default:
		panic(fmt.Errorf("unsupport state %T", state))
	}
}

func (e *Executor) processStartBranchCommand(cmd *StartBranchCommand) {
	e.ev <- &BranchStartedEvent{
		BranchID: cmd.BranchID,
	}
	err := Set(context.Background(), e.store, cmd.BranchID, &BranchContext{
		BranchID:   cmd.BranchID,
		Index:      cmd.Index,
		ActivityID: cmd.ActivityID,
		Input:      cmd.Input,
	})
	if err != nil {
		panic(err)
	}

	at, err := Get[ActivityContext](context.Background(), e.store, cmd.ActivityID)
	if err != nil {
		panic(err)
	}

	state := findState(e.sm, at.StateName).(*fsl.ParallelState)
	e.ev <- &StartStateCommand{
		StateName:         state.Branches[cmd.Index].StartAt,
		ActivityID:        uuid.NewString(),
		ParentBranchID:    &cmd.BranchID,
		ParentIterationID: nil,
		Input:             cmd.Input,
	}
}

func (e *Executor) processCompleteBranchCommand(cmd *CompleteBranchCommand) {
	e.ev <- &BranchCompletedEvent{
		BranchID: cmd.BranchID,
	}
	branchCtx, err := Get[BranchContext](context.Background(), e.store, cmd.BranchID)
	if err != nil {
		panic(err)
	}

	at, err := Get[ActivityContext](context.Background(), e.store, branchCtx.ActivityID)
	if err != nil {
		panic(err)
	}

	at.BranchStatus.Done++
	// 此处还需要考虑顺序问题
	at.BranchStatus.Output = append(at.BranchStatus.Output, cmd.Output)
	err = Set(context.Background(), e.store, branchCtx.ActivityID, at)
	if err != nil {
		panic(err)
	}

	if at.BranchStatus.Done == at.BranchStatus.Max {
		// 拼接所有的 output 为列表
		var output []interface{}
		for _, oo := range at.BranchStatus.Output {
			output = append(output, utils.MustUnmarshalJSON(oo))
		}
		outputBytes, err := json.Marshal(output)
		if err != nil {
			panic(err)
		}

		e.ev <- &CompleteStateCommand{
			ActivityID: branchCtx.ActivityID,
			Output:     outputBytes,
		}
	}
}

func (e *Executor) processStartIterationCommand(cmd *StartIterationCommand) {
	e.ev <- &IterationStartedEvent{
		IterationID: cmd.IterationID,
	}
	err := Set(context.Background(), e.store, cmd.IterationID, &IterationContext{
		IterationID: cmd.IterationID,
		Index:       cmd.Index,
		ActivityID:  cmd.ActivityID,
	})
	if err != nil {
		panic(err)
	}

	at, err := Get[ActivityContext](context.Background(), e.store, cmd.ActivityID)
	if err != nil {
		panic(err)
	}

	state := findState(e.sm, at.StateName).(*fsl.MapState)
	e.ev <- &StartStateCommand{
		StateName:         state.ItemProcessor.StartAt,
		ActivityID:        uuid.NewString(),
		ParentBranchID:    nil,
		ParentIterationID: &cmd.IterationID,
		Input:             cmd.Input,
	}
}

func (e *Executor) processCompleteIterationCommand(cmd *CompleteIterationCommand) {
	e.ev <- &IterationCompletedEvent{
		IterationID: cmd.IterationID,
	}
	iterCtx, err := Get[IterationContext](context.Background(), e.store, cmd.IterationID)
	if err != nil {
		panic(err)
	}

	at, err := Get[ActivityContext](context.Background(), e.store, iterCtx.ActivityID)
	if err != nil {
		panic(err)
	}

	at.IterationStatus.Done++
	at.IterationStatus.Output = append(at.IterationStatus.Output, cmd.Output)
	err = Set(context.Background(), e.store, iterCtx.ActivityID, at)
	if err != nil {
		panic(err)
	}

	if at.IterationStatus.Done == at.IterationStatus.Max {
		// 拼接所有的 output 为列表
		var output []interface{}
		for _, oo := range at.IterationStatus.Output {
			output = append(output, utils.MustUnmarshalJSON(oo))
		}
		outputBytes, err := json.Marshal(output)
		if err != nil {
			panic(err)
		}

		e.ev <- &CompleteStateCommand{
			ActivityID: iterCtx.ActivityID,
			Output:     outputBytes,
		}
	}
}

func (e *Executor) processCompleteExecutionCommand(cmd *CompleteExecutionCommand) {
	slog.InfoContext(context.Background(),
		"complete execution",
		slog.Any("command", cmd),
		slog.String("output", string(cmd.Output)),
	)

	close(e.done)
}
