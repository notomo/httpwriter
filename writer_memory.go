package httpwriter

import (
	"bytes"
	"io"
	"net/http"
	"sync"
)

func NewMemoryWriter(memory *Memory) func(*http.Request) (io.WriteCloser, error) {
	return func(req *http.Request) (io.WriteCloser, error) {
		return memory.CreateBuffer(req), nil
	}
}

type Memory struct {
	Buffers []*Buffer

	mu sync.Mutex
}

func (m *Memory) CreateBuffer(req *http.Request) *Buffer {
	m.mu.Lock()
	defer m.mu.Unlock()

	b := &Buffer{}
	m.Buffers = append(m.Buffers, b)
	return b
}

type Buffer struct {
	bytes.Buffer
}

func (b *Buffer) Close() error {
	return nil
}
