package bsonproto

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"reflect"
	"testing"
)

func TestScalars(t *testing.T) {
	for i, tc := range []struct {
		v any
		b []byte
	}{{
		v: 42.13,
		b: []byte{0x71, 0x3d, 0xa, 0xd7, 0xa3, 0x10, 0x45, 0x40},
	}, {
		v: math.Inf(-1),
		b: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xf0, 0xff},
	}, {
		v: "foo",
		b: []byte{0x08, 0x00, 0x00, 0x00, 0x66, 0x6f, 0x6f, 0x00},
	}, {
		v: int32(123456789),
		b: []byte{0x15, 0xcd, 0x5b, 0x07},
	}, {
		v: int64(1234567890123456789),
		b: []byte{0x15, 0x81, 0xe9, 0x7d, 0xf4, 0x10, 0x22, 0x11},
	}} {
		t.Run(fmt.Sprintf("%[1]d_%[2]T(%[2]v)", i, tc.v), func(t *testing.T) {
			s := SizeAny(tc.v)
			if s != len(tc.b) {
				t.Fatalf("Size(%[1]T(%[1]v)) = %[2]d, want %[3]d", tc.v, s, len(tc.b))
			}

			actualB := make([]byte, s)
			EncodeAny(actualB, tc.v)
			if !bytes.Equal(actualB, tc.b) {
				t.Errorf("Encode(%[1]T(%[1]v)) = %#[2]v, want %#[3]v", tc.v, actualB, tc.b)
			}

			actualV := reflect.New(reflect.TypeOf(tc.v)).Interface() // actualV := new(T)
			err := DecodeAny(actualB, actualV)
			if err != nil {
				t.Fatalf("Decode(%v): %s", actualB, err)
			}

			actualV = reflect.ValueOf(actualV).Elem().Interface() // *actualV
			if !reflect.DeepEqual(actualV, tc.v) {
				t.Errorf("Decode(%v) = %v, want %v", actualB, actualV, tc.v)
			}
		})
	}
}

func TestFloat64(t *testing.T) {
	t.Run("NegativeZero", func(t *testing.T) {
		v := math.Copysign(0, -1)
		b := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80}

		actualB := make([]byte, 8)
		EncodeFloat64(actualB, v)
		if !bytes.Equal(actualB, b) {
			t.Errorf("Encode(%[1]T(%[1]v)) = %#[2]v, want %#[3]v", v, actualB, b)
		}

		actualV, err := DecodeFloat64(actualB)
		if err != nil {
			t.Fatalf("Decode(%v): %s", actualB, err)
		}
		if !reflect.DeepEqual(actualV, v) || !math.Signbit(actualV) {
			t.Errorf("Decode(%v) = %v, want %v", actualB, actualV, v)
		}
	})

	t.Run("NaNBits", func(t *testing.T) {
		for name, tc := range map[string]struct {
			v float64
			b []byte
		}{
			"qNaN": {
				v: math.Float64frombits(0b_11111111_11111000_00000000_00000000_00000000_00000000_00000000_00000000),
				b: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xf8, 0xff},
			},
			"qNaNpayload": {
				v: math.Float64frombits(0b_11111111_11110000_00000000_00000000_00000000_01110111_00010011_01000010),
				b: []byte{0x42, 0x13, 0x77, 0x00, 0x00, 0x00, 0xf0, 0xff},
			},
			"sNaN": {
				v: math.Float64frombits(0b_01111111_11111000_00000000_00000000_00000000_00000000_00000000_00000000),
				b: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xf8, 0x7f},
			},
			"sNaNpayload": {
				v: math.Float64frombits(0b_01111111_11110000_00000000_00000000_00000000_01110111_00010011_01000010),
				b: []byte{0x42, 0x13, 0x77, 0x00, 0x00, 0x00, 0xf0, 0x7f},
			},
		} {
			t.Run(name, func(t *testing.T) {
				actualB := make([]byte, 8)
				EncodeFloat64(actualB, tc.v)
				if !bytes.Equal(actualB, tc.b) {
					t.Errorf("Encode(%[1]T(%[1]v)) = %#[2]v, want %#[3]v", tc.v, actualB, tc.b)
				}

				actualV, err := DecodeFloat64(actualB)
				if err != nil {
					t.Fatalf("Decode(%v): %s", actualB, err)
				}
				if !math.IsNaN(actualV) {
					t.Errorf("Decode(%v) = %v, want NaN", actualB, actualV)
				}
			})
		}
	})
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
			err := DecodeAny(tc.b, v)
			if !errors.Is(err, tc.err) {
				t.Errorf("Decode(%v): %v, want %v", tc.b, err, tc.err)
			}
		})
	}
}
