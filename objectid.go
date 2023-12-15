package bson

import (
	"crypto/rand"
	"encoding"
	"encoding/hex"
	"encoding/json"
	"errors"
	mathrand "math/rand"
	"sync/atomic"
	"time"
)

var ErrBadObjectID = errors.New("not a valid ObjectID")

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
// Example: ObjectID('64d526fa37931c1e97eea90f').
func (oid ObjectID) String() string {
	return "ObjectID('" + hex.EncodeToString(oid[:]) + "')"
}

// MarshalBSON implements [bson.Marshaler].
func (oid *ObjectID) MarshalBSON() ([]byte, error) {
	b := make([]byte, len(oid))
	copy(b, oid[:])
	return b, nil
}

// UnmarshalBSON implements [bson.Unmarshaler].
func (oid *ObjectID) UnmarshalBSON(b []byte) error {
	switch len(b) {
	case 12:
		copy(oid[:], b)
		return nil
	case 24:
		n, err := hex.Decode(oid[:], b)
		if n != 12 {
			panic("unreachable")
		}
		return err
	default:
		return ErrBadObjectID
	}
}

// MarshalText implements [encoding.TextMarshaler].
func (oid ObjectID) MarshalText() ([]byte, error) {
	return oid.MarshalBSON()
}

// UnmarshalText implements [encoding.TextUnmarshaler].
func (oid *ObjectID) UnmarshalText(b []byte) error {
	return oid.UnmarshalBSON(b)
}

// MarshalBinary implements [encoding.BinaryMarshaler].
func (oid ObjectID) MarshalBinary() ([]byte, error) {
	return oid.MarshalBSON()
}

// UnmarshalBinary implements [encoding.BinaryUnmarshaler].
func (oid *ObjectID) UnmarshalBinary(b []byte) error {
	return oid.UnmarshalBSON(b)
}

// MarshalJSON implements [json.Marshaler].
func (oid ObjectID) MarshalJSON() ([]byte, error) {
	return oid.MarshalText()
}

// UnmarshalJSON implements [json.Unmarshaler].
func (oid *ObjectID) UnmarshalJSON(b []byte) error {
	return oid.UnmarshalBSON(b)
}

var (
	procUniqueID    [5]byte
	objectIDCounter atomic.Uint32
)

func init() {
	must(rand.Read(procUniqueID[:]))
	objectIDCounter.Store(mathrand.Uint32())
}

var (
	_ Marshaler                  = &ObjectID{}
	_ Unmarshaler                = &ObjectID{}
	_ encoding.TextMarshaler     = &ObjectID{}
	_ encoding.TextUnmarshaler   = &ObjectID{}
	_ encoding.BinaryMarshaler   = &ObjectID{}
	_ encoding.BinaryUnmarshaler = &ObjectID{}
	_ json.Marshaler             = &ObjectID{}
	_ json.Unmarshaler           = &ObjectID{}
)
