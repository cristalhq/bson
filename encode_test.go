package bson

import (
	"bytes"
	"testing"
)

func TestEncodeA(t *testing.T) {
	testCases := []struct {
		arr  A
		want string
	}{
		{
			arr:  A{},
			want: "0500000000",
		},
		{
			arr:  A{int32(30)},
			want: "0c0000001030001e00000000",
		},
		{
			arr:  A{"a", int32(10), "c"},
			want: "1e0000000230000200000061001031000a00000002320002000000630000",
		},
		{
			arr:  A{"a", int32(10), "c", true, "b", int64(10203040)},
			want: "360000000230000200000061001031000a00000002320002000000630008330001023400020000006200123500a0af9b000000000000",
		},
		{
			arr: A{
				bool(true),
				int(123),
				int8(123),
				int16(123),
				int32(123),
				int64(123),
				uint(123),
				uint8(123),
				uint16(123),
				uint32(123),
				uint64(123),
				float32(123),
				float64(123),
			},
			want: "70000000083000011031007b0000001032007b0000001033007b0000001034007b0000001235007b000000000000001036007b0000001037007b0000001038007b0000001039007b000000123130007b00000000000000013131000000000000c05e40013132000000000000c05e4000",
		},
	}

	for _, tc := range testCases {
		var buf bytes.Buffer

		err := NewEncoder(&buf).Encode(tc.arr)
		mustOk(t, err)
		wantBytes(t, buf.Bytes(), tc.want)
	}
}

func TestEncodeD(t *testing.T) {
	testCases := []struct {
		doc  D
		want string
	}{
		{
			doc:  D{},
			want: "0500000000",
		},

		{
			doc:  D{{"a", int32(10)}},
			want: "0c0000001061000a00000000",
		},

		{
			doc:  D{{"a", int32(10)}, {"c", true}},
			want: "100000001061000a0000000863000100",
		},

		{
			doc:  D{{"a", int32(10)}, {"c", true}, {"b", int64(10203040)}},
			want: "1b0000001061000a00000008630001126200a0af9b000000000000",
		},
	}

	for _, tc := range testCases {
		var buf bytes.Buffer

		err := NewEncoder(&buf).Encode(tc.doc)
		mustOk(t, err)
		wantBytes(t, buf.Bytes(), tc.want)
	}
}

func TestEncodeM(t *testing.T) {
	testCases := []struct {
		doc  M
		want string
	}{
		{
			doc:  M{},
			want: "0500000000",
		},
		{
			doc:  M{"a": int32(10)},
			want: "0c0000001061000a00000000",
		},
		{
			doc:  M{"a": int32(10), "c": true},
			want: "100000001061000a0000000863000100",
		},
		{
			doc:  M{"a": int32(10), "c": true, "b": int64(10203040)},
			want: "1b0000001061000a000000126200a0af9b00000000000863000100",
		},
	}

	for _, tc := range testCases {
		var buf bytes.Buffer

		err := NewEncoder(&buf).Encode(tc.doc)
		mustOk(t, err)
		wantBytes(t, buf.Bytes(), tc.want)
	}
}

func TestEncodeReflectMap(t *testing.T) {
	testCases := []struct {
		m    any
		want string
	}{
		{
			m:    map[string]string{},
			want: "0500000000",
		},
		{
			m:    map[string]int{"abc": int(123)},
			want: "0e00000010616263007b00000000",
		},
		{
			m:    map[string]any{"hello": "world", "foo": int32(123)},
			want: "1f00000010666f6f007b0000000268656c6c6f0006000000776f726c640000",
		},
	}

	for _, tc := range testCases {
		var buf bytes.Buffer

		err := NewEncoder(&buf).Encode(tc.m)
		mustOk(t, err)
		wantBytes(t, buf.Bytes(), tc.want)
	}
}

func TestEncodeReflectStruct(t *testing.T) {
	type foo struct {
		A string `bson:"a-field"`
		B int    `bson:"-"`
		C bool   `bson:"c-bool-field,omitempty"`
		D map[string]int32
	}

	testCases := []struct {
		s    any
		want string
	}{
		{
			s:    foo{},
			want: "1b00000002612d6669656c64000100000000036400050000000000",
		},
		{
			s: foo{
				A: "123",
				B: 456,
				C: true,
				D: map[string]int32{"a": int32(42)},
			},
			want: "3400000002612d6669656c6400040000003132330008632d626f6f6c2d6669656c6400010364000c0000001061002a0000000000",
		},
	}

	for _, tc := range testCases {
		var buf bytes.Buffer

		err := NewEncoder(&buf).Encode(tc.s)
		mustOk(t, err)
		wantBytes(t, buf.Bytes(), tc.want)
	}
}

func TestEncodeReflectSlice(t *testing.T) {
	testCases := []struct {
		slice any
		want  string
	}{
		{
			slice: []string{},
			want:  "0500000000",
		},

		{
			slice: []string{"abc"},
			want:  "10000000023000040000006162630000",
		},

		{
			slice: []string{"hello", "world"},
			want:  "1f0000000230000600000068656c6c6f0002310006000000776f726c640000",
		},
	}

	for _, tc := range testCases {
		var buf bytes.Buffer

		err := NewEncoder(&buf).Encode(tc.slice)
		mustOk(t, err)
		wantBytes(t, buf.Bytes(), tc.want)
	}
}

func TestEncodeReflectArray(t *testing.T) {
	testCases := []struct {
		arr  any
		want string
	}{
		{
			arr:  [0]string{},
			want: "0500000000",
		},
		{
			arr:  [1]string{"abc"},
			want: "10000000023000040000006162630000",
		},
		{
			arr:  [...]string{"hello", "world"},
			want: "1f0000000230000600000068656c6c6f0002310006000000776f726c640000",
		},
	}

	for _, tc := range testCases {
		var buf bytes.Buffer

		err := NewEncoder(&buf).Encode(tc.arr)
		mustOk(t, err)
		wantBytes(t, buf.Bytes(), tc.want)
	}
}
