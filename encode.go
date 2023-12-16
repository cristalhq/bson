package bson

import (
	"fmt"
	"io"
	"math"
	"reflect"
	"sort"
	"strconv"
)

// Encoder writes BSON values to an output stream.
type Encoder struct {
	w   io.Writer
	buf []byte
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		w:   w,
		buf: make([]byte, 0, 512),
	}
}

// Encode writes the BSON encoding of v to the stream.
func (enc *Encoder) Encode(v any) error {
	enc.buf = enc.buf[:0]
	if err := enc.marshal(v); err != nil {
		return fmt.Errorf("encode failed: %w", err)
	}
	_, err := enc.w.Write(enc.buf)
	return err
}

func (enc *Encoder) marshal(v any) error {
	var err error
	switch v := v.(type) {
	case D:
		_, err = enc.writeD(v)
	case M:
		_, err = enc.writeD(v.AsD())
	case map[string]any:
		_, err = enc.writeD(M(v).AsD())
	case A:
		_, err = enc.writeA(v)
	case []any:
		_, err = enc.writeA(v)
	case RawObject:
		enc.buf = append(enc.buf, v...)
	case RawArray:
		enc.buf = append(enc.buf, v...)

	default:
		switch rv := reflect.ValueOf(v); rv.Kind() {
		case reflect.Struct:
			_, err = enc.writeStruct(rv)
		case reflect.Map:
			_, err = enc.writeMap(rv)
		case reflect.Array, reflect.Slice:
			_, err = enc.writeSlice(rv)
		default:
			return fmt.Errorf("type %T is not supported yet", v)
		}
	}
	return err
}

func (enc *Encoder) writeD(d D) (int, error) {
	start := len(enc.buf)
	enc.buf = append(enc.buf, 0, 0, 0, 0)
	count := 4 + 1 // sizeof(int) + sizeof(\0)

	for i := 0; i < len(d); i++ {
		n, err := enc.writeAny(d[i].K, d[i].V)
		if err != nil {
			return 0, err
		}
		count += n
	}

	enc.buf = append(enc.buf, 0)
	enc.buf[start] = byte(count)
	enc.buf[start+1] = byte(count >> 8)
	enc.buf[start+2] = byte(count >> 16)
	enc.buf[start+3] = byte(count >> 24)
	return count, nil
}

func (enc *Encoder) writeA(a A) (int, error) {
	start := len(enc.buf)
	enc.buf = append(enc.buf, 0, 0, 0, 0)
	count := 4 + 1 // sizeof(int) + sizeof(\0)

	for i := range a {
		n, err := enc.writeAny(strconv.Itoa(i), a[i])
		if err != nil {
			return 0, err
		}
		count += n
	}

	enc.buf = append(enc.buf, 0)
	enc.buf[start] = byte(count)
	enc.buf[start+1] = byte(count >> 8)
	enc.buf[start+2] = byte(count >> 16)
	enc.buf[start+3] = byte(count >> 24)
	return count, nil
}

func (enc *Encoder) writeMap(v reflect.Value) (int, error) {
	if m, ok := v.Interface().(map[string]any); ok {
		return enc.writeD(M(m).AsD())
	}

	d := make(D, v.Len())

	for i, key := range v.MapKeys() {
		d[i] = e{
			K: key.String(),
			V: v.MapIndex(key).Interface(),
		}
	}
	sort.Sort(d)
	return enc.writeD(d)
}

func (enc *Encoder) writeStruct(v reflect.Value) (int, error) {
	start := len(enc.buf)
	enc.buf = append(enc.buf, 0, 0, 0, 0)
	count := 4 + 1 // sizeof(int) + sizeof(\0)

	d := getStruct(v).asDoc(v)

	for i := 0; i < len(d); i++ {
		n, err := enc.writeAny(d[i].Key, d[i].Val)
		if err != nil {
			return 0, err
		}
		count += n
	}

	enc.buf = append(enc.buf, 0)
	enc.buf[start] = byte(count)
	enc.buf[start+1] = byte(count >> 8)
	enc.buf[start+2] = byte(count >> 16)
	enc.buf[start+3] = byte(count >> 24)
	return count, nil
}

func (enc *Encoder) writeSlice(v reflect.Value) (int, error) {
	if a, ok := v.Interface().([]any); ok {
		return enc.writeA(a)
	}

	start := len(enc.buf)
	enc.buf = append(enc.buf, 0, 0, 0, 0)
	count := 4 + 1 // sizeof(int) + sizeof(\0)

	n := v.Len()
	for i := 0; i < n; i++ {
		n, err := enc.writeAny(strconv.Itoa(i), v.Index(i).Interface())
		if err != nil {
			return 0, err
		}
		count += n
	}

	enc.buf = append(enc.buf, 0)
	enc.buf[start] = byte(count)
	enc.buf[start+1] = byte(count >> 8)
	enc.buf[start+2] = byte(count >> 16)
	enc.buf[start+3] = byte(count >> 24)
	return count, nil
}

