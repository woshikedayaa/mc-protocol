package tcpc

type DataPack interface {
	Unpack() []byte
	Len() int
	String() string
}

type ByteDataPack []byte

func (d ByteDataPack) Unpack() []byte {
	return d
}

func (d ByteDataPack) Len() int {
	return len(d)
}

func (d ByteDataPack) String() string {
	return string(d)
}

func NewByteDataPack(b []byte) DataPack {
	if b == nil {
		return nil
	}
	return ByteDataPack(b)
}
