package executor

// IterationContext ...
type IterationContext struct {
	// /{tenant}/{namespace}/execution/{id}/activity/{id}/iteration/{id}
	IterationID string

	Index int
	//
	ActivityID string
}
