package query

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

type Response interface {
	Encode([]byte) error
}

type EmptyResponse struct {
	typ       queryType
	sessionID int32
}

func (e *EmptyResponse) Encode(bs []byte) error {
	if len(bs) < 5 {
		return errors.New("response bytes length to short")
	}
	// parse the typ and sessionID
	e.typ = queryType(bs[0])
	// big-Ending
	e.sessionID = int32(uint32(bs[1])<<24 | uint32(bs[2])<<16 | uint32(bs[3])<<8 | uint32(bs[4]))
	return nil
}

type BasicResponse struct {
	EmptyResponse
	_MOTD      string
	gameType   string
	_map       string
	curPlayers int // current-player
	maxPlayer  int
	port       uint16
	ip         string // alias hostname
}

type FullResponse struct {
	EmptyResponse
	_MOTD      string
	gameType   string
	_map       string
	curPlayers int // current-player
	maxPlayer  int
	port       uint16
	ip         string // alias hostname
	// extend
	player  []string
	version string
}

type HandleShakeResponse struct {
	EmptyResponse
	token int32
}

// Encode BasicResponse
func (r *BasicResponse) Encode(bs []byte) error {
	var err error
	// 5 for sessionID and queryType
	buffer := bytes.NewBuffer(bs[5:])
	// motd
	r._MOTD, err = buffer.ReadString(0x00)
	if err != nil {
		return err
	}
	// game-type
	r.gameType, err = buffer.ReadString(0x00)
	if err != nil {
		return err
	}
	// map
	r._map, err = buffer.ReadString(0x00)
	if err != nil {
		return err
	}
	// playerNum for numPlayer and maxPlayer
	var playerNum []byte
	// curPlayer
	playerNum, err = buffer.ReadBytes(0x00)
	if err != nil {
		return err
	}
	// 这个 playerNum 其实是长度不定的
	// 1 <= len(playerNum) <= 4    [1,4]
	// 要使用位运算来算
	// r.curPlayers = int(binary.LittleEndian.Uint32(playerNum))
	for i := 0; i < len(playerNum); i++ {
		r.curPlayers = r.curPlayers<<8 | int(playerNum[i])
	}
	if err != nil {
		return err
	}
	playerNum, err = buffer.ReadBytes(0x00)
	if err != nil {
		return err
	}
	// 同理
	// r.maxPlayer = int(binary.LittleEndian.Uint32(playerNum))
	for i := 0; i < len(playerNum); i++ {
		r.maxPlayer = r.maxPlayer<<8 | int(playerNum[i])
	}
	// port
	port := []byte{0xDD, 0x63} // 25565
	_, err = buffer.Read(port)
	if err != nil {
		return err
	}
	r.port = binary.LittleEndian.Uint16(port)
	// hostname
	r.ip, err = buffer.ReadString(0x00)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	return r.EmptyResponse.Encode(bs)
}

func (r *FullResponse) Encode(bs []byte) error {
	return r.EmptyResponse.Encode(bs)
}

func (r *HandleShakeResponse) Encode(bs []byte) error {
	buffer := bytes.NewBuffer(bs[5:])
	token, err := buffer.ReadString(0x00)
	if err != nil {
		return err
	}
	r.token = parseTokenString(token)
	return r.EmptyResponse.Encode(bs)
}
