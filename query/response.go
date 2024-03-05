package query

type Response interface {
	Encode([]byte) error
}

type EmptyResponse struct {
	typ       queryType
	sessionID int32
}

func (e *EmptyResponse) Encode([]byte) error {
	return nil
}

type BasicResponse struct {
	EmptyResponse
	_MOTD      string
	gameType   string
	_map       string
	curPlayers int
	maxPlayer  int
	port       uint16
	ip         string
}

type FullResponse struct {
	EmptyResponse
	_MOTD      string
	gameType   string
	_map       string
	curPlayers int
	maxPlayer  int
	port       uint16
	ip         string
	// extend
	player  []string
	version string
}

type HandleShakeResponse struct {
	EmptyResponse
	token int32
}

func (r *BasicResponse) Encode([]byte) error {

}

func (r *FullResponse) Encode([]byte) error {

}

func (r *HandleShakeResponse) Encode([]byte) error {

}
