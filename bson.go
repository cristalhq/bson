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

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
)

type rawObject []byte

type rawArray []byte

type field struct {
	name  string
	value any // BSON type, rawObject, or rawArray
}

type Object struct {
	fields []field
}

type Array struct {
	elements []any // BSON type, rawObject, or rawArray
}

func (o *Object) MarshalBinary() ([]byte, error) {
	res := make([]byte, 0, 5)

	for _, f := range o.fields {
		var t byte
		var b []byte
		var err error

		switch value := f.value.(type) {
		case rawObject:
			t = 0x03
			b = value
		case rawArray:
			t = 0x04
			b = value
		case float64:
			t = 0x01
			b, err = marshalDouble(value)
			if err != nil {
				return nil, fmt.Errorf("marshal double: %w", err)
			}

			// ...
		}

		res = append(res, t)
		res = append(res, f.name...)
		res = append(res, 0x00)
		res = append(res, b...)
	}

	res = append(res, 0x00)

	return res, nil
}

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

// Unmarshal parses the BSON data and stores the result
// in the value pointed to by v.
func Unmarshal(data []byte, v any) error {
	d := NewDecodeBytes(data)
	if err := d.Decode(v); err != nil {
		return err
	}
	return nil
}

// A is a BSON array.
//
// Example:
//
//	bson.A{"hello", "world", 3.14159, bson.D{{"foo", 12345}}}
type A []any

func (a A) AsD() D {
	d := make(D, len(a))
	for i, v := range a {
		d[i] = e{K: strconv.Itoa(i), V: v}
	}
	return d
}

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

func (d D) Len() int           { return len(d) }
func (d D) Less(i, j int) bool { return d[i].K < d[j].K }
func (d D) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }

func (d D) AsM() M {
	m := make(M, len(d))
	for _, pair := range d {
		m[pair.K] = pair.V
	}
	return m
}

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

// RawArray represents a raw array which will be encoded or decoded as is.
type RawArray []byte

// RawObject represents a raw object which will be encoded or decoded as is.
type RawObject []byte
