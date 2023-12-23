package bsonproto

import (
	"encoding/binary"
	"fmt"
)

// SizeInt64 returns a size of the encoding of int64 in bytes - 8.
//
// The argument is unused.
func SizeInt64(int64) int {
	return 8
}

// EncodeInt64 encodes int64 value v into b.
//
// b must be at least 8 bytes long; otherwise, EncodeInt64 will panic.
// Only b[0:8] bytes are modified.
func EncodeInt64(b []byte, v int64) {
	binary.LittleEndian.PutUint64(b, uint64(v))
}

// DecodeInt64 decodes int64 value from b.
//
// If there is not enough bytes, DecodeInt64 will return a wrapped ErrDecodeShortInput.
func DecodeInt64(b []byte) (int64, error) {
	if len(b) < 8 {
		return 0, fmt.Errorf("DecodeInt64: expected at least 8 bytes, got %d: %w", len(b), ErrDecodeShortInput)
	}

	return int64(binary.LittleEndian.Uint64(b)), nil
}
