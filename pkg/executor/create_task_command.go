package executor

// CreateTaskCommand ...
type CreateTaskCommand struct {
	// /{tenant}/{namespace}/task/{id}
	ID string

	ActivityID string
	Resource   string
}
