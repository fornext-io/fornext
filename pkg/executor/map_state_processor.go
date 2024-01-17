package executor

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/fornext-io/fornext/pkg/fsl"
)

type mapStateProcessor struct {
}

func (p *mapStateProcessor) StartState(ctx context.Context, cmd *StartStateCommand, e *Executor, s *fsl.MapState) error {
	slog.InfoContext(ctx, "start map state",
		slog.String("StateName", cmd.StateName),
		slog.String("ActivityID", cmd.ActivityID),
	)

	at, err := Get[ActivityContext](context.Background(), e.store, cmd.ActivityID)
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

	var inputObjs []interface{}
	err = json.Unmarshal(input, &inputObjs)
	if err != nil {
		return err
	}

	at.IterationStatus = &ActivityIterationStatus{
		Max:  len(inputObjs),
		Done: 0,
	}
	err = Set(context.Background(), e.store, cmd.ActivityID, at)
	if err != nil {
		panic(err)
	}

	// 需要处理 0 个的情况
	for i, obj := range inputObjs {
		// 再应用 ItemSelector & ItemPath
		data, err := json.Marshal(obj)
		if err != nil {
			return err
		}

		e.ev <- &StartIterationCommand{
			IterationID: fmt.Sprintf("%s/i%v", cmd.ActivityID, i),
			ActivityID:  cmd.ActivityID,
			ExecutionID: cmd.ExecutionID,
			Index:       i,
			Input:       data,
		}
	}

	return nil
}

func (p *mapStateProcessor) CompleteState(ctx context.Context, cmd *CompleteStateCommand, e *Executor, s *fsl.MapState) error {
	slog.InfoContext(ctx, "complete map state",
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
		contextData:    []byte(`{}`),
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
