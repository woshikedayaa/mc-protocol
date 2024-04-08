package ping

import (
	"github.com/woshikedayaa/mc-protocol/internal/ver"
	"net"
)

type Client struct {
	impl IPing
	conn net.Conn
	op   *options

	server string
}

func NewClient(server string, ops ...option) (*Client, error) {
	var (
		c   = new(Client)
		err error
	)
	c.server = server
	c.op = new(options)
	c.op.ops = ops
	err = c.op.check(c)
	if err != nil {
		return nil, err
	}
	c.conn, err = net.Dial(c.op.network, server)
	if err != nil {
		return nil, err
	}
	c.impl = chooseImpl(c.op.version)
	return c, nil
}

func chooseImpl(version ver.Version) IPing {
	if version.Minor() >= 7 {
		return new(c17)
	} else if version.Minor() >= 6 {
		return new(c16)
	} else if version.Minor() >= 4 {
		return new(c1415)
	} else {
		return new(cNoop)
	}
}
