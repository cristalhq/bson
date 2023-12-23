package bsonproto

import (
	"encoding/binary"
	"fmt"
	"math"
)

// SizeFloat64 returns a size of the encoding of float64 in bytes - 8.
//
// The argument is unused.
func SizeFloat64(float64) int {
	return 8
}

// EncodeFloat64 encodes float64 value v into b.
//
// b must be at least 8 bytes long; otherwise, EncodeFloat64 will panic.
// Only b[0:8] bytes are modified.
//
// Infinities, NaNs, negative zeros are preserved.
func EncodeFloat64(b []byte, v float64) {
	binary.LittleEndian.PutUint64(b, math.Float64bits(float64(v)))
}

// DecodeFloat64 decodes float64 value from b.
//
// If there is not enough bytes, DecodeFloat64 will return a wrapped ErrDecodeShortInput.
//
// Infinities, NaNs, negative zeros are preserved.
func DecodeFloat64(b []byte) (float64, error) {
	if len(b) < 8 {
		return 0, fmt.Errorf("DecodeFloat64: expected at least 8 bytes, got %d: %w", len(b), ErrDecodeShortInput)
	}

	return math.Float64frombits(binary.LittleEndian.Uint64(b)), nil
}
