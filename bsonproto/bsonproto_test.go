package bsonproto

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"reflect"
	"testing"
	"time"
)

func TestScalars(t *testing.T) {
	for i, tc := range []struct {
		v any
		b []byte
	}{{
		v: float64(42.13),
		b: []byte{0x71, 0x3d, 0xa, 0xd7, 0xa3, 0x10, 0x45, 0x40},
	}, {
		v: math.Inf(-1),
		b: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xf0, 0xff},
	}, {
		v: "foo",
		b: []byte{0x04, 0x00, 0x00, 0x00, 0x66, 0x6f, 0x6f, 0x00},
	}, {
		v: "f",
		b: []byte{0x02, 0x00, 0x00, 0x00, 0x66, 0x00},
	}, {
		v: "",
		b: []byte{0x01, 0x00, 0x00, 0x00, 0x00},
	}, {
		v: Binary{B: []byte("foo"), Subtype: BinaryUser},
		b: []byte{0x03, 0x00, 0x00, 0x00, 0x80, 0x66, 0x6f, 0x6f},
	}, {
		v: Binary{B: []byte("f"), Subtype: BinaryUser},
		b: []byte{0x01, 0x00, 0x00, 0x00, 0x80, 0x66},
	}, {
		v: Binary{Subtype: BinaryUser},
		b: []byte{0x00, 0x00, 0x00, 0x00, 0x80},
	}, {
		v: ObjectID{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c},
		b: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c},
	}, {
		v: false,
		b: []byte{0x00},
	}, {
		v: true,
		b: []byte{0x01},
	}, {
		v: time.Date(2023, 12, 26, 13, 22, 42, 123000000, time.UTC),
		b: []byte{0x4b, 0xb1, 0x4a, 0xa6, 0x8c, 0x01, 0x00, 0x00},
	}, {
		v: time.Unix(0, 0).UTC(),
		b: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
	}, {
		v: time.Time{},
		b: []byte{0x00, 0x28, 0xd3, 0xed, 0x7c, 0xc7, 0xff, 0xff},
	}, {
		v: Regex{Pattern: "foo", Options: "m"},
		b: []byte{0x66, 0x6f, 0x6f, 0x00, 0x6d, 0x00},
	}, {
		v: Regex{Pattern: "", Options: ""},
		b: []byte{0x00, 0x00},
	}, {
		v: int32(123456789),
		b: []byte{0x15, 0xcd, 0x5b, 0x07},
	}, {
		v: Timestamp(1234567890123456789),
		b: []byte{0x15, 0x81, 0xe9, 0x7d, 0xf4, 0x10, 0x22, 0x11},
	}, {
		v: int64(1234567890123456789),
		b: []byte{0x15, 0x81, 0xe9, 0x7d, 0xf4, 0x10, 0x22, 0x11},
	}} {
		t.Run(fmt.Sprintf("%[1]d_%[2]T(%[2]v)", i, tc.v), func(t *testing.T) {
			s := SizeAny(tc.v)
			if s != len(tc.b) {
				t.Fatalf("SizeAny(%[1]T(%[1]v)) = %[2]d, expected %[3]d", tc.v, s, len(tc.b))
			}

			actualB := make([]byte, s)
			EncodeAny(actualB, tc.v)
			if !bytes.Equal(actualB, tc.b) {
				t.Errorf("EncodeAny(%[1]T(%[1]v))\n actual   %#[2]v\n expected %#[3]v", tc.v, actualB, tc.b)
			}

			actualV := reflect.New(reflect.TypeOf(tc.v)).Interface() // actualV := new(T)
			err := DecodeAny(actualB, actualV)
			if err != nil {
				t.Fatalf("DecodeAny(%v): %s", actualB, err)
			}

			actualV = reflect.ValueOf(actualV).Elem().Interface() // *actualV
			if !reflect.DeepEqual(actualV, tc.v) {
				t.Errorf("DecodeAny(%v)\n actual   %v\n expected %v", actualB, actualV, tc.v)
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
			t.Errorf("EncodeFloat64(%[1]T(%[1]v)) = %#[2]v, expected %#[3]v", v, actualB, b)
		}

		actualV, err := DecodeFloat64(actualB)
		if err != nil {
			t.Fatalf("DecodeFloat64(%v): %s", actualB, err)
		}
		if !reflect.DeepEqual(actualV, v) || !math.Signbit(actualV) {
			t.Errorf("DecodeFloat64(%v) = %v, expected %v", actualB, actualV, v)
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
					t.Errorf("EncodeFloat64(%[1]T(%[1]v)) = %#[2]v, expected %#[3]v", tc.v, actualB, tc.b)
				}

				actualV, err := DecodeFloat64(actualB)
				if err != nil {
					t.Fatalf("DecodeFloat64(%v): %s", actualB, err)
				}
				if !math.IsNaN(actualV) {
					t.Errorf("DecodeFloat64(%v) = %v, expected NaN", actualB, actualV)
				}
			})
		}
	})
}

