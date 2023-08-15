package bson

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
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
func (e *Encoder) Encode(v any) error {
	e.buf.Reset()
	if err := e.marshal(v); err != nil {
		return fmt.Errorf("encode failed: %w", err)
	}
	_, err := e.w.Write(e.buf.Bytes())
	return err
}

func (e *Encoder) marshal(v any) error {
	switch v := v.(type) {
	case Marshaler:
		raw, err := v.MarshalBSON()
		if err != nil {
			return err
		}
		e.buf.Write(raw)

	case []byte:
		var b [4]byte
		binary.LittleEndian.PutUint32(b[:], uint32(len(v)+1))
		e.buf.Write(b[:])
		e.buf.WriteByte(0x80) // TODO(cristaloleg): better binary type?
		e.buf.Write(v)

	case string:
		var b [4]byte
		binary.LittleEndian.PutUint32(b[:], uint32(len(v)+1))
		e.buf.Write(b[:])
		e.buf.Write([]byte(v))
		e.buf.WriteByte(0)

	case int32:
		var b [4]byte
		binary.LittleEndian.PutUint32(b[:], uint32(v))
		e.buf.Write(b[:])

	case int64:
		var b [8]byte
		binary.LittleEndian.PutUint64(b[:], uint64(v))
		e.buf.Write(b[:])

	case float64:
		var b [8]byte
		binary.LittleEndian.PutUint64(b[:], math.Float64bits(float64(v)))
		e.buf.Write(b[:])

	case bool:
		var b [1]byte
		if v {
			b[0] = 1
		}
		e.buf.Write(b[:])

	default:
		return fmt.Errorf("type %T is not supported yet", v)
	}
	return nil
}
