package query

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"strconv"
)

type Response interface {

	// decode bind a byte array to Response
	decode([]byte) error

	// JSON return jsonify data
	JSON() ([]byte, error)

	// SessionID return the response's sessionID
	SessionID() int32

	IsStatQuery() bool
}

type EmptyResponse struct {
	typ       queryType
	sessionID int32
}

func (e *EmptyResponse) SessionID() int32 {
	return e.sessionID
}

func (e *EmptyResponse) IsStatQuery() bool {
	return e.typ == StatQueryType
}

func (e *EmptyResponse) decode(bs []byte) error {
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
	port       int
	ip         string // alias hostname
}

type FullResponse struct {
	EmptyResponse
	_MOTD      string
	gameType   string
	_map       string
	curPlayers int // current-player
	maxPlayer  int
	port       int
	ip         string // alias hostname
	// extend than BasicResponse
	plugins string
	gameID  string
	player  []string
	version string
}

type HandleShakeResponse struct {
	EmptyResponse
	token int32
}

func (r *BasicResponse) decode(bs []byte) error {
	var (
		lastRead = 5 // start
		n        = len(bs)
		find     = func() []byte {
			for i := lastRead; i < n; i++ {
				if bs[i] == 0x00 {
					if i == lastRead {
						lastRead++
						return nil
					}
					res := bs[lastRead:i]
					lastRead = i + 1 // skip 0x00
					return res
				}
			}
			return nil
		}
	)
	r._MOTD = string(find())
	r.gameType = string(find())
	r._map = string(find())
	r.curPlayers = byteArrayToInt(find())
	r.maxPlayer = byteArrayToInt(find())
	tempLastRead := lastRead
	r.port = int(binary.LittleEndian.Uint16(find())) // 这里应该只读2个Byte
	lastRead = tempLastRead + 2
	r.ip = string(find())
	return r.EmptyResponse.decode(bs)
}

func (r *FullResponse) decode(bs []byte) error {
	Skip1 := 1 + 4 + 11 // head and padding
	Skip2 := 10         // padding-2
	Skip2 += r.parseKVString(bs[Skip1:]) + Skip1
	r.parsePlayerString(bs[Skip2:])
	return r.EmptyResponse.decode(bs)
}

func (r *HandleShakeResponse) decode(bs []byte) error {
	var err error
	err = r.EmptyResponse.decode(bs)
	if err != nil {
		return err
	}

	if r.IsStatQuery() {
		return errors.New("except QueryType is HandShakeType,but current QueryType is StatQueryType")
	}
	buffer := bytes.NewBuffer(bs[5:])
	token, err := buffer.ReadBytes(0x00)
	if err != nil {
		return err
	}

	r.token = r.parseTokenString(string(token))
	return nil
}

func (r *HandleShakeResponse) parseTokenString(s string) int32 {
	atoi, _ := strconv.Atoi(s[:len(s)-1])
	return int32(atoi)
}

// parseKVString return the length of KVString
func (r *FullResponse) parseKVString(bs []byte) int {
	n := len(bs)
	var lastRead int
	var find = func() []byte {
		keyFind := false
		ValueStart := lastRead
		for i := lastRead; i < n; i++ {
			if bs[i] == 0x00 {
				if lastRead == i { // end
					lastRead = i + 1 // skip
					return nil
				}
				if keyFind {
					lastRead = i + 1 // skip 0x00
					return bs[ValueStart:i]
				}
				keyFind = true
				// skip 0x00
				ValueStart = i + 1
			}
		}
		return nil
	}

	r._MOTD = string(find())
	r.gameType = string(find())
	r.gameID = string(find())
	r.version = string(find())
	r.plugins = string(find())
	r._map = string(find())
	r.curPlayers = byteArrayToInt(find())
	r.maxPlayer = byteArrayToInt(find())
	r.port = byteArrayToInt(find())
	r.ip = string(find())

	for len(find()) != 0 {
	} // find last
	return lastRead
}

func (r *FullResponse) parsePlayerString(bs []byte) {
	var (
		err error
		buf = bytes.NewBuffer(bs)
		s   = ""
	)
	for len(s) != 0 || !errors.Is(err, io.EOF) {
		s, err = buf.ReadString(0x00)
		if len(s) != 0 {
			r.player = append(r.player, s)
		}
	}
	// avoid memory-leak
	r.player = append([]string{}, r.player[:len(r.player)-1]...)
	if err != nil && !errors.Is(err, io.EOF) {
		panic(err)
	}
}

// ascii code to int
func byteArrayToInt(bs []byte) (ans int) {
	for i := 0; i < len(bs); i++ {
		ans = ans*10 + int(bs[i]-'0')
	}
	return
}
