package executor

// BranchContext ...
type BranchContext struct {
	// /{tenant}/{namespace}/execution/{id}/activity/{id}/branch/{id}
	BranchID string

	Index int

	//
	ActivityID string
}
