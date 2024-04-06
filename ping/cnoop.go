package ping

import (
	"net"
)

type cNoop struct{}

func (c *cNoop) Ping(conn net.Conn) Response {
	return nil
}
