package query

import (
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
	FullRequest() (*FullResponse, error)
	BasicRequest() (*BasicResponse, error)
	HandShake() (*HandleShakeResponse, error)
	Close() error
}

type BaseClient struct {
	conn   *net.UDPConn
	server string
	port   int

	lastRefresh int64
	sessionID   int32
	cachedToken int32
	// options

	timeout time.Duration
}

func NewQueryClient(server string, ops ...Option) (Client, error) {
	c := &BaseClient{}
	c.server = server
	addr, err := net.ResolveUDPAddr("udp", server)
	if err != nil {
		return nil, err
	}
	c.sessionID = int32(time.Now().Unix()) & 0x0F0F0F0F
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}
	c.conn = conn
	for _, v := range append(defaultOptions, ops...) {
		v.apply(c)
	}
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
	recv, err := b.sendAndRecv(StatType, true)
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
	recv, err := b.sendAndRecv(StatType, false)
	if err != nil {
		return nil, err
	}
	err = res.decode(recv)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (b *BaseClient) HandShake() (*HandleShakeResponse, error) {
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
