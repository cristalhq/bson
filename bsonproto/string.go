package bsonproto

import (
	"encoding/binary"
	"io"
)

func SizeString(s string) int {
	return len(s) + 5
}

func DecodeString(b []byte) (string, error) {
	if len(b) < 5 {
		return "", io.ErrUnexpectedEOF
	}

	l := int(binary.LittleEndian.Uint32(b))
	if len(b) < l {
		return "", io.ErrUnexpectedEOF
	}
	if b[l-1] != 0 {
		return "", io.ErrUnexpectedEOF
	}

	return string(b[4 : l-1]), nil
}

func EncodeString(b []byte, v string) {
	l := len(v) + 5

	b[l-1] = 0
	binary.LittleEndian.PutUint32(b, uint32(l))
	copy(b[4:], v)
}
