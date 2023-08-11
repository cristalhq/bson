package bson

import "testing"

func TestTimestamp(t *testing.T) {
	ts := Timestamp(4294967298)

	mustEqual(t, ts.String(), "Timestamp(4294967298, 1)")
}
