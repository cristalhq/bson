package bson

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strconv"
)

// Encoder writes BSON values to an output stream.
type Encoder struct {
	w   io.Writer
	buf *bytes.Buffer
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		w:   w,
		buf: bytes.NewBuffer(make([]byte, 0, 512)),
	}
}

// Encode writes the BSON encoding of v to the stream.
func (enc *Encoder) Encode(v any) error {
	enc.buf.Reset()
	if err := enc.marshal(v); err != nil {
		return fmt.Errorf("encode failed: %w", err)
	}
	_, err := enc.w.Write(enc.buf.Bytes())
	return err
}

func (enc *Encoder) marshal(v any) error {
	switch v := v.(type) {
	case Marshaler:
		raw, err := v.MarshalBSON()
		if err != nil {
			return err
		}
		enc.buf.Write(raw)
	case A:
		enc.marshalArray(enc.buf, v)
	case D:
		enc.marshalDoc(enc.buf, v)
	case M:
		enc.marshalDoc(enc.buf, v.AsD())
	default:
		return enc.marshalReflect(enc.buf, v)
	}
	return nil
}

func (enc *Encoder) marshalArray(w io.Writer, arr A) error {
	doc := make(D, len(arr))
	for i := range arr {
		doc[i] = e{
			K: strconv.Itoa(i),
			V: arr[i],
		}
	}
	return enc.marshalDoc(w, doc)
}

func (enc *Encoder) marshalDoc(w io.Writer, doc D) error {
	// TODO(cristaloleg): prealloc or smarter way.
	elist := bytes.NewBuffer(make([]byte, 0, 128))

	for i := range doc {
		key := doc[i].K
		val := doc[i].V

		switch v := val.(type) {
		case string:
			enc.writeKey(elist, TypeString, key)
			b := putUint32(uint32(len(v) + 1))
			elist.Write(b[:])
			elist.WriteString(v)
			elist.WriteByte(0)

		case int32:
			enc.writeKey(elist, TypeInt32, key)
			b := putUint32(uint32(v))
			elist.Write(b[:])

		case int64:
			enc.writeKey(elist, TypeInt64, key)
			b := putUint64(uint64(v))
			elist.Write(b[:])

		case bool:
			enc.writeKey(elist, TypeBool, key)
			elist.WriteByte(putBool(v))

		default:
			var err error
			switch rv := reflect.ValueOf(val); rv.Kind() {
			case reflect.Map:
				enc.writeKey(elist, TypeDocument, key)
				err = enc.marshalMap(elist, rv)

			case reflect.Array, reflect.Slice:
				enc.writeKey(elist, TypeArray, key)
				err = enc.marshalSlice(elist, rv)

			default:
				return fmt.Errorf("type %T is not supported yet", v)
			}
			if err != nil {
				return err
			}
		}
	}

	size := 4 + elist.Len() + 1 // header + len + null.
	b := putUint32(uint32(size))
	w.Write(b[:])

	io.Copy(w, elist)
	w.Write([]byte{0})
	return nil
}

func (enc *Encoder) writeKey(buf *bytes.Buffer, t Type, s string) {
	buf.WriteByte(byte(t))
	buf.WriteString(s)
	buf.WriteByte(0)
}

func (enc *Encoder) marshalReflect(w io.Writer, v any) error {
	switch rv := reflect.ValueOf(v); rv.Kind() {
	// TODO(cristaloleg): add reflect.Struct
	case reflect.Map:
		return enc.marshalMap(w, rv)
	case reflect.Array, reflect.Slice:
		return enc.marshalSlice(w, rv)
	default:
		return fmt.Errorf("type %T is not supported yet", v)
	}
}

func (enc *Encoder) marshalMap(w io.Writer, v reflect.Value) error {
	doc := make(D, v.Len())
	for i, iter := 0, v.MapRange(); iter.Next(); i++ {
		doc[i] = e{
			K: iter.Key().String(),
			V: iter.Value().Interface(),
		}
	}
	return enc.marshalDoc(w, doc)
}

func (enc *Encoder) marshalSlice(w io.Writer, v reflect.Value) error {
	doc := make(D, v.Len())
	for i := 0; i < v.Len(); i++ {
		doc[i] = e{
			K: strconv.Itoa(i),
			V: v.Index(i).Interface(),
		}
	}
	return enc.marshalDoc(w, doc)
}
