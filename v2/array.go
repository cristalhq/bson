package bson

// RawArray represents a BSON array in the binary encoded form.
type RawArray []byte

// Array represents a BSON array in the (partially) decoded form.
type Array struct {
	elements []any
}
