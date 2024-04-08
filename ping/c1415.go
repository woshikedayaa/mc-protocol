package ping

import (
	"github.com/woshikedayaa/mc-protocol/internal/pkgb"
	"net"
)

type c1415 struct{ cNoop }

func (c *c1415) Ping(conn net.Conn) (Response, error) {
	//TODO implement me
	panic("implement me")
}

func (c *c1415) GetLatency(conn net.Conn) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (c *c1415) GetPackage() *pkgb.PKG {
	//TODO implement me
	panic("implement me")
}

func (c *c1415) HandShake(conn net.Conn) error {
	//TODO implement me
	panic("implement me")
}

func (c *c1415) Close(conn net.Conn) error {
	//TODO implement me
	panic("implement me")
}
