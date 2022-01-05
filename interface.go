package log

import "time"

// Interface represents the API of both Logger and Entry and exposes 3 types of
// functions:
// - functions named like WithXX and Watch return entries that can be used in chained call
// - xxf are in printf style with message format and arguments
// - logging functions - like Info - log a message and optional kv arguments that
//   are expected to be key/value pairs exception made of error values that can
//   be passed alone and will automatically be assigned to a key 'error'.
type Interface interface {
	// WithFields returns a new entry with the given Fields appended
	WithFields(Fielder) *Entry
	// WithField returns a new entry with the given name and value appended to fields
	WithField(name string, value interface{}) *Entry
	// WithDuration returns a new entry with the given duration appended as a 'duration' field
	WithDuration(time.Duration) *Entry
	// WithError returns a new entry with the given error appended as an 'error' field
	WithError(error) *Entry

	Trace(msg string, kv ...interface{}) // Trace is a Trace level message with KV values.
	Debug(msg string, kv ...interface{}) // Debug is a Debug level message with KV values.
	Info(msg string, kv ...interface{})  // Info is a Info level message with KV values.
	Warn(msg string, kv ...interface{})  // Warn is a Warn level message with KV values.
	Error(msg string, kv ...interface{}) // Error is a Error level message with KV values.
	Fatal(msg string, kv ...interface{}) // Fatal is a Fatal level message with KV values.

	Tracef(string, ...interface{}) // Tracef is a Trace level formatted message.
	Debugf(string, ...interface{}) // Debugf is a Debug level formatted message.
	Infof(string, ...interface{})  // Infof is a Info level formatted message.
	Warnf(string, ...interface{})  // Warnf is a Warn level formatted message.
	Errorf(string, ...interface{}) // Errorf is a Error level formatted message.
	Fatalf(string, ...interface{}) // Tracef is a Fatal level formatted message.

	// Watch returns a new entry whose Stop method can be used to fire off a
	// corresponding log that will include the duration taken for completion:
	// useful with defer.
	Watch(string) *Entry
}
