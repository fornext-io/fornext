package fsl

import (
	"context"
	"encoding/json"

	"github.com/ohler55/ojg/jp"
)

// ReferencePath is a Path with syntax limited as below:
//  1. Can only identify a single node in a JSON structure
//  2. The operators '@'、','、':'、'?' are not supported
//
// ReferencePath MUST be unambiguous references to a single value, array or object.
type ReferencePath string

type referencePathContext struct {
	Input  []byte
	Output []byte
}

// Apply will apply this ReferencePath on provided context
func (p ReferencePath) Apply(_ context.Context, pc referencePathContext) ([]byte, error) {
	pp := string(p)
	switch {
	case pp == "":
		// NOTE: this should change to state's raw input
		return pc.Input, nil
	case pp == "$":
		// We do nothing, the default `$` means use raw output
		return pc.Output, nil
	}

	expr, err := jp.ParseString(pp)
	if err != nil {
		return nil, err
	}

	var inputObj interface{}
	err = json.Unmarshal(pc.Input, &inputObj)
	if err != nil {
		return nil, err
	}

	var outputObj interface{}
	err = json.Unmarshal(pc.Output, &outputObj)
	if err != nil {
		return nil, err
	}

	err = expr.Set(inputObj, outputObj)
	if err != nil {
		return nil, err
	}

	result, err := json.Marshal(inputObj)
	if err != nil {
		return nil, err
	}
	return result, nil
}
