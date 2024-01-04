package utils

import (
	"encoding/json"
)

// MustUnmarshalJSON will unmarshal data to interface{}, otherwise
// it will panic.
func MustUnmarshalJSON(data []byte) interface{} {
	var v interface{}
	err := json.Unmarshal(data, &v)
	if err != nil {
		panic(err)
	}

	return v
}
