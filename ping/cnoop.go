package ping

import (
	"errors"
	"github.com/woshikedayaa/mc-protocol/internal/pkgb"
	"github.com/woshikedayaa/mc-protocol/internal/ver"
	"net"
)

var errNOOP = errors.New("noop implementation client")

type cNoop struct {
	handShaken bool
	v          ver.Version
}

func (c *cNoop) StatusRequest(conn net.Conn) (Response, error) {
	return nil, errNOOP
}

func (c *cNoop) Latency(conn net.Conn) (int, error) {
	return 0, errNOOP
}

func (c *cNoop) GetPackage() *pkgb.PKG {
	return nil
}

func (c *cNoop) HandShake(conn net.Conn) error {
	return errNOOP
}

func (c *cNoop) Close(conn net.Conn) error {
	return errNOOP
}

func (c *cNoop) IsHandShaken() bool {
	return c.handShaken
}

func (c *cNoop) SetVersion(v ver.Version) {
	c.v = v
}
