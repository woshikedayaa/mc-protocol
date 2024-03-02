package query

import (
	"encoding/binary"
	"errors"
	"net"
	"time"
)

const (
	magic uint16 = 0xFEFD
)

type queryType byte

const (
	HandShakeType queryType = 9
	StatType      queryType = 0
)

type Client interface {
	Token() int32
	RefreshToken() error
	IsTokenExpire() bool
	SessionID() int32
	FullRequest() (Response, error)
	BasicRequest() (Response, error)
	HandShake() (Response, error)
}

type BaseClient struct {
	conn   net.Conn
	server string

	lastRefresh int64
	sessionID   int32
	cachedToken int32
}

func (b *BaseClient) Token() int32 {
	return b.cachedToken
}

func (b *BaseClient) RefreshToken() error {
	resp, err := b.sendAndRecv(HandShakeType, false)
	if err != nil {
		return err
	}
	hresp, ok := resp.(*HandleShakeResponse)
	if !ok {
		return errors.New("can not cast Response interface to HandleShakeResponse")
	}
	b.lastRefresh = time.Now().Unix()
	b.cachedToken = hresp.newToken
	return nil
}

func (b *BaseClient) IsTokenExpire() bool {
	return time.Now().Unix()-b.lastRefresh >= 30
}

func (b *BaseClient) SessionID() int32 {
	return b.sessionID
}

func (b *BaseClient) FullRequest() (Response, error) {
	return b.sendAndRecv(StatType, true)
}

func (b *BaseClient) BasicRequest() (Response, error) {
	return b.sendAndRecv(StatType, false)
}

func (b *BaseClient) HandShake() (Response, error) {
	return b.sendAndRecv(HandShakeType, false)
}

func (b *BaseClient) sendAndRecv(qt queryType, isFull bool) (Response, error) {

	var (
		err      error
		pkg      Package
		response Response
	)
	if b.IsTokenExpire() && qt != HandShakeType {
		err = b.RefreshToken()
		if err != nil {
			return nil, err
		}
	}

	if isFull {
		pkg = NewFullPackage(b, qt)
	} else {
		pkg = NewPackage(b, qt)
	}
	err = b.send(pkg)
	if err != nil {
		return nil, err
	}
	var recv []byte
	recv, err = b.recv()
	if err != nil {
		return nil, err
	}
	if qt == HandShakeType {
		response = &HandleShakeResponse{}
	} else {
		if isFull {
			response = &FullResponse{}
		} else {
			response = &BasicResponse{}
		}
	}
	err = response.Encode(recv, qt)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (b *BaseClient) send(p Package) error {
	var (
		err   error
		bytes []byte
	)
	bytes, err = p.Encode()
	if err != nil {
		return err
	}
	_, err = b.conn.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

func (b *BaseClient) recv() ([]byte, error) {
	var (
		bufSize = 4096
		buf     = make([]byte, bufSize)
	)
	n, err := b.conn.Read(buf)
	if err != nil {
		return nil, err
	}
	return append([]byte{}, buf[:n]...), nil
}

func parseTokenString(s string) int32 {
	return int32(binary.BigEndian.Uint32([]byte(s)))
}
