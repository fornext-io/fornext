package executor

// StartBranchCommand ...
type StartBranchCommand struct {
	// /{tenant}/{namespace}/execution/{id}/activity/{id}/branch/{id}
	BranchID string

	Index int

	//
	ActivityID string
}
