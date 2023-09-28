package asl

import (
	"context"
	"time"
)

// ChoiceState add branching logic to a state machine.
type ChoiceState struct {
	// Type is the type's name of ChoiceState
	// +Required
	// MUST be `Choice`
	Type StateType `json:"Type"`

	// Comment provided for human-readable description of the state.
	// +Optional
	Comment string `json:"Comment,omitempty"`

	// InputPath is an Path, which is applied to a State's raw input to select some
	// or all of it, that selection is used by the state.
	// +Optional
	// Defaults to '$'
	InputPath *Path `json:"InputPath,omitempty"`

	// OutputPath is an Path, which is applied to the state's output after the application of `ResultPath`,
	// producing the effective output which serves as the raw input for the next state.
	// +Optional
	// Defaults to '$'
	OutputPath *Path `json:"OutputPath,omitempty"`

	// ------ State specified field

	// Choices is an array of rules that determines which state transitions to next.
	// +Required
	// MUST non-empty
	Choices []Choice `json:"Choices"`

	// Default is the name of state to transition to if none of the transitions in `Choices` is taken.
	// +Optional
	Default *string `json:"Default,omitempty"`
}

// Choice is a JSON object may be evaluated to boolean value.
type Choice struct {
	// Next is the state name transition to.
	// +Required
	Next string `json:"Next"`

	// ChoiceRule is the rule for this choice.
	ChoiceRule `json:",inline"`
}

