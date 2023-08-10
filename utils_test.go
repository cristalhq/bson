package bson

import "testing"

func mustEqual[T comparable](tb testing.TB, have, want T) {
	tb.Helper()

	if have != want {
		tb.Fatalf("\nhave: %+v\nwant: %+v\n", have, want)
	}
}
