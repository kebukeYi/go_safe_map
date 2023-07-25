package two

import "time"

func (db *database) operator() {
	exitChan := make(chan struct{}, db.workers)
	go listenExitChan(db.exit, exitChan, db.workers)

	for i := uint(0); i < db.workers; i++ {
		go startupOperator(db, exitChan)
	}
}

func startupOperator(db *database, exit chan struct{}) {
	for {
		select {
		case request := <-db.buffer:
			switch request.requestType {
			case getValue:
				request.callback <- getValueOperator(db, request)
			case setValue:
				request.callback <- setValueOperator(db, request)
			case deleteKey:
				request.callback <- deleteKeyOperator(db, request)
			case rangeDB:
				request.callback <- rangeDBOperator(db, request)
			}
		case <-exit:
			return
		}
	}
}
func getValueOperator(db *database, req *dbRequest) dbResponse {
	key := req.args[requestKey]
	index := db.hash(key) % db.size
	map_ := db.data[index]
	map_.lock.RLock()
	defer map_.lock.RUnlock()
	resp := newDbResponse(withError)
	if entity, ok := map_.dirty[key]; !ok {
		resp.returns[responseError] = notFoundError{key: key}
	} else {
		if entity.timeout.Before(time.Now()) {
			map_.lock.Lock()
			defer map_.lock.Unlock()
			delete(map_.dirty, key)
			resp.returns[responseError] = keyExpiredError{key: key}
		} else {
			resp.responseType = withValue
			resp.returns[responseValue] = entity.data
			resp.returns[responseTimeout] = entity.timeout.Sub(time.Now())
		}
	}
	return resp
}
func setValueOperator(db *database, request *dbRequest) dbResponse {
	key := request.args[requestKey]
	value := request.args[requestValue]
	timeout := request.args[requestTimeout].(time.Duration)
	index := db.hash(key) % db.size
	map_ := db.data[index]
	en := &entity{
		data:    value,
		timeout: time.Now().Add(timeout),
	}
	map_.lock.Lock()
	defer map_.lock.Unlock()
	temp := map_.dirty[key]
	map_.dirty[key] = en
	resp := newDbResponse(withValue)
	resp.returns[responseValue] = temp
	resp.returns[responseTimeout] = en.timeout.Sub(time.Now())
	return resp
}
func deleteKeyOperator(db *database, request *dbRequest) dbResponse {
	key := request.args[requestKey]
	index := db.hash(key) % db.size
	map_ := db.data[index]
	map_.lock.RLock()
	defer map_.lock.RUnlock()
	resp := newDbResponse(withError)
	if _, ok := map_.dirty[key]; ok {
		map_.lock.Lock()
		defer map_.lock.Unlock()
		delete(map_.dirty, key)
		resp.responseType = withValue
	} else {
		resp.returns[responseError] = notFoundError{key: key}
	}
	return resp
}

func rangeDBOperator(db *database, request *dbRequest) dbResponse {
	f := request.args[requestFunc].(func(interface{}, interface{}))
	for i := uint(0); i < db.size; i++ {
		current := db.data[i]
		current.lock.Lock()
		for k, v := range current.dirty {
			f(k, v.data)
		}
		current.lock.Unlock()
	}
	resp := newDbResponse(withValue)
	return resp
}

func listenExitChan(from, to chan struct{}, workers uint) {
	select {
	case <-from:
		for i := uint(0); i < workers; i++ {
			to <- struct{}{}
		}
	}
}
