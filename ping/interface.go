package ping

import "net"

type Response interface {
	GetBody() []byte
}

type IPing interface {
	Ping(conn net.Conn) Response
}
