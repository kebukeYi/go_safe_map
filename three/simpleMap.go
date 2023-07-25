package three

import "sync"

// SimpleMap is a simple map with atomic operations.
type SimpleMap[Key, Value any] struct {
	mu *sync.RWMutex
	mp map[any]Value
}

func (m *SimpleMap[Key, Value]) Load(key Key) (value Value, ok bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	got, gotten := m.mp[key]
	return got, gotten
}

func (m *SimpleMap[Key, Value]) Store(key Key, value Value) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.mp[key] = value
}

func (m *SimpleMap[Key, Value]) Delete(key Key) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.mp, key)
}

func (m *SimpleMap[Key, Value]) Range(f func(key Key, value Value) bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, v := range m.mp {
		if !f(k.(Key), v) {
			return
		}
	}
}

func newSimpleMap[Key, Value any]() IMap[Key, Value] {
	return &SimpleMap[Key, Value]{
		mu: &sync.RWMutex{},
		mp: map[any]Value{},
	}
}
