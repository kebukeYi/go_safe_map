package two

import (
	"sync"
	"time"
)

type entity struct {
	data    interface{}
	timeout time.Time
}

type dirtyMap struct {
	lock  sync.RWMutex
	dirty map[interface{}]*entity
}
