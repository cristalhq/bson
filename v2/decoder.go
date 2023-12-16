package bson

import (
	"bufio"
	"io"
)

type Decoder struct {
	r *bufio.Reader
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		r: bufio.NewReader(r),
	}
}

func (dec *Decoder) Decode() (any, error) {
	panic("TODO")
}
