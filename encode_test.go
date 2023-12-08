package bson

import (
	"bytes"
	"testing"
)

func TestEncode(t *testing.T) {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)

	var err error
	var arr A

	arr = A{
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
	}
	err = enc.Encode(arr)
	mustOk(t, err)
	wantBytes(t, buf.Bytes(), "70000000083000011031007b0000001032007b0000001033007b0000001034007b0000001235007b000000000000001036007b0000001037007b0000001038007b0000001039007b000000123130007b00000000000000013131000000000000c05e40013132000000000000c05e4000")
	buf.Reset()
}

func TestEncodeA(t *testing.T) {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)

	var err error
	var arr A

	arr = A{}
	err = enc.Encode(arr)
	mustOk(t, err)
	wantBytes(t, buf.Bytes(), "0500000000")
	buf.Reset()

	arr = A{int32(30)}
	err = enc.Encode(arr)
	mustOk(t, err)
	wantBytes(t, buf.Bytes(), "0c0000001030001e00000000")
	buf.Reset()

	arr = A{"a", int32(10), "c"}
	err = enc.Encode(arr)
	mustOk(t, err)
	wantBytes(t, buf.Bytes(), "1e0000000230000200000061001031000a00000002320002000000630000")
	buf.Reset()

	arr = A{"a", int32(10), "c", true, "b", int64(10203040)}
	err = enc.Encode(arr)
	mustOk(t, err)
	wantBytes(t, buf.Bytes(), "360000000230000200000061001031000a00000002320002000000630008330001023400020000006200123500a0af9b000000000000")
	buf.Reset()
}

func TestEncodeD(t *testing.T) {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)

	var err error
	var doc D

	doc = D{}
	err = enc.Encode(doc)
	mustOk(t, err)
	wantBytes(t, buf.Bytes(), "0500000000")
	buf.Reset()

	doc = D{{"a", int32(10)}}
	err = enc.Encode(doc)
	mustOk(t, err)
	wantBytes(t, buf.Bytes(), "0c0000001061000a00000000")
	buf.Reset()

	doc = D{{"a", int32(10)}, {"c", true}}
	err = enc.Encode(doc)
	mustOk(t, err)
	wantBytes(t, buf.Bytes(), "100000001061000a0000000863000100")
	buf.Reset()

	doc = D{{"a", int32(10)}, {"c", true}, {"b", int64(10203040)}}
	err = enc.Encode(doc)
	mustOk(t, err)
	wantBytes(t, buf.Bytes(), "1b0000001061000a00000008630001126200a0af9b000000000000")
	buf.Reset()
}

func TestEncodeM(t *testing.T) {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)

	var err error
	var maa M

	maa = M{"a": int32(10)}
	err = enc.Encode(maa)
	mustOk(t, err)
	wantBytes(t, buf.Bytes(), "0c0000001061000a00000000")
	buf.Reset()

	maa = M{"a": int32(10), "c": true}
	err = enc.Encode(maa)
	mustOk(t, err)
	wantBytes(t, buf.Bytes(), "100000001061000a0000000863000100")
	buf.Reset()

	maa = M{"a": int32(10), "c": true, "b": int64(10203040)}
	err = enc.Encode(maa)
	mustOk(t, err)
	wantBytes(t, buf.Bytes(), "1b0000001061000a000000126200a0af9b00000000000863000100")
	buf.Reset()
}

func TestEncodeReflectMap(t *testing.T) {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)

	var err error
	var m map[string]any

	m = map[string]any{}
	err = enc.Encode(m)
	mustOk(t, err)
	wantBytes(t, buf.Bytes(), "0500000000")
	buf.Reset()

	m = map[string]any{"abc": int32(123)}
	err = enc.Encode(m)
	mustOk(t, err)
	wantBytes(t, buf.Bytes(), "0e00000010616263007b00000000")
	buf.Reset()

	m = map[string]any{"hello": "world", "foo": int32(123)}
	err = enc.Encode(m)
	mustOk(t, err)
	wantBytes(t, buf.Bytes(), "1f00000010666f6f007b0000000268656c6c6f0006000000776f726c640000")
	buf.Reset()
}

func TestEncodeReflectStruct(t *testing.T) {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)

	var err error
	type foo struct {
		A string `bson:"a-field"`
		B int32  `bson:"-"`
		C bool   `bson:"c-bool-field,omitempty"`
		D map[string]int32
	}
	var x foo

	err = enc.Encode(x)
	mustOk(t, err)
	wantBytes(t, buf.Bytes(), "1b00000002612d6669656c64000100000000036400050000000000")
	buf.Reset()

	x = foo{
		A: "123",
		B: 456,
		C: true,
		D: map[string]int32{"a": int32(42)},
	}
	err = enc.Encode(x)
	mustOk(t, err)
	wantBytes(t, buf.Bytes(), "3400000002612d6669656c6400040000003132330008632d626f6f6c2d6669656c6400010364000c0000001061002a0000000000")
	buf.Reset()
}

func TestEncodeReflectSlice(t *testing.T) {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)

	var err error
	var s []string

	s = []string{}
	err = enc.Encode(s)
	mustOk(t, err)
	wantBytes(t, buf.Bytes(), "0500000000")
	buf.Reset()

	s = []string{"abc"}
	err = enc.Encode(s)
	mustOk(t, err)
	wantBytes(t, buf.Bytes(), "10000000023000040000006162630000")
	buf.Reset()

	s = []string{"hello", "world"}
	err = enc.Encode(s)
	mustOk(t, err)
	wantBytes(t, buf.Bytes(), "1f0000000230000600000068656c6c6f0002310006000000776f726c640000")
	buf.Reset()
}

func TestEncodeReflectArray(t *testing.T) {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)

	var err error

	err = enc.Encode([0]string{})
	mustOk(t, err)
	wantBytes(t, buf.Bytes(), "0500000000")
	buf.Reset()

	err = enc.Encode([1]string{"abc"})
	mustOk(t, err)
	wantBytes(t, buf.Bytes(), "10000000023000040000006162630000")
	buf.Reset()

	err = enc.Encode([...]string{"hello", "world"})
	mustOk(t, err)
	wantBytes(t, buf.Bytes(), "1f0000000230000600000068656c6c6f0002310006000000776f726c640000")
	buf.Reset()
}
