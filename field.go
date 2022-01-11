package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"
)

var (
	fieldsPool        sync.Pool
	encBufferPool     = newBufferPool()
	defaultBufferSize = 1024
)

type Field struct {
	pool  bool
	Name  string
	Value interface{}
}

// newField returns a Field initialized with the given name and value.
// The field is retrieved from a pool and will be released back to the pool once
// logging occurred.
func newField(name string, value interface{}) *Field {
	var e *Field
	if v := fieldsPool.Get(); v != nil {
		e = v.(*Field)
	} else {
		e = new(Field)
	}
	e.pool = true
	e.Reset(name, value)
	return e
}

func (f *Field) release() {
	if f.pool {
		f.Reset("", nil)
		fieldsPool.Put(f)
	}
}

func (f *Field) Reset(name string, value interface{}) {
	f.Name = name
	f.Value = value
}

func (f *Field) _toJSON() ([]byte, error) {
	val, err := json.Marshal(f.Value)
	if err != nil {
		return nil, err
	}
	s := fmt.Sprintf("\"%s\":%v", f.Name, string(val))
	return []byte(s), nil
}

func (f *Field) toJSON() ([]byte, error) {
	bp := encBufferPool.get()
	defer bp.release()
	buf := bytes.NewBuffer(bp.bs)

	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(f.Value)
	if err != nil {
		return nil, err
	}
	s := fmt.Sprintf("\"%s\": %v", f.Name, string(buf.Bytes()))
	return []byte(s), nil
}

// bufferPool is a pool of buffers
type bufferPool struct {
	p *sync.Pool
}

// newBufferPool constructs a new BufferPool.
func newBufferPool() bufferPool {
	return bufferPool{p: &sync.Pool{
		New: func() interface{} {
			return &buffer{bs: make([]byte, 0, defaultBufferSize)}
		},
	}}
}

// Get retrieves a buffer from the pool, creating one if necessary.
func (p bufferPool) get() *buffer {
	buf := p.p.Get().(*buffer)
	buf.reset()
	buf.pool = p
	return buf
}

func (p bufferPool) put(buf *buffer) {
	p.p.Put(buf)
}

// buffer is a wrapper around a byte slice.
type buffer struct {
	bs   []byte
	pool bufferPool
}

// reset resets the underlying byte slice.
func (b *buffer) reset() {
	b.bs = b.bs[:0]
}

// release returns the buffer to its Pool.
func (b *buffer) release() {
	b.pool.put(b)
}
