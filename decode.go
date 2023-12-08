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
	iter, err := newReader(data)
	if err != nil {
		return err
	}

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
			v.SetString(trimlast(element))

		case TypeDocument:
			m := make(map[string]any)
			vv := reflect.ValueOf(m)
			if err := decodeMap(element, vv); err != nil {
				return err
			}
			v.Set(vv)

		case TypeArray:
			s := make([]any, 0)
			if err := decodeSlice(element, &s); err != nil {
				return err
			}
			vv := reflect.ValueOf(s)
			v.Set(vv)

		case TypeObjectID:
			var oid ObjectID
			copy(oid[:], element)
			v.Set(reflect.ValueOf(oid))

		case TypeBool:
			v.SetBool(element[0] == 1)

		case TypeInt32:
			element := int64(element[0]) |
				int64(element[1])<<8 |
				int64(element[2])<<16 |
				int64(element[3])<<24
			v.SetInt(element)

		case TypeTimestamp:
			ts := Timestamp(element[0]) |
				Timestamp(element[1])<<8 |
				Timestamp(element[2])<<16 |
				Timestamp(element[3])<<24 |
				Timestamp(element[4])<<32 |
				Timestamp(element[5])<<40 |
				Timestamp(element[6])<<48 |
				Timestamp(element[7])<<56
			v.Set(reflect.ValueOf(ts))

		case TypeInt64:
			element := int64(element[0]) |
				int64(element[1])<<8 |
				int64(element[2])<<16 |
				int64(element[3])<<24 |
				int64(element[4])<<32 |
				int64(element[5])<<40 | int64(element[6])<<48 | int64(element[7])<<56
			v.SetInt(element)

		case TypeBinary,
			TypeUndefined,
			TypeDateTime,
			TypeNull,
			TypeRegex,
			TypeDBPointer,
			TypeCodeWithScope,
			TypeSymbol,
			TypeJavaScriptScope,
			TypeDecimal,
			TypeMinKey,
			TypeMaxKey:
			return fmt.Errorf("unsupported type %x", typ)
		default:
			return fmt.Errorf("unknown element type %x", typ)

		}
	}
	return iter.Err()
}

func decodeMap(data []byte, v reflect.Value) error {
	iter, err := newReader(data)
	if err != nil {
		return err
	}

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
			vv := reflect.ValueOf(trimlast(element))
			v.SetMapIndex(key, vv)

		case TypeDocument:
			m := make(map[string]any)
			vv := reflect.ValueOf(m)
			if err := decodeMap(element, vv); err != nil {
				return err
			}
			v.SetMapIndex(key, vv)

		case TypeArray:
			s := make([]any, 0)
			if err := decodeSlice(element, &s); err != nil {
				return err
			}
			vv := reflect.ValueOf(s)
			v.SetMapIndex(key, vv)

		case TypeObjectID:
			var oid ObjectID
			copy(oid[:], element)
			v.SetMapIndex(key, reflect.ValueOf(oid))

		case TypeBool:
			v.SetMapIndex(key, reflect.ValueOf(element[0] == 1))

		case TypeInt32:
			element := int32(element[0]) |
				int32(element[1])<<8 |
				int32(element[2])<<16 |
				int32(element[3])<<24
			v.SetMapIndex(key, reflect.ValueOf(element))

		case TypeTimestamp:
			ts := Timestamp(element[0]) |
				Timestamp(element[1])<<8 |
				Timestamp(element[2])<<16 |
				Timestamp(element[3])<<24 |
				Timestamp(element[4])<<32 |
				Timestamp(element[5])<<40 |
				Timestamp(element[6])<<48 |
				Timestamp(element[7])<<56
			v.SetMapIndex(key, reflect.ValueOf(ts))

		case TypeInt64:
			element := int64(element[0]) |
				int64(element[1])<<8 |
				int64(element[2])<<16 |
				int64(element[3])<<24 |
				int64(element[4])<<32 |
				int64(element[5])<<40 |
				int64(element[6])<<48 |
				int64(element[7])<<56
			v.SetMapIndex(key, reflect.ValueOf(element))

		case TypeBinary,
			TypeUndefined,
			TypeDateTime,
			TypeRegex,
			TypeDBPointer,
			TypeCodeWithScope,
			TypeSymbol,
			TypeJavaScriptScope,
			TypeDecimal,
			TypeMinKey,
			TypeMaxKey:
			return fmt.Errorf("unsupported type %x", typ)

		default:
			return fmt.Errorf("unknown element type %x", typ)
		}
	}
	return iter.Err()
}

