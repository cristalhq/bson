package bsonproto

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestScalars(t *testing.T) {
	for i, tc := range []struct {
		v any
		b []byte
	}{{
		v: "foo",
		b: []byte{0x08, 0x00, 0x00, 0x00, 0x66, 0x6f, 0x6f, 0x00},
	}, {
		v: int32(123456789),
		b: []byte{0x15, 0xcd, 0x5b, 0x07},
	}, {
		v: int64(1234567890123456789),
		b: []byte{0x15, 0x81, 0xe9, 0x7d, 0xf4, 0x10, 0x22, 0x11},
	}} {
		t.Run(fmt.Sprintf("%d_%[1]T(%[1]v)", i, tc.v), func(t *testing.T) {
			s := SizeAny(tc.v)
			if s != len(tc.b) {
				t.Fatalf("Size(%[1]T(%[1]v)) = %[2]d, want %[3]d", tc.v, s, len(tc.b))
			}

			b := make([]byte, s)
			EncodeAny(b, tc.v)
			if !bytes.Equal(b, tc.b) {
				t.Errorf("Encode(%[1]T(%[1]v)) = %#[2]v, want %#[3]v", tc.v, b, tc.b)
			}

			v := reflect.New(reflect.TypeOf(tc.v)).Interface() // v := new(T)
			err := DecodeAny(v, b)
			if err != nil {
				t.Fatalf("Decode(%v): %s", b, err)
			}

			v = reflect.ValueOf(v).Elem().Interface() // *v
			if !reflect.DeepEqual(v, tc.v) {
				t.Errorf("Decode(%v) = %v, want %v", b, v, tc.v)
			}
		})
	}
}

func TestScalarsDecodeErrors(t *testing.T) {
	for i, tc := range []struct {
		b   []byte
		v   any
		err error
	}{{
		b:   []byte{0x42},
		v:   string(""),
		err: ErrDecodeShortInput,
	}, {
		b:   []byte{0x42},
		v:   int32(0),
		err: ErrDecodeShortInput,
	}, {
		b:   []byte{0x42},
		v:   int64(0),
		err: ErrDecodeShortInput,
	}} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			v := reflect.New(reflect.TypeOf(tc.v)).Interface() // v := new(T)
			err := DecodeAny(v, tc.b)
			if !errors.Is(err, tc.err) {
				t.Errorf("Decode(%v): %v, want %v", tc.b, err, tc.err)
			}
		})
	}
}
