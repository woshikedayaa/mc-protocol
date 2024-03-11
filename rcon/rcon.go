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
	Close() error
	SendCommand(command string) (*Response, error)
	Auth(pwd string) error
	ClientID() int32
	Connection() net.Conn
	Reconnect() error
}

type BaseClient struct {
	conn   net.Conn
	server string
	isAuth bool
	// options

	timeout time.Duration
}

func encode(payload string, pt PackageType) ([]byte, error) {
	size := int32(4 + 4 + len(payload) + 2)
	buf := bytes.NewBuffer([]byte{})
	// use an array to extend in the future
	/*
		Copy from : https://wiki.vg/RCON
		Field name	Field-type	Notes
		Length		int32		Length of remainder of packet
		Request-ID	int32		Client-generated ID
		Type		int32		3 for login, 2 to run a command, 0 for a multi-packet response
		Payload		byte[]		NULL-terminated ASCII text
		1-byte pad	byte		NULL
	*/
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

func (b *BaseClient) recv() ([]byte, error) {
	var (
		bufSize = 4096
		buf     = make([]byte, bufSize)
		err     error
		n       int
	)
	n, err = b.conn.Read(buf)
	if err != nil {
		return nil, err
	}
	// avoid memory leak
	return append([]byte{}, buf[:n]...), nil
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
	conn, err = net.Dial("tcp", b.server)
	if err != nil {
		return err
	}
	if b.conn.Close() != nil {
		return errors.New("error when close old TCP connection")
	}
	b.conn = conn
	return nil
}

func NewRconClient(server string, ops ...Option) (Client, error) {
	var (
		conn net.Conn
		err  error
		c    = &BaseClient{}
	)
	conn, err = net.Dial("tcp", server)
	if err != nil {
		return nil, err
	}
	c.conn = conn
	c.server = server
	c.isAuth = false
	for _, v := range append(defaultOptions, ops...) {
		v.apply(c)
	}
	return c, nil
}
