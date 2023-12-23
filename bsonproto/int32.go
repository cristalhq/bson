package bsonproto

import (
	"encoding/binary"
	"fmt"
)

// SizeInt32 returns a size of the encoding of int32 in bytes - 4.
//
// The argument is unused.
func SizeInt32(int32) int {
	return 4
}

// EncodeInt32 encodes int32 value v into b.
//
// b must be at least 4 bytes long; otherwise, EncodeInt32 will panic.
// Only b[0:4] bytes are modified.
func EncodeInt32(b []byte, v int32) {
	binary.LittleEndian.PutUint32(b, uint32(v))
}

// DecodeInt32 decodes int32 value from b.
//
// If there is not enough bytes, DecodeInt32 will return a wrapped ErrDecodeShortInput.
func DecodeInt32(b []byte) (int32, error) {
	if len(b) < 4 {
		return 0, fmt.Errorf("DecodeInt32: expected at least 4 bytes, got %d: %w", len(b), ErrDecodeShortInput)
	}

	return int32(binary.LittleEndian.Uint32(b)), nil
}
