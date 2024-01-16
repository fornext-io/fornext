package executor

import (
	"context"
	"log/slog"

	"github.com/fornext-io/fornext/pkg/fsl"
	"github.com/google/uuid"
)

type taskStateProcessor struct {
}

func (p *taskStateProcessor) StartState(ctx context.Context, cmd *StartStateCommand, e *Executor, s *fsl.TaskState) error {
	slog.InfoContext(ctx, "start task state",
		slog.String("StateName", cmd.StateName),
		slog.String("ActivityID", cmd.ActivityID),
	)

	t := e.c.Next()
	taskInput, err := s.ApplyInput(ctx, &stateContextHolder{
		input:       cmd.Input,
		contextData: []byte{},
	})
	if err != nil {
		return err
	}

	e.ev <- &CreateTaskCommand{
		ID:         "/tenant/namespace/t/" + t.AsString(),
		ActivityID: cmd.ActivityID,
		Resource:   s.Resource,
		Input:      taskInput,
		Timestamp:  t.AsUint64(),
	}

	return nil
}

func (p *taskStateProcessor) CompleteState(ctx context.Context, cmd *CompleteStateCommand, e *Executor, s *fsl.TaskState) error {
	slog.InfoContext(ctx, "start task state",
		slog.String("ActivityID", cmd.ActivityID),
	)

	at, err := Get[ActivityContext](context.Background(), e.store, cmd.ActivityID)
	if err != nil {
		panic(err)
	}

	output, err := s.ApplyOutput(ctx, &stateContextHolder{
		input:          at.Input,
		effectiveInput: nil,
		output:         cmd.Output,
	})
	if err != nil {
		return err
	}

	if s.End {
		switch {
		case at.ParentBranchID != nil:
			e.ev <- &CompleteBranchCommand{
				BranchID: *at.ParentBranchID,
				Output:   output,
			}
		case at.ParentIterationID != nil:
			e.ev <- &CompleteIterationCommand{
				IterationID: *at.ParentIterationID,
				Output:      output,
			}
		default:
			e.ev <- &CompleteExecutionCommand{
				ID:     "",
				Output: output,
			}
		}
		return nil
	}

	e.ev <- &StartStateCommand{
		ActivityID:        uuid.NewString(), // TODO: add executionid in context
		StateName:         s.Next,
		ParentBranchID:    at.ParentBranchID,
		ParentIterationID: at.ParentIterationID,
		Input:             output,
	}
	return nil
}
