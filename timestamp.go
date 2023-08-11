package bson

import (
	"encoding/binary"
	"fmt"
	"sync/atomic"
	"time"
)

var (
	_ Marshaler   = Timestamp(0)
	_ Unmarshaler = new(Timestamp)
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
// Example: Timestamp('64d526fa37931c1e97eea90f').
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

func (ts *Timestamp) UnmarshalBSON(b []byte) error {
	v := binary.LittleEndian.Uint64(b)
	*ts = Timestamp(v)
	return nil
}

var timestampCounter atomic.Uint32
