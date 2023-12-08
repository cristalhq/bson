package bson

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
)

// Decoder reads and decodes BSON values from an input stream.
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
	if len(d.data) < 4 {
		return errors.New("not enough data") // TODO(cristaloleg): static error?
	}

	rv := reflect.ValueOf(v)
	switch {
	case rv.Kind() != reflect.Ptr:
		return errors.New("unmarshal non-pointer: " + rv.Type().String())
	case rv.IsNil():
		return errors.New("unmarshal nil: " + rv.Type().String())
	}

	switch rv := rv.Elem(); rv.Kind() {
	case reflect.Struct:
		return decodeStruct(d.data, rv)
	case reflect.Map:
		return decodeMap(d.data, rv)
	default:
		return errors.New("unmarshal unsupported: " + rv.Type().String())
	}
}

func decodeStruct(data []byte, v reflect.Value) error {
	iter := newReader(data)

	for iter.Next() {
		typ, name, element := iter.Peek()

		v := v.FieldByName(trimlast(name))
		if !v.IsValid() {
			continue
		}

		switch typ {
		case TypeDouble:
			val := uint64(element[0]) |
				uint64(element[1])<<8 |
				uint64(element[2])<<16 |
				uint64(element[3])<<24 |
				uint64(element[4])<<32 |
				uint64(element[5])<<40 |
				uint64(element[6])<<48 |
				uint64(element[7])<<56
			v.SetFloat(math.Float64frombits(val))

		case TypeString:
		case TypeDocument:
		case TypeArray:
		case TypeBinary:
		case TypeUndefined:
		case TypeObjectID:
		case TypeBool:
		case TypeDateTime:
		case TypeNull:
		case TypeRegex:
		case TypeDBPointer:
		case TypeCodeWithScope:
		case TypeSymbol:
		case TypeJavaScriptScope:
		case TypeInt32:
		case TypeTimestamp:
		case TypeInt64:
		case TypeDecimal:
		case TypeMinKey:
		case TypeMaxKey:
		}
	}
	return iter.Err()
}

func decodeMap(data []byte, v reflect.Value) error {
	iter := newReader(data)

	for iter.Next() {
		typ, name, element := iter.Peek()

		key := reflect.ValueOf(trimlast(name))

		switch typ {
		case TypeDouble:
			val := uint64(element[0]) |
				uint64(element[1])<<8 |
				uint64(element[2])<<16 |
				uint64(element[3])<<24 |
				uint64(element[4])<<32 |
				uint64(element[5])<<40 |
				uint64(element[6])<<48 |
				uint64(element[7])<<56
			v.SetMapIndex(key, reflect.ValueOf(math.Float64frombits(val)))

		case TypeString:
		case TypeDocument:
		case TypeArray:
		case TypeBinary:
		case TypeUndefined:
		case TypeObjectID:
		case TypeBool:
		case TypeDateTime:
		case TypeNull:
		case TypeRegex:
		case TypeDBPointer:
		case TypeCodeWithScope:
		case TypeSymbol:
		case TypeJavaScriptScope:
		case TypeInt32:
		case TypeTimestamp:
		case TypeInt64:
		case TypeDecimal:
		case TypeMinKey:
		case TypeMaxKey:
		}
	}
	return iter.Err()
}

func decodeSlice(data []byte, v *[]any) error {
	iter := newReader(data)

	for iter.Next() {
		typ, _, element := iter.Peek()

		switch typ {
		case TypeDouble:
			bits := uint64(element[0]) |
				uint64(element[1])<<8 |
				uint64(element[2])<<16 |
				uint64(element[3])<<24 |
				uint64(element[4])<<32 |
				uint64(element[5])<<40 |
				uint64(element[6])<<48 |
				uint64(element[7])<<56
			*v = append(*v, bits)

		case TypeString:
		case TypeDocument:
		case TypeArray:
		case TypeBinary:
		case TypeUndefined:
		case TypeObjectID:
		case TypeBool:
		case TypeDateTime:
		case TypeNull:
		case TypeRegex:
		case TypeDBPointer:
		case TypeCodeWithScope:
		case TypeSymbol:
		case TypeJavaScriptScope:
		case TypeInt32:
		case TypeTimestamp:
		case TypeInt64:
		case TypeDecimal:
		case TypeMinKey:
		case TypeMaxKey:
		}
	}
	return iter.Err()
}

type reader struct {
	data    []byte // data to process.
	name    []byte // name of the current element.
	element []byte // current element.
	err     error  // err during reading.
}

func newReader(b []byte) reader {
	return reader{
		data: b[4 : len(b)-1], // TODO(cristaloleg): this can easily panic.
	}
}

func (r *reader) Err() error {
	return r.err
}

func (r *reader) Peek() (Type, []byte, []byte) {
	return Type(r.name[0]), r.name[1:], r.element
}

func (r *reader) Next() bool {
	if r.err != nil || len(r.data) == 0 {
		return false
	}
	if len(r.data) == 1 {
		return r.setErr(errors.New("corrupt BSON, only 1 byte remains"))
	}

	i := bytes.IndexByte(r.data[1:], 0)
	if i < 0 {
		return r.setErr(errors.New("corrupt BSON name"))
	}
	i += 2

	name, rest := r.data[:i], r.data[i:]
	var element []byte

	switch typ := name[0]; Type(typ) {
	case TypeDouble:
		if len(rest) < 8 {
			return r.setErr(errors.New("corrupt BSON reading double"))
		}
		element, rest = rest[:8], rest[8:]

	case TypeString:
	case TypeDocument:
	case TypeArray:
	case TypeBinary:
	case TypeUndefined:
	case TypeObjectID:
	case TypeBool:
	case TypeDateTime:
	case TypeNull:
	case TypeRegex:
	case TypeDBPointer:
	case TypeCodeWithScope:
	case TypeSymbol:
	case TypeJavaScriptScope:
	case TypeInt32:
	case TypeTimestamp:
	case TypeInt64:
	case TypeDecimal:
	case TypeMinKey:
	case TypeMaxKey:
	default:
		return r.setErr(fmt.Errorf("unknown element type %x", typ))
	}

	r.data, r.name, r.element = rest, name, element
	return true
}

func (r *reader) setErr(err error) bool {
	if r.err == nil {
		r.err = err
	}
	return false
}

func readInt32(buf []byte) (int, []byte) {
	v := int(buf[0]) |
		int(buf[1])<<8 |
		int(buf[2])<<16 |
		int(buf[3])<<24
	return v, buf[4:]
}

func readCstring(buf []byte) ([]byte, []byte, error) {
	i := bytes.IndexByte(buf, 0)
	if i == -1 {
		return nil, nil, errors.New("cstring missing \\0")
	}
	i++
	return buf[:i], buf[i:], nil
}

func trimlast(s []byte) string { return string(s[:len(s)-1]) }
