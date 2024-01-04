package executor

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"regexp"
	"time"

	"github.com/fornext-io/fornext/pkg/fsl"
	"github.com/google/uuid"
	"github.com/ohler55/ojg/jp"
)

type choiceStateProcessor struct {
}

func (p *choiceStateProcessor) StartState(ctx context.Context, cmd *StartStateCommand, e *Executor, s *fsl.ChoiceState) error {
	slog.InfoContext(ctx, "start choice state",
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

func (p *choiceStateProcessor) isMatchedAllChoiceRule(ctx context.Context, output interface{}, rules []fsl.ChoiceRule) bool {
	for _, rule := range rules {
		if !p.isMatchedChoiceRule(ctx, output, rule) {
			return false
		}
	}

	return true
}

func (p *choiceStateProcessor) isMatchedAnyChoiceRule(ctx context.Context, output interface{}, rules []fsl.ChoiceRule) bool {
	for _, rule := range rules {
		if p.isMatchedChoiceRule(ctx, output, rule) {
			return true
		}
	}

	return false
}

func getElementWithPath[T any](path fsl.Path, output interface{}) (T, error) {
	expr, err := jp.ParseString(string(path))
	if err != nil {
		return *new(T), err
	}

	result := expr.Get(output)
	if len(result) != 1 {
		return *new(T), fmt.Errorf("not found")
	}

	vv, ok := result[0].(T)
	if !ok {
		return *new(T), fmt.Errorf("wrong type")
	}

	return vv, nil
}

func variableCompare[T any](variablePath fsl.Path, output interface{}, fn func(variable T) bool) bool {
	variable, err := getElementWithPath[T](variablePath, output)
	if err != nil {
		panic(err)
	}

	return fn(variable)
}

func variableComparePathValue[T any](variablePath fsl.Path, p fsl.Path, output interface{}, fn func(variable T, pValue T) bool) bool {
	pValue, err := getElementWithPath[T](p, output)
	if err != nil {
		panic(err)
	}

	return variableCompare[T](variablePath, output, func(variable T) bool {
		return fn(variable, pValue)
	})
}

func variableTypeMatch[T any](variablePath fsl.Path, output interface{}) bool {
	expr, err := jp.ParseString(string(variablePath))
	if err != nil {
		panic(err)
	}

	result := expr.Get(output)
	if len(result) != 1 {
		return false
	}

	_, ok := result[0].(T)
	return ok
}

func variableIsNull(variablePath fsl.Path, output interface{}) bool {
	expr, err := jp.ParseString(string(variablePath))
	if err != nil {
		panic(err)
	}

	result := expr.Get(output)
	if len(result) != 1 {
		return false
	}

	return result == nil
}

func variableIsPresent(variablePath fsl.Path, output interface{}) bool {
	expr, err := jp.ParseString(string(variablePath))
	if err != nil {
		panic(err)
	}

	result := expr.Get(output)
	return len(result) != 0
}

func (p *choiceStateProcessor) isMatchedChoiceRule(ctx context.Context, output interface{}, rule fsl.ChoiceRule) bool {
	if len(rule.And) > 0 {
		return p.isMatchedAllChoiceRule(ctx, output, rule.And)
	}

	if len(rule.Or) > 0 {
		return p.isMatchedAnyChoiceRule(ctx, output, rule.Or)
	}

	if rule.Not != nil {
		return !p.isMatchedChoiceRule(ctx, output, *rule.Not)
	}

	// At now, this must be an explicit rule
	switch {
	case rule.StringEquals != nil:
		return variableCompare[string](*rule.Variable, output, func(variable string) bool {
			return variable == *rule.StringEquals
		})
	case rule.StringEqualsPath != nil:
		return variableComparePathValue[string](*rule.Variable, *rule.StringEqualsPath, output, func(variable, pValue string) bool {
			return variable == pValue
		})
	case rule.StringLessThan != nil:
		return variableCompare[string](*rule.Variable, output, func(variable string) bool {
			return variable < *rule.StringLessThan
		})
	case rule.StringLessThanPath != nil:
		return variableComparePathValue[string](*rule.Variable, *rule.StringLessThanPath, output, func(variable, pValue string) bool {
			return variable < pValue
		})
	case rule.StringGreaterThan != nil:
		return variableCompare[string](*rule.Variable, output, func(variable string) bool {
			return variable > *rule.StringGreaterThan
		})
	case rule.StringGreaterThanPath != nil:
		return variableComparePathValue[string](*rule.Variable, *rule.StringGreaterThanPath, output, func(variable, pValue string) bool {
			return variable > pValue
		})
	case rule.StringLessThanEquals != nil:
		return variableCompare[string](*rule.Variable, output, func(variable string) bool {
			return variable <= *rule.StringLessThanEquals
		})
	case rule.StringLessThanEqualsPath != nil:
		return variableComparePathValue[string](*rule.Variable, *rule.StringLessThanEqualsPath, output, func(variable, pValue string) bool {
			return variable <= pValue
		})
	case rule.StringGreaterThanEquals != nil:
		return variableCompare[string](*rule.Variable, output, func(variable string) bool {
			return variable >= *rule.StringGreaterThanEquals
		})
	case rule.StringGreaterThanEqualsPath != nil:
		return variableComparePathValue[string](*rule.Variable, *rule.StringGreaterThanEqualsPath, output, func(variable, pValue string) bool {
			return variable >= pValue
		})
	case rule.StringMatches != nil:
		return variableCompare[string](*rule.Variable, output, func(variable string) bool {
			v, err := regexp.MatchString(*rule.StringMatches, variable)
			if err != nil {
				panic(err)
			}

			return v
		})
	case rule.NumericEquals != nil:
		return variableCompare[float64](*rule.Variable, output, func(variable float64) bool {
			return variable == *rule.NumericEquals
		})
	case rule.NumericEqualsPath != nil:
		return variableComparePathValue[float64](*rule.Variable, *rule.NumericEqualsPath, output, func(variable, pValue float64) bool {
			return variable == pValue
		})
	case rule.NumericLessThan != nil:
		return variableCompare[float64](*rule.Variable, output, func(variable float64) bool {
			return variable < *rule.NumericLessThan
		})
	case rule.NumericLessThanPath != nil:
		return variableComparePathValue[float64](*rule.Variable, *rule.NumericLessThanPath, output, func(variable, pValue float64) bool {
			return variable < pValue
		})
	case rule.NumericGreaterThan != nil:
		return variableCompare[float64](*rule.Variable, output, func(variable float64) bool {
			return variable > *rule.NumericGreaterThan
		})
	case rule.NumericGreaterThanPath != nil:
		return variableComparePathValue[float64](*rule.Variable, *rule.NumericGreaterThanPath, output, func(variable, pValue float64) bool {
			return variable > pValue
		})
	case rule.NumericLessThanEquals != nil:
		return variableCompare[float64](*rule.Variable, output, func(variable float64) bool {
			return variable <= *rule.NumericLessThanEquals
		})
	case rule.NumericLessThanEqualsPath != nil:
		return variableComparePathValue[float64](*rule.Variable, *rule.NumericLessThanEqualsPath, output, func(variable, pValue float64) bool {
			return variable <= pValue
		})
	case rule.NumericGreaterThanEquals != nil:
		return variableCompare[float64](*rule.Variable, output, func(variable float64) bool {
			return variable >= *rule.NumericGreaterThanEquals
		})
	case rule.NumericGreaterThanEqualsPath != nil:
		return variableComparePathValue[float64](*rule.Variable, *rule.NumericGreaterThanEqualsPath, output, func(variable, pValue float64) bool {
			return variable >= pValue
		})
	case rule.BooleanEquals != nil:
		return variableCompare[bool](*rule.Variable, output, func(variable bool) bool {
			return variable == *rule.BooleanEquals
		})
	case rule.BooleanEqualsPath != nil:
		return variableComparePathValue[bool](*rule.Variable, *rule.BooleanEqualsPath, output, func(variable, pValue bool) bool {
			return variable == pValue
		})
	case rule.TimestampEquals != nil:
		return variableCompare[time.Time](*rule.Variable, output, func(variable time.Time) bool {
			return variable.Equal(*rule.TimestampEquals)
		})
	case rule.TimestampEqualsPath != nil:
		return variableComparePathValue[time.Time](*rule.Variable, *rule.TimestampEqualsPath, output, func(variable, pValue time.Time) bool {
			return variable.Equal(pValue)
		})
	case rule.TimestampLessThan != nil:
		return variableCompare[time.Time](*rule.Variable, output, func(variable time.Time) bool {
			return variable.Before(*rule.TimestampLessThan)
		})
	case rule.TimestampLessThanPath != nil:
		return variableComparePathValue[time.Time](*rule.Variable, *rule.TimestampLessThanPath, output, func(variable, pValue time.Time) bool {
			return variable.Before(pValue)
		})
	case rule.TimestampGreaterThan != nil:
		return variableCompare[time.Time](*rule.Variable, output, func(variable time.Time) bool {
			return variable.After(*rule.TimestampGreaterThan)
		})
	case rule.TimestampGreaterThanPath != nil:
		return variableComparePathValue[time.Time](*rule.Variable, *rule.TimestampGreaterThanPath, output, func(variable, pValue time.Time) bool {
			return variable.After(pValue)
		})
	case rule.TimestampLessThanEquals != nil:
		return variableCompare[time.Time](*rule.Variable, output, func(variable time.Time) bool {
			return variable.Before(*rule.TimestampLessThanEquals) || variable.Equal(*rule.TimestampLessThanEquals)
		})
	case rule.TimestampLessThanEqualsPath != nil:
		return variableComparePathValue[time.Time](*rule.Variable, *rule.TimestampLessThanEqualsPath, output, func(variable, pValue time.Time) bool {
			return variable.Before(pValue) || variable.Equal(pValue)
		})
	case rule.TimestampGreaterThanEquals != nil:
		return variableCompare[time.Time](*rule.Variable, output, func(variable time.Time) bool {
			return variable.After(*rule.TimestampGreaterThanEquals) || variable.Equal(*rule.TimestampGreaterThanEquals)
		})
	case rule.TimestampGreaterThanEqualsPath != nil:
		return variableComparePathValue[time.Time](*rule.Variable, *rule.TimestampGreaterThanEqualsPath, output, func(variable, pValue time.Time) bool {
			return variable.After(pValue) || variable.Equal(pValue)
		})
	case rule.IsNull != nil:
		return variableIsNull(*rule.Variable, output)
	case rule.IsPresent != nil:
		return variableIsPresent(*rule.Variable, output)
	case rule.IsNumeric != nil:
		return variableTypeMatch[float64](*rule.Variable, output)
	case rule.IsString != nil:
		return variableTypeMatch[string](*rule.Variable, output)
	case rule.IsBoolean != nil:
		return variableTypeMatch[bool](*rule.Variable, output)
	case rule.IsTimestamp != nil:
		return variableTypeMatch[time.Time](*rule.Variable, output)
	}

	return false
}

func (p *choiceStateProcessor) findNextState(ctx context.Context, output []byte, s *fsl.ChoiceState) (string, error) {
	var outputObj interface{}
	err := json.Unmarshal(output, &outputObj)
	if err != nil {
		return "", err
	}

	for _, choice := range s.Choices {
		if p.isMatchedChoiceRule(ctx, outputObj, choice.ChoiceRule) {
			return choice.Next, nil
		}
	}

	if s.Default != nil {
		return *s.Default, nil
	}

	return "", fmt.Errorf("no next state is choiced")
}

func (p *choiceStateProcessor) CompleteState(ctx context.Context, cmd *CompleteStateCommand, e *Executor, s *fsl.ChoiceState) error {
	slog.InfoContext(ctx, "complete choice state",
		slog.String("ActivityID", cmd.ActivityID),
	)

	at := e.activityContextes[cmd.ActivityID]
	// TODO: found the next state
	next, err := p.findNextState(ctx, cmd.Output, s)
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

	e.ev <- &StartStateCommand{
		ActivityID:        uuid.NewString(), // TODO: add executionid in context
		StateName:         next,
		ParentBranchID:    at.ParentBranchID,
		ParentIterationID: at.ParentIterationID,
		Input:             output,
	}
	return nil
}
