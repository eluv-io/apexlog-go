// Package logfmt implements a "logfmt" format handler.
package logfmt

import (
	"io"
	"os"
	"sync"

	"github.com/eluv-io/apexlog-go"
	"github.com/go-logfmt/logfmt"
)

// Default handler outputting to stderr.
var Default = New(os.Stderr)

// Handler implementation.
type Handler struct {
	mu  sync.Mutex
	enc *logfmt.Encoder
}

// New handler.
func New(w io.Writer) *Handler {
	return &Handler{
		enc: logfmt.NewEncoder(w),
	}
}

// HandleLog implements log.Handler.
func (h *Handler) HandleLog(e *log.Entry) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	_ = h.enc.EncodeKeyval("timestamp", e.Timestamp)
	_ = h.enc.EncodeKeyval("level", e.Level.String())
	_ = h.enc.EncodeKeyval("message", e.Message)

	for _, field := range e.Fields {
		_ = h.enc.EncodeKeyval(field.Name, field.Value)
	}

	_ = h.enc.EndRecord()

	return nil
}
