package ping

import "net"

type Client struct {
	impl IPing
	conn net.Conn
	op   *optionType
}

func NewClient(server string, ops ...option) (*Client, error) {
	var (
		c   = new(Client)
		err error
	)
	c.op.ops = ops
	err = c.op.check(c)
	if err != nil {
		return nil, err
	}
	c.conn, err = net.Dial(c.op.network, server)
	if err != nil {
		return nil, err
	}

}
