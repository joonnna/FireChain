package blockchain

import (
	"bytes"
	"hash"
)

type entry struct {
	data []byte
	h    hash.Hash
}

func (e *entry) CalculateHash() []byte {
	return e.h.Sum256([]byte{e.data})
}

func (e *entry) Equals(other *entry) bool {
	return bytes.Equal(e.data, other.data)
}
