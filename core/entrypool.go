package core

import (
	"errors"

	"github.com/joonnna/blocks/protobuf"
)

var (
	errFullPending = errors.New("Pending queue is full")
)

/*
// Allow buffering of block entries equal to 2 full blocks
const (
	maxSize = maxBlockSize * 2
)
*/

type entryPool struct {
	/*
		cap      uint32
		currSize uint32
		full     bool
	*/
	pending map[string]*entry

	confirmed map[string]*entry

	missing map[string]bool
}

func newEntryPool() *entryPool {
	return &entryPool{
		pending:   make(map[string]*entry),
		confirmed: make(map[string]*entry),
		missing:   make(map[string]bool),
		//cap:       uint32(maxSize),
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

func (e *entryPool) isMissing(key string) bool {
	_, exists := e.missing[key]

	return exists
}

func (e *entryPool) getPending(key string) *entry {
	return e.pending[key]
}

func (e *entryPool) getConfirmed(key string) *entry {
	return e.confirmed[key]
}

func (e *entryPool) addPending(new *entry) {
	/*
		if e.full {
			return errFullPending
		}

		size := len(new.data)

		if (e.currSize + uint32(size)) > e.cap {
			e.full = true
			return errFullPending
		}

		e.currSize += uint32(size)
	*/
	e.pending[string(new.hash)] = new
}

func (e *entryPool) addConfirmed(new *entry) {
	e.confirmed[string(new.hash)] = new
}

func (e *entryPool) addMissing(key string) {
	e.missing[key] = true
}

func (e *entryPool) removePending(key string) {
	/*
		if entry, exists := e.pending[key]; exists {
			e.currSize -= uint32(len(entry.data))

			if e.currSize < e.cap {
				e.full = false
			}
		}
	*/
	delete(e.pending, key)
}

func (e *entryPool) removeConfirmed(key string) {
	delete(e.confirmed, key)
}

func (e *entryPool) removeMissing(key string) {
	delete(e.missing, key)
}

func (e *entryPool) resetMissing() {
	e.missing = make(map[string]bool)
}

func (e *entryPool) getAllMissing() [][]byte {
	ret := make([][]byte, len(e.missing))

	idx := 0
	for k, _ := range e.missing {
		ret[idx] = []byte(k)
		idx++
	}

	return ret
}

func (e *entryPool) diff(other [][]byte) []*blockchain.BlockEntry {
	if other == nil {
		return nil
	}

	var ret []*blockchain.BlockEntry

	for _, k := range other {
		key := string(k)
		if v, exists := e.pending[key]; exists {
			ret = append(ret, &blockchain.BlockEntry{Content: v.data, Hash: v.hash})
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

/*
func (e *entryPool) isFull() bool {
	return e.full
}

*/
