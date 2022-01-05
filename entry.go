package log

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

// assert interface compliance.
var _ Interface = (*Entry)(nil)

var entryPool sync.Pool

// Now returns the current time.
var Now = time.Now

// Entry represents a single log entry.
type Entry struct {
	Logger    *Logger   `json:"-"`
	Fields    Fields    `json:"fields"`
	Level     Level     `json:"level"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
	start     time.Time
	fields    []Fields
	pool      bool
}

// newEntry returns a new entry for the given `log`.
// The entry is retrieved from the pool.
func newEntry(log *Logger) *Entry {
	var e *Entry
	if v := entryPool.Get(); v != nil {
		e = v.(*Entry)
	} else {
		e = new(Entry)
	}
	e.pool = true
	e.reset(log)
	return e
}

func NewEntry(log *Logger) *Entry {
	return &Entry{
		Logger: log,
	}
}

func (e *Entry) reset(l *Logger) {
	e.Logger = l
	e.Fields = nil
	e.fields = nil
}

// Release the entry to the pool if it was retrieved from it.
// The function is public such that asynchronous handlers can call it.
func (e *Entry) Release() {
	if e.pool {
		e.releaseFields()
		e.reset(nil)
		entryPool.Put(e)
	}
}

func (e *Entry) appendFields(fields Fielder) []Fields {
	f := make([]Fields, 0)
	f = append(f, e.fields...)
	if fields != nil {
		f = append(f, fields.Fields())
	}
	return f
}

// WithFields returns a new entry with `fields` set.
func (e *Entry) WithFields(fields Fielder) *Entry {
	return &Entry{
		Logger: e.Logger,
		fields: e.appendFields(fields),
	}
}

// withFields returns a new entry from the pool with `fields` set.
func (e *Entry) withFields(fields Fielder) *Entry {
	ret := newEntry(e.Logger)
	ret.fields = e.appendFields(fields)
	return ret
}

// WithField returns a new entry with the `key` and `value` set.
func (e *Entry) WithField(key string, value interface{}) *Entry {
	return e.WithFields(Fields{&Field{Name: key, Value: value}})
}

// WithDuration returns a new entry with the "duration" field set
// to the given duration in milliseconds.
func (e *Entry) WithDuration(d time.Duration) *Entry {
	return e.WithField("duration", d.Milliseconds())
}

// WithError returns a new entry with the "error" set to `err`.
//
// The given error may implement .Fielder, if it does the method
// will add all its `.Fields()` into the returned entry.
func (e *Entry) WithError(err error) *Entry {
	if err == nil {
		return e
	}

	ctx := e.WithField("error", err.Error())

	if s, ok := err.(stackTracer); ok {
		frame := s.StackTrace()[0]

		name := fmt.Sprintf("%n", frame)
		file := fmt.Sprintf("%+s", frame)
		line := fmt.Sprintf("%d", frame)

		parts := strings.Split(file, "\n\t")
		if len(parts) > 1 {
			file = parts[1]
		}

		ctx = ctx.WithField("source", fmt.Sprintf("%s: %s:%s", name, file, line))
	}

	if f, ok := err.(Fielder); ok {
		ctx = ctx.WithFields(f.Fields())
	}

	return ctx
}

// convert converts fields depending on their type.
// For example, it converts instances of "error" to strings, since errors are
// marshalled to "{}" by the standard json library...
func convert(val interface{}) interface{} {
	if err, ok := val.(error); ok {
		if _, ok := val.(json.Marshaler); !ok {
			return err.Error()
		}
	}
	if err, ok := val.(Sanitizer); ok {
		return err.Sanitize()
	}
	return val
}

func (e *Entry) withKvFields(args ...interface{}) *Entry {
	count := len(args)
	if args == nil || count == 0 {
		return e
	}
	if count == 1 {
		if slice, ok := args[0].([]interface{}); ok {
			// there is a single argument, and it's an []interface{}... most
			// probably the caller forgot to specify the ellipsis in the call
			// invocation: log.Info(msg, slice...). Hence we treat the slice as
			// the fields.
			args = slice
			count = len(slice)
		}
	}

	f := make(Fields, 0, (count+1)/2)
	for idx := 0; idx < count; idx++ {
		_, ok := args[idx].(error)
		if ok {
			// an error value without key
			f = append(f, newField("error", convert(args[idx])))
		} else if fields, ok := args[idx].(Fielder); ok {
			f = append(f, fields.Fields()...)
		} else if field, ok := args[idx].(Field); ok {
			f = append(f, &field)
		} else if field, ok := args[idx].(*Field); ok {
			f = append(f, field)
		} else if idx+1 < count {
			// there are (at least) two args left
			key, ok := args[idx].(string)
			if !ok {
				key = fmt.Sprintf("%v", args[idx])
			}
			f = append(f, newField(key, convert(args[idx+1])))
			idx++
		} else {
			f = append(f, newField("unknown", convert(args[idx])))
		}
	}
	return e.withFields(f)
}

// Trace level message.
func (e *Entry) Trace(msg string, fields ...interface{}) {
	e.Logger.log(TraceLevel, e.withKvFields(fields...), msg)
}

// Debug level message.
func (e *Entry) Debug(msg string, fields ...interface{}) {
	e.Logger.log(DebugLevel, e.withKvFields(fields...), msg)
}

// Info level message.
func (e *Entry) Info(msg string, fields ...interface{}) {
	e.Logger.log(InfoLevel, e.withKvFields(fields...), msg)
}

// Warn level message.
func (e *Entry) Warn(msg string, fields ...interface{}) {
	e.Logger.log(WarnLevel, e.withKvFields(fields...), msg)
}

// Error level message.
func (e *Entry) Error(msg string, fields ...interface{}) {
	e.Logger.log(ErrorLevel, e.withKvFields(fields...), msg)
}

// Fatal level message, followed by an exit.
func (e *Entry) Fatal(msg string, fields ...interface{}) {
	e.Logger.log(FatalLevel, e.withKvFields(fields...), msg)
	os.Exit(1)
}

// Tracef level formatted message.
func (e *Entry) Tracef(msg string, v ...interface{}) {
	e.Trace(fmt.Sprintf(msg, v...))
}

// Debugf level formatted message.
func (e *Entry) Debugf(msg string, v ...interface{}) {
	e.Debug(fmt.Sprintf(msg, v...))
}

// Infof level formatted message.
func (e *Entry) Infof(msg string, v ...interface{}) {
	e.Info(fmt.Sprintf(msg, v...))
}

// Warnf level formatted message.
func (e *Entry) Warnf(msg string, v ...interface{}) {
	e.Warn(fmt.Sprintf(msg, v...))
}

// Errorf level formatted message.
func (e *Entry) Errorf(msg string, v ...interface{}) {
	e.Error(fmt.Sprintf(msg, v...))
}

// Fatalf level formatted message, followed by an exit.
func (e *Entry) Fatalf(msg string, v ...interface{}) {
	e.Fatal(fmt.Sprintf(msg, v...))
}

// Watch returns a new entry with a Stop method to fire off
// a corresponding completion log, useful with defer.
func (e *Entry) Watch(msg string) *Entry {
	e.Info(msg)
	v := e.WithFields(e.Fields)
	v.Message = msg
	v.start = time.Now()
	return v
}

// Stop should be used with Trace, to fire off the completion message. When
// an `err` is passed the "error" field is set, and the log level is error.
func (e *Entry) Stop(err *error) {
	if err == nil || *err == nil {
		e.WithDuration(time.Since(e.start)).Info(e.Message)
	} else {
		e.WithDuration(time.Since(e.start)).WithError(*err).Error(e.Message)
	}
}

// mergedFields returns the fields list collapsed into a single one.
func (e *Entry) mergedFields() Fields {
	f := Fields{}

	for _, fields := range e.fields {
		for _, v := range fields {
			f = append(f, v)
		}
	}

	return f
}

// finalize returns a copy of the Entry with Fields merged.
func (e *Entry) finalize(level Level, msg string, pool bool) *Entry {
	if pool {
		// note: async entry cannot be taken from the pool since some handlers
		//       (e.g. memory handler) keep entries
		ret := newEntry(e.Logger)
		ret.Fields = e.mergedFields()
		ret.Level = level
		ret.Message = msg
		ret.Timestamp = Now()
		return ret
	}
	return &Entry{
		Logger:    e.Logger,
		Fields:    e.mergedFields(),
		Level:     level,
		Message:   msg,
		Timestamp: Now(),
	}
}

func (e *Entry) releaseFields() {
	for _, fields := range e.fields {
		for _, f := range fields {
			f.release()
		}
	}
}
