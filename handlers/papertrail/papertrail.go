// Package papertrail implements a papertrail logfmt format handler.
package papertrail

import (
	"bytes"
	"fmt"
	"log/syslog"
	"net"
	"os"
	"sync"
	"time"

	"github.com/eluv-io/apexlog-go"
	"github.com/go-logfmt/logfmt"
)

// TODO: syslog portion is ad-hoc for my serverless use-case,
// I don't really need hostnames etc, but this should be improved

// Config for Papertrail.
type Config struct {
	// Papertrail settings.
	Host string // Host subdomain such as "logs4"
	Port int    // Port number

	// Application settings
	Hostname string // Hostname value
	Tag      string // Tag value
}

// Handler implementation.
type Handler struct {
	*Config

	mu   sync.Mutex
	conn net.Conn
}

// New handler.
func New(config *Config) *Handler {
	conn, err := net.Dial("udp", fmt.Sprintf("%s.papertrailapp.com:%d", config.Host, config.Port))
	if err != nil {
		panic(err)
	}

	return &Handler{
		Config: config,
		conn:   conn,
	}
}

// HandleLog implements log.Handler.
func (h *Handler) HandleLog(e *log.Entry) error {
	ts := time.Now().Format(time.Stamp)

	var buf bytes.Buffer

	enc := logfmt.NewEncoder(&buf)
	_ = enc.EncodeKeyval("level", e.Level.String())
	_ = enc.EncodeKeyval("message", e.Message)

	for _, field := range e.Fields {
		_ = enc.EncodeKeyval(field.Name, field.Value)
	}

	_ = enc.EndRecord()

	msg := []byte(fmt.Sprintf("<%d>%s %s %s[%d]: %s\n", syslog.LOG_KERN, ts, h.Hostname, h.Tag, os.Getpid(), buf.String()))

	h.mu.Lock()
	_, err := h.conn.Write(msg)
	h.mu.Unlock()

	return err
}
