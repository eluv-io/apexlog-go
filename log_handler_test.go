package log_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/apex/log"
	hjson "github.com/apex/log/handlers/json"
)

func TestLog_text(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0))

	l := &log.Logger{
		Handler: NewText(buf),
		Level:   log.InfoLevel,
	}

	ctx := l.WithField("file", "sloth.png").WithField("user", "Tobi")
	ctx.Debug("uploading")
	ctx.Info("upload complete")

	ctx = ctx.WithError(fmt.Errorf("comparison failed"))
	ctx.Error("bad file")
	l.Error("bad file", "file", "sloth.png", "user", "Tobi", fmt.Errorf("comparison failed"))
	fmt.Println(string(buf.Bytes()))

	// Output:
	//
	//   INFO[0000] upload complete           file=sloth.png user=Tobi
	//  ERROR[0000] bad file                  file=sloth.png user=Tobi error=comparison failed
	//  ERROR[0000] bad file                  file=sloth.png user=Tobi error=comparison failed
	//
}

func ExampleLog_text() {
	buf := bytes.NewBuffer(make([]byte, 0))

	l := &log.Logger{
		Handler: NewText(buf),
		Level:   log.InfoLevel,
	}

	ctx := l.WithField("file", "sloth.png").WithField("user", "Tobi")
	ctx.Debug("uploading")
	ctx.Info("upload complete")

	ctx = ctx.WithError(fmt.Errorf("comparison failed1"))
	ctx.Error("bad file")
	l.Error("bad file", "file", "sloth.png", "user", "Tobi", fmt.Errorf("comparison failed2"))
	fmt.Println(string(buf.Bytes()))

	// Output:
	//
	//   INFO[0000] upload complete           file=sloth.png user=Tobi
	//  ERROR[0000] bad file                  file=sloth.png user=Tobi error=comparison failed1
	//  ERROR[0000] bad file                  file=sloth.png user=Tobi error=comparison failed2
	//
}

func ExampleLog_fields_text() {
	buf := bytes.NewBuffer(make([]byte, 0))

	l := &log.Logger{
		Handler: NewText(buf),
		Level:   log.InfoLevel,
	}

	var err error
	err = fmt.Errorf("comparison failed")
	l.Error("bad file", "file", "sloth.png", "user", "Tobi", err)
	l.Error("bad file", err, "file", "sloth.png", "user", "Tobi")
	err = &structuredError{
		Errno:  25,
		Reason: "bad descriptor1",
	}
	l.Error("bad file", "file", "sloth.png", "user", "Tobi", err)
	err = &jstructuredError{
		Errno:  26,
		Reason: "bad descriptor2",
	}
	l.Error("bad file", "file", "sloth.png", "user", "Tobi", err)

	fmt.Println(string(buf.Bytes()))

	// Output:
	//
	//  ERROR[0000] bad file                  file=sloth.png user=Tobi error=comparison failed
	//  ERROR[0000] bad file                  error=comparison failed file=sloth.png user=Tobi
	//  ERROR[0000] bad file                  file=sloth.png user=Tobi error=25/bad descriptor1
	//  ERROR[0000] bad file                  file=sloth.png user=Tobi error=26/bad descriptor2
	//
}

func ExampleLog_json() {
	log.Now = func() time.Time {
		return time.Unix(0, 0).UTC()
	}
	defer func() { log.Now = time.Now }()

	buf := bytes.NewBuffer(make([]byte, 0))

	l := &log.Logger{
		Handler: hjson.New(buf, false),
		Level:   log.InfoLevel,
	}

	ctx := l.WithField("file", "sloth.png").WithField("user", "Tobi")
	ctx.Debug("uploading")
	ctx.Info("upload complete")

	ctx = ctx.WithError(fmt.Errorf("comparison failed"))
	ctx.Error("bad file")
	l.Error("bad file", "file", "sloth.png", "user", "Tobi", fmt.Errorf("comparison failed"))

	fmt.Println(string(buf.Bytes()))

	// Output:
	// {"fields":{"file":"sloth.png","user":"Tobi"},"level":"info","timestamp":"1970-01-01T00:00:00Z","message":"upload complete"}
	// {"fields":{"file":"sloth.png","user":"Tobi","error":"comparison failed"},"level":"error","timestamp":"1970-01-01T00:00:00Z","message":"bad file"}
	// {"fields":{"file":"sloth.png","user":"Tobi","error":"comparison failed"},"level":"error","timestamp":"1970-01-01T00:00:00Z","message":"bad file"}
	//
}

