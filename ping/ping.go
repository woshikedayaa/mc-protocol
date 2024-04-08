package ping

import (
	"github.com/woshikedayaa/mc-protocol/internal/ver"
	"net"
)

type Client struct {
	impl IPing
	conn net.Conn
	op   *optionType

	server string
}

func NewClient(server string, ops ...option) (*Client, error) {
	var (
		c   = new(Client)
		err error
	)
	c.server = server
	c.op = new(optionType)
	c.op.ops = ops
	err = c.op.check(c)
	if err != nil {
		return nil, err
	}
	c.conn, err = net.Dial(c.op.network, server)
	if err != nil {
		return nil, err
	}
	c.impl = ver.chooseImpl(c.op.version)
	return c, nil
}
