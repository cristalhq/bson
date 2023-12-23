// Package bsonproto provides primitives for encoding and decoding of BSON.
package bsonproto

import (
	"errors"
	"fmt"
	"time"
)

// ScalarType represents a BSON scalar type.
type ScalarType interface {
	float64 | string | Binary | ObjectID | bool | time.Time | NullType | Regex | int32 | Timestamp | int64
}

// Size returns a size of the encoding of value v in bytes.
func Size[T ScalarType](v T) int {
	return SizeAny(v)
}

// SizeAny returns a size of the encoding of value v in bytes.
//
// It panics if v is not a ScalarType.
func SizeAny(v any) int {
	switch v := v.(type) {
	case float64:
		return SizeFloat64(v)
	case string:
		return SizeString(v)
	case int32:
		return SizeInt32(v)
	case int64:
		return SizeInt64(v)
	default:
		panic(fmt.Sprintf("unsupported type %T", v))
	}
}

// Encode encodes value v into b.
//
// b must be at least Size(v) bytes long; otherwise, Encode will panic.
// Only b[0:Size(v)] bytes are modified.
func Encode[T ScalarType](b []byte, v T) {
	EncodeAny(b, v)
}

// EncodeAny encodes value v into b.
//
// b must be at least Size(v) bytes long; otherwise, EncodeAny will panic.
// Only b[0:Size(v)] bytes are modified.
//
// It panics if v is not a ScalarType.
func EncodeAny(b []byte, v any) {
	switch v := v.(type) {
	case float64:
		EncodeFloat64(b, v)
	case string:
		EncodeString(b, v)
	case int32:
		EncodeInt32(b, v)
	case int64:
		EncodeInt64(b, v)
	default:
		panic(fmt.Sprintf("unsupported type %T", v))
	}
}

// Decode decodes value from b into v.
//
// If there is not enough bytes, Decode will return a wrapped ErrDecodeShortInput.
// If the input is otherwise invalid, a wrapped ErrDecodeInvalidInput is returned.
//
// If the value can't be decoded, a wrapped ErrDecodeInvalidInput is returned.
func Decode[T ScalarType](b []byte, v *T) error {
	return DecodeAny(b, v)
}

// DecodeAny decodes value from b into v.
//
// If there is not enough bytes, DecodeAny will return a wrapped ErrDecodeShortInput.
// If the input is otherwise invalid, a wrapped ErrDecodeInvalidInput is returned.
//
// It panics if v is not a pointer to ScalarType.
func DecodeAny(b []byte, v any) error {
	var err error
	switch v := v.(type) {
	case *float64:
		*v, err = DecodeFloat64(b)
	case *string:
		*v, err = DecodeString(b)
	case *int32:
		*v, err = DecodeInt32(b)
	case *int64:
		*v, err = DecodeInt64(b)
	default:
		panic(fmt.Sprintf("unsupported type %T", v))
	}

	return err
}

var (
	// ErrDecodeShortInput is returned wrapped by Decode functions if the input bytes slice is too short.
	ErrDecodeShortInput = errors.New("bsonproto: short input")

	// ErrDecodeInvalidInput is returned wrapped by Decode functions if the input bytes slice is invalid.
	ErrDecodeInvalidInput = errors.New("bsonproto: invalid input")
)
