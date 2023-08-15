package bson

import (
	"bytes"
	"math"
	"testing"
)

func TestEncode(t *testing.T) {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)

	err := enc.Encode(int32(123456789))
	mustOk(t, err)
	wantBytes(t, buf.Bytes(), "15cd5b07")
	buf.Reset()

	err = enc.Encode(int64(123456789123456789))
	mustOk(t, err)
	wantBytes(t, buf.Bytes(), "155fd0ac4b9bb601")
	buf.Reset()

	err = enc.Encode(true)
	mustOk(t, err)
	wantBytes(t, buf.Bytes(), "01")
	buf.Reset()
}

func TestEncodeBytes(t *testing.T) {
	var err error
	var buf bytes.Buffer
	enc := NewEncoder(&buf)

	err = enc.Encode([]byte("foo"))
	mustOk(t, err)
	wantBytes(t, buf.Bytes(), "0400000080666f6f")
	buf.Reset()

	err = enc.Encode([]byte{0x00})
	mustOk(t, err)
	wantBytes(t, buf.Bytes(), "020000008000")
	buf.Reset()
}

func TestEncodeString(t *testing.T) {
	var err error
	var buf bytes.Buffer
	enc := NewEncoder(&buf)

	err = enc.Encode("foo")
	mustOk(t, err)
	wantBytes(t, buf.Bytes(), "04000000666f6f00")
	buf.Reset()

	err = enc.Encode("")
	mustOk(t, err)
	wantBytes(t, buf.Bytes(), "0100000000")
	buf.Reset()
}

func TestEncodeNumbers(t *testing.T) {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)

	err := enc.Encode(float64(3.14159))
	mustOk(t, err)
	wantBytes(t, buf.Bytes(), "6e861bf0f9210940")
	buf.Reset()

	err = enc.Encode(float64(0))
	mustOk(t, err)
	wantBytes(t, buf.Bytes(), "0000000000000000")
	buf.Reset()

	err = enc.Encode(math.NaN())
	mustOk(t, err)
	wantBytes(t, buf.Bytes(), "010000000000f87f")
	buf.Reset()

	err = enc.Encode(math.Inf(+1))
	mustOk(t, err)
	wantBytes(t, buf.Bytes(), "000000000000f07f")
	buf.Reset()

	err = enc.Encode(math.Inf(-1))
	mustOk(t, err)
	wantBytes(t, buf.Bytes(), "000000000000f0ff")
	buf.Reset()

	err = enc.Encode(42.13)
	mustOk(t, err)
	wantBytes(t, buf.Bytes(), "713d0ad7a3104540")
	buf.Reset()
}
