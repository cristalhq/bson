// Package bson implements encoding and decoding of BSON as defined by https://bsonspec.org/spec.html
// and https://www.mongodb.com/docs/manual/reference/bson-types/.
//
// # Types
//
// The following BSON types are supported:
//
//	BSON                Go
//
//	Object              *bson.Object or RawObject
//	Array               *bson.Array  or RawArray
//
//	Double              float64
//	String              string
//	Binary data         bson.Binary
//	ObjectId            bson.ObjectID
//	Boolean             bool
//	Date                time.Time
//	Null                bson.NullType
//	Regular Expression  bson.Regex
//	32-bit integer      int32
//	Timestamp           bson.Timestamp
//	64-bit integer      int64
//
// Composite types (Object and Array) are passed by pointers.
// Raw composite type and scalars are passed by values.
package bson

import (
	"github.com/cristalhq/bson/bsonproto"
)

// Type represents a BSON type.
type Type interface {
	bsonproto.ScalarType | CompositeType
}

// CompositeType represents a BSON composite type (including raw types).
type CompositeType interface {
	*Object | *Array | RawObject | RawArray
}

// validType checks if v is a valid BSON type (including raw types).
func validType(v any) bool {
	switch v.(type) {
	case *Object:
	case RawObject:
	case *Array:
	case RawArray:
	case float64:
	case string:
	case int32:
	case int64:
	default:
		return false
	}

	return true
}

const (
	TagFloat64   = byte(0x01)
	TagString    = byte(0x02)
	TagObject    = byte(0x03)
	TagArray     = byte(0x04)
	TagBinary    = byte(0x05)
	TagObjectID  = byte(0x07)
	TagBool      = byte(0x08)
	TagTime      = byte(0x09)
	TagNullType  = byte(0x0a)
	TagRegex     = byte(0x0b)
	TagInt32     = byte(0x10)
	TagTimestamp = byte(0x11)
	TagInt64     = byte(0x12)
)
