package bsonproto

import (
	"encoding/binary"
	"fmt"
)

// SizeString returns a size of the encoding of s string in bytes - len(s)+5.
func SizeString(s string) int {
	return len(s) + 5
}

// EncodeString encodes string value v into b.
//
// b must be at least len(v)+5 bytes long; otherwise, EncodeString will panic.
// Only b[0:len(v)+5] bytes are modified.
func EncodeString(b []byte, v string) {
	l := len(v) + 5

	b[l-1] = 0
	binary.LittleEndian.PutUint32(b, uint32(l))
	copy(b[4:], v)
}

// DecodeString decodes string value from b.
//
// If there is not enough bytes, DecodeString will return a wrapped ErrDecodeShortInput.
// If the input is otherwise invalid, a wrapped ErrDecodeInvalidInput is returned.
func DecodeString(b []byte) (string, error) {
	if len(b) < 5 {
		return "", fmt.Errorf("DecodeString: expected at least 5 bytes, got %d: %w", len(b), ErrDecodeShortInput)
	}

	l := int(binary.LittleEndian.Uint32(b))
	if len(b) < l {
		return "", fmt.Errorf("DecodeString: expected at least %d bytes, got %d: %w", l, len(b), ErrDecodeShortInput)
	}
	if b[l-1] != 0 {
		return "", fmt.Errorf("DecodeString: expected the last byte to be 0: %w", ErrDecodeInvalidInput)
	}

	return string(b[4 : l-1]), nil
}
