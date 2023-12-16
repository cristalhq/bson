package bson

// RawObject represents a BSON object a.k.a document in the binary encoded form.
type RawObject []byte

// field represents a single Object field in the (partially) decoded form.
type field struct {
	value any // BSON type, RawObject, or RawArray
	name  string
}

// Object represents a BSON object a.k.a document in the (partially) decoded form.
//
// It may contain duplicate field names.
type Object struct {
	fields []field
}
