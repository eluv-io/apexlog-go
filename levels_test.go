package log

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseLevel(t *testing.T) {
	cases := []struct {
		String string
		Level  Level
		Num    int
	}{
		{"trace", TraceLevel, 0},
		{"debug", DebugLevel, 1},
		{"info", InfoLevel, 2},
		{"warn", WarnLevel, 3},
		{"warning", WarnLevel, 4},
		{"error", ErrorLevel, 5},
		{"fatal", FatalLevel, 6},
	}

	for _, c := range cases {
		t.Run(c.String, func(t *testing.T) {
			l, err := ParseLevel(c.String)
			assert.NoError(t, err, "parse")
			assert.Equal(t, c.Level, l)
		})
	}

	t.Run("invalid", func(t *testing.T) {
		l, err := ParseLevel("something")
		assert.Equal(t, ErrInvalidLevel, err)
		assert.Equal(t, InvalidLevel, l)
	})
}

func TestLevel_MarshalJSON(t *testing.T) {
	e := Entry{
		Level:   InfoLevel,
		Message: "hello",
		Fields:  Fields{{Name: "name", Value: "bob"}, {Name: "foo", Value: "bar"}},
	}

	expect := `{"fields":{"name":"bob","foo":"bar"},"level":"info","timestamp":"0001-01-01T00:00:00Z","message":"hello"}`

	b, err := json.Marshal(e)
	assert.NoError(t, err)
	assert.Equal(t, expect, string(b))
}

func TestLevel_UnmarshalJSON(t *testing.T) {
	s := `{"fields":{"name":"bob"},"level":"info","timestamp":"0001-01-01T00:00:00Z","message":"hello"}`
	e := new(Entry)

	err := json.Unmarshal([]byte(s), e)
	assert.NoError(t, err)
	assert.Equal(t, InfoLevel, e.Level)
	assert.Equal(t, "hello", e.Message)
	assert.Equal(t, "bob", e.Fields.Get("name"))
}
