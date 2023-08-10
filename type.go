package bson

// Type represents a BSON type.
type Type byte

// BSON element types as described in https://bsonspec.org/spec.html.
const (
	TypeDouble          Type = 0x01
	TypeString          Type = 0x02
	TypeDocument        Type = 0x03
	TypeArray           Type = 0x04
	TypeBinary          Type = 0x05
	TypeUndefined       Type = 0x06
	TypeObjectID        Type = 0x07
	TypeBool            Type = 0x08
	TypeDateTime        Type = 0x09
	TypeNull            Type = 0x0a
	TypeRegex           Type = 0x0b
	TypeDBPointer       Type = 0x0c
	TypeCodeWithScope   Type = 0x0d
	TypeSymbol          Type = 0x0e
	TypeJavaScriptScope Type = 0x0f
	TypeInt32           Type = 0x10
	TypeTimestamp       Type = 0x11
	TypeInt64           Type = 0x12
	TypeDecimal         Type = 0x13
	TypeMinKey          Type = 0xff
	TypeMaxKey          Type = 0x7f
)
