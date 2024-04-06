package ping

import (
	"net"
)

type c16 struct{}

func (c *c16) Ping(conn net.Conn) Response {
}
