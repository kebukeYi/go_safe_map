package one

import (
	"sync"
	"time"
)

type IDataBase interface {
	// Get 查
	Get(key interface{}) (value interface{}, duration time.Duration, err error)
	// Set 增改
	Set(key, value interface{}, duration time.Duration) (v interface{}, t time.Duration, err error)
	// Delete 删
	Delete(key interface{}) (ok bool, err error)

	Range(f func(key, value interface{}))

	Close()
}

type database struct {
	mu    sync.RWMutex
	buff  chan *dbRequest
	exit  chan struct{}
	dirty map[interface{}]*entity
}

func NewDatabase(workers, bufferSize int) IDataBase {
	db := &database{
		buff:  make(chan *dbRequest, bufferSize),
		exit:  make(chan struct{}, 1),
		dirty: make(map[interface{}]*entity),
		mu:    sync.RWMutex{},
	}
	db.operator(workers)
	return db
}

func (db database) Get(key interface{}) (value interface{}, timeout time.Duration, err error) {
	req := newDbRequest(getValue)
	defer close(req.callback)
	req.args[requestKey] = key
	db.buff <- req
	resp := <-req.callback
	switch resp.typeResponse {
	case withError:
		err := resp.returns[responseError].(error)
		return nil, time.Duration(0), err
	default:
		value := resp.returns[responseValue]
		times := resp.returns[responseTimeout].(time.Time)
		out := times.Sub(time.Now())
		return value, out, nil
	}
}

func (db database) Set(key, value interface{}, duration time.Duration) (v interface{}, t time.Duration, err error) {
	req := newDbRequest(setValue)
	defer close(req.callback)
	req.args[requestKey] = key
	req.args[requestValue] = value
	req.args[requestTimeout] = time.Now().Add(duration)
	db.buff <- req
	resp := <-req.callback
	switch resp.typeResponse {
	case withError:
		err := resp.returns[responseError].(error)
		return nil, time.Duration(0), err
	default:
		value := resp.returns[responseValue]
		timeout := resp.returns[responseTimeout].(*time.Time)
		out := timeout.Sub(time.Now())
		return value, out, nil
	}
}

func (db database) Delete(key interface{}) (ok bool, err error) {
	req := newDbRequest(deleteKey)
	defer close(req.callback)
	db.buff <- req
	resp := <-req.callback
	switch resp.typeResponse {
	case withError:
		err := resp.returns[responseError].(error)
		return false, err
	default:
		ok = true
		return ok, err
	}
}

func (db database) Range(f func(key, value interface{})) {
	req := newDbRequest(rangDB)
	defer close(req.callback)
	req.args[requestFunc] = f
	db.buff <- req
	<-req.callback
	return
}

func (db database) Close() {
	db.exit <- struct{}{}
}

func (db *database) operator(workers int) {
	if workers == 1 {
		go func(da *database) {
			for {
				select {
				case request := <-db.buff:
					switch request.typeRequest {
					case getValue:
						key := request.args[requestKey]
						if entity, ok := db.dirty[key]; !ok {
							callback := newDbResponse(withError)
							callback.returns[responseError] = notFoundError{key: key}
							request.callback <- callback
						} else {
							if entity.timeout.Before(time.Now()) {
								delete(db.dirty, key)
								callback := newDbResponse(withError)
								callback.returns[responseError] = notFoundError{key: key}
								request.callback <- callback
							} else {
								callback := newDbResponse(withValue)
								callback.returns[responseValue] = &entity.data
								callback.returns[responseTimeout] = &entity.timeout
								request.callback <- callback
							}
						}
					case setValue:
						key := request.args[requestKey]
						value := request.args[requestValue]
						timeout := request.args[requestTimeout].(time.Time)
						entity := &entity{data: value, timeout: timeout, mu: sync.RWMutex{}}
						db.dirty[key] = entity
						callback := newDbResponse(withValue)
						callback.returns[responseValue] = &value
						callback.returns[responseTimeout] = &timeout
						request.callback <- callback
					case deleteKey:
						key := request.args[requestKey]
						if _, ok := db.dirty[key]; ok {
							delete(db.dirty, key)
						}
						callback := newDbResponse(withValue)
						callback.returns = nil
						request.callback <- callback
					case rangDB:
						function := request.args[requestFunc].(func(key, value interface{}))
						for i, e := range db.dirty {
							function(i, e.data)
						}
						callback := newDbResponse(withValue)
						callback.returns = nil
						request.callback <- callback
					} // switch request.typeRequest
				case <-db.exit:
					return
				} // select chan chanRequest
			} // for
		}(db) // go func
	} else { // workers != 1
		//
		exitChan := make(chan struct{}, workers)
		for i := 0; i < workers; i++ {
			go func(e chan struct{}, db *database) {
				for {
					select {
					case request := <-db.buff:
						switch request.typeRequest {
						case getValue:
							db.mu.RLock()
							key := request.args[requestKey]
							if entity, ok := db.dirty[key]; !ok {
								callback := newDbResponse(withError)
								callback.returns[responseError] = notFoundError{key: key}
								request.callback <- callback
							} else {
								entity.read()
								if entity.timeout.Before(time.Now()) {
									delete(db.dirty, key)
									callback := newDbResponse(withError)
									callback.returns[responseError] = keyExpired{key: key}
									request.callback <- callback
								} else {
									callback := newDbResponse(withValue)
									callback.returns[responseValue] = &entity.data
									callback.returns[responseTimeout] = &entity.timeout
									request.callback <- callback
								}
								entity.readDone()
							}
							db.mu.RUnlock()
						case setValue:
							db.mu.Lock()
							key := request.args[requestKey]
							value := request.args[requestValue]
							timeout_ := request.args[requestTimeout].(time.Time)
							en := &entity{
								timeout: timeout_,
								data:    value,
								mu:      sync.RWMutex{},
							}
							db.dirty[key] = en
							callback := newDbResponse(withValue)
							callback.returns[responseTimeout] = en.timeout
							en.writeDone()
							request.callback <- callback
							db.mu.Unlock()
						case deleteKey:
							db.mu.Lock()
							key := request.args[requestKey]
							callback := newDbResponse(withError)
							if entity, ok := db.dirty[key]; ok {
								entity.write()
								delete(db.dirty, key)
								callback.typeResponse = withValue
								entity.writeDone()
							}
							callback.returns = nil
							request.callback <- callback
							db.mu.Unlock()
						case rangDB:
							db.mu.Lock()
							function := request.args[requestFunc].(func(key, value interface{}))
							for i2, e2 := range db.dirty {
								e2.write()
								function(i2, e2.data)
								e2.writeDone()
							}
							callback := newDbResponse(withValue)
							callback.returns = nil
							request.callback <- callback
							db.mu.Unlock()
						} // switch request.typeRequest
					case <-exitChan:
						return
					} // select
				} // for
			}(exitChan, db)
		}

		//
		go func(e chan struct{}, d *database, w int) {
			select {
			case <-db.exit:
				for i := 0; i < workers; i++ {
					e <- struct{}{}
				}
			}
		}(exitChan, db, workers)
	}
}
