package bsonproto

import (
	"encoding/binary"
	"io"
)

func SizeInt32(int32) int {
	return 4
}

func DecodeInt32(b []byte) (int32, error) {
	if len(b) < 4 {
		return 0, io.ErrUnexpectedEOF
	}

	return int32(binary.LittleEndian.Uint32(b)), nil
}

func EncodeInt32(b []byte, v int32) {
	binary.LittleEndian.PutUint32(b, uint32(v))
}
