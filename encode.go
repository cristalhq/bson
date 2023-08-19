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
	if doc, ok := v.(D); ok {
		_, err := enc.writeDoc(doc)
		return err
	}

	var err error
	switch rv := reflect.ValueOf(v); rv.Kind() {
	case reflect.Map:
		_, err = enc.writeMap(rv)
	case reflect.Struct:
		_, err = enc.writeStruct(rv)
	case reflect.Array, reflect.Slice:
		_, err = enc.writeSlice(rv)
	default:
		return fmt.Errorf("type %T is not supported yet", v)
	}
	return err
}

func (enc *Encoder) writeDoc(doc D) (int, error) {
	start := len(enc.buf)
	enc.buf = append(enc.buf, 0, 0, 0, 0)
	count := 4 + 1 // sizeof(int) + sizeof(\0)

	for i := 0; i < len(doc); i++ {
		n, err := enc.writeValue(doc[i].K, reflect.ValueOf(doc[i].V))
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

// TODO(cristaloleg): doc[i] value box-unbox can be omitted.
func (enc *Encoder) writeMap(v reflect.Value) (int, error) {
	start := len(enc.buf)
	enc.buf = append(enc.buf, 0, 0, 0, 0)
	count := 4 + 1 // sizeof(int) + sizeof(\0)

	doc := make(docRefl, v.Len())
	if m, ok := v.Interface().(map[string]any); ok {
		i := 0
		for k, v := range m {
			doc[i] = pairRefl{
				Key: k,
				Val: reflect.ValueOf(v),
			}
			i++
		}
	} else {
		keys := v.MapKeys()
		for i := range keys {
			key := keys[i]
			doc[i] = pairRefl{
				Key: key.String(),
				Val: v.MapIndex(key),
			}
		}
	}
	sortPairRefl(doc)

	for i := 0; i < len(doc); i++ {
		n, err := enc.writeValue(doc[i].Key, doc[i].Val)
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

func (enc *Encoder) writeStruct(v reflect.Value) (int, error) {
	start := len(enc.buf)
	enc.buf = append(enc.buf, 0, 0, 0, 0)
	count := 4 + 1 // sizeof(int) + sizeof(\0)

	doc := getStruct(v).asDoc(v)

	sort.Sort(doc)

	for i := 0; i < len(doc); i++ {
		n, err := enc.writeValue(doc[i].Key, doc[i].Val)
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
	start := len(enc.buf)
	enc.buf = append(enc.buf, 0, 0, 0, 0)
	count := 4 + 1 // sizeof(int) + sizeof(\0)

	if a, ok := v.Interface().([]any); ok {
		for i := range a {
			n, err := enc.writeValue(strconv.Itoa(i), reflect.ValueOf(a[i]))
			if err != nil {
				return 0, err
			}
			count += n
		}
	} else {
		n := v.Len()
		for i := 0; i < n; i++ {
			val := v.Index(i)

			n, err := enc.writeValue(strconv.Itoa(i), val)
			if err != nil {
				return 0, err
			}
			count += n
		}
	}

	enc.buf = append(enc.buf, 0)
	enc.buf[start] = byte(count)
	enc.buf[start+1] = byte(count >> 8)
	enc.buf[start+2] = byte(count >> 16)
	enc.buf[start+3] = byte(count >> 24)
	return count, nil
}

// TODO(cristaloleg): probably split into simple & compound types.
func (enc *Encoder) writeValue(ename string, v reflect.Value) (int, error) {
	if v.Kind() == reflect.Interface {
		return enc.writeValue(ename, v.Elem())
	}

	var count int
	switch v.Kind() {
	case reflect.String:
		count += enc.writeElem(TypeString, ename)
		count += enc.writeString(v.String())
	case reflect.Int32:
		count += enc.writeElem(TypeInt32, ename)
		count += enc.writeInt32(int32(v.Int()))
	case reflect.Int64:
		count += enc.writeElem(TypeInt64, ename)
		count += enc.writeInt64(int64(v.Int()))
	case reflect.Bool:
		count += enc.writeElem(TypeBool, ename)
		count += enc.writeBool(v.Bool())
	case reflect.Float64:
		count += enc.writeElem(TypeDouble, ename)
		count += enc.writeInt64(int64(math.Float64bits(v.Float())))

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

func (enc *Encoder) writeString(s string) int {
	size := len(s) + 1
	enc.writeInt32(int32(size))
	enc.buf = append(enc.buf, s...)
	enc.buf = append(enc.buf, 0)
	return 4 + size
}

func (enc *Encoder) writeInt32(v int32) int {
	enc.buf = append(enc.buf,
		byte(v),
		byte(v>>8),
		byte(v>>16),
		byte(v>>24),
	)
	return 4
}

func (enc *Encoder) writeInt64(v int64) int {
	enc.buf = append(enc.buf,
		byte(v),
		byte(v>>8),
		byte(v>>16),
		byte(v>>24),
		byte(v>>32),
		byte(v>>40),
		byte(v>>48),
		byte(v>>56),
	)
	return 8
}

func (enc *Encoder) writeBool(b bool) int {
	var v byte
	if b {
		v = 1
	}
	enc.buf = append(enc.buf, v)
	return 1
}
