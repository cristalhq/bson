package bson

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"time"
)

// Encoder writes BSON values to an output stream.
type Encoder struct {
	w *bufio.Writer
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		w: bufio.NewWriter(w),
	}
}

// Encode encodes the given value to BSON.
func (enc *Encoder) Encode(v any) error {
	var err error
	switch v := v.(type) {
	case *Object:
		err = enc.encodeObject(v)
	case RawObject:
		err = enc.encodeRawObject(v)
	case *Array:
		err = enc.encodeArray(v)
	case RawArray:
		err = enc.encodeRawArray(v)
	case float64:
		err = enc.encodeFloat64(v)
	case string:
		err = enc.encodeString(v)
	case Binary:
		err = enc.encodeBinary(v)
	case ObjectID:
		err = enc.encodeObjectID(v)
	case bool:
		err = enc.encodeBool(v)
	case time.Time:
		err = enc.encodeTime(v)
	case NullType:
		err = enc.encodeNullType(v)
	case Regex:
		err = enc.encodeRegex(v)
	case int32:
		err = enc.encodeInt32(v)
	case Timestamp:
		err = enc.encodeTimestamp(v)
	case int64:
		err = enc.encodeInt64(v)
	default:
		return fmt.Errorf("unsupported type %T", v)
	}

	if err != nil {
		return err
	}

	if err = enc.w.Flush(); err != nil {
		return err
	}

	return nil
}

func (enc *Encoder) encodeObject(v *Object) error {
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
			if err = oEnc.encodeObject(v); err != nil {
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

		case Binary:
			if err = oEnc.w.WriteByte(TagBinary); err != nil {
				return err
			}
			if err = oEnc.encodeString(f.name); err != nil {
				return err
			}
			if err = oEnc.encodeBinary(v); err != nil {
				return err
			}

		case ObjectID:
			if err = oEnc.w.WriteByte(TagObjectID); err != nil {
				return err
			}
			if err = oEnc.encodeString(f.name); err != nil {
				return err
			}
			if err = oEnc.encodeObjectID(v); err != nil {
				return err
			}

		case bool:
			if err = oEnc.w.WriteByte(TagBool); err != nil {
				return err
			}
			if err = oEnc.encodeString(f.name); err != nil {
				return err
			}
			if err = oEnc.encodeBool(v); err != nil {
				return err
			}

		case time.Time:
			if err = oEnc.w.WriteByte(TagTime); err != nil {
				return err
			}
			if err = oEnc.encodeString(f.name); err != nil {
				return err
			}
			if err = oEnc.encodeTime(v); err != nil {
				return err
			}

		case NullType:
			if err = oEnc.w.WriteByte(TagNullType); err != nil {
				return err
			}
			if err = oEnc.encodeString(f.name); err != nil {
				return err
			}
			if err = oEnc.encodeNullType(v); err != nil {
				return err
			}

		case Regex:
			if err = oEnc.w.WriteByte(TagRegex); err != nil {
				return err
			}
			if err = oEnc.encodeString(f.name); err != nil {
				return err
			}
			if err = oEnc.encodeRegex(v); err != nil {
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

		case Timestamp:
			if err = oEnc.w.WriteByte(TagTimestamp); err != nil {
				return err
			}
			if err = oEnc.encodeString(f.name); err != nil {
				return err
			}
			if err = oEnc.encodeTimestamp(v); err != nil {
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

func (enc *Encoder) encodeBinary(v Binary) error {
	panic("TODO")
}

func (enc *Encoder) encodeObjectID(v ObjectID) error {
	panic("TODO")
}

func (enc *Encoder) encodeBool(v bool) error {
	panic("TODO")
}

func (enc *Encoder) encodeTime(v time.Time) error {
	panic("TODO")
}

func (enc *Encoder) encodeNullType(v NullType) error {
	panic("TODO")
}

func (enc *Encoder) encodeRegex(v Regex) error {
	panic("TODO")
}

func (enc *Encoder) encodeInt32(v int32) error {
	panic("TODO")
}

func (enc *Encoder) encodeTimestamp(v Timestamp) error {
	panic("TODO")
}

func (enc *Encoder) encodeInt64(v int64) error {
	panic("TODO")
}
