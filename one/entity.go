package one

import (
	"sync"
	"time"
)

type entity struct {
	mu      sync.RWMutex
	timeout time.Time
	data    interface{}
}

func (e *entity) write() {
	e.mu.Lock()
}

func (e *entity) writeDone() {
	e.mu.Unlock()
}

func (e *entity) read() {
	e.mu.RUnlock()
}

func (e *entity) readDone() {
	e.mu.RUnlock()
}
