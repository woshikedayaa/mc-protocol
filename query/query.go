// Package query
// It was completed under  document https://wiki.vg/Query
package query

import (
	"net"
	"time"
)

const (
	// magic each package will start with the magic number
	magic uint16 = 0xFEFD
)

// see below
type queryType byte

const (
	// HandShakeType for get token
	// see Client.RefreshToken
	HandShakeType queryType = 9
	// StatQueryType for get server stats
	StatQueryType queryType = 0
)

type Client interface {
	// Token return current token.
	// warning: token maybe has been expired.
	// you can use RefreshToken to refresh token.
	// use IsTokenExpire to check does token expired.
	Token() int32

	RefreshToken() error

	// IsTokenExpire return true when token has expired.
	IsTokenExpire() bool

	// SessionID return the client unique-sessionID
	// usually , it is time-based.
	// see NewQueryClient
	SessionID() int32

	// FullRequest return a full server stat,
	// compared with BasicRequest, there are more players,plugins,version and gameID.
	// see FullResponse and BasicResponse
	FullRequest() (*FullResponse, error)

	// BasicRequest return a basic server stat.
	// effective than FullRequest
	BasicRequest() (*BasicResponse, error)

	// HandShakeRequest to get challenge-token
	// about challenge-token: https://wiki.vg/Query
	HandShakeRequest() (*HandleShakeResponse, error)

	// Close is optional.
	// Because query network is based on UDP.
	// UDP is connectionless.
	Close() error
}

type BaseClient struct {
	conn   *net.UDPConn
	server string

	lastRefresh int64 // second
	sessionID   int32
	cachedToken int32
	// options

	network string
	timeout time.Duration
}

func NewQueryClient(server string, ops ...Option) (Client, error) {
	c := &BaseClient{}
	for _, v := range append(defaultOptions, ops...) {
		v.apply(c)
	}
	addr, err := net.ResolveUDPAddr(c.network, server)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP(c.network, nil, addr)
	if err != nil {
		return nil, err
	}
	c.sessionID = int32(time.Now().Unix()) & 0x0F0F0F0F
	c.conn = conn
	c.server = server
	return c, nil
}

func (b *BaseClient) Close() error {
	return b.conn.Close()
}

func (b *BaseClient) Token() int32 {
	return b.cachedToken
}

func (b *BaseClient) RefreshToken() error {
	resp, err := b.sendAndRecv(HandShakeType, false)
	if err != nil {
		return err
	}
	hresp := &HandleShakeResponse{}
	err = hresp.decode(resp)
	if err != nil {
		return err
	}
	b.lastRefresh = time.Now().Unix()
	b.cachedToken = hresp.token
	return nil
}

func (b *BaseClient) IsTokenExpire() bool {
	return time.Now().Unix()-b.lastRefresh >= 30
}

func (b *BaseClient) SessionID() int32 {
	return b.sessionID
}

func (b *BaseClient) FullRequest() (*FullResponse, error) {
	res := &FullResponse{}
	recv, err := b.sendAndRecv(StatQueryType, true)
	if err != nil {
		return nil, err
	}
	err = res.decode(recv)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (b *BaseClient) BasicRequest() (*BasicResponse, error) {
	res := &BasicResponse{}
	recv, err := b.sendAndRecv(StatQueryType, false)
	if err != nil {
		return nil, err
	}
	err = res.decode(recv)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (b *BaseClient) HandShakeRequest() (*HandleShakeResponse, error) {
	res := &HandleShakeResponse{}
	recv, err := b.sendAndRecv(HandShakeType, true)
	if err != nil {
		return nil, err
	}
	err = res.decode(recv)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (b *BaseClient) sendAndRecv(qt queryType, isFull bool) ([]byte, error) {
	var (
		err  error
		pkg  Package
		recv []byte
	)
	err = b.conn.SetDeadline(time.Now().Add(b.timeout)) // set timeout
	if err != nil {
		return nil, err
	}

	if b.IsTokenExpire() && qt != HandShakeType {
		err = b.RefreshToken()
		if err != nil {
			return nil, err
		}
	}

	if isFull {
		pkg = newFullPackage(b, qt)
	} else {
		pkg = newPackage(b, qt)
	}
	err = b.send(pkg)
	if err != nil {
		return nil, err
	}
	recv, err = b.recv()
	if err != nil {
		return nil, err
	}
	return recv, nil
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
	n, _, err := b.conn.ReadFromUDP(buf)
	if err != nil {
		return nil, err
	}
	return append([]byte{}, buf[:n]...), nil
}
