package bson

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func putBool(b bool) byte {
	if b {
		return 1
	}
	return 0
}

func putUint32(v uint32) [4]byte {
	var b [4]byte
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
	return b
}

func putUint64(v uint64) [8]byte {
	var b [8]byte
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
	b[4] = byte(v >> 32)
	b[5] = byte(v >> 40)
	b[6] = byte(v >> 48)
	b[7] = byte(v >> 56)
	return b
}
