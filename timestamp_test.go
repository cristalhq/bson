package bson

import "testing"

func TestTimestamp(t *testing.T) {
	ts := Timestamp(4294967298 << 10)

	mustEqual(t, ts.String(), "Timestamp(4398046513152, 1024)")
}
