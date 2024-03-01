package rcon

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
	"sync/atomic"
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
}

type BaseClient struct {
	conn   net.Conn
	server string
	isAuth bool
}

func (b *BaseClient) encode(payload string, pt PackageType) ([]byte, error) {
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

func (b *BaseClient) decode(bs []byte) (*Response, error) {
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
	encode, err := b.encode(command, TypeCommand)
	if err != nil {
		return nil, err
	}
	err = b.send(encode)
	if err != nil {
		return nil, err
	}
	var recv []byte
	recv, err = b.recv()
	if err != nil {
		return nil, err
	}
	return b.decode(recv)
}

func (b *BaseClient) Auth(pwd string) error {
	encode, err := b.encode(pwd, TypeAuthorize)
	if err != nil {
		return err
	}
	err = b.send(encode)
	if err != nil {
		return err
	}
	var recv []byte
	recv, err = b.recv()
	if err != nil {
		return err
	}
	response, err := b.decode(recv)
	if err != nil {
		return err
	}
	if response.id == -1 {
		return errors.New("auth fail with password=" + pwd)
	}
	return nil
}

func (b *BaseClient) ClientID() int32 {
	return uniqueID.Load()
}

func (b *BaseClient) Connection() net.Conn {
	return b.conn
}

func NewClient(server string) (Client, error) {
	var (
		conn net.Conn
		err  error
	)
	conn, err = net.Dial("tcp4", server)
	if err != nil {
		return nil, err
	}

	return &BaseClient{
		server: server,
		conn:   conn,
	}, nil
}
