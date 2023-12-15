package bson

import (
	"testing"
	"time"
)

func TestObjectID(t *testing.T) {
	oid := ObjectID([12]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12})

	t.Run("NewObjectID", func(t *testing.T) {
		now := time.Now()
		id := NewObjectID()
		counter := objectIDCounter.Load()
		b := [12]byte(id)
		ts := uint32(now.UTC().Unix())

		mustEqual(t, b[0], byte(ts>>24))
		mustEqual(t, b[1], byte(ts>>16))
		mustEqual(t, b[2], byte(ts>>8))
		mustEqual(t, b[3], byte(ts>>0))

		mustEqual(t, b[4], procUniqueID[0])
		mustEqual(t, b[5], procUniqueID[1])
		mustEqual(t, b[6], procUniqueID[2])
		mustEqual(t, b[7], procUniqueID[3])
		mustEqual(t, b[8], procUniqueID[4])

		mustEqual(t, b[9], byte(counter>>16))
		mustEqual(t, b[10], byte(counter>>8))
		mustEqual(t, b[11], byte(counter>>0))
	})

	t.Run("Marshal", func(t *testing.T) {
		mustEqual(t, oid.String(), "ObjectID('0102030405060708090a0b0c')")
		wantBytes(t, must(oid.MarshalText()), "0102030405060708090a0b0c")
		wantBytes(t, must(oid.MarshalBinary()), "0102030405060708090a0b0c")
		wantBytes(t, must(oid.MarshalJSON()), "0102030405060708090a0b0c")
		wantBytes(t, must(oid.MarshalBSON()), "0102030405060708090a0b0c")
	})

	buf := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}

	t.Run("UnmarshalBSON", func(t *testing.T) {
		var id ObjectID
		err := id.UnmarshalBSON(buf)
		mustOk(t, err)
		mustEqual(t, oid, id)
	})
	t.Run("UnmarshalJSON", func(t *testing.T) {
		var id ObjectID
		err := id.UnmarshalJSON(buf)
		mustOk(t, err)
		mustEqual(t, oid, id)
	})
	t.Run("UnmarshalBinary", func(t *testing.T) {
		var id ObjectID
		err := id.UnmarshalBinary(buf)
		mustOk(t, err)
		mustEqual(t, oid, id)
	})
	t.Run("UnmarshalText", func(t *testing.T) {
		var id ObjectID
		err := id.UnmarshalText(buf)
		mustOk(t, err)
		mustEqual(t, oid, id)
	})
	t.Run("UnmarshalBSON_hex", func(t *testing.T) {
		buf := []byte("0102030405060708090a0b0c")
		var id ObjectID
		err := id.UnmarshalBSON(buf)
		mustOk(t, err)
		mustEqual(t, oid, id)
	})
	t.Run("UnmarshalBSON_bad", func(t *testing.T) {
		buf := []byte("01020__30405060708090a0b0c")
		var id ObjectID
		err := id.UnmarshalBSON(buf)
		mustFail(t, err)
	})
}
