package executor

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/fornext-io/fornext/pkg/fsl"
)

type parallelStateProcessor struct {
}

func (p *parallelStateProcessor) StartState(ctx context.Context, cmd *StartStateCommand, e *Executor, s *fsl.ParallelState) error {
	slog.InfoContext(ctx, "start parallel state",
		slog.String("StateName", cmd.StateName),
		slog.String("ActivityID", cmd.ActivityID),
	)

	at, err := Get[ActivityContext](context.Background(), e.store, cmd.ActivityID)
	if err != nil {
		panic(err)
	}

	at.BranchStatus = &ActivityBranchStatus{
		Max:  len(s.Branches),
		Done: 0,
	}
	err = Set[*ActivityContext](context.Background(), e.store, cmd.ActivityID, at)
	if err != nil {
		panic(err)
	}

	input, err := s.ApplyInput(ctx, &stateContextHolder{
		input:       cmd.Input,
		contextData: []byte{},
	})
	if err != nil {
		return err
	}

	for i := range s.Branches {
		e.ev <- &StartBranchCommand{
			BranchID:    fmt.Sprintf("%s/b%v", cmd.ActivityID, i),
			ActivityID:  cmd.ActivityID,
			ExecutionID: cmd.ExecutionID,
			Index:       i,
			Input:       input,
		}
	}

	return nil
}

func (p *parallelStateProcessor) CompleteState(ctx context.Context, cmd *CompleteStateCommand, e *Executor, s *fsl.ParallelState) error {
	slog.InfoContext(ctx, "complete parallel state",
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
		ActivityID:        at.ExecutionID + "/" + e.c.Next().AsString(),
		StateName:         s.Next,
		ExecutionID:       at.ExecutionID,
		ParentBranchID:    at.ParentBranchID,
		ParentIterationID: at.ParentIterationID,
		Input:             output,
	}

	return nil
}
