package encoder

import "encoding/binary"

var ByteOrder = binary.LittleEndian

const (
	_ = (2 << (iota + 1)) >> 3
	Uint8Size
	Uint16Size
	Uint32Size
	Uint64Size
)

const (
	_ = (2 << (iota + 1)) >> 3
	Int8Size
	Int16Size
	Int32Size
	Int64Size
)

func MakeUint32() []byte {
	return make([]byte, Uint32Size)
}

func BinaryUint32(input uint32) []byte {
	buff := MakeUint32()
	ByteOrder.PutUint32(buff, input)
	return buff
}

func Uint32Binary(input []byte) uint32 {
	return ByteOrder.Uint32(input)
}
