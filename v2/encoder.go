package bson

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"time"
)

// Encoder writes BSON values to an output stream.
type Encoder struct {
	w io.Writer
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		w: w,
	}
}

// encode encodes the given value to BSON.
func (enc *Encoder) Encode(v any) error {
	switch v := v.(type) {
	case *Object:
		return enc.encodeObject(v)
	case RawObject:
		return enc.encodeRawObject(v)
	case *Array:
		return enc.encodeArray(v)
	case RawArray:
		return enc.encodeRawArray(v)
	case float64:
		return enc.encodeFloat64(v)
	case string:
		return enc.encodeString(v)
	case Binary:
		return enc.encodeBinary(v)
	case ObjectID:
		return enc.encodeObjectID(v)
	case bool:
		return enc.encodeBool(v)
	case time.Time:
		return enc.encodeTime(v)
	case NullType:
		return enc.encodeNullType(v)
	case Regex:
		return enc.encodeRegex(v)
	case int32:
		return enc.encodeInt32(v)
	case Timestamp:
		return enc.encodeTimestamp(v)
	case int64:
		return enc.encodeInt64(v)
	default:
		panic(fmt.Sprintf("unsupported type %T", v))
	}
}

func (enc *Encoder) encodeObject(v *Object) error {
	buf := bytes.NewBuffer(make([]byte, 0, 256))
	oEnc := NewEncoder(buf)

	var err error
	for _, f := range v.fields {
		switch v := f.value.(type) {
		case *Object:
			if _, err = oEnc.w.Write([]byte{TagObject}); err != nil {
				return err
			}
			if err = oEnc.encodeString(f.name); err != nil {
				return err
			}
			if err = oEnc.encodeObject(v); err != nil {
				return err
			}

		case RawObject:
			if _, err = oEnc.w.Write([]byte{TagObject}); err != nil {
				return err
			}
			if err = oEnc.encodeString(f.name); err != nil {
				return err
			}
			if err = oEnc.encodeRawObject(v); err != nil {
				return err
			}

		case *Array:
			if _, err = oEnc.w.Write([]byte{TagArray}); err != nil {
				return err
			}
			if err = oEnc.encodeString(f.name); err != nil {
				return err
			}
			if err = oEnc.encodeArray(v); err != nil {
				return err
			}

		case RawArray:
			if _, err = oEnc.w.Write([]byte{TagArray}); err != nil {
				return err
			}
			if err = oEnc.encodeString(f.name); err != nil {
				return err
			}
			if err = oEnc.encodeRawArray(v); err != nil {
				return err
			}

		case float64:
			if _, err = oEnc.w.Write([]byte{TagFloat64}); err != nil {
				return err
			}
			if err = oEnc.encodeString(f.name); err != nil {
				return err
			}
			if err = oEnc.encodeFloat64(v); err != nil {
				return err
			}

		case string:
			if _, err = oEnc.w.Write([]byte{TagString}); err != nil {
				return err
			}
			if err = oEnc.encodeString(f.name); err != nil {
				return err
			}
			if err = oEnc.encodeString(v); err != nil {
				return err
			}

		case Binary:
			if _, err = oEnc.w.Write([]byte{TagBinary}); err != nil {
				return err
			}
			if err = oEnc.encodeString(f.name); err != nil {
				return err
			}
			if err = oEnc.encodeBinary(v); err != nil {
				return err
			}

		case ObjectID:
			if _, err = oEnc.w.Write([]byte{TagObjectID}); err != nil {
				return err
			}
			if err = oEnc.encodeString(f.name); err != nil {
				return err
			}
			if err = oEnc.encodeObjectID(v); err != nil {
				return err
			}

		case bool:
			if _, err = oEnc.w.Write([]byte{TagBool}); err != nil {
				return err
			}
			if err = oEnc.encodeString(f.name); err != nil {
				return err
			}
			if err = oEnc.encodeBool(v); err != nil {
				return err
			}

		case time.Time:
			if _, err = oEnc.w.Write([]byte{TagTime}); err != nil {
				return err
			}
			if err = oEnc.encodeString(f.name); err != nil {
				return err
			}
			if err = oEnc.encodeTime(v); err != nil {
				return err
			}

		case NullType:
			if _, err = oEnc.w.Write([]byte{TagNullType}); err != nil {
				return err
			}
			if err = oEnc.encodeString(f.name); err != nil {
				return err
			}
			if err = oEnc.encodeNullType(v); err != nil {
				return err
			}

		case Regex:
			if _, err = oEnc.w.Write([]byte{TagRegex}); err != nil {
				return err
			}
			if err = oEnc.encodeString(f.name); err != nil {
				return err
			}
			if err = oEnc.encodeRegex(v); err != nil {
				return err
			}

		case int32:
			if _, err = oEnc.w.Write([]byte{TagInt32}); err != nil {
				return err
			}
			if err = oEnc.encodeString(f.name); err != nil {
				return err
			}
			if err = oEnc.encodeInt32(v); err != nil {
				return err
			}

		case Timestamp:
			if _, err = oEnc.w.Write([]byte{TagTimestamp}); err != nil {
				return err
			}
			if err = oEnc.encodeString(f.name); err != nil {
				return err
			}
			if err = oEnc.encodeTimestamp(v); err != nil {
				return err
			}

		case int64:
			if _, err = oEnc.w.Write([]byte{TagInt64}); err != nil {
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

	l := buf.Len() + 5
	if l > MaxObjectSize {
		return fmt.Errorf("object size %d exceeds maximum allowed %d", l, MaxObjectSize)
	}

	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], uint32(l))
	if _, err = enc.w.Write(b[:]); err != nil {
		return err
	}

	if _, err = enc.w.Write(buf.Bytes()); err != nil {
		return err
	}

	if _, err = enc.w.Write([]byte{0x00}); err != nil {
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
	panic("TODO")
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