func TestCString(t *testing.T) {
	for _, tc := range []struct {
		v string
		b []byte
	}{{
		v: "foo",
		b: []byte{0x66, 0x6f, 0x6f, 0x00},
	}, {
		v: "f",
		b: []byte{0x66, 0x00},
	}, {
		v: "",
		b: []byte{0x00},
	}} {
		t.Run(tc.v, func(t *testing.T) {
			s := SizeCString(tc.v)
			if s != len(tc.b) {
				t.Fatalf("SizeCString(%[1]T(%[1]v)) = %[2]d, expected %[3]d", tc.v, s, len(tc.b))
			}

			actualB := make([]byte, s)
			EncodeCString(actualB, tc.v)
			if !bytes.Equal(actualB, tc.b) {
				t.Errorf("EncodeCString(%[1]q)\n actual   %#[2]v\n expected %#[3]v", tc.v, actualB, tc.b)
			}

			actualV, err := DecodeCString(actualB)
			if err != nil {
				t.Fatalf("DecodeCString(%v): %s", actualB, err)
			}
			if actualV != tc.v {
				t.Errorf("DecodeCString(%v) = %q, expected %q", actualB, actualV, tc.v)
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
		v:   float64(0),
		err: ErrDecodeShortInput,
	}, {
		b:   []byte{0x42},
		v:   string(""),
		err: ErrDecodeShortInput,
	}, {
		b:   []byte{0x42, 0x42, 0x42, 0x42, 0x42},
		v:   string(""),
		err: ErrDecodeShortInput,
	}, {
		b:   []byte{0x00, 0x00, 0x00, 0x00, 0x42},
		v:   string(""),
		err: ErrDecodeInvalidInput,
	}, {
		b:   []byte{0x01, 0x00, 0x00, 0x00, 0x42},
		v:   string(""),
		err: ErrDecodeInvalidInput,
	}, {
		b:   []byte{0x42},
		v:   Binary{},
		err: ErrDecodeShortInput,
	}, {
		b:   []byte{0x01, 0x00, 0x00, 0x00, 0x80},
		v:   Binary{},
		err: ErrDecodeShortInput,
	}, {
		b:   []byte{},
		v:   ObjectID{},
		err: ErrDecodeShortInput,
	}, {
		b:   []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b},
		v:   ObjectID{},
		err: ErrDecodeShortInput,
	}, {
		b:   []byte{},
		v:   false,
		err: ErrDecodeShortInput,
	}, {
		b:   []byte{0x42},
		v:   false,
		err: ErrDecodeInvalidInput,
	}, {
		b:   []byte{0x42},
		v:   time.Time{},
		err: ErrDecodeShortInput,
	}, {
		b:   []byte{0x00},
		v:   Regex{},
		err: ErrDecodeShortInput,
	}, {
		b:   []byte{0x00, 0x42},
		v:   Regex{},
		err: ErrDecodeShortInput,
	}, {
		b:   []byte{0x42},
		v:   int32(0),
		err: ErrDecodeShortInput,
	}, {
		b:   []byte{0x42},
		v:   Timestamp(0),
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
				t.Errorf("DecodeAny(%v): %v, expected %v", tc.b, err, tc.err)
			}
		})
	}
}
