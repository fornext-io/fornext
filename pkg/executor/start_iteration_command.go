package executor

// StartIterationCommand is triggered when an MapState is started,
// and this command will been replicated from leader to followers.
type StartIterationCommand struct {
	// IterationID is the identify of this iteration, it formated as
	// `/{ActivityID}/i{id}`
	// in which:
	//  1. {ActivityID}: is the referenced activity
	//  2. {id}: is the iteration index of current iteration, begin with 0
	IterationID string

	ExecutionID string

	// Index is serial of this iteration in MapState.
	Index int

	// ActivityID is the identify of activity which this branch belongs.
	ActivityID string

	// Input is the arguments of this iteration. It MUST be an valid JSON-Object.
	Input []byte

	// Timestamp is `HLC` value of this command request time.
	Timestamp uint64
}
