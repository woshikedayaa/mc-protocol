package ping

import (
	"net"
)

type c16 struct{ cNoop }

func (c *c16) Ping(conn net.Conn) Response {
	return nil
}
