package bson

import (
	"bytes"
	"sort"
)

// Marshaler is the interface implemented by types that
// can marshal themselves into valid BSON.
type Marshaler interface {
	MarshalBSON() ([]byte, error)
}

// Marshal returns BSON encoding of v.
func Marshal(v any) ([]byte, error) {
	buf := &bytes.Buffer{}
	if err := NewEncoder(buf).Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// MarshalTo returns BSON encoding of v written to dst.
func MarshalTo(dst []byte, v any) ([]byte, error) {
	enc := &Encoder{buf: dst}
	if err := enc.marshal(v); err != nil {
		return nil, err
	}
	return enc.buf, nil
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

func (doc D) Len() int           { return len(doc) }
func (doc D) Less(i, j int) bool { return doc[i].K < doc[j].K }
func (doc D) Swap(i, j int)      { doc[i], doc[j] = doc[j], doc[i] }

// M is an unordered representation of a BSON document.
//
// Example usage:
//
//	bson.M{"hello": "world", "foo": "bar", "pi": 3.14159}
type M map[string]any

func (m M) AsD() D {
	d := make(D, len(m))
	i := 0
	for k, v := range m {
		d[i] = e{K: k, V: v}
		i++
	}

	sort.Sort(d)
	return d
}
