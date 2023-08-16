package bson

import (
	"bytes"
	"fmt"
	"io"
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
		enc.marshalArray(v)
	case D:
		enc.marshalDoc(v)
	case M:
		enc.marshalDoc(v.AsD())
	default:
		return fmt.Errorf("type %T is not supported yet", v)
	}
	return nil
}

func (enc *Encoder) marshalArray(arr A) error {
	doc := make(D, len(arr))
	for i := range arr {
		doc[i] = e{
			K: strconv.Itoa(i),
			V: arr[i],
		}
	}
	return enc.marshalDoc(doc)
}

func (enc *Encoder) marshalDoc(doc D) error {
	// TODO(cristaloleg): prealloc or smarter way.
	var elist bytes.Buffer

	for i := range doc {
		pair := doc[i]
		key := doc[i].K

		switch v := pair.V.(type) {
		case string:
			enc.writeKey(&elist, TypeString, key)
			b := putUint32(uint32(len(v) + 1))
			elist.Write(b[:])
			elist.WriteString(v)
			elist.WriteByte(0)

		case int32:
			enc.writeKey(&elist, TypeInt32, key)
			b := putUint32(uint32(v))
			elist.Write(b[:])

		case int64:
			enc.writeKey(&elist, TypeInt64, key)
			b := putUint64(uint64(v))
			elist.Write(b[:])

		case bool:
			enc.writeKey(&elist, TypeBool, key)
			elist.WriteByte(putBool(v))
		}
	}

	size := 4 + elist.Len() + 1 // header + len + null.
	b := putUint32(uint32(size))
	enc.buf.Write(b[:])

	io.Copy(enc.buf, &elist)
	enc.buf.WriteByte(0)
	return nil
}

func (enc *Encoder) writeKey(buf *bytes.Buffer, t Type, s string) {
	buf.WriteByte(byte(t))
	buf.WriteString(s)
	buf.WriteByte(0)
}
