package bson

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strconv"

	"github.com/cristalhq/bson/bsonproto"
)

// Encoder writes BSON values to an output stream.
type Encoder struct {
	p *pool
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		p: new(pool),
	}
}

func encodeObjectSize(v *Object) int {
	size := 4 // length

	for _, f := range v.fields {
		size += 1 + len(f.name) + 1 // tag + name + terminating 0x00

		switch v := f.value.(type) {
		case *Object:
			size += encodeObjectSize(v)

		case RawObject:
			size += len(v)

		case *Array:
			size += encodeArraySize(v)

		case RawArray:
			size += len(v)

		default:
			size += bsonproto.SizeAny(v)
		}
	}

	return size + 1 // terminating 0x00
}

func encodeArraySize(v *Array) int {
	size := 4 // length

	for i, f := range v.elements {
		size += 1 + len(strconv.Itoa(i)) + 1 // tag + name + terminating 0x00

		switch v := f.(type) {
		case *Object:
			size += encodeObjectSize(v)

		case RawObject:
			size += len(v)

		case *Array:
			size += encodeArraySize(v)

		case RawArray:
			size += len(v)

		default:
			size += bsonproto.SizeAny(v)
		}
	}

	return size + 1 // terminating 0x00
}

func (enc *Encoder) EncodeObject(v *Object) error {
	buf := bytes.NewBuffer(make([]byte, 0, 256))
	oEnc := NewEncoder(buf)

	var err error
	for _, f := range v.fields {
		switch v := f.value.(type) {
		case *Object:
			if err = oEnc.w.WriteByte(TagObject); err != nil {
				return err
			}
			if err = oEnc.encodeString(f.name); err != nil {
				return err
			}
			if err = oEnc.EncodeObject(v); err != nil {
				return err
			}

		case RawObject:
			if err = oEnc.w.WriteByte(TagObject); err != nil {
				return err
			}
			if err = oEnc.encodeString(f.name); err != nil {
				return err
			}
			if err = oEnc.encodeRawObject(v); err != nil {
				return err
			}

		case *Array:
			if err = oEnc.w.WriteByte(TagArray); err != nil {
				return err
			}
			if err = oEnc.encodeString(f.name); err != nil {
				return err
			}
			if err = oEnc.encodeArray(v); err != nil {
				return err
			}

		case RawArray:
			if err = oEnc.w.WriteByte(TagArray); err != nil {
				return err
			}
			if err = oEnc.encodeString(f.name); err != nil {
				return err
			}
			if err = oEnc.encodeRawArray(v); err != nil {
				return err
			}

		case float64:
			if err = oEnc.w.WriteByte(TagFloat64); err != nil {
				return err
			}
			if err = oEnc.encodeString(f.name); err != nil {
				return err
			}
			if err = oEnc.encodeFloat64(v); err != nil {
				return err
			}

		case string:
			if err = oEnc.w.WriteByte(TagString); err != nil {
				return err
			}
			if err = oEnc.encodeString(f.name); err != nil {
				return err
			}
			if err = oEnc.encodeString(v); err != nil {
				return err
			}

		case int32:
			if err = oEnc.w.WriteByte(TagInt32); err != nil {
				return err
			}
			if err = oEnc.encodeString(f.name); err != nil {
				return err
			}
			if err = oEnc.encodeInt32(v); err != nil {
				return err
			}

		case int64:
			if err = oEnc.w.WriteByte(TagInt64); err != nil {
				return err
			}
			if err = oEnc.encodeString(f.name); err != nil {
				return err
			}
			if err = oEnc.encodeInt64(v); err != nil {
				return err
			}

		default:
			panic(fmt.Sprintf("unsupported type %T", v))
		}
	}

	if err = oEnc.w.Flush(); err != nil {
		return err
	}

	l := buf.Len() + 5

	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], uint32(l))
	if _, err = enc.w.Write(b[:]); err != nil {
		return err
	}

	if _, err = enc.w.Write(buf.Bytes()); err != nil {
		return err
	}

	if err = enc.w.WriteByte(0x00); err != nil {
		return err
	}

	return nil
}

func (enc *Encoder) encodeRawObject(v RawObject) error {
	_, err := enc.w.Write(v)
	return err
}

func (enc *Encoder) encodeArray(v *Array) error {
	panic("TODO")
}

func (enc *Encoder) encodeRawArray(v RawArray) error {
	_, err := enc.w.Write(v)
	return err
}

func (enc *Encoder) encodeFloat64(v float64) error {
	panic("TODO")
}

func (enc *Encoder) encodeString(v string) error {
	var err error

	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], uint32(len(v)+1))
	if _, err = enc.w.Write(b[:]); err != nil {
		return err
	}

	if _, err = io.WriteString(enc.w, v); err != nil {
		return err
	}

	if err = enc.w.WriteByte(0x00); err != nil {
		return err
	}

	return nil
}

func (enc *Encoder) encodeInt32(v int32) error {
	panic("TODO")
}

func (enc *Encoder) encodeInt64(v int64) error {
	panic("TODO")
}
