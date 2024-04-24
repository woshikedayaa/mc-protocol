package query

import (
	"bytes"
	"encoding/binary"
)

type Package interface {
	Encode() ([]byte, error)
}

type PackageQuery struct {
	typ       queryType
	isFull    bool
	sessionID int32
	token     int32
}

// Encode parse the package to a byte array for send udp package
func (p *PackageQuery) Encode() ([]byte, error) {
	/*
			Field-name	Field-Type	Notes
			Magic		uint16		Always 65277 (0xFEFD)
			Type		byte		9 for handshake, 0 for stat
			Session-ID	int32
		// stats only
			Token		Varies
		// full stats only
			Padding		Value:00 00 00 00
	*/

	buf := bytes.NewBuffer(make([]byte, 0))
	err := binary.Write(buf, binary.BigEndian, magic)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.LittleEndian, p.typ)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.BigEndian, p.sessionID)
	if err != nil {
		return nil, err
	}
	switch p.typ {
	case HandShakeType:
		// do nothing
	case StatQueryType:
		err = binary.Write(buf, binary.BigEndian, p.token)
		if err != nil {
			return nil, err
		}
		// padding
		if p.isFull {
			err = binary.Write(buf, binary.BigEndian, int32(0))
			if err != nil {
				return nil, err
			}
		}
	}
	return buf.Bytes(), nil
}

func newPackage(c Client, typ queryType) Package {
	return &PackageQuery{
		typ:       typ,
		isFull:    false,
		sessionID: c.SessionID(),
		token:     c.Token(),
	}
}

func newFullPackage(c Client, typ queryType) Package {
	return &PackageQuery{
		typ:       typ,
		isFull:    true,
		sessionID: c.SessionID(),
		token:     c.Token(),
	}
}
