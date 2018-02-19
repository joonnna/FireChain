package core

import (
	"bytes"
	"crypto/sha256"

	"github.com/cbergoon/merkletree"
)

type entry struct {
	data []byte
	hash []byte
}

func (e *entry) CalculateHash() []byte {
	return sha256.New().Sum(e.data)
}

func (e *entry) Equals(other merkletree.Content) bool {
	inst, _ := other.(*entry)
	return e.equal(inst)
}

func (e *entry) cmp(other *entry) int {
	return bytes.Compare(e.data, other.data)
}

func (e *entry) equal(other *entry) bool {
	return bytes.Equal(e.data, other.data)
}