func (enc *Encoder) writeAny(ename string, v any) (int, error) {
	// TODO(cristaloleg): support bson.Marshaler.
	switch v := v.(type) {
	case string:
		return enc.writeString(ename, v), nil
	case bool:
		return enc.writeBool(ename, v), nil

	// TODO(cristaloleg): force int64 for int and uint to prevent bit truncation.
	case int:
		return enc.writeInt32(ename, int32(v)), nil
	case uint:
		return enc.writeInt32(ename, int32(v)), nil

	case int8:
		return enc.writeInt32(ename, int32(v)), nil
	case uint8:
		return enc.writeInt32(ename, int32(v)), nil

	case int16:
		return enc.writeInt32(ename, int32(v)), nil
	case uint16:
		return enc.writeInt32(ename, int32(v)), nil

	case int32:
		return enc.writeInt32(ename, int32(v)), nil
	case uint32:
		return enc.writeInt32(ename, int32(v)), nil

	case int64:
		return enc.writeInt64(ename, v), nil
	case uint64:
		return enc.writeInt64(ename, int64(v)), nil

	case float64:
		return enc.writeFloat64(ename, v), nil
	case float32:
		return enc.writeFloat64(ename, float64(v)), nil

	default:
		return enc.writeValue(ename, reflect.ValueOf(v))
	}
}

func (enc *Encoder) writeValue(ename string, v reflect.Value) (int, error) {
	if v.Kind() == reflect.Interface {
		return enc.writeValue(ename, v.Elem())
	}

	var count int
	switch v.Kind() {

	case reflect.Map:
		count += enc.writeElem(TypeDocument, ename)
		n, err := enc.writeMap(v)
		if err != nil {
			return 0, err
		}
		count += n

	case reflect.Struct:
		count += enc.writeElem(TypeDocument, ename)
		n, err := enc.writeStruct(v)
		if err != nil {
			return 0, err
		}
		count += n

	case reflect.Array, reflect.Slice:
		count += enc.writeElem(TypeArray, ename)
		n, err := enc.writeSlice(v)
		if err != nil {
			return 0, err
		}
		count += n

	default:
		return 0, fmt.Errorf("type %v is not supported", v.Kind())
	}
	return count, nil
}

func (enc *Encoder) writeElem(typ Type, key string) int {
	enc.buf = append(enc.buf, byte(typ))
	enc.buf = append(enc.buf, key...)
	enc.buf = append(enc.buf, 0)
	return 1 + len(key) + 1
}

func (enc *Encoder) writeString(key, s string) int {
	size := len(s) + 1
	enc.buf = append(enc.buf, byte(TypeString))
	enc.buf = append(enc.buf, key...)
	enc.buf = append(enc.buf,
		0,
		byte(size),
		byte(size>>8),
		byte(size>>16),
		byte(size>>24),
	)
	enc.buf = append(enc.buf, s...)
	enc.buf = append(enc.buf, 0)
	return (1 + len(key) + 1) + (4 + size)
}

func (enc *Encoder) writeBool(key string, b bool) int {
	var v byte
	if b {
		v = 1
	}
	enc.buf = append(enc.buf, byte(TypeBool))
	enc.buf = append(enc.buf, key...)
	enc.buf = append(enc.buf, 0, v)
	return 1 + len(key) + 1 + 1
}

func (enc *Encoder) writeInt32(key string, v int32) int {
	enc.buf = append(enc.buf, byte(TypeInt32))
	enc.buf = append(enc.buf, key...)
	enc.buf = append(enc.buf,
		0,
		byte(v),
		byte(v>>8),
		byte(v>>16),
		byte(v>>24),
	)
	return (1 + len(key) + 1) + 4
}

func (enc *Encoder) writeInt64(key string, v int64) int {
	enc.buf = append(enc.buf, byte(TypeInt64))
	enc.buf = append(enc.buf, key...)
	enc.buf = append(enc.buf,
		0,
		byte(v),
		byte(v>>8),
		byte(v>>16),
		byte(v>>24),
		byte(v>>32),
		byte(v>>40),
		byte(v>>48),
		byte(v>>56),
	)
	return (1 + len(key) + 1) + 8
}

func (enc *Encoder) writeFloat64(key string, f float64) int {
	v := int64(math.Float64bits(f))
	enc.buf = append(enc.buf, byte(TypeDouble))
	enc.buf = append(enc.buf, key...)
	enc.buf = append(enc.buf,
		0,
		byte(v),
		byte(v>>8),
		byte(v>>16),
		byte(v>>24),
		byte(v>>32),
		byte(v>>40),
		byte(v>>48),
		byte(v>>56),
	)
	return (1 + len(key) + 1) + 8
}
