package tcpc

import (
	"context"
	"errors"
	"io"
	"net"
)

type SimpleTcpClient struct {
	conn net.Conn
	addr *net.TCPAddr
}

func (c *SimpleTcpClient) Connection() net.Conn {
	return c.conn
}

func (c *SimpleTcpClient) Write(ctx context.Context, pack DataPack) (int, error) {
	ch := c.writeAsChannel(pack)
	for {
		select {
		case e := <-ch:
			return pack.Len(), e
		case <-ctx.Done():
			return 0, errors.New("context done")
		}
	}
}

func (c *SimpleTcpClient) writeAsChannel(pack DataPack) chan error {
	ch := make(chan error)
	go func() {
		_, err := c.Connection().Write(pack.Unpack())
		ch <- err
	}()
	return ch
}

func (c *SimpleTcpClient) Read() (DataPack, error) {
	bytes, err := io.ReadAll(c.reader())
	return NewByteDataPack(bytes), err
}

func (c *SimpleTcpClient) reader() io.Reader {
	return io.Reader(c.Connection())
}

func (c *SimpleTcpClient) ReadN(n int) (DataPack, error) {
	if n <= 0 {
		return nil, nil
	}
	buf := make([]byte, n)
	count, err := io.ReadAtLeast(c.reader(), buf, n)
	return NewByteDataPack(buf[:count]), err
}

func (c *SimpleTcpClient) Close() error {
	return c.Connection().Close()
}

func newClient(network string, addrString string) (Client, error) {
	var (
		addr *net.TCPAddr
		conn net.Conn
		err  error
	)
	addr, err = net.ResolveTCPAddr(network, addrString)
	if err != nil {
		return nil, err
	}
	// TODO 完成属性
	conn, err = net.Dial(network, addr.String())
	return &SimpleTcpClient{addr: addr, conn: conn}, err
}

func NewIPV4Client(addrString string) (Client, error) {
	return newClient("tcp4", addrString)
}

func NewIPV6Client(addrString string) (Client, error) {
	return newClient("tcp6", addrString)
}
