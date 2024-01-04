package executor

// CompleteExecutionCommand is triggered when an execution is completed,
// and this command will been replicated from leader to followers.
type CompleteExecutionCommand struct {
	// ID is the identify of this execution.
	ID string

	// Output is the results of this execution. it MUST be an valid JSON-Object.
	Output []byte

	// Timestamp is `HLC` value of this command request time.
	Timestamp uint64
}
