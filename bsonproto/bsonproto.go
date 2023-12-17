package bsonproto

import (
	"fmt"
	"time"
)

// ScalarType represents a BSON scalar type.
type ScalarType interface {
	float64 | string | Binary | ObjectID | bool | time.Time | NullType | Regex | int32 | Timestamp | int64
}

type Binary struct{}

type ObjectID struct{}

type NullType struct{}

type Regex struct{}

type Timestamp struct{}

func Size[T ScalarType](v T) int {
	return SizeAny(v)
}

func SizeAny(v any) int {
	switch v := any(v).(type) {
	case string:
		return len(v) + 5
	case int32:
		return SizeInt32(v)
	case int64:
		return SizeInt64(v)
	default:
		panic(fmt.Sprintf("unsupported type %T", v))
	}
}

func Encode[T ScalarType](b []byte, v T) {
	EncodeAny(b, v)
}

func EncodeAny(b []byte, v any) {
	switch v := any(v).(type) {
	case int32:
		EncodeInt32(b, v)
	case int64:
		EncodeInt64(b, v)
	default:
		panic(fmt.Sprintf("unsupported type %T", v))
	}
}

func Decode[T ScalarType](v *T, b []byte) error {
	return DecodeAny(v, b)
}

func DecodeAny(v any, b []byte) error {
	var err error
	switch v := any(v).(type) {
	case *int32:
		*v, err = DecodeInt32(b)
	case *int64:
		*v, err = DecodeInt64(b)
	default:
		panic(fmt.Sprintf("unsupported type %T", v))
	}

	return err
}
