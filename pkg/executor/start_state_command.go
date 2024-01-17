package executor

// StartStateCommand is triggered when an State is ready for
// execution, and this command will been replicated from leader to
// followers.
type StartStateCommand struct {
	// ActivityID is the identify of this state's context, it formated as
	// `/{ExecutionID}/{id}`
	// in which:
	//  1. {ExecutionID}: is the referenced execution of current state's context
	//  2. {id}: is the 8-length hex string of current `hybrid logical clock` value.
	//
	// An execution's all activity have the same prefix, and will in same partition.
	ActivityID string

	// ExecutionID is the identify of execution which this state belong to.
	ExecutionID string

	// StateName is current state's name which ready for execution, and because all state's
	// name must be unique (include in each nest state), so only name is enough.
	StateName string

	// ParentBranchID is exists when this state is invoked in ParallelState.
	ParentBranchID *string

	// ParentIterationID is exists when this state is invoked in MapState.
	ParentIterationID *string

	// Input is the arguments of this state. It MUST be an valid JSON-Object.
	Input []byte

	// Timestamp is `HLC` value of this command request time.
	Timestamp uint64
}
