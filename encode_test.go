package bson

import (
	"bytes"
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
