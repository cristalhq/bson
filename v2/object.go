package bson

import "fmt"

// RawObject represents a BSON object a.k.a document in the binary encoded form.
type RawObject []byte

// field represents a single Object field in the (partially) decoded form.
type field struct {
	value any
	name  string
}

// Object represents a BSON object a.k.a document in the (partially) decoded form.
//
// It may contain duplicate field names.
type Object struct {
	fields []field
}

func NewObject(pairs ...any) (*Object, error) {
	l := len(pairs)
	if l%2 != 0 {
		return nil, fmt.Errorf("invalid number of arguments: %d", l)
	}

	res := &Object{
		fields: make([]field, l/2),
	}

	for i := 0; i < l; i += 2 {
		name, ok := pairs[i].(string)
		if !ok {
			return nil, fmt.Errorf("invalid field name type: %T", pairs[i])
		}

		value := pairs[i+1]
		if !checkType(value) {
			return nil, fmt.Errorf("invalid field value type: %T", value)
		}

		res.fields[i/2] = field{
			name:  name,
			value: value,
		}
	}

	return res, nil
}

func (o *Object) All(yield func(name string, value any) bool) bool {
	for _, f := range o.fields {
		if !yield(f.name, f.value) {
			return false
		}
	}

	return true
}
