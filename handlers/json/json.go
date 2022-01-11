// Package json implements a JSON handler.
package json

import (
	j "encoding/json"
	"io"
	"os"
	"sync"

	"github.com/apex/log"
)

// Default handler outputting to stderr.
var Default = New(os.Stderr)

// Handler implementation.
type Handler struct {
	*j.Encoder
	mu sync.Mutex
}

// New returns a new handler. By default, the json encoder used by the handler
// has SetEscapeHTML(false). The first escapeHtml optional params can be used
// to change this behavior.
func New(w io.Writer, escapeHtml ...bool) *Handler {
	ret := &Handler{
		Encoder: j.NewEncoder(w),
	}
	ret.Encoder.SetEscapeHTML(false)
	if len(escapeHtml) > 0 {
		ret.Encoder.SetEscapeHTML(escapeHtml[0])
	}
	return ret
}

// HandleLog implements log.Handler.
func (h *Handler) HandleLog(e *log.Entry) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.Encoder.Encode(e)
}
