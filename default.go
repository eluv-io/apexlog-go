package log

import (
	"bytes"
	"fmt"
	"log"
)

// by sorts fields by name.
type byName []Field

func (a byName) Len() int           { return len(a) }
func (a byName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byName) Less(i, j int) bool { return a[i].Name < a[j].Name }

// handleStdLog outpouts to the stlib log.
func handleStdLog(e *Entry) error {
	level := levelNames[e.Level]

	var fields []Field

	for _, f := range e.Fields {
		fields = append(fields, *f)
	}

	//sort.Sort(byName(fields))

	var b bytes.Buffer
	_, _ = fmt.Fprintf(&b, "%5s %-25s", level, e.Message)

	for _, f := range fields {
		_, _ = fmt.Fprintf(&b, " %s=%v", f.Name, f.Value)
	}

	log.Println(b.String())

	return nil
}
