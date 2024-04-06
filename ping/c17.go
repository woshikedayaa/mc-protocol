package ping

import (
	"net"
)

type c17 struct{ cNoop }

func (c *c17) Ping(conn net.Conn) Response {
	return nil
}
