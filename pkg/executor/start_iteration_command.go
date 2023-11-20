package executor

// StartIterationCommand ...
type StartIterationCommand struct {
	// /{tenant}/{namespace}/execution/{id}/activity/{id}/iteration/{id}
	IterationID string

	Index int

	//
	ActivityID string
}
