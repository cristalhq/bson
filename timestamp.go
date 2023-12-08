package bson

import (
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
	"sync/atomic"
	"time"
)

// Timestamp represents BSON type Timestamp.
type Timestamp int64

// NowTimestamp returns a timestamp with current UTC time.
func NowTimestamp() Timestamp {
	c := timestampCounter.Add(1)
	return NewTimestampWithCounter(time.Now().UTC(), c)
}

// NewTimestamp returns a timestamp with a given time.
func NewTimestamp(t time.Time) Timestamp {
	c := timestampCounter.Add(1)
	return NewTimestampWithCounter(t, c)
}

// NewTimestampWithCounter returns a timestamp with a given time and counter.
func NewTimestampWithCounter(t time.Time, c uint32) Timestamp {
	return Timestamp((t.UTC().Unix() << 32) | int64(c))
}

// String returns a hex string representation of the id.
// Example: Timestamp('4398046513152, 1024').
func (ts Timestamp) String() string {
	return fmt.Sprintf(`Timestamp(%d, %d)`, ts, ts>>32)
}

// Time returns time.Time ignoring increment.
func (ts Timestamp) Time() time.Time {
	return time.Unix(int64(ts>>32), 0)
}

// Counter returns timestamp counter.
func (ts Timestamp) Counter() uint32 {
	return uint32(ts)
}

// MarshalBSON implements [Marshaler].
func (ts Timestamp) MarshalBSON() ([]byte, error) {
	var b [8]byte
	b[0] = byte(ts)
	b[1] = byte(ts >> 8)
	b[2] = byte(ts >> 16)
	b[3] = byte(ts >> 24)
	b[4] = byte(ts >> 32)
	b[5] = byte(ts >> 40)
	b[6] = byte(ts >> 48)
	b[7] = byte(ts >> 56)
	return b[:], nil
}

// UnmarshalBSON implements [Unmarshaler].
func (ts *Timestamp) UnmarshalBSON(b []byte) error {
	if len(b) < 8 {
		return errors.New("not enough bytes for timestamp")
	}
	v := uint64(b[0]) |
		uint64(b[1])<<8 |
		uint64(b[2])<<16 |
		uint64(b[3])<<24 |
		uint64(b[4])<<32 |
		uint64(b[5])<<40 |
		uint64(b[6])<<48 |
		uint64(b[7])<<56
	*ts = Timestamp(v)
	return nil
}

// MarshalText implements [encoding.TextMarshaler].
func (ts Timestamp) MarshalText() ([]byte, error) {
	return ts.MarshalBSON()
}

// UnmarshalText implements [encoding.TextUnmarshaler].
func (ts *Timestamp) UnmarshalText(b []byte) error {
	return ts.UnmarshalBSON(b)
}

// MarshalBinary implements [encoding.BinaryMarshaler].
func (ts Timestamp) MarshalBinary() ([]byte, error) {
	return ts.MarshalBSON()
}

// UnmarshalBinary implements [encoding.BinaryUnmarshaler].
func (ts *Timestamp) UnmarshalBinary(b []byte) error {
	return ts.UnmarshalBSON(b)
}

// MarshalJSON implements [json.Marshaler].
func (ts Timestamp) MarshalJSON() ([]byte, error) {
	return ts.MarshalBSON()
}

// UnmarshalJSON implements [json.Unmarshaler].
func (ts *Timestamp) UnmarshalJSON(b []byte) error {
	return ts.UnmarshalBSON(b)
}

var timestampCounter atomic.Uint32

var (
	_ Marshaler                  = Timestamp(0)
	_ Unmarshaler                = new(Timestamp)
	_ encoding.TextMarshaler     = Timestamp(0)
	_ encoding.TextUnmarshaler   = new(Timestamp)
	_ encoding.BinaryMarshaler   = Timestamp(0)
	_ encoding.BinaryUnmarshaler = new(Timestamp)
	_ json.Marshaler             = Timestamp(0)
	_ json.Unmarshaler           = new(Timestamp)
)
