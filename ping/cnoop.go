package ping

import (
	"errors"
	"github.com/woshikedayaa/mc-protocol/internal/pkgb"
	"net"
)

var errNOOP = errors.New("noop implementation client")

type cNoop struct{}

func (c *cNoop) Ping(conn net.Conn) (Response, error) {
	return nil, errNOOP
}

func (c *cNoop) GetLatency(conn net.Conn) (int, error) {
	return 0, errNOOP
}

func (c *cNoop) GetPackage() *pkgb.PKG {
	return nil
}

func (c *cNoop) HandleShake(conn net.Conn) error {
	return errNOOP
}

func (c *cNoop) Close(conn net.Conn) error {
	return errNOOP
}
