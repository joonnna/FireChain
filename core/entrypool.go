package core

import (
	"github.com/joonnna/blocks/protobuf"
)

type entryPool struct {
	pending map[string]*entry

	confirmed map[string]*entry
}

func newEntryPool() *entryPool {
	return &entryPool{
		pending:   make(map[string]*entry),
		confirmed: make(map[string]*entry),
	}
}

func (e *entryPool) exists(key string) bool {
	return e.isPending(key) || e.isConfirmed(key)
}

func (e *entryPool) isPending(key string) bool {
	_, exists := e.pending[key]

	return exists
}

func (e *entryPool) isConfirmed(key string) bool {
	_, exists := e.confirmed[key]

	return exists
}

func (e *entryPool) getPending(key string) *entry {
	return e.pending[key]
}

func (e *entryPool) getConfirmed(key string) *entry {
	return e.confirmed[key]
}

func (e *entryPool) addPending(new *entry) {
	e.pending[string(new.hash)] = new
}

func (e *entryPool) addConfirmed(new *entry) {
	e.confirmed[string(new.hash)] = new
}

func (e *entryPool) removePending(key string) {
	delete(e.pending, key)
}

func (e *entryPool) removeConfirmed(key string) {
	delete(e.confirmed, key)
}

func (e *entryPool) diff(other map[string]bool) []*blockchain.BlockEntry {
	var ret []*blockchain.BlockEntry

	if other == nil {
		for _, v := range e.pending {
			ret = append(ret, &blockchain.BlockEntry{Content: v.data, Hash: v.hash})
		}
	} else {
		for k, v := range e.pending {
			if _, exists := other[k]; !exists {
				ret = append(ret, &blockchain.BlockEntry{Content: v.data, Hash: v.hash})
			}
		}
	}
	return ret
}

func (e *entryPool) fillWithPending(b *block) bool {
	for _, ent := range e.pending {
		err := b.add(ent)
		if err == errFullBlock {
			return true
		}
	}

	return false
}

func (e *entryPool) getAllPending() []*entry {
	idx := 0
	ret := make([]*entry, len(e.pending))

	for _, v := range e.pending {
		ret[idx] = v
		idx++
	}
	return ret
}
