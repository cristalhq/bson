package bsonproto

import (
	"encoding/binary"
	"io"
)

func SizeInt64(int64) int {
	return 8
}

func DecodeInt64(b []byte) (int64, error) {
	if len(b) < 8 {
		return 0, io.ErrUnexpectedEOF
	}

	return int64(binary.LittleEndian.Uint64(b)), nil
}

func EncodeInt64(b []byte, v int64) {
	binary.LittleEndian.PutUint64(b, uint64(v))
}
