package pkgb

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/woshikedayaa/mc-protocol/internal/conv"
)

type PKG struct {
	defaultOrder binary.ByteOrder
	buf          *bytes.Buffer
	isRO         bool // read-only
}

func NewPackage(defaultOrder binary.ByteOrder) *PKG {
	p := new(PKG)
	p.defaultOrder = defaultOrder
	p.isRO = false
	// add length
	p.buf.Write(conv.Uint32(0, binary.LittleEndian))
	// package id
	p.buf.Write(conv.Uint32(0, binary.LittleEndian))
	return p
}

func (p *PKG) Write(bs []byte) (n int, err error) {
	p.check()
	return p.buf.Write(bs)
}

func (p *PKG) Grow(n int) {
	p.check()
	p.buf.Grow(n)
}

func (p *PKG) WriteOnIndex(idx int, bs []byte) {
	copy(p.Bytes()[idx:idx+len(bs)], bs)
}

func (p *PKG) Bytes() []byte {
	return p.buf.Bytes()
}

func (p *PKG) READONLY() bool {
	p.isRO = true
	return true
}

func (p *PKG) Len() int {
	return p.buf.Len()
}

func (p *PKG) build() {
	p.check()
	// calc length
	p.WriteOnIndex(0, conv.Uint32(uint32(p.Len()-4), binary.LittleEndian))
}

func (p *PKG) check() {
	if p.isRO {
		panic(errors.New("panic when try to write data on a read-only package"))
	}
}
