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

	// SessionID return the response's session
	SessionID() int32

	IsStatQuery() bool
}

type EmptyResponse struct {
	Typ     queryType `json:"Typ"`
	Session int32     `json:"session"`
}

func (e *EmptyResponse) SessionID() int32 {
	return e.Session
}

func (e *EmptyResponse) IsStatQuery() bool {
	return e.Typ == StatQueryType
}

func (e *EmptyResponse) decode(bs []byte) error {
	if len(bs) < 5 {
		return errors.New("response bytes length to short")
	}
	// parse the Typ and session
	e.Typ = queryType(bs[0])
	// big-Ending
	e.Session = int32(uint32(bs[1])<<24 | uint32(bs[2])<<16 | uint32(bs[3])<<8 | uint32(bs[4]))
	return nil
}

type BasicResponse struct {
	EmptyResponse
	MOTD       string `json:"motd"`
	GameType   string `json:"gameType"`
	Map        string `json:"map"`
	NumPlayers int    `json:"numPlayers"`
	MaxPlayer  int    `json:"maxPlayer"`
	Port       int    `json:"port"`
	HostName   string `json:"hostName"`
}

type FullResponse struct {
	BasicResponse
	// extend than BasicResponse
	Plugins string   `json:"plugins"`
	GameID  string   `json:"gameID"`
	Players []string `json:"players"`
	Version string   `json:"version"`
}

type HandleShakeResponse struct {
	EmptyResponse
	Token int32 `json:"token"`
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
	r.MOTD = string(find())
	r.GameType = string(find())
	r.Map = string(find())
	r.NumPlayers = byteArrayToInt(find())
	r.MaxPlayer = byteArrayToInt(find())
	tempLastRead := lastRead
	r.Port = int(binary.LittleEndian.Uint16(find())) // 这里应该只读2个Byte
	lastRead = tempLastRead + 2
	r.HostName = string(find())
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

	r.Token = r.parseTokenString(string(token))
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

	r.MOTD = string(find())
	r.GameType = string(find())
	r.GameID = string(find())
	r.Version = string(find())
	r.Plugins = string(find())
	r.Map = string(find())
	r.NumPlayers = byteArrayToInt(find())
	r.MaxPlayer = byteArrayToInt(find())
	r.Port = byteArrayToInt(find())
	r.HostName = string(find())

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
			r.Players = append(r.Players, s)
		}
	}
	// avoid memory-leak
	r.Players = append([]string{}, r.Players[:len(r.Players)-1]...)
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
