package asl

// PathResolver will return the Path
type PathResolver interface {
	// Resolve will return the object which is selected by JSONPath
	Resolve(p string) (interface{}, error)
}
