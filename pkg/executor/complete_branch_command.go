package executor

// CompleteBranchCommand is triggered when an branch is complete,
// and this command will been replicated from leader to followers.
type CompleteBranchCommand struct {
	// BranchID is the identify of branch.
	BranchID string

	// Output is the results of this branch. it MUST be an valid JSON-Object.
	Output []byte

	// Timestamp is `HLC` value of this command request time.
	Timestamp uint64
}
