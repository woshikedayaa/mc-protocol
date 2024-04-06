package ping

import "net"

type Client struct {
	impl IPing
	conn net.Conn
}

func NewClient(v string) (*Client, error) {
	c := &Client{}
	vs, err := newVersion(v)
	if err != nil {
		return nil, err
	}
	c.impl = chooseImpl(vs)
	// connect
	return c, nil
}
