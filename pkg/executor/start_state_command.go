package executor

// StartStateCommand ...
type StartStateCommand struct {
	// /{tenant}/{namespace}/execution/{id}/activity/{id}
	ActivityID string

	StateName         string
	ParentBranchID    *string
	ParentIterationID *string
}
