package asl

import (
	"github.com/lsytj0413/ena/xerrors"
	"github.com/ohler55/ojg/jp"
)

// ReferencePath is a Path with syntax limited as below:
//  1. Can only identify a single node in a JSON structure
//  2. The operators '@'、','、':'、'?' are not supported
//
// ReferencePath MUST be unambiguous references to a single value, array or object.
type ReferencePath string

// ApplyReferencePath ...
func ApplyReferencePath(path *string, object interface{}, value interface{}) (interface{}, error) {
	if path == nil {
		return object, nil
	}

	p := *path
	if p == "" || p[0] != '$' {
		return nil, xerrors.Errorf("path must beginning with '$'")
	}

	expr, err := jp.ParseString(p)
	if err != nil {
		return nil, xerrors.Wrapf(err, "cannot parse path with JSONPath format")
	}

	err = expr.Set(object, value)
	if err != nil {
		return nil, err
	}
	return object, nil
}
