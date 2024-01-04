package executor

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/fornext-io/fornext/pkg/fsl"
	"github.com/google/uuid"
)

type mapStateProcessor struct {
}

func (p *mapStateProcessor) StartState(ctx context.Context, cmd *StartStateCommand, e *Executor, s *fsl.MapState) error {
	slog.InfoContext(ctx, "start map state",
		slog.String("StateName", cmd.StateName),
		slog.String("ActivityID", cmd.ActivityID),
	)

	at := e.activityContextes[cmd.ActivityID]
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
	// 需要处理 0 个的情况
	for i, obj := range inputObjs {
		// 再应用 ItemSelector & ItemPath
		data, err := json.Marshal(obj)
		if err != nil {
			return err
		}

		e.ev <- &StartIterationCommand{
			IterationID: uuid.NewString(),
			ActivityID:  cmd.ActivityID,
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

	at := e.activityContextes[cmd.ActivityID]
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
		ActivityID:        uuid.NewString(),
		StateName:         s.Next,
		ParentBranchID:    at.ParentBranchID,
		ParentIterationID: at.ParentIterationID,
		Input:             output,
	}

	return nil
}
