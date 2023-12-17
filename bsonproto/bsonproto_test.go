package bsonproto

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

func TestScalars(t *testing.T) {
	for _, tc := range []struct {
		v any
		b []byte
	}{
		{
			v: int32(123456),
			b: []byte{0x40, 0xe2, 0x01, 0x00},
		},
	} {
		t.Run(fmt.Sprintf("%[1]T %[1]v", tc.v), func(t *testing.T) {
			s := SizeAny(tc.v)
			if s != len(tc.b) {
				t.Fatalf("Size(%[1]T %[1]v) = %[2]d, want %[3]d", tc.v, s, len(tc.b))
			}

			b := make([]byte, s)
			EncodeAny(b, tc.v)
			if !bytes.Equal(b, tc.b) {
				t.Fatalf("Encode(%[1]T %[1]v) = %#[2]v, want %#[3]v", tc.v, b, tc.b)
			}

			v := reflect.New(reflect.TypeOf(tc.v)).Interface() // v := new(T)
			err := DecodeAny(v, b)
			if err != nil {
				t.Fatalf("Decode(%v): %s", b, err)
			}

			v = reflect.ValueOf(v).Elem().Interface() // *v
			if !reflect.DeepEqual(v, tc.v) {
				t.Fatalf("Decode(%v) = %v, want %v", b, v, tc.v)
			}
		})
	}
}
