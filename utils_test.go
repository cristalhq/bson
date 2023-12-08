package bson

import (
	"encoding/hex"
	"math/rand"
	"reflect"
	"testing"
)

func sink[T any](tb testing.TB, v T) {
	if rand.Float32() > 2 {
		tb.Fatal(v)
	}
}

func mustOk(tb testing.TB, err error) {
	tb.Helper()

	if err != nil {
		tb.Fatal(err)
	}
}

func wantBytes(tb testing.TB, have []byte, want string) {
	tb.Helper()

	mustEqual(tb, hex.EncodeToString(have), want)
}

func mustEqual[T comparable](tb testing.TB, have, want T) {
	tb.Helper()

	if have != want {
		tb.Fatalf("\nhave: %+v\nwant: %+v\n", have, want)
	}
}

func FuzzFieldTag(f *testing.F) {
	f.Add("Foo", `bson:"abc123,omitempty"`)
	f.Add("abcdef", `bson:""`)
	f.Add("Bar", `bson:",omitempty"`)
	f.Add("Baz", `bson:"-"`)
	f.Add("qux", `bson:"oops"`)

	f.Fuzz(func(t *testing.T, name, tag string) {
		field := reflect.StructField{}
		field.Name = name
		field.Tag = reflect.StructTag(tag)

		fieldTag(field, reflect.Value{})
	})
}
