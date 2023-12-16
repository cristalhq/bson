package bson

//go:generate stringer -linecomment -output scalars_string.go -type BinarySubtype

// BinarySubtype represents BSON Binary's subtype.
type BinarySubtype byte

const (
	// BinaryGeneric represents a BSON generic binary subtype.
	BinaryGeneric = BinarySubtype(0x00) // generic

	// BinaryFunction represents a BSON function binary subtype.
	BinaryFunction = BinarySubtype(0x01) // function

	// BinaryBinaryOld represents a BSON binary-old binary subtype.
	BinaryBinaryOld = BinarySubtype(0x02) // binary-old

	// BinaryUUIDOld represents a BSON uuid-old binary subtype.
	BinaryUUIDOld = BinarySubtype(0x03) // uuid-old

	// BinaryUUID represents a BSON UUID binary subtype.
	BinaryUUID = BinarySubtype(0x04) // uuid

	// BinaryMD5 represents a BSON md5 binary subtype.
	BinaryMD5 = BinarySubtype(0x05) // md5

	// BinaryEncrypted represents a Encrypted BSON value binary subtype.
	BinaryEncrypted = BinarySubtype(0x06) // encrypted

	// BinaryEncrypted represents a Compressed BSON column binary subtype.
	BinaryCompressed = BinarySubtype(0x07) // compressed

	// BinaryUser represents a User defined binary subtype.
	BinaryUser = BinarySubtype(0x80) // user
)

// Binary represents a BSON binary data in the decoded form.
type Binary struct {
	B       []byte
	Subtype BinarySubtype
}

// ObjectID represents a BSON ObjectID data in the decoded form.
type ObjectID [12]byte

// NullType represents a BSON null type.
type NullType struct{}

// Null represents a BSON null value.
var Null NullType

// Regex represents a BSON ObjectID regular expression in the decoded form.
type Regex struct {
	// TODO
}

// Timestamp represents a BSON timestamp in the decoded form.
type Timestamp uint64
