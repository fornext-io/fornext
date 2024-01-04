package executor

// StartExecutionCommand is triggered when an workflow is ready for
// execution, and this command will been replicated from leader to
// followers.
type StartExecutionCommand struct {
	// ID is the identify of this execution, it formated as
	// `/{tenant}/{namespace}/e{xxxxxxx}/{id}`
	// in which:
	//  1. {tenant}: is the tenant name of current execution
	//  2. {namespace}: is the namespace name of current execution
	//  3. {xxxxxxx}: is 7-length hex string, when start with idempotent,
	//     it was the hash of idempotent key; when without idempotent, it's
	//     the hash of last id.
	//  4. {id}: is the 16-length hex string of current `hybrid logical clock` value.
	//
	// All execution will partitioned within namespace, the partition prefix is
	// without last {id} to support idempotent without cross-partition transaction.
	// So there maybe large partition.
	ID string

	// WorkflowID is the identify of workflow definition which this execution belongs.
	WorkflowID string

	// Input is the arguments of this execution. It MUST be an valid JSON-Object.
	Input []byte

	// Timestamp is `HLC` value of this command request time.
	Timestamp uint64
}
