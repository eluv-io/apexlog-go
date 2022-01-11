package log

// Sanitizer is an interface for types that may contain sensitive information
// (like passwords), which shouldn't be printed to the log.
type Sanitizer interface {
	Sanitize() interface{}
}
