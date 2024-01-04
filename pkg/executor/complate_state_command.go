package executor

// CompleteStateCommand is triggered when an State is complete,
// and this command will been replicated from leader to followers.
type CompleteStateCommand struct {
	// ActivityID is the identify of this state's context.
	ActivityID string

	// Output is the results of this state. it MUST be an valid JSON-Object.
	Output []byte

	// Status is the result of this state, which could be one of
	// 1. `Succeeded`: the state is executed successful.
	// 2. `Failed`: the state is executed failed, it always will be caused by workflow interpreter.
	Status string

	// Reason is the errcode string when `Status` is `Failed`.
	// It's a brief CamelCase message.
	Reason string

	// Message is the err message string when `Status` is `Failed`.
	// It's a human readable message indicating details.
	Message string

	// Timestamp is `HLC` value of this command request time.
	Timestamp uint64
}
