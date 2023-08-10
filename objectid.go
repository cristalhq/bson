package bson

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	mathrand "math/rand"
	"sync/atomic"
	"time"
)

var (
	_ Marshaler   = &ObjectID{}
	_ Unmarshaler = &ObjectID{}
)

var ErrBadObjectID = errors.New("provided input is not a valid ObjectID")

// ObjectID represents BSON object ID.
type ObjectID [12]byte

// NewObjectID returns a new ObjectID.
func NewObjectID() ObjectID {
	return NewObjectIDWithTime(time.Now())
}

// NewObjectIDWithTime returns a new ObjectID.
func NewObjectIDWithTime(t time.Time) ObjectID {
	ts := uint32(t.UTC().Unix())
	c := objectIDCounter.Add(1)

	var oid ObjectID
	oid[0] = byte(ts >> 24)
	oid[1] = byte(ts >> 16)
	oid[2] = byte(ts >> 8)
	oid[3] = byte(ts)

	oid[4] = procUniqueID[0]
	oid[5] = procUniqueID[1]
	oid[6] = procUniqueID[2]
	oid[7] = procUniqueID[3]
	oid[8] = procUniqueID[4]

	oid[9] = byte(c >> 16)
	oid[10] = byte(c >> 8)
	oid[11] = byte(c)
	return oid
}

// String returns a hex string representation of the id.
// Example: ObjectIdHex('64d526fa37931c1e97eea90f').
func (oid ObjectID) String() string {
	return fmt.Sprintf(`ObjectIdHex('%x')`, string(oid[:]))
}

func (oid *ObjectID) MarshalBSON() ([]byte, error) {
	b := make([]byte, len(oid))
	copy(b, oid[:])
	return b, nil
}

func (oid *ObjectID) UnmarshalBSON(b []byte) error {
	switch len(b) {
	case 12:
		copy(oid[:], b)
		return nil
	case 24:
		n, err := hex.Decode(oid[:], b)
		if n != 24 {
			panic("unreachable")
		}
		return err
	default:
		return ErrBadObjectID
	}
}

var (
	procUniqueID    [5]byte
	objectIDCounter atomic.Uint32
)

func init() {
	must(rand.Read(procUniqueID[:]))
	objectIDCounter.Store(mathrand.Uint32())
}
