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
