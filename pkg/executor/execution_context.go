package executor

// ExecutionContext ...
type ExecutionContext struct {
	ID string

	WorkflowID string

	Input []byte

	Timestamp uint64

	Status string
}
