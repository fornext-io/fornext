package fsl

import (
	"context"
	"encoding/json"

	"github.com/lsytj0413/ena/xerrors"
)

// Catcher represent transitions when State's reports an error
// and there is no Retrier or retries have failed to resolve the error.
// When the error appears in the value of ErrroEquals field, transition
// the machine to the state named in the value of the 'Next' field.
type Catcher struct {
	// ErrorEquals is a non-empty array of Strings, which match Error Names.
	// +Required
	ErrorEquals []string `json:"ErrorEquals"`

	// Next is a string exactly matching a State Name.
	// +Required
	Next string `json:"Next"`

	// ResultPath must be a Reference Path, which specifies the
	// raw input's combination with or replacement by the state's result.
	// +Optional
	// Defaults to $
	ResultPath string `json:"ResultPath"`
}

// UnmarshalJSON implement json.Unmarshaler and will set defaults
// value for Catcher.
func (c *Catcher) UnmarshalJSON(data []byte) error {
	c.ResultPath = "$"

	type CatcherForUnmarshal Catcher
	return json.Unmarshal(data, (*CatcherForUnmarshal)(c))
}

// Validate will validate the Catcher configuration
func (c *Catcher) Validate(_ context.Context) error {
	err := ValidateErrorEquals(c.ErrorEquals)
	if err != nil {
		return err
	}

	if len(c.Next) == 0 {
		return xerrors.New("Next cannot be empty")
	}

	return nil
}

// IsAnyErrorWildcard represent where is Catcher will matches any Error Names
func (c *Catcher) IsAnyErrorWildcard() bool {
	return IsAnyErrorWildcard(c.ErrorEquals)
}

// Catchers must be an array of Catcher.
type Catchers []Catcher

// Validate will validate the Catchers configuration
func (c Catchers) Validate(ctx context.Context) error {
	for i, catcher := range c {
		if err := catcher.Validate(ctx); err != nil {
			return err
		}

		if i != len(c)-1 && catcher.IsAnyErrorWildcard() {
			return xerrors.New("States.ALL must be the last catcher")
		}
	}

	return nil
}
