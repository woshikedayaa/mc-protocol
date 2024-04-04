// Package pkgb is a package builder for minecraft protocol
package pkgb

import (
	"encoding/binary"
	"github.com/woshikedayaa/mc-protocol/internal/conv"
)

type Builder struct {
	p *PKG
}

func B() *Builder {
	b := &Builder{}
	b.p = NewPackage(binary.LittleEndian)
	return b
}

func (b *Builder) PackageID(v uint32) *Builder {
	b.p.WriteOnIndex(4, conv.LUint32(v))
	return b
}

func (b *Builder) AppendString(s string) *Builder {
	b.Append([]byte(s))
	return b
}

func (b *Builder) Append(bs []byte) *Builder {
	_, _ = b.p.Write(bs)
	return b
}

func (b *Builder) Build() *PKG {
	// if u want to write more data to package
	// write code on here
	b.p.build()
	//
	b.p.READONLY()
	return b.p
}