func decodeSlice(data []byte, v *[]any) error {
	iter, err := newReader(data)
	if err != nil {
		return err
	}

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
			*v = append(*v, trimlast(element))

		case TypeDocument:
			m := make(map[string]any)
			vv := reflect.ValueOf(m)
			if err := decodeMap(element, vv); err != nil {
				return err
			}
			*v = append(*v, m)

		case TypeArray:
			s := make([]any, 0)
			if err := decodeSlice(element, &s); err != nil {
				return err
			}
			*v = append(*v, s)

		case TypeObjectID:
			var oid ObjectID
			copy(oid[:], element)
			*v = append(*v, oid)

		case TypeBool:
			b := element[0] == 1
			*v = append(*v, b)

		// case TypeDateTime:
		// 	dt := Datetime(element[0]) | Datetime(element[1])<<8 | Datetime(element[2])<<16 | Datetime(element[3])<<24 | Datetime(element[4])<<32 | Datetime(element[5])<<40 | Datetime(element[6])<<48 | Datetime(element[7])<<56
		// 	*v = append(*v, dt)

		case TypeNull:
			*v = append(*v, nil)

		case TypeInt32:
			element := int32(element[0]) | int32(element[1])<<8 | int32(element[2])<<16 | int32(element[3])<<24
			*v = append(*v, element)

		case TypeTimestamp:
			ts := Timestamp(element[0]) | Timestamp(element[1])<<8 | Timestamp(element[2])<<16 | Timestamp(element[3])<<24 | Timestamp(element[4])<<32 | Timestamp(element[5])<<40 | Timestamp(element[6])<<48 | Timestamp(element[7])<<56
			*v = append(*v, ts)

		case TypeInt64:
			element := int64(element[0]) | int64(element[1])<<8 | int64(element[2])<<16 | int64(element[3])<<24 | int64(element[4])<<32 | int64(element[5])<<40 | int64(element[6])<<48 | int64(element[7])<<56
			*v = append(*v, element)

		case TypeBinary,
			TypeUndefined,
			TypeDateTime,
			TypeRegex,
			TypeDBPointer,
			TypeCodeWithScope,
			TypeSymbol,
			TypeJavaScriptScope,
			TypeDecimal,
			TypeMinKey,
			TypeMaxKey:
			return fmt.Errorf("unsupported type %x", typ)

		default:
			return fmt.Errorf("unknown element type %x", typ)
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

func newReader(b []byte) (reader, error) {
	if len(b) < 5 {
		return reader{}, errors.New("not enough data")
	}
	return reader{
		data: b[4 : len(b)-1], // TODO(cristaloleg): this can easily panic.
	}, nil
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
		if len(rest) < 5 {
			return r.setErr(errors.New("corrupt BSON reading string len"))
		}
		elen, _ := readInt32(rest)
		if len(rest) < elen {
			return r.setErr(errors.New("corrupt BSON reading string"))
		}
		element = rest[4 : 4+elen]
		if len(rest) < 4+elen {
			return r.setErr(errors.New("corrupt BSON reading string"))
		}
		rest = rest[4+elen:]

	case TypeDocument, TypeArray:
		if len(rest) < 5 {
			return r.setErr(errors.New("corrupt BSON reading string len"))
		}
		elen, _ := readInt32(rest)
		if len(rest) < elen {
			return r.setErr(fmt.Errorf("corrupt document: want %x bytes, have %x", elen, len(rest)))
		}
		element, rest = rest[:elen], rest[elen:]

	case TypeObjectID:
		if len(rest) < 12 {
			return r.setErr(errors.New("corrupt BSON reading object id"))
		}
		element, rest = rest[:12], rest[12:]

	case TypeBool:
		if len(rest) < 1 {
			return r.setErr(errors.New("corrupt BSON reading boolean"))
		}
		element, rest = rest[:1], rest[1:]

	case TypeDateTime:
		if len(rest) < 8 {
			return r.setErr(errors.New("corrupt BSON reading utc datetime"))
		}
		element, rest = rest[:8], rest[8:]

	case TypeNull:
		element, rest = rest[:0], rest[0:]

	// TODO(cristaloleg): fuzzer fails on 2nd IndexByte.
	// case TypeRegex:
	// 	if len(rest) < 2 {
	// 		return r.setErr(errors.New("corrupt BSON reading regex"))
	// 	}
	// 	i := bytes.IndexByte(rest, 0)
	// 	if i < 0 {
	// 		return r.setErr(errors.New("corrupt BSON regex 1"))
	// 	}
	// 	i++
	// 	j := bytes.IndexByte(rest[i+1:], 0)
	// 	if j < 0 {
	// 		return r.setErr(errors.New("corrupt BSON regex 2"))
	// 	}
	// 	j++
	// 	element, rest = rest[:i+j+1], rest[i+j+1:]

	case TypeInt32:
		if len(rest) < 4 {
			return r.setErr(errors.New("corrupt BSON reading int32"))
		}
		element, rest = rest[:4], rest[4:]

	case TypeTimestamp:
		fallthrough

	case TypeInt64:
		if len(rest) < 8 {
			return r.setErr(errors.New("corrupt BSON reading int64"))
		}
		element, rest = rest[:8], rest[8:]

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

// readInt32 the 4 bytes in little endian and return tail.
// Panics if less than 4 bytes is passed.
func readInt32(buf []byte) (int, []byte) {
	v := int(buf[0]) |
		int(buf[1])<<8 |
		int(buf[2])<<16 |
		int(buf[3])<<24
	return v, buf[4:]
}

// readCstring returns CString including \0.
func readCstring(buf []byte) ([]byte, []byte, error) {
	i := bytes.IndexByte(buf, 0)
	if i == -1 {
		return nil, nil, errors.New("cstring missing \\0")
	}
	i++
	return buf[:i], buf[i:], nil
}

func trimlast(s []byte) string {
	if len(s) == 0 {
		return ""
	}
	return string(s[:len(s)-1])
}
