package ping

import (
	"github.com/woshikedayaa/mc-protocol/internal/pkgb"
	"github.com/woshikedayaa/mc-protocol/internal/ver"
	"net"
)

type Response interface {
	GetLength() int
	GetBody() []byte
}

type IPing interface {
	StatusRequest(conn net.Conn) (Response, error)
	Latency(conn net.Conn) (int, error)

	GetPackage() *pkgb.PKG
	HandShake(conn net.Conn) error
	Close(conn net.Conn) error
	IsHandShaken() bool
	SetVersion(version ver.Version)
}
