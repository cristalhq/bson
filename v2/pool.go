package bson

type pool struct{}

func (p *pool) Get(l int) []byte {
	return make([]byte, l)
}

func (p *pool) Put(b []byte) {
}
