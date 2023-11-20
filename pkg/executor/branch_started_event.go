package executor

// BranchStartedEvent ...
type BranchStartedEvent struct {
	// /{tenant}/{namespace}/execution/{id}/activity/{id}/branch/{id}
	BranchID string
}
