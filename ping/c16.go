package ping

import (
	"github.com/woshikedayaa/mc-protocol/internal/pkgb"
	"net"
)

type c16 struct{ cNoop }

func (c *c16) StatusRequest(conn net.Conn) (Response, error) {
	//TODO implement me
	panic("implement me")
}

func (c *c16) Latency(conn net.Conn) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (c *c16) GetPackage() *pkgb.PKG {
	//TODO implement me
	panic("implement me")
}

func (c *c16) HandShake(conn net.Conn) error {
	//TODO implement me
	panic("implement me")
}

func (c *c16) Close(conn net.Conn) error {
	//TODO implement me
	panic("implement me")
}
