package conv

import "encoding/binary"

func LUint32(u uint32) []byte {
	return Uint32(u, binary.LittleEndian)
}

func BUint32(u uint32) []byte {
	return Uint32(u, binary.BigEndian)
}

func Uint32(u uint32, order binary.ByteOrder) []byte {
	res := make([]byte, 4)
	order.PutUint32(res, u)
	return res
}
