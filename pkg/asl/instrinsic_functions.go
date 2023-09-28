package asl

const (
	// FuncStatesFormat takes one or more arguments and return first-argument string with each '{}'
	// replaced by the value of positionally-corresponding argument.
	FuncStatesFormat string = "States.Format"

	// FuncStatesStringToJSON takes a single argument whose Value MUST be a string.
	FuncStatesStringToJSON string = "States.StringToJson"

	// FuncStatesJSONToString takes a single aargument which MUST be a Path and return a string.
	FuncStatesJSONToString string = "States.JsonToString"

	// FuncStatesArray takes zero or more arguments, returns a JSON array
	FuncStatesArray string = "States.Array"

	// FuncStatesArrayPartition to partition a large array:
	//   1. The first argument is an array
	//   2. The second argument defines the chunk size
	FuncStatesArrayPartition string = "States.ArrayPartition"

	// FuncStatesArrayContains to determine if a specific value is present in an array.
	FuncStatesArrayContains string = "States.ArrayContains"

	// FuncStatesArrayRange to create a new array containing a specific range of elements.
	// up to 1000 elements.
	FuncStatesArrayRange string = "States.ArrayRange"

	// FuncStatesArrayGetItem returns a specified index's value.
	FuncStatesArrayGetItem string = "States.ArrayGetItem"

	// FuncStatesArrayLength returns the length of an array.
	FuncStatesArrayLength string = "States.ArrayLength"

	// FuncStatesArrayUnique removes duplicate values from an array and returns an array containing
	// only unique elements.
	FuncStatesArrayUnique string = "States.ArrayUnique"

	// FuncStatesBase64Encode to encode data based on MIME Base64 scheme.
	FuncStatesBase64Encode string = "States.Base64Encode"

	// FuncStatesBase64Decode to decode data base on MIME Base64 decoding scheme.
	FuncStatesBase64Decode string = "States.Base64Decode"

	// FuncStatesHash to calculate the hash value of a given input.
	// Support hashing algorithm: MD5、SHA-1、SHA-256、SHA-384、SHA-512
	FuncStatesHash string = "States.Hash"

	// FuncStatesJSONMerge to merge two JSON objects into a single object.
	FuncStatesJSONMerge string = "States.JsonMerge"

	// FuncStatesMathRandom to return a random number between the specified start and end number.
	FuncStatesMathRandom string = "States.MathRandom"

	// FuncStatesMathAdd to return the sum of two numbers.
	FuncStatesMathAdd string = "States.MathAdd"

	// FuncStatesStringSplit to split a string into an array of values.
	FuncStatesStringSplit string = "States.StringSplit"

	// FuncStatesUUID to return a version 4 universally unique identifier.
	FuncStatesUUID string = "States.UUID"
)

var (
	// AllFuncs contains all pre defined funcs
	AllFuncs = map[string]bool{
		FuncStatesFormat:         true,
		FuncStatesStringToJSON:   true,
		FuncStatesJSONToString:   true,
		FuncStatesArray:          true,
		FuncStatesArrayPartition: true,
		FuncStatesArrayContains:  true,
		FuncStatesArrayRange:     true,
		FuncStatesArrayGetItem:   true,
		FuncStatesArrayLength:    true,
		FuncStatesArrayUnique:    true,
		FuncStatesBase64Encode:   true,
		FuncStatesBase64Decode:   true,
		FuncStatesHash:           true,
		FuncStatesJSONMerge:      true,
		FuncStatesMathRandom:     true,
		FuncStatesMathAdd:        true,
		FuncStatesStringSplit:    true,
		FuncStatesUUID:           true,
	}
)
