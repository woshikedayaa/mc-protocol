// Package rcon
// It was completed under document https://wiki.vg/RCON
package rcon

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
	"sync/atomic"
	"time"
)

var (
	uniqueID = atomic.Int32{}
)

type Client interface {

	// Close the RCON client
	Close() error

	// SendCommand will send the command to the target server
	// it will return a Response and an error
	SendCommand(command string) (*Response, error)

	// Auth will auth with pwd
	// if success, it will return nil
	Auth(pwd string) error

	// ClientID return the current ClientID
	// see global value uniqueID
	ClientID() int32

	// Connection will return the tcp connection
	Connection() net.Conn

	// Reconnect will create a new tcp connection
	Reconnect() error
}

type BaseClient struct {
	conn   net.Conn
	server string
	isAuth bool
	// options

	timeout time.Duration
	network string
}

// encode payload to an available package
// PackageType is defined in response.go
func encode(payload string, pt PackageType) ([]byte, error) {
	size := int32(4 + 4 + len(payload) + 2)
	buf := bytes.NewBuffer([]byte{})
	writeArray := []interface{}{
		size, uniqueID.Load(), pt, []byte(payload), []byte{0, 0},
	}
	var err error
	for i := 0; i < len(writeArray); i++ {
		err = binary.Write(buf, binary.LittleEndian, writeArray[i])
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

// decode the response from target server to a Response
func decode(bs []byte) (*Response, error) {
	resp := &Response{}
	buffer := bytes.NewBuffer(bs)
	err := binary.Read(buffer, binary.LittleEndian, &resp.size)
	if err != nil {
		return nil, err
	}
	err = binary.Read(buffer, binary.LittleEndian, &resp.id)
	if err != nil {
		return nil, err
	}
	err = binary.Read(buffer, binary.LittleEndian, &resp.typ)
	if err != nil {
		return nil, err
	}
	resp.body = make([]byte, len(bs)-12)
	err = binary.Read(buffer, binary.LittleEndian, &resp.body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// recv
// at first , read 4 byte to get the length of this package
// and read-all from the connection
func (b *BaseClient) recv() ([]byte, error) {
	var (
		res = make([]byte, 4) // 4 for length(int32)
		err error
	)
	_, err = b.conn.Read(res)
	if err != nil {
		return nil, err
	}
	length := binary.LittleEndian.Uint32(res)
	res = make([]byte, length+4)
	binary.LittleEndian.PutUint32(res, length)
	_, err = b.conn.Read(res[4:])
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (b *BaseClient) send(bs []byte) error {
	uniqueID.Add(1)
	_, err := b.conn.Write(bs)
	return err
}

func (b *BaseClient) Close() error {
	return b.conn.Close()
}

func (b *BaseClient) SendCommand(command string) (*Response, error) {
	var (
		encodeRes []byte
		err       error
		recv      []byte
	)
	encodeRes, err = encode(command, TypeCommand)
	if err != nil {
		return nil, err
	}
	err = b.conn.SetDeadline(time.Now().Add(b.timeout))
	if err != nil {
		return nil, err
	}
	err = b.send(encodeRes)
	if err != nil {
		return nil, err
	}
	recv, err = b.recv()
	if err != nil {
		return nil, err
	}
	return decode(recv)
}

func (b *BaseClient) Auth(pwd string) error {
	if b.isAuth {
		return nil
	}
	var (
		encodeRes []byte
		err       error
		response  *Response
		recv      []byte
	)
	encodeRes, err = encode(pwd, TypeAuthorize)
	if err != nil {
		return err
	}
	err = b.conn.SetDeadline(time.Now().Add(b.timeout))
	if err != nil {
		return err
	}

	err = b.send(encodeRes)
	if err != nil {
		return err
	}
	recv, err = b.recv()
	if err != nil {
		return err
	}

	response, err = decode(recv)
	if err != nil {
		return err
	}
	if response.id == -1 {
		return errors.New("auth fail with password=" + pwd)
	}

	b.isAuth = true
	return nil
}

func (b *BaseClient) ClientID() int32 {
	return uniqueID.Load()
}

func (b *BaseClient) Connection() net.Conn {
	return b.conn
}

func (b *BaseClient) Reconnect() error {
	var (
		conn net.Conn
		err  error
	)
	conn, err = net.Dial(b.network, b.server)
	if err != nil {
		return err
	}
	_ = b.conn.Close()
	b.conn = conn
	return nil
}

// NewRconClient will enable the ops
// create a new tcp connection for RCON
// notice: it will not authorize to the target server
// Use: Client.Auth  to authorize to the target server
func NewRconClient(server string, ops ...Option) (Client, error) {
	var (
		conn net.Conn
		err  error
		c    = &BaseClient{}
	)
	for _, v := range append(defaultOptions, ops...) {
		v.apply(c)
	}
	conn, err = net.Dial(c.network, server)
	if err != nil {
		return nil, err
	}
	c.conn = conn
	c.server = server
	c.isAuth = false

	return c, nil
}
