package asl

import (
	"context"
	"fmt"
	"time"
)

func RunStateMachine(ctx context.Context, sm *StateMachine) {
	s := sm.States[sm.StartAt]
	runState(ctx, sm.States, s, sm.StartAt)
}

func runState(ctx context.Context, states States, s State, n string) {
	fmt.Printf("run state %s\n", n)

	switch ss := s.(type) {
	case *WaitState:
		time.Sleep(time.Duration(*ss.Seconds))
		if ss.End {
			return
		}

		runState(ctx, states, states[ss.Next], ss.Next)
	case *TaskState:
		if ss.End {
			return
		}

		runState(ctx, states, states[ss.Next], ss.Next)
	case *SucceedState:
		return
	case *FailState:
		return
	case *PassState:
		if ss.End {
			return
		}

		runState(ctx, states, states[ss.Next], ss.Next)
	case *ChoiceState:
		runState(ctx, states, states[*ss.Default], *ss.Default)
	case *MapState:
		runState(ctx, ss.ItemProcessor.States, ss.ItemProcessor.States[ss.ItemProcessor.StartAt], ss.ItemProcessor.StartAt)
		if ss.End {
			return
		}

		runState(ctx, states, states[ss.Next], ss.Next)
	case *ParallelState:
		for _, branch := range ss.Branches {
			runState(ctx, branch.States, branch.States[branch.StartAt], branch.StartAt)
		}
		if ss.End {
			return
		}

		runState(ctx, states, states[ss.Next], ss.Next)
	}
}
