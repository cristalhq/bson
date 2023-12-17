package bson

import (
	"bytes"
	"encoding/hex"
	"io"
	"testing"
)

func TestEncodeDecode(t *testing.T) {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)

	for name, tc := range map[string]struct {
		v any
		b []byte
	}{
		"SimpleObject": {
			v: must(NewObject(
				"foo", "bar",
			)),
			b: []byte{
				0x16, 0x00, 0x00, 0x00, // 22 bytes
				0x02,                   // string field
				0x04, 0x00, 0x00, 0x00, // 4 bytes of "foo"
				0x66, 0x6f, 0x6f, 0x00, // "foo"
				0x04, 0x00, 0x00, 0x00, // 4 bytes of "bar"
				0x62, 0x61, 0x72, 0x00, // "bar"
				0x00,
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			buf.Reset()

			if err := enc.Encode(tc.v); err != nil {
				t.Fatal(err)
			}

			actualB, err := io.ReadAll(&buf)
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(actualB, tc.b) {
				t.Fatalf("expected:\n%#v\n%s\ngot:\n%#v\n%s", tc.b, hex.Dump(tc.b), actualB, hex.Dump(actualB))
			}
		})
	}
}
