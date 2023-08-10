package bson

import (
	"testing"
)

func TestObjectID(t *testing.T) {
	oid := ObjectID([12]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12})

	mustEqual(t, oid.String(), "ObjectIdHex('0102030405060708090a0b0c')")
}