func ExampleLog_fields_json() {
	log.Now = func() time.Time {
		return time.Unix(0, 0).UTC()
	}
	defer func() { log.Now = time.Now }()

	buf := bytes.NewBuffer(make([]byte, 0))

	l := &log.Logger{
		Handler: hjson.New(buf, false),
		Level:   log.InfoLevel,
	}

	var err error
	err = fmt.Errorf("comparison failed")
	l.Error("bad file", "file", "sloth.png", "user", "Tobi", err)
	l.Error("bad file", err, "file", "sloth.png", "user", "Tobi")
	err = &structuredError{
		Errno:  25,
		Reason: "bad descriptor",
	}
	l.Error("bad file", "file", "sloth.png", "user", "Tobi", err)
	err = &jstructuredError{
		Errno:  25,
		Reason: "bad descriptor",
	}
	l.Error("bad file", "file", "sloth.png", "user", "Tobi", err)

	fmt.Println(string(buf.Bytes()))

	// Output:
	// {"fields":{"file":"sloth.png","user":"Tobi","error":"comparison failed"},"level":"error","timestamp":"1970-01-01T00:00:00Z","message":"bad file"}
	// {"fields":{"error":"comparison failed","file":"sloth.png","user":"Tobi"},"level":"error","timestamp":"1970-01-01T00:00:00Z","message":"bad file"}
	// {"fields":{"file":"sloth.png","user":"Tobi","error":"25/bad descriptor"},"level":"error","timestamp":"1970-01-01T00:00:00Z","message":"bad file"}
	// {"fields":{"file":"sloth.png","user":"Tobi","error":{"errno":25,"reason":"bad descriptor"}},"level":"error","timestamp":"1970-01-01T00:00:00Z","message":"bad file"}
	//
}

// Strings mapping.
var Strings = [...]string{
	log.DebugLevel: "DEBUG",
	log.InfoLevel:  "INFO",
	log.WarnLevel:  "WARN",
	log.ErrorLevel: "ERROR",
	log.FatalLevel: "FATAL",
}

// TextHandler implementation.
type TextHandler struct {
	mu     sync.Mutex
	start  time.Time
	Writer io.Writer
}

// New handler.
func NewText(w io.Writer) *TextHandler {
	return &TextHandler{
		start:  time.Now(),
		Writer: w,
	}
}

// HandleLog implements log.Handler.
func (h *TextHandler) HandleLog(e *log.Entry) error {
	level := Strings[e.Level]

	h.mu.Lock()
	defer h.mu.Unlock()

	ts := time.Since(h.start) / time.Second
	_, _ = fmt.Fprintf(h.Writer, "%6s[%04d] %-25s", level, ts, e.Message)

	for _, field := range e.Fields {
		_, _ = fmt.Fprintf(h.Writer, " %s=%v", field.Name, field.Value)
	}
	_, _ = fmt.Fprintln(h.Writer)

	return nil
}

type structuredError struct {
	Errno  int    `json:"errno"`
	Reason string `json:"reason"`
}

func (e *structuredError) Error() string {
	return fmt.Sprintf("%d/%s", e.Errno, e.Reason)
}

type jstructuredError struct {
	Errno  int    `json:"errno"`
	Reason string `json:"reason"`
}

func (e *jstructuredError) Error() string {
	return fmt.Sprintf("%d/%s", e.Errno, e.Reason)
}

func (e *jstructuredError) MarshalJSON() ([]byte, error) {
	type jj jstructuredError
	ej := jj(*e)
	return json.Marshal(&ej)
}
