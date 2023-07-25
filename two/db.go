package two

import (
	"sync"
	"time"
)

type IDataBase interface {
	Get(key interface{}) (value interface{}, timeout time.Duration, err error)
	Set(key interface{}, value interface{}, timeout time.Duration) (v interface{}, t time.Duration, e error)
	Delete(key interface{}) (ok bool, err error)
	Range(f func(key, value interface{}))
	Close()
}

type database struct {
	size    uint
	workers uint
	buffer  chan *dbRequest
	exit    chan struct{}
	data    map[uint]*dirtyMap
	hash    func(interface{}) uint
}

func NewIDataBase(workers, bufferSize, dirtyCount uint, hash func(key interface{}) uint) IDataBase {
	db := &database{
		buffer:  make(chan *dbRequest, bufferSize),
		exit:    make(chan struct{}),
		data:    make(map[uint]*dirtyMap, dirtyCount),
		size:    dirtyCount,
		workers: workers,
		hash:    hash,
	}
	for i := uint(0); i < dirtyCount; i++ {
		db.data[i] = &dirtyMap{
			lock:  sync.RWMutex{},
			dirty: make(map[interface{}]*entity),
		}
	}
	db.operator()
	return db
}

func (db *database) Get(key interface{}) (value interface{}, timeout time.Duration, err error) {
	req := newDbRequest(getValue)
	defer close(req.callback)
	req.args[requestKey] = key
	db.buffer <- req
	resp := <-req.callback
	switch resp.responseType {
	case withValue:
		result := resp.returns[responseValue]
		timeout := resp.returns[responseTimeout].(time.Duration)
		return result, timeout, nil
	case withError:
		err := resp.returns[responseError].(error)
		return nil, time.Duration(0), err
	default:
		return nil, time.Duration(0), nil
	}
}
func (db *database) Set(key interface{}, value interface{}, timeout time.Duration) (v interface{}, t time.Duration, e error) {
	req := newDbRequest(setValue)
	defer close(req.callback)
	req.args[requestKey] = key
	req.args[requestValue] = value
	req.args[requestTimeout] = timeout
	db.buffer <- req
	resp := <-req.callback
	switch resp.responseType {
	case withValue:
		return resp.returns[responseValue], resp.returns[responseTimeout].(time.Duration), nil
	case withError:
		err := resp.returns[responseError].(error)
		return nil, time.Duration(0), err
	default:
		return nil, time.Duration(0), nil
	}
}
func (db *database) Delete(key interface{}) (ok bool, err error) {
	req := newDbRequest(deleteKey)
	defer close(req.callback)
	req.args[requestKey] = key
	db.buffer <- req
	resp := <-req.callback
	switch resp.responseType {
	case withValue:
		return true, nil
	case withError:
		return false, resp.returns[responseError].(error)
	default:
		return false, nil
	}
}
func (db *database) Range(f func(key, value interface{})) {
	req := newDbRequest(rangeDB)
	defer close(req.callback)
	req.args[requestFunc] = f
	db.buffer <- req
	<-req.callback
	return
}
func (db *database) Close() {
	db.exit <- struct{}{}
}
