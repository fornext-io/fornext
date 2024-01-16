package executor

import (
	"context"
	"log/slog"

	"github.com/fornext-io/fornext/pkg/fsl"
)

type succeedStateProcessor struct {
}

func (p *succeedStateProcessor) StartState(ctx context.Context, cmd *StartStateCommand, e *Executor, s *fsl.SucceedState) error {
	slog.InfoContext(ctx, "start succeed state",
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

	e.ev <- &CompleteStateCommand{
		ActivityID: cmd.ActivityID,
		Output:     taskInput,
		Timestamp:  t.AsUint64(),
	}

	return nil
}

func (p *succeedStateProcessor) CompleteState(ctx context.Context, cmd *CompleteStateCommand, e *Executor, s *fsl.SucceedState) error {
	slog.InfoContext(ctx, "complete succeed state",
		slog.String("ActivityID", cmd.ActivityID),
	)

	// 1. find the next state
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
