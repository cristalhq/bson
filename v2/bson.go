// Package bson implements encoding and decoding of BSON as defined by https://bsonspec.org/spec.html
// and https://www.mongodb.com/docs/manual/reference/bson-types/.
//
// # Types
//
// The following BSON types are supported:
//
//	BSON                Go
//
//	Double              float64
//	String              string
//	Object              *bson.Object
//	Array               *bson.Array
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
// Composite types (Object and Array) are passed by pointers. Scalar types are passed by values.
package bson
