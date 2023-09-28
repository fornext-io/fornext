package asl

import (
	"context"

	"github.com/lsytj0413/ena/xerrors"
)

// Retrier represent when State's reports an error the interpreter can do
type Retrier struct {
	// ErrorEquals is a non-empty array of Strings, which match Error Names.
	// +Required
	ErrorEquals []string `json:"ErrorEquals"`

	// IntervalSeconds representing the number of seconds before the first retry attempt.
	// +Optional
	// Defaults to 1
	IntervalSeconds *int `json:"IntervalSeconds"`

	// MaxAttempts representing the maximum number of retry attempts.
	// +Optional
	// Defaults to 3
	MaxAttempts *int `json:"MaxAttempts"`

	// BackoffRate is a number which is the multiplier that increases the retry interval on each attempt.
	// +Optional
	// >= 1.0
	// default 2.0
	BackoffRate *float32 `json:"BackoffRate"`
}

// Validate will validate the Retirer configuration
func (r *Retrier) Validate(_ context.Context) error {
	err := ValidateErrorEquals(r.ErrorEquals)
	if err != nil {
		return err
	}

	if r.IntervalSeconds != nil && *r.IntervalSeconds <= 0 {
		return xerrors.New("IntervalSeconds must be an positive integer")
	}

	if r.MaxAttempts != nil && *r.MaxAttempts < 0 {
		return xerrors.New("MaxAttempts must be a non-negative integer")
	}

	if r.BackoffRate != nil && *r.BackoffRate < 1.0 {
		return xerrors.New("BackoffRate must be greater than or equal to 1.0")
	}

	return nil
}

// IsAnyErrorWildcard represent where is Catcher will matches any Error Names
func (r *Retrier) IsAnyErrorWildcard() bool {
	return IsAnyErrorWildcard(r.ErrorEquals)
}

// Retriers must be an array of Retrier.
type Retriers []Retrier

// Validate will validate the Catchers configuration
func (c Retriers) Validate(ctx context.Context) error {
	for i, retrier := range c {
		if err := retrier.Validate(ctx); err != nil {
			return err
		}

		if i != len(c)-1 && retrier.IsAnyErrorWildcard() {
			return xerrors.New("States.ALL must be the last retrier")
		}
	}

	return nil
}
