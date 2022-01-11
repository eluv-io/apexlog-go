package log_test

import (
	"fmt"
	"testing"

	"github.com/apex/log"
	"github.com/apex/log/handlers/discard"
	"github.com/apex/log/handlers/memory"
	"github.com/stretchr/testify/assert"
)

func TestLogger_printf(t *testing.T) {
	h := memory.New()

	l := &log.Logger{
		Handler: h,
		Level:   log.InfoLevel,
	}

	l.Infof("logged in %s", "Tobi")

	e := h.Entries[0]
	assert.Equal(t, e.Message, "logged in Tobi")
	assert.Equal(t, e.Level, log.InfoLevel)
}

func TestLogger_levels(t *testing.T) {
	h := memory.New()

	l := &log.Logger{
		Handler: h,
		Level:   log.InfoLevel,
	}

	l.Debug("uploading")
	l.Info("upload complete")

	assert.Equal(t, 1, len(h.Entries))

	e := h.Entries[0]
	assert.Equal(t, e.Message, "upload complete")
	assert.Equal(t, e.Level, log.InfoLevel)
}

func TestLogger_WithFields(t *testing.T) {
	h := memory.New()

	l := &log.Logger{
		Handler: h,
		Level:   log.InfoLevel,
	}

	ctx := l.WithFields(log.Fields{{Name: "file", Value: "sloth.png"}})
	ctx.Debug("uploading")
	ctx.Info("upload complete")

	assert.Equal(t, 1, len(h.Entries))

	e := h.Entries[0]
	assert.Equal(t, e.Message, "upload complete")
	assert.Equal(t, e.Level, log.InfoLevel)
	assert.Equal(t, log.Fields{{Name: "file", Value: "sloth.png"}}, e.Fields)
}

func TestLogger_WithField(t *testing.T) {
	h := memory.New()

	l := &log.Logger{
		Handler: h,
		Level:   log.InfoLevel,
	}

	ctx := l.WithField("file", "sloth.png").WithField("user", "Tobi")
	ctx.Debug("uploading")
	ctx.Info("upload complete")

	assert.Equal(t, 1, len(h.Entries))

	e := h.Entries[0]
	assert.Equal(t, e.Message, "upload complete")
	assert.Equal(t, e.Level, log.InfoLevel)
	assert.Equal(t, log.Fields{{Name: "file", Value: "sloth.png"}, {Name: "user", Value: "Tobi"}}, e.Fields)
}

func TestLogger_Watch_info(t *testing.T) {
	h := memory.New()

	l := &log.Logger{
		Handler: h,
		Level:   log.InfoLevel,
	}

	func() (err error) {
		defer l.WithField("file", "sloth.png").Watch("upload").Stop(&err)
		return nil
	}()

	assert.Equal(t, 2, len(h.Entries))

	{
		e := h.Entries[0]
		assert.Equal(t, e.Message, "upload")
		assert.Equal(t, e.Level, log.InfoLevel)
		assert.Equal(t, log.Fields{{Name: "file", Value: "sloth.png"}}, e.Fields)
	}

	{
		e := h.Entries[1]
		assert.Equal(t, e.Message, "upload")
		assert.Equal(t, e.Level, log.InfoLevel)
		assert.Equal(t, "sloth.png", e.Fields.Get("file"))
		assert.IsType(t, int64(0), e.Fields.Get("duration"))
	}
}

func TestLogger_Watch_error(t *testing.T) {
	h := memory.New()

	l := &log.Logger{
		Handler: h,
		Level:   log.InfoLevel,
	}

	func() (err error) {
		defer l.WithField("file", "sloth.png").Watch("upload").Stop(&err)
		return fmt.Errorf("boom")
	}()

	assert.Equal(t, 2, len(h.Entries))

	{
		e := h.Entries[0]
		assert.Equal(t, e.Message, "upload")
		assert.Equal(t, e.Level, log.InfoLevel)
		assert.Equal(t, "sloth.png", e.Fields.Get("file"))
	}

	{
		e := h.Entries[1]
		assert.Equal(t, e.Message, "upload")
		assert.Equal(t, e.Level, log.ErrorLevel)
		assert.Equal(t, "sloth.png", e.Fields.Get("file"))
		assert.Equal(t, "boom", e.Fields.Get("error"))
		assert.IsType(t, int64(0), e.Fields.Get("duration"))
	}
}

func TestLogger_Watch_nil(t *testing.T) {
	h := memory.New()

	l := &log.Logger{
		Handler: h,
		Level:   log.InfoLevel,
	}

	func() {
		defer l.WithField("file", "sloth.png").Watch("upload").Stop(nil)
	}()

	assert.Equal(t, 2, len(h.Entries))

	{
		e := h.Entries[0]
		assert.Equal(t, e.Message, "upload")
		assert.Equal(t, e.Level, log.InfoLevel)
		assert.Equal(t, log.Fields{{Name: "file", Value: "sloth.png"}}, e.Fields)
	}

	{
		e := h.Entries[1]
		assert.Equal(t, e.Message, "upload")
		assert.Equal(t, e.Level, log.InfoLevel)
		assert.Equal(t, "sloth.png", e.Fields.Get("file"))
		assert.IsType(t, int64(0), e.Fields.Get("duration"))
	}
}

func TestLogger_HandlerFunc(t *testing.T) {
	h := memory.New()
	f := func(e *log.Entry) error {
		return h.HandleLog(e)
	}

	l := &log.Logger{
		Handler: log.HandlerFunc(f),
		Level:   log.InfoLevel,
	}

	l.Infof("logged in %s", "Tobi")

	e := h.Entries[0]
	assert.Equal(t, e.Message, "logged in Tobi")
	assert.Equal(t, e.Level, log.InfoLevel)
}

func TestFields_Append(t *testing.T) {
	var f log.Fields
	f = f.Append("k1", 1).Append("k2", "v2")
	assert.Equal(t, 2, len(f))
	assert.Equal(t, f[0], &log.Field{Name: "k1", Value: 1})
	assert.Equal(t, f[1], &log.Field{Name: "k2", Value: "v2"})
}

func BenchmarkLogger_small(b *testing.B) {
	l := &log.Logger{
		Handler: discard.New(),
		Level:   log.InfoLevel,
	}

	for i := 0; i < b.N; i++ {
		l.Info("login")
	}
}

func BenchmarkLogger_medium(b *testing.B) {
	l := &log.Logger{
		Handler: discard.New(),
		Level:   log.InfoLevel,
	}

	for i := 0; i < b.N; i++ {
		l.WithFields(log.Fields{
			{Name: "file", Value: "sloth.png"},
			{Name: "type", Value: "image/png"},
			{Name: "size", Value: 1 << 20},
		}).Info("upload")
	}
}

func BenchmarkLogger_large(b *testing.B) {
	l := &log.Logger{
		Handler: discard.New(),
		Level:   log.InfoLevel,
	}

	err := fmt.Errorf("boom")

	for i := 0; i < b.N; i++ {
		l.WithFields(log.Fields{
			{Name: "file", Value: "sloth.png"},
			{Name: "type", Value: "image/png"},
			{Name: "size", Value: 1 << 20},
		}).WithFields(log.Fields{
			{Name: "some", Value: "more"},
			{Name: "data", Value: "here"},
			{Name: "whatever", Value: "blah blah"},
			{Name: "more", Value: "stuff"},
			{Name: "context", Value: "such useful"},
			{Name: "much", Value: "fun"},
		}).WithError(err).Error("upload failed")
	}
}
