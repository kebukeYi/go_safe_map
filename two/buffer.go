package two

type typeRequest int

const (
	getValue typeRequest = iota
	setValue
	deleteKey
	rangeDB
)

type typeResponse int

const (
	withError typeResponse = iota
	withValue
)

type requestFields int

const (
	requestKey requestFields = iota
	requestValue
	requestTimeout
	requestFunc
)

type responseFields int

const (
	responseKey responseFields = iota
	responseValue
	responseTimeout
	responseError
)

type dbRequest struct {
	requestType typeRequest
	args        map[requestFields]interface{}
	callback    chan dbResponse
}

type dbResponse struct {
	responseType typeResponse
	returns      map[responseFields]interface{}
}

func newDbRequest(rType typeRequest) *dbRequest {
	return &dbRequest{
		requestType: rType,
		args:        make(map[requestFields]interface{}),
		callback:    make(chan dbResponse, 1),
	}
}

func newDbResponse(rType typeResponse) dbResponse {
	return dbResponse{
		responseType: rType,
		returns:      make(map[responseFields]interface{}),
	}
}