// ChoiceRule may be evaluated to return a boolean value.
type ChoiceRule struct {
	// ------ Boolean expression

	// And is an non-empty array of ChoiceRule.
	// +Optional
	And []ChoiceRule `json:"And,omitempty"`

	// Or is an non-empty array of ChoiceRule.
	// +optional
	Or []ChoiceRule `json:"Or,omitempty"`

	// Not is a single ChoiceRule.
	// +Optional
	Not *ChoiceRule `json:"Not,omitempty"`

	// ------ Data-test expression

	// Variable MUST be a Path to field.
	// +Optional
	Variable *Path `json:"Variable,omitempty"`

	// StringEquals is the value field comparison to.
	// +Optional
	StringEquals *string `json:"StringEquals,omitempty"`
	// StringEqualsPath is an Refercence Path to the StringEquals data.
	// +Optional
	StringEqualsPath *Path `json:"StringEqualsPath,omitempty"`

	// StringLessThan is the value field comparison to.
	// +Optional
	StringLessThan *string `json:"StringLessThan,omitempty"`
	// StringLessThanPath is an Refercence Path to the StringLessThan data.
	// +Optional
	StringLessThanPath *Path `json:"StringLessThanPath,omitempty"`

	// StringGreaterThan is the value field comparison to.
	// +Optional
	StringGreaterThan *string `json:"StringGreaterThan,omitempty"`
	// StringGreaterThanPath is an Refercence Path to the StringGreaterThan data.
	// +Optional
	StringGreaterThanPath *Path `json:"StringGreaterThanPath,omitempty"`

	// StringLessThanEquals is the value field comparison to.
	// +Optional
	StringLessThanEquals *string `json:"StringLessThanEquals,omitempty"`
	// StringLessThanEqualsPath is an Refercence Path to the StringLessThanEquals data.
	// +Optional
	StringLessThanEqualsPath *Path `json:"StringLessThanEqualsPath,omitempty"`

	// StringGreaterThanEquals is the value field comparison to.
	// +Optional
	StringGreaterThanEquals *string `json:"StringGreaterThanEquals,omitempty"`
	// StringGreaterThanEqualsPath is an Refercence Path to the StringGreaterThanEquals data.
	// +Optional
	StringGreaterThanEqualsPath *Path `json:"StringGreaterThanEqualsPath,omitempty"`

	// StringMatches MUST be a string which MAY contain one or more '*' characters.
	// +Optional
	StringMatches *string `json:"StringMatches,omitempty"`

	// NumericEquals is the value field comparison to.
	// +Optional
	NumericEquals *float64 `json:"NumericEquals,omitempty"`
	// NumericEqualsPath is an Refercence Path to the NumericEquals data.
	// +Optional
	NumericEqualsPath *Path `json:"NumericEqualsPath,omitempty"`

	// NumericLessThan is the value field comparison to.
	// +Optional
	NumericLessThan *float64 `json:"NumericLessThan,omitempty"`
	// NumericLessThanPath is an Refercence Path to the NumericLessThan data.
	// +Optional
	NumericLessThanPath *Path `json:"NumericLessThanPath,omitempty"`

	// NumericGreaterThan is the value field comparison to.
	// +Optional
	NumericGreaterThan *float64 `json:"NumericGreaterThan,omitempty"`
	// NumericGreaterThanPath is an Refercence Path to the NumericGreaterThan data.
	// +Optional
	NumericGreaterThanPath *Path `json:"NumericGreaterThanPath,omitempty"`

	// NumericLessThanEquals is the value field comparison to.
	// +Optional
	NumericLessThanEquals *float64 `json:"NumericLessThanEquals,omitempty"`
	// NumericLessThanEqualsPath is an Refercence Path to the NumericLessThanEquals data.
	// +Optional
	NumericLessThanEqualsPath *Path `json:"NumericLessThanEqualsPath,omitempty"`

	// NumericGreaterThanEquals is the value field comparison to.
	// +Optional
	NumericGreaterThanEquals *float64 `json:"NumericGreaterThanEquals,omitempty"`
	// NumericGreaterThanEqualsPath is an Refercence Path to the NumericGreaterThanEquals data.
	// +Optional
	NumericGreaterThanEqualsPath *Path `json:"NumericGreaterThanEqualsPath,omitempty"`

	// BooleanEquals is the value field comparison to.
	// +Optional
	BooleanEquals *bool `json:"BooleanEquals,omitempty"`
	// BooleanEqualsPath is an Refercence Path to the BooleanEquals data.
	// +Optional
	BooleanEqualsPath *Path `json:"BooleanEqualsPath,omitempty"`

	// TimestampEquals is the value field comparison to.
	// +Optional
	TimestampEquals *time.Time `json:"TimestampEquals,omitempty"`
	// TimestampEqualsPath is an Refercence Path to the TimestampEquals data.
	// +Optional
	TimestampEqualsPath *Path `json:"TimestampEqualsPath,omitempty"`

	// TimestampLessThan is the value field comparison to.
	// +Optional
	TimestampLessThan *time.Time `json:"TimestampLessThan,omitempty"`
	// TimestampLessThanPath is an Refercence Path to the TimestampLessThan data.
	// +Optional
	TimestampLessThanPath *Path `json:"TimestampLessThanPath,omitempty"`

	// TimestampGreaterThan is the value field comparison to.
	// +Optional
	TimestampGreaterThan *time.Time `json:"TimestampGreaterThan,omitempty"`
	// TimestampGreaterThanPath is an Refercence Path to the TimestampGreaterThan data.
	// +Optional
	TimestampGreaterThanPath *Path `json:"TimestampGreaterThanPath,omitempty"`

	// TimestampLessThanEquals is the value field comparison to.
	// +Optional
	TimestampLessThanEquals *time.Time `json:"TimestampLessThanEquals,omitempty"`
	// TimestampLessThanEqualsPath is an Refercence Path to the TimestampLessThanEquals data.
	// +Optional
	TimestampLessThanEqualsPath *Path `json:"TimestampLessThanEqualsPath,omitempty"`

	// TimestampGreaterThanEquals is the value field comparison to.
	// +Optional
	TimestampGreaterThanEquals *time.Time `json:"TimestampGreaterThanEquals,omitempty"`
	// TimestampGreaterThanEqualsPath is an Refercence Path to the TimestampGreaterThanEquals data.
	// +Optional
	TimestampGreaterThanEqualsPath *Path `json:"TimestampGreaterThanEqualsPath,omitempty"`

	// IsNull comparison the value `null`
	// +Optional
	IsNull *bool `json:"IsNull,omitempty"`

	// IsPresent check the Variable-field Path is exists.
	// +Optional
	IsPresent *bool `json:"IsPresent,omitempty"`

	// IsNumeric check the value if numeric.
	// +Optional
	IsNumeric *bool `json:"IsNumeric,omitempty"`

	// IsString check the value is string.
	// +Optional
	IsString *bool `json:"IsString,omitempty"`

	// IsBoolean check the value is boolean.
	// +Optional
	IsBoolean *bool `json:"IsBoolean,omitempty"`

	// IsTimestamp check the value is timestamp.
	// +Optional
	IsTimestamp *bool `json:"IsTimestamp,omitempty"`
}

// Validate ChoiceState will validate the ChoiceState configuration
func (s *ChoiceState) Validate(_ context.Context) error {
	return nil
}
