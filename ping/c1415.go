package ping

import (
	"net"
)

type c1415 struct{ cNoop }

func (c *c1415) Ping(conn net.Conn) Response {
	return nil
}
