package core

// Read/write access to peer objects is regulated by the state mutex.
// We do not support concurrent state merging, hence, holding state mutex regulates
// access to all peers.

type peer struct {
	id        string
	entries   [][]byte
	rootHash  []byte
	prevHash  []byte
	epoch     uint64
	signature []byte
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

func (p *peer) update(epoch uint64, entries [][]byte, rootHash, prevHash []byte) {
	p.epoch = epoch
	p.entries = entries
	p.rootHash = rootHash
	p.prevHash = prevHash
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
	p.signature = nil
}

func (p *peer) haveFavourite() bool {
	if p.entries == nil {
		return false
	}

	return true
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
