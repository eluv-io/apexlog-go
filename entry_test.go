package log

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEntry_WithFields(t *testing.T) {
	a := NewEntry(nil)
	assert.Nil(t, a.Fields)

	b := a.WithFields(Fields{{Name: "foo", Value: "bar"}})
	assert.Equal(t, Fields{}, a.mergedFields())
	assert.Equal(t, Fields{{Name: "foo", Value: "bar"}}, b.mergedFields())

	c := a.WithFields(Fields{{Name: "foo", Value: "hello"}, {Name: "bar", Value: "world"}})

	e := c.finalize(InfoLevel, "upload", false)
	assert.Equal(t, e.Message, "upload")
	assert.Equal(t, e.Fields, Fields{{Name: "foo", Value: "hello"}, {Name: "bar", Value: "world"}})
	assert.Equal(t, e.Level, InfoLevel)
	assert.NotEmpty(t, e.Timestamp)
}

func TestEntry_WithField(t *testing.T) {
	a := NewEntry(nil)
	b := a.WithField("foo", "bar")
	assert.Equal(t, Fields{}, a.mergedFields())
	assert.Equal(t, Fields{{Name: "foo", Value: "bar"}}, b.mergedFields())
}

func TestEntry_WithError(t *testing.T) {
	a := NewEntry(nil)
	b := a.WithError(fmt.Errorf("boom"))
	assert.Equal(t, Fields{}, a.mergedFields())
	assert.Equal(t, Fields{{Name: "error", Value: "boom"}}, b.mergedFields())
}

func TestEntry_WithError_fields(t *testing.T) {
	a := NewEntry(nil)
	b := a.WithError(errFields("boom"))
	assert.Equal(t, Fields{}, a.mergedFields())
	assert.Equal(t, Fields{
		{Name: "error", Value: "boom"},
		{Name: "reason", Value: "timeout"},
	}, b.mergedFields())
}

func TestEntry_WithError_nil(t *testing.T) {
	a := NewEntry(nil)
	b := a.WithError(nil)
	assert.Equal(t, Fields{}, a.mergedFields())
	assert.Equal(t, Fields{}, b.mergedFields())
}

func TestEntry_WithDuration(t *testing.T) {
	a := NewEntry(nil)
	b := a.WithDuration(time.Second * 2)
	assert.Equal(t, Fields{{Name: "duration", Value: int64(2000)}}, b.mergedFields())
}

type errFields string

func (ef errFields) Error() string {
	return string(ef)
}

func (ef errFields) Fields() Fields {
	return Fields{{Name: "reason", Value: "timeout"}}
}
