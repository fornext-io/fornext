package executor

// StartBranchCommand is triggered when an ParallelState is started,
// and this command will been replicated from leader to followers.
type StartBranchCommand struct {
	// BranchID is the identify of this branch, it formated as
	// `/{ActivityID}/b{id}`
	// in which:
	//  1. {ActivityID}: is the referenced activity
	//  2. {id}: is the branch index of current branch, begin with 0
	BranchID string

	ExecutionID string

	// Index is serial of this branch in ParallelState.
	Index int

	// ActivityID is the identify of activity which this branch belongs.
	ActivityID string

	// Input is the arguments of this branch. It MUST be an valid JSON-Object.
	Input []byte

	// Timestamp is `HLC` value of this command request time.
	Timestamp uint64
}
