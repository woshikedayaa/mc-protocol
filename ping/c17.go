package ping

import (
	"encoding/binary"
	"errors"
	"github.com/woshikedayaa/mc-protocol/internal/pkgb"
	"net"
	"net/netip"
	"time"
)

type c17 struct{ cNoop }

func (c *c17) Ping(conn net.Conn) (Response, error) {
	if !c.IsHandShaken() {
		return nil, errors.New("before ping, handshake required")
	}
	var err error
	// em...
	// maybe impossible
	pkg := c.GetPackage()
	if pkg == nil {
		return nil, errors.New("nil package")
	}
	_, err = conn.Write(pkg.Bytes())
	if err != nil {
		return nil, err
	}
	// todo Response struct
	return nil, nil
}

func (c *c17) GetLatency(conn net.Conn) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (c *c17) GetPackage() *pkgb.PKG {
	return pkgb.B().
		PackageID(0x00).
		Build()
}

func (c *c17) HandShake(conn net.Conn) error {
	// todo handShake
	if c.IsHandShaken() {
		return errors.New("already handshake")
	}
	// addr
	var (
		addr netip.AddrPort
		err  error
	)
	addr, err = netip.ParseAddrPort(conn.RemoteAddr().String())
	if err != nil {
		return err
	}
	// pkg
	pkg := pkgb.B().
		PackageID(0x00).
		AppendUint32()

	// finally
	c.handShaken = true
	return nil
}

func (c *c17) Close(conn net.Conn) error {
	if !c.IsHandShaken() {
		return errors.New("before close, handshake required")
	}
	var err error
	_, err = conn.Write(pkgb.B().
		PackageID(0x01).
		AppendUint64(uint64(time.Now().UnixMilli()), binary.LittleEndian).
		Build().Bytes())
	if err != nil {
		return err
	}
	// todo valid timestamp from server
	return conn.Close()
}
