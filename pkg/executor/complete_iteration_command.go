package executor

// CompleteIterationCommand is triggered when an iteration is complete,
// and this command will been replicated from leader to followers.
type CompleteIterationCommand struct {
	// IterationID is the identify of iteration.
	IterationID string

	// Output is the results of this iteration. it MUST be an valid JSON-Object.
	Output []byte

	// Timestamp is `HLC` value of this command request time.
	Timestamp uint64
}
