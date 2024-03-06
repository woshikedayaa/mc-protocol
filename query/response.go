package query

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"strconv"
)

type Response interface {
	encode([]byte) error
}

type EmptyResponse struct {
	typ       queryType
	sessionID int32
}

func (e *EmptyResponse) encode(bs []byte) error {
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
	curPlayers string // current-player
	maxPlayer  string
	port       string
	ip         string // alias hostname
	// extend
	plugins string
	gameID  string
	player  []string
	version string
}

type HandleShakeResponse struct {
	EmptyResponse
	token int32
}

// Encode BasicResponse
func (r *BasicResponse) encode(bs []byte) error {
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
	port := []byte{0xDD, 0x63} // default 25565
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
	return r.EmptyResponse.encode(bs)
}

func (r *FullResponse) encode(bs []byte) error {
	Skip1 := 1 + 4 + 11
	Skip2 := 10
	Skip2 += r.parseKVString(bs[Skip1:]) + Skip1
	r.parsePlayerString(bs[Skip2:])
	return r.EmptyResponse.encode(bs)
}

func (r *HandleShakeResponse) encode(bs []byte) error {
	// todo  验证是不是 HandleShake
	buffer := bytes.NewBuffer(bs[5:])
	token, err := buffer.ReadString(0x00)
	if err != nil {
		return err
	}
	r.token = r.parseTokenString(token)
	return r.EmptyResponse.encode(bs)
}

func (r *HandleShakeResponse) parseTokenString(s string) int32 {
	atoi, _ := strconv.Atoi(s)
	return int32(atoi)
}

func (r *FullResponse) parseKVString(bs []byte) int {
	n := len(bs)
	var lastRead int
	var find = func() string {
		keyFind := false
		ValueStart := lastRead
		for i := lastRead; i < n; i++ {
			if bs[i] == 0x00 {
				if lastRead == i { // end
					return ""
				}
				if keyFind {
					lastRead = i + 1 // skip 0x00
					return string(bs[ValueStart:i])
				}
				keyFind = true
				// skip 0x00
				ValueStart = i + 1
			}
		}
		return ""
	}

	r._MOTD = find()
	r.gameType = find()
	r.gameID = find()
	r.version = find()
	r.plugins = find()
	r._map = find()
	r.curPlayers = find()
	r.maxPlayer = find()
	r.port = find()
	r.ip = find()

	for len(find()) != 0 {
	} // find last
	return lastRead + 1 // skip
}

func (r *FullResponse) parsePlayerString(bs []byte) {
	var (
		err error
		buf = bytes.NewBuffer(bs)
		s   = ""
	)
	for !errors.Is(err, io.EOF) || len(s) != 0 {
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
