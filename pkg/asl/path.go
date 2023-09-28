package asl

import (
	"strings"

	"github.com/lsytj0413/ena/xerrors"
	"github.com/ohler55/ojg/jp"
)

// Path is a string, beginning with '$', used to identify components with a JSON text.
// When a Path begins with '$$', this signals that it is intended to identify content within the Context Object.
// The syntax is https://github.com/json-path/JsonPath
type Path string

// ApplyPath ...
func ApplyPath(path *string, obj interface{}, contextObj interface{}) (interface{}, error) {
	if path == nil {
		return obj, nil
	}

	p := *path
	if p == "" {
		return nil, xerrors.Errorf("cannot apply path with empty expression")
	}

	inputObj := obj
	if strings.HasPrefix(p, "$$") {
		// The path is used with Context Object
		p = p[1:]
		inputObj = contextObj
	}

	return ApplyPathOnObject(p, inputObj)
}

// ApplyPathOnObject ...
func ApplyPathOnObject(path string, object interface{}) (interface{}, error) {
	if path == "" || path[0] != '$' {
		return nil, xerrors.Errorf("path must beginning with '$'")
	}

	expr, err := jp.ParseString(path)
	if err != nil {
		return nil, xerrors.Wrapf(err, "cannot parse path with JSONPath format")
	}

	result := expr.Get(object)
	if result == nil {
		// If the result is null (which means no element selected by this expression), we return an empty object
		return map[string]interface{}{}, nil
	}

	if len(result) == 1 {
		// TODO: This will be wrong if the expression selected only one element such as a[0]
		return result[0], nil
	}

	return result, nil
}
