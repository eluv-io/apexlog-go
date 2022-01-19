// Package multi implements a handler which invokes a number of handlers.
package multi

import (
	"github.com/eluv-io/apexlog-go"
)

// Handler implementation.
type Handler struct {
	Handlers []log.Handler
	async    bool
}

// New handler.
func New(h ...log.Handler) *Handler {
	async := false
	for _, l := range h {
		if as, ok := l.(log.Asynchronous); ok && as.Asynchronous() {
			async = true
		}
	}
	return &Handler{
		Handlers: h,
		async:    async,
	}
}

// HandleLog implements log.Handler.
func (h *Handler) HandleLog(e *log.Entry) error {
	for _, handler := range h.Handlers {
		// TODO(tj): maybe just write to stderr here, definitely not ideal
		// to miss out logging to a more critical handler if something
		// goes wrong
		if err := handler.HandleLog(e); err != nil {
			return err
		}
	}

	return nil
}

func (h *Handler) Asynchronous() bool {
	return h.async
}
