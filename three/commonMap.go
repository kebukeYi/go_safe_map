package three

import (
	"sync"
)

// CommonMap  is a type that implements several common map operations, including set, get, and remove.
type CommonMap[Key, Value any] struct {
	mu *sync.Mutex
	m  map[any]Value
}

func (c *CommonMap[Key, Value]) Load(key Key) (value Value, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	got, gotten := c.m[key]
	return got, gotten
}

func (c *CommonMap[Key, Value]) Store(key Key, value Value) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[key] = value
}

func (c *CommonMap[Key, Value]) Delete(key Key) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.m, key)
}

func (c *CommonMap[Key, Value]) Range(f func(key Key, value Value) bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, v := range c.m {
		if !f(k.(Key), v) {
			return
		}
	}
}

func newCommonMap[Key, Value any]() IMap[Key, Value] {
	return &CommonMap[Key, Value]{
		mu: &sync.Mutex{},
		m:  make(map[any]Value),
	}
}
