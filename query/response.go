package query

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"unsafe"
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
	e.sessionID = int32(binary.LittleEndian.Uint32(bs[1 : unsafe.Sizeof(int32(0))+1]))
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

func (r *BasicResponse) Encode(bs []byte) error {
	var err error
	buffer := bytes.NewBuffer(bs[5:])
	r._MOTD, err = buffer.ReadString(0x00)
	if err != nil {
		return err
	}
	r.gameType, err = buffer.ReadString(0x00)
	if err != nil {
		return err
	}
	r._map, err = buffer.ReadString(0x00)
	if err != nil {
		return err
	}
	var num32 []byte
	num32, err = buffer.ReadBytes(0x00)
	if err != nil {
		return err
	}
	r.curPlayers = int(binary.LittleEndian.Uint32(num32))
	if err != nil {
		return err
	}
	num32, err = buffer.ReadBytes(0x00)
	if err != nil {
		return err
	}
	r.maxPlayer = int(binary.LittleEndian.Uint32(num32))
	port := []byte{0xDD, 0x63} // 25565
	_, err = buffer.Read(port)
	if err != nil {
		return err
	}
	r.port = binary.LittleEndian.Uint16(port)
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
	return r.EmptyResponse.Encode(bs)
}
