package executor

import (
	"context"
	"log/slog"

	"github.com/fornext-io/fornext/pkg/fsl"
)

type failStateProcessor struct {
}

func (p *failStateProcessor) StartState(ctx context.Context, cmd *StartStateCommand, e *Executor, s *fsl.FailState) error {
	slog.InfoContext(ctx, "start fail state",
		slog.String("StateName", cmd.StateName),
		slog.String("ActivityID", cmd.ActivityID),
	)

	t := e.c.Next()
	e.ev <- &CompleteStateCommand{
		ActivityID: cmd.ActivityID,
		Output:     cmd.Input,
		Timestamp:  t.AsUint64(),
	}

	return nil
}

func (p *failStateProcessor) CompleteState(ctx context.Context, cmd *CompleteStateCommand, e *Executor, s *fsl.FailState) error {
	slog.InfoContext(ctx, "complete fail state",
		slog.String("ActivityID", cmd.ActivityID),
	)

	// 1. find the next state
	at, err := Get[ActivityContext](context.Background(), e.store, cmd.ActivityID)
	if err != nil {
		panic(err)
	}

	switch {
	case at.ParentBranchID != nil:
		e.ev <- &CompleteBranchCommand{
			BranchID: *at.ParentBranchID,
			Output:   cmd.Output,
		}
	case at.ParentIterationID != nil:
		e.ev <- &CompleteIterationCommand{
			IterationID: *at.ParentIterationID,
			Output:      cmd.Output,
		}
	default:
		e.ev <- &CompleteExecutionCommand{
			ID:     "",
			Output: cmd.Output,
		}
	}

	return nil
}
