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

	activityContextes  map[string]*ActivityContext
	branchContextes    map[string]*BranchContext
	iterationContextes map[string]*IterationContext
	taskContextes      map[string]*TaskContext

	done chan interface{}
	ev   chan interface{}

	c *hybridLogicalClock
}

// NewExecutor ...
func NewExecutor(sm *fsl.StateMachine) *Executor {
	v := &Executor{
		sm:                 sm,
		done:               make(chan interface{}),
		ev:                 make(chan interface{}, 100),
		activityContextes:  map[string]*ActivityContext{},
		branchContextes:    map[string]*BranchContext{},
		iterationContextes: map[string]*IterationContext{},
		taskContextes:      map[string]*TaskContext{},
		c:                  newHybridLogicalClock(),
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

	e.activityContextes[cmd.ActivityID] = &ActivityContext{
		ID:                cmd.ActivityID,
		StateName:         cmd.StateName,
		ParentBranchID:    cmd.ParentBranchID,
		ParentIterationID: cmd.ParentIterationID,
		Input:             cmd.Input,
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

func (e *Executor) processCreateTaskCommand(cmd *CreateTaskCommand) {
	e.ev <- &TaskCreatedEvent{
		ID: cmd.ID,
	}
	e.taskContextes[cmd.ID] = &TaskContext{
		ActivityID: cmd.ActivityID,
	}

	go func() {
		if cmd.Resource == "__sleep" {
			time.Sleep(5 * time.Second)
			e.ev <- &CompleteTaskCommand{
				TaskID: cmd.ID,
				Output: cmd.Input,
			}
		} else {
			e.ev <- &CompleteTaskCommand{
				TaskID: cmd.ID,
				Output: cmd.Input,
			}
		}
	}()
}

func (e *Executor) processActivateTaskCommand(_ *ActivateTaskCommand) {

}

func (e *Executor) processCompleteTaskCommand(cmd *CompleteTaskCommand) {
	fmt.Printf("complete task command: %+v\n", cmd)

	taskCtx := e.taskContextes[cmd.TaskID]

	e.ev <- &TaskCompletedEvent{
		TaskID: cmd.TaskID,
	}
	e.ev <- &CompleteStateCommand{
		ActivityID: taskCtx.ActivityID,
		Output:     cmd.Output,
	}
}

// 本地状态在什么时候修改？
//  1. 发送消息之前，那么可能出现消息没有发送成功，需要回滚（其实不会出现这个问题，因为 command 是 determin 的）
//  2. 关键是什么时候回复给客户？如果是 determin 的，那么无所谓
//  3. 收到一个 cmd，然后记录到本地存储中（标记为还没有处理），然后发送给 scheduler + processor，调度（不同的 execution 可以并行）
//     然后执行并记录本地状态，发送新的 command，删除旧的 command，处理下一个（不用等待发送成功？应该无需，可以直接开始处理下一个，不然要等待发送完成才能进行下一个的处理，延迟会很高；而且这里不等待是没有问题的，因为 cmd 是 determin 的，假设有新的 leader 启动，得到的状态也会是一致的），但是这里删除、修改是需要 transaction 的，避免修改成功，删除没成功，导致重启时 cmd 的重复处理
//     另外，可能出现本地应用成功，然后重启重新成为 leader，但是 cmd & event 没有发送成功的情况，则需要将要发送的 cmd 保存到 kv，然后重新成为 leader 的时候重发
//     对于 follower，则在收到 event 之后删除旧的 command
func (e *Executor) processCompleteStateCommand(cmd *CompleteStateCommand) {
	slog.InfoContext(context.Background(), "complete state command", slog.Any("Command", cmd), slog.String("Output", string(cmd.Output)))

	e.ev <- &StateCompletedEvent{
		ActivityID: cmd.ActivityID,
	}
	state := findState(e.sm, e.activityContextes[cmd.ActivityID].StateName)

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
	e.branchContextes[cmd.BranchID] = &BranchContext{
		BranchID:   cmd.BranchID,
		Index:      cmd.Index,
		ActivityID: cmd.ActivityID,
		Input:      cmd.Input,
	}
	at := e.activityContextes[cmd.ActivityID]

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
	at := e.activityContextes[e.branchContextes[cmd.BranchID].ActivityID]
	at.BranchStatus.Done++
	// 此处还需要考虑顺序问题
	at.BranchStatus.Output = append(at.BranchStatus.Output, cmd.Output)

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
			ActivityID: e.branchContextes[cmd.BranchID].ActivityID,
			Output:     outputBytes,
		}
	}
}

func (e *Executor) processStartIterationCommand(cmd *StartIterationCommand) {
	e.ev <- &IterationStartedEvent{
		IterationID: cmd.IterationID,
	}
	e.iterationContextes[cmd.IterationID] = &IterationContext{
		IterationID: cmd.IterationID,
		Index:       cmd.Index,
		ActivityID:  cmd.ActivityID,
	}
	at := e.activityContextes[cmd.ActivityID]

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
	at := e.activityContextes[e.iterationContextes[cmd.IterationID].ActivityID]
	at.IterationStatus.Done++
	at.IterationStatus.Output = append(at.IterationStatus.Output, cmd.Output)

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
			ActivityID: e.iterationContextes[cmd.IterationID].ActivityID,
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
