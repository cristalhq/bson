package bson

// Marshaler is the interface implemented by types that
// can marshal themselves into valid BSON.
type Marshaler interface {
	MarshalBSON() ([]byte, error)
}

// Unmarshaler is the interface implemented by types that
// can unmarshal a BSON representation of themselves.
//
// The input can be assumed to be a valid encoding of a BSON.
// UnmarshalBSON must copy the BSON data if it wishes to retain the data after returning.
type Unmarshaler interface {
	UnmarshalBSON([]byte) error
}

// A is a BSON array.
//
// Example:
//
//	bson.A{"hello", "world", 3.14159, bson.D{{"foo", 12345}}}
type A []any

// D is an ordered representation of a BSON document.
//
// Example usage:
//
//	bson.D{{"hello", "world"}, {"foo", "bar"}, {"pi", 3.14159}}
type D []e

// e represents a BSON element for a D. It is usually used inside a D.
type e struct {
	K string
	V any
}

// M is an unordered representation of a BSON document.
//
// Example usage:
//
//	bson.M{"hello": "world", "foo": "bar", "pi": 3.14159}
type M map[string]any
