package core

import (
	"errors"
	"math/rand"
	"time"

	log "github.com/inconshreveable/log15"
	"github.com/joonnna/blocks/protobuf"
)

var (
	errFullPending = errors.New("Pending queue is full")
)

type entryPool struct {
	cap      uint32
	currSize uint32
	full     bool

	pending   map[string]*entry
	confirmed map[string]*entry
	missing   map[string]bool
	favorite  map[string]bool
}

func newEntryPool(cap uint32) *entryPool {
	return &entryPool{
		pending:   make(map[string]*entry),
		confirmed: make(map[string]*entry),
		missing:   make(map[string]bool),
		favorite:  make(map[string]bool),
		cap:       cap,
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

func (e *entryPool) addFavorite(key string) {
	e.favorite[key] = true
}

func (e *entryPool) addPending(new *entry) error {
	if e.full {
		return errFullPending
	}

	size := len(new.data)

	newSize := e.currSize + uint32(size)

	if newSize > e.cap {
		e.full = true
		return errFullPending
	} else if newSize == e.cap {
		e.full = true
	}

	e.currSize += uint32(size)

	e.pending[string(new.hash)] = new

	return nil
}

func (e *entryPool) missingToPending(ent *entry) {
	key := string(ent.hash)
	e.removeMissing(key)
	if e.isPending(key) {
		return
	}

	err := e.addPending(ent)
	if err != nil {
		log.Info("Full pending, removing random entry")
		e.replaceRandomPending(ent)
	}
}

func (e *entryPool) addConfirmed(new *entry) {
	e.confirmed[string(new.hash)] = new
}

func (e *entryPool) addMissing(key string) {
	e.missing[key] = true
}

func (e *entryPool) replaceRandomPending(ent *entry) {
	for k, _ := range e.pending {
		if _, exist := e.favorite[k]; !exist {
			e.removePending(k)
		}

		if err := e.addPending(ent); err == nil {
			break
		}
	}
}

func (e *entryPool) removePending(key string) {
	if entry, exists := e.pending[key]; exists {
		e.currSize -= uint32(len(entry.data))

		if e.currSize < e.cap {
			e.full = false
		}
	}

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

func (e *entryPool) resetFavorite() {
	e.favorite = make(map[string]bool)
}

func (e *entryPool) getAllMissing() [][]byte {
	ret := make([][]byte, 0, len(e.missing))

	for k, _ := range e.missing {
		ret = append(ret, []byte(k))
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

func (e *entryPool) fillWithPending(b *block) {
	for _, ent := range e.pending {
		err := b.add(ent)
		if err == errFullBlock {
			return
		}
		e.addFavorite(string(ent.hash))
	}

	rand.Seed(time.Now().UnixNano())
	for {
		buf := make([]byte, entrySize)
		_, err := rand.Read(buf)
		if err != nil {
			log.Error(err.Error())
			continue
		}

		ent := &entry{
			data: buf,
			hash: hashBytes(buf),
		}

		err = b.add(ent)
		if err == errFullBlock {
			return
		} else if err == nil {
			e.replaceRandomPending(ent)
			e.addFavorite(string(ent.hash))
		} else {
			log.Error(err.Error())
		}
	}
}

func (e *entryPool) getAllPending() []*entry {
	ret := make([]*entry, 0, len(e.pending))

	for _, v := range e.pending {
		ret = append(ret, v)
	}
	return ret
}

/*
func (e *entryPool) isFull() bool {
	return e.full
}

*/
