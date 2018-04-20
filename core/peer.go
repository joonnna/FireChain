package core

import "bytes"

// Read/write access to peer objects is regulated by the state mutex.
// We do not support concurrent state merging, hence, holding state mutex regulates
// access to all peers.

type peer struct {
	id       string
	entries  [][]byte
	rootHash []byte
	prevHash []byte
	epoch    uint64

	// Elliptic signature
	r, s []byte

	// Only for worm interaction
	// TODO solve this differently
	httpAddr string
}

func (p *peer) getRootHash() []byte {
	return p.rootHash
}

func (p *peer) getPrevHash() []byte {
	return p.prevHash
}

func (p *peer) getEpoch() uint64 {
	return p.epoch
}

func (p *peer) getEntries() [][]byte {
	return p.entries
}

func (p *peer) getR() []byte {
	return p.r
}
func (p *peer) getS() []byte {
	return p.s
}

func (p *peer) addSignature(r, s []byte) {
	p.r = r
	p.s = s
}

func (p *peer) update(epoch uint64, entries [][]byte, rootHash, prevHash, r, s []byte) {
	p.epoch = epoch
	p.entries = entries
	p.rootHash = rootHash
	p.prevHash = prevHash
	p.r = r
	p.s = s
}

func (p *peer) addBlock(b *block) {
	p.rootHash = b.getRootHash()
	p.prevHash = b.getPrevHash()
	p.entries = b.getHashes()
	p.epoch++
}

func (p *peer) increment(rootHash, prevHash []byte, entries [][]byte) {
	p.epoch++
	p.entries = entries
	p.rootHash = rootHash
	p.prevHash = prevHash
}

func (p *peer) reset() {
	p.entries = nil
	p.rootHash = nil
	p.prevHash = nil
	p.r = nil
	p.s = nil
}

func (p *peer) hasFavourite() bool {
	if p.entries == nil {
		return false
	}

	return true
}

func (p *peer) isEqual(other *peer) bool {
	return bytes.Equal(p.rootHash, other.rootHash) && bytes.Equal(p.prevHash, other.prevHash)
}

/*
func (p *peer) getBlock() *block {
	p.RLock()
	defer p.RUnlock()

	return p.block
}

func (p *peer) setFirstBlock(newBlock *block) {
	p.Lock()
	defer p.Unlock()

	if p.block == nil {
		p.block = newBlock
	}
}
*/
