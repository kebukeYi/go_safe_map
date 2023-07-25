package three

import "sync"

type Table[Key, Value any] struct {
	mu    sync.RWMutex
	lines map[any]Value
}

func (t *Table[Key, Value]) get(key Key) (value Value, ok bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	value, ok = t.lines[key]
	return value, ok
}
func (t *Table[Key, Value]) set(key Key, value Value) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.lines[key] = value
}
func (t *Table[Key, Value]) delete(key Key) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.lines, key)
}
func (t *Table[Key, Value]) rangeDB(f func(key Key, value Value) bool) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	for k, v := range t.lines {
		if !f(k.(Key), v) {
			return false
		}
	}
	return true
}

func newTable[Key, Value any]() *Table[Key, Value] {
	return &Table[Key, Value]{
		lines: make(map[any]Value),
		mu:    sync.RWMutex{},
	}
}
