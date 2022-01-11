package log

import (
	"bytes"
	"encoding/json"
	stdlog "log"
	"sort"
	"time"
)

// assert interface compliance.
var _ Interface = (*Logger)(nil)

// Fielder is an interface for providing fields to custom types.
type Fielder interface {
	Fields() Fields
}

// Fields represents a slice of entry level data used for structured logging.
type Fields []*Field

// Fields implements Fielder.
func (f Fields) Fields() Fields {
	return f
}

// Append adds the given name/value pair as a Field and returns the updated Fields
// ff := log.Fields(nil).Append("count", 1).Append("name", "bob")
func (f Fields) Append(name string, value interface{}) Fields {
	ret := append(f, &Field{
		Name:  name,
		Value: value,
	})
	return ret
}

// Get field value by name.
func (f Fields) Get(name string) interface{} {
	for _, f := range f {
		if f.Name == name {
			return f.Value
		}
	}
	return nil
}

// Names returns field names sorted.
func (f Fields) Names() (v []string) {
	for _, k := range f {
		v = append(v, k.Name)
	}

	sort.Strings(v)
	return
}

func (f Fields) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	for i, field := range f {
		bb, err := field.toJSON()
		if err != nil {
			return nil, err
		}
		_, err = buf.WriteString(string(bb))
		if err != nil {
			return nil, err
		}
		if i < len(f)-1 {
			buf.WriteString(",")
		}
	}
	return []byte("{" + string(buf.Bytes()) + "}"), nil
}

func (f *Fields) UnmarshalJSON(b []byte) error {
	m := make(map[string]interface{})
	err := json.Unmarshal(b, &m)
	if err != nil {
		return err
	}
	ff := *f
	for k, v := range m {
		ff = append(ff, &Field{Name: k, Value: v})
	}
	*f = ff
	return nil
}

func (f Fields) Map() map[string]interface{} {
	if len(f) == 0 {
		return nil
	}
	ret := make(map[string]interface{})
	for _, field := range f {
		ret[field.Name] = field.Value
	}
	return ret
}

// The HandlerFunc type is an adapter to allow the use of ordinary functions as
// log handlers. If f is a function with the appropriate signature,
// HandlerFunc(f) is a Handler object that calls f.
type HandlerFunc func(*Entry) error

// HandleLog calls f(e).
func (f HandlerFunc) HandleLog(e *Entry) error {
	return f(e)
}

// Handler is used to handle log events, outputting them to
// stdio or sending them to remote services. See the "handlers"
// directory for implementations.
//
// It is left up to Handlers to implement thread-safety.
type Handler interface {
	HandleLog(*Entry) error
}

// Asynchronous is an optional interface for handlers that wish to keep log
// entries (known implementation are in multi, es, memory)
type Asynchronous interface {
	// Asynchronous returns true if the handler takes responsibility for
	// releasing log entries
	Asynchronous() bool
}

// Logger represents a logger with configurable Level and Handler.
type Logger struct {
	Handler Handler
	Level   Level
}

// WithFields returns a new entry with `fields` set.
func (l *Logger) WithFields(fields Fielder) *Entry {
	ret := l.newEntry()
	defer ret.Release()
	return ret.WithFields(fields.Fields())
}

// WithField returns a new entry with the `key` and `value` set.
//
// Note that the `key` should not have spaces in it - use camel
// case or underscores
func (l *Logger) WithField(key string, value interface{}) *Entry {
	ret := l.newEntry()
	defer ret.Release()
	return ret.WithField(key, value)
}

// WithDuration returns a new entry with the "duration" field set
// to the given duration in milliseconds.
func (l *Logger) WithDuration(d time.Duration) *Entry {
	ret := l.newEntry()
	defer ret.Release()
	return ret.WithDuration(d)
}

// WithError returns a new entry with the "error" set to `err`.
func (l *Logger) WithError(err error) *Entry {
	if err == nil {
		return NewEntry(l)
	}
	ret := l.newEntry()
	defer ret.Release()
	return ret.WithError(err)
}

func (l *Logger) Trace(msg string, fields ...interface{}) {
	e := l.newEntry()
	defer e.Release()
	e.Trace(msg, fields...)
}

// Debug level message.
func (l *Logger) Debug(msg string, fields ...interface{}) {
	e := l.newEntry()
	defer e.Release()
	e.Debug(msg, fields...)
}

// Info level message.
func (l *Logger) Info(msg string, fields ...interface{}) {
	e := l.newEntry()
	defer e.Release()
	e.Info(msg, fields...)
}

// Warn level message.
func (l *Logger) Warn(msg string, fields ...interface{}) {
	e := l.newEntry()
	defer e.Release()
	e.Warn(msg, fields...)
}

// Error level message.
func (l *Logger) Error(msg string, fields ...interface{}) {
	e := l.newEntry()
	defer e.Release()
	e.Error(msg, fields...)
}

// Fatal level message, followed by an exit.
func (l *Logger) Fatal(msg string, fields ...interface{}) {
	e := l.newEntry()
	defer e.Release()
	e.Fatal(msg, fields...)
}

// Tracef level formatted message.
func (l *Logger) Tracef(msg string, v ...interface{}) {
	e := l.newEntry()
	defer e.Release()
	e.Tracef(msg, v...)
}

// Debugf level formatted message.
func (l *Logger) Debugf(msg string, v ...interface{}) {
	e := l.newEntry()
	defer e.Release()
	e.Debugf(msg, v...)
}

// Infof level formatted message.
func (l *Logger) Infof(msg string, v ...interface{}) {
	e := l.newEntry()
	defer e.Release()
	e.Infof(msg, v...)
}

// Warnf level formatted message.
func (l *Logger) Warnf(msg string, v ...interface{}) {
	e := l.newEntry()
	defer e.Release()
	e.Warnf(msg, v...)
}

// Errorf level formatted message.
func (l *Logger) Errorf(msg string, v ...interface{}) {
	e := l.newEntry()
	defer e.Release()
	e.Errorf(msg, v...)
}

// Fatalf level formatted message, followed by an exit.
func (l *Logger) Fatalf(msg string, v ...interface{}) {
	e := l.newEntry()
	defer e.Release()
	e.Fatalf(msg, v...)
}

// Watch returns a new entry with a Stop method to fire off
// a corresponding completion log, useful with defer.
func (l *Logger) Watch(msg string) *Entry {
	return NewEntry(l).Watch(msg)
}

// log the message, invoking the handler. We clone the entry here
// to bypass the overhead in Entry methods when the level is not met.
func (l *Logger) log(level Level, e *Entry, msg string) {
	if l == nil {
		return
	}
	if level < l.Level {
		return
	}
	entry := e.finalize(level, msg, l.usePool())
	defer entry.Release()

	if err := l.Handler.HandleLog(entry); err != nil {
		stdlog.Printf("error logging: %s", err)
	}
}

func (l *Logger) newEntry() *Entry {
	if l.usePool() {
		return newEntry(l)
	}
	return NewEntry(l)
}

func (l *Logger) usePool() bool {
	async := false
	if fin, ok := l.Handler.(Asynchronous); ok {
		async = fin.Asynchronous()
	}
	return !async
}
