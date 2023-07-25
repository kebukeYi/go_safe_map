package one

type requestType int

const (
	getValue requestType = iota
	setValue
	deleteKey
	rangDB
)

type responseType int

const (
	withError responseType = iota
	withValue
)

type requestFieldIndex int

const (
	requestKey requestFieldIndex = iota
	requestValue
	requestTimeout
	requestFunc
)

type responseFieldIndex int

const (
	responseKey responseFieldIndex = iota
	responseValue
	responseTimeout
	responseError
)

type dbRequest struct {
	typeRequest requestType
	args        map[requestFieldIndex]interface{}
	callback    chan *dbResponse
}

type dbResponse struct {
	typeResponse responseType
	returns      map[responseFieldIndex]interface{}
}

func newDbRequest(rType requestType) *dbRequest {
	return &dbRequest{
		typeRequest: rType,
		args:        make(map[requestFieldIndex]interface{}),
		callback:    make(chan *dbResponse, 1),
	}
}

func newDbResponse(rType responseType) *dbResponse {
	return &dbResponse{
		typeResponse: rType,
		returns:      make(map[responseFieldIndex]interface{}),
	}
}
