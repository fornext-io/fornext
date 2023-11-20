package executor

// ActivityContext ...
type ActivityContext struct {
	// /{tenant}/{namespace}/execution/{id}/activity/{id}
	ID string

	StateName string

	ParentBranchID    *string
	ParentIterationID *string

	BranchStatus    *ActivityBranchStatus
	IterationStatus *ActivityIterationStatus
}

// ActivityBranchStatus ...
type ActivityBranchStatus struct {
	Max  int
	Done int
}

// ActivityIterationStatus ...
type ActivityIterationStatus struct {
	Max  int
	Done int
}
