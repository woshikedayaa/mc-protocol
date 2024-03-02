package query

type Response interface {
	Encode([]byte, queryType) error
}

type EmptyResponse struct {
	typ       queryType
	SessionID int32
}

func (e *EmptyResponse) Encode([]byte, queryType) error {
	return nil
}

type BasicResponse struct {
	EmptyResponse
}

type FullResponse struct {
	EmptyResponse
}

type HandleShakeResponse struct {
	EmptyResponse
	newToken int32
}
