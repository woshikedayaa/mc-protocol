package ping

import (
	"github.com/woshikedayaa/mc-protocol/internal/pkgb"
	"net"
)

type Response interface {
	GetLength() int
	GetBody() []byte
}

type IPing interface {
	Ping(conn net.Conn) (Response, error)
	GetLatency(conn net.Conn) (int, error)
	GetPackage() *pkgb.PKG
	HandShake(conn net.Conn) error
	Close(conn net.Conn) error
	IsHandShaken() bool
}
