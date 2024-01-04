package fsl

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/lsytj0413/ena/xerrors"
	"github.com/ohler55/ojg/jp"
)

// Path is a string, beginning with '$', used to identify components with a JSON text.
// When a Path begins with '$$', this signals that it is intended to identify content within the Context Object.
// The syntax is https://github.com/json-path/JsonPath
type Path string

type pathContext struct {
	Input       []byte
	ContextData []byte
}

// Apply will apply this Path on provided context
func (p Path) Apply(_ context.Context, pc pathContext) ([]byte, error) {
	pp := string(p)
	switch {
	case pp == "":
		// If the user set InputPath to empty explicit, that means raw
		// input is discarded, and the effective input for the state is
		// and tempty JSON object, {}.
		return []byte(`{}`), nil
	case pp == "$":
		// We do nothing, the default `$` means use raw input
		return pc.Input, nil
	}

	input := pc.Input
	if strings.HasPrefix(pp, "$$") {
		input = pc.ContextData
		pp = pp[1:]
	}

	var inputObj interface{}
	err := json.Unmarshal(input, &inputObj)
	if err != nil {
		return nil, err
	}

	expr, err := jp.ParseString(pp)
	if err != nil {
		return nil, err
	}
	result := expr.Get(inputObj)
	switch {
	case len(result) == 0:
		return nil, xerrors.Errorf("must at least one element InputPath, got %v", len(result))
	case len(result) == 1:
		// This expr only select on element
		inputObj = result[0]
	default:
		inputObj = result
	}

	input, err = json.Marshal(inputObj)
	if err != nil {
		return nil, err
	}
	return input, nil
}
