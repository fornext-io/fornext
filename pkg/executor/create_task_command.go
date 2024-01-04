package executor

// CreateTaskCommand is triggered when an `TaskState`、 `WaitState` or other task is ready
// to fired. And this command will been replicated from leader to followers.
type CreateTaskCommand struct {
	// ID is the identify of this task, it formated as
	// `/{tenant}/{namespace}/t{xxxxxxx}/{id}`
	//  1. {tenant}: is the tenant name of current execution
	//  2. {namespace}: is the namespace name of current execution
	//  3. {xxxxxxx}: is 7-length hex string, it the hash of task's type.
	//  4. {id}: is the 8-length hex string of current `hybrid logical clock` value.
	ID string

	// ActivityID is the identify of activity which this task belongs.
	ActivityID string

	// Resource is the task's reference resource definitions.
	Resource string

	// Input is the arguments of this task. It MUST be an valid JSON-Object.
	Input []byte

	// Timestamp is `HLC` value of this command request time.
	Timestamp uint64
}
