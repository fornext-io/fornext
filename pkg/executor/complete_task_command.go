package executor

// CompleteTaskCommand is triggered when an Task is complete,
// and this command will been replicated from leader to followers.
type CompleteTaskCommand struct {
	// TaskID is the identify of this task.
	TaskID string

	// Output is the results of this task. it MUST be an valid JSON-Object.
	Output []byte

	// Timestamp is `HLC` value of this command request time.
	Timestamp uint64
}
