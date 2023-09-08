package bson

import (
	"errors"
	"io"
	"reflect"
)

// A Decoder reads and decodes BSON values from an input stream.
type Decoder struct {
	r    io.Reader
	data []byte
}

func NewDecodeBytes(buf []byte) *Decoder {
	d := &Decoder{
		data: buf,
	}
	return d
}

func (d *Decoder) Decode(v any) error {
	return d.decode(d.data, v)
}

func (d *Decoder) decode(data []byte, v any) error {
	rv := reflect.ValueOf(v)
	switch {
	case rv.Kind() != reflect.Ptr:
		return errors.New("unmarshal non-pointer: " + rv.Type().String())
	case rv.IsNil():
		return errors.New("unmarshal nil: " + rv.Type().String())
	}

	switch rv := rv.Elem(); rv.Kind() {
	// TODO(cristaloleg): support these types.
	// case reflect.Struct:
	// 	return decodeStruct(data, rv)
	// case reflect.Map:
	// 	return decodeMap(data, rv)
	default:
		return errors.New("unmarshal unsupported: " + rv.Type().String())
	}
}
