// Code generated by protoc-gen-go. DO NOT EDIT.
// source: chainstate.proto

/*
Package blockchain is a generated protocol buffer package.

It is generated from these files:
	chainstate.proto

It has these top-level messages:
	EntryHeader
	BlockEntry
	BlockHeader
	BlockContent
	Signature
	PeerState
	PeerInfo
	State
	StateResponse
	ChainState
	ChainContent
	WormPayload
*/
package blockchain

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type EntryHeader struct {
	Hash []byte `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
}

func (m *EntryHeader) Reset()                    { *m = EntryHeader{} }
func (m *EntryHeader) String() string            { return proto.CompactTextString(m) }
func (*EntryHeader) ProtoMessage()               {}
func (*EntryHeader) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *EntryHeader) GetHash() []byte {
	if m != nil {
		return m.Hash
	}
	return nil
}

type BlockEntry struct {
	Hash    []byte `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	Content []byte `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"`
}

func (m *BlockEntry) Reset()                    { *m = BlockEntry{} }
func (m *BlockEntry) String() string            { return proto.CompactTextString(m) }
func (*BlockEntry) ProtoMessage()               {}
func (*BlockEntry) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *BlockEntry) GetHash() []byte {
	if m != nil {
		return m.Hash
	}
	return nil
}

func (m *BlockEntry) GetContent() []byte {
	if m != nil {
		return m.Content
	}
	return nil
}

type BlockHeader struct {
	RootHash []byte `protobuf:"bytes,1,opt,name=rootHash,proto3" json:"rootHash,omitempty"`
	PrevHash []byte `protobuf:"bytes,2,opt,name=prevHash,proto3" json:"prevHash,omitempty"`
	BlockNum uint64 `protobuf:"varint,3,opt,name=blockNum" json:"blockNum,omitempty"`
}

func (m *BlockHeader) Reset()                    { *m = BlockHeader{} }
func (m *BlockHeader) String() string            { return proto.CompactTextString(m) }
func (*BlockHeader) ProtoMessage()               {}
func (*BlockHeader) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *BlockHeader) GetRootHash() []byte {
	if m != nil {
		return m.RootHash
	}
	return nil
}

func (m *BlockHeader) GetPrevHash() []byte {
	if m != nil {
		return m.PrevHash
	}
	return nil
}

func (m *BlockHeader) GetBlockNum() uint64 {
	if m != nil {
		return m.BlockNum
	}
	return 0
}

type BlockContent struct {
	RootHash []byte   `protobuf:"bytes,1,opt,name=rootHash,proto3" json:"rootHash,omitempty"`
	PrevHash []byte   `protobuf:"bytes,2,opt,name=prevHash,proto3" json:"prevHash,omitempty"`
	Content  [][]byte `protobuf:"bytes,3,rep,name=content,proto3" json:"content,omitempty"`
	BlockNum uint64   `protobuf:"varint,4,opt,name=blockNum" json:"blockNum,omitempty"`
}

func (m *BlockContent) Reset()                    { *m = BlockContent{} }
func (m *BlockContent) String() string            { return proto.CompactTextString(m) }
func (*BlockContent) ProtoMessage()               {}
func (*BlockContent) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *BlockContent) GetRootHash() []byte {
	if m != nil {
		return m.RootHash
	}
	return nil
}

func (m *BlockContent) GetPrevHash() []byte {
	if m != nil {
		return m.PrevHash
	}
	return nil
}

func (m *BlockContent) GetContent() [][]byte {
	if m != nil {
		return m.Content
	}
	return nil
}

func (m *BlockContent) GetBlockNum() uint64 {
	if m != nil {
		return m.BlockNum
	}
	return 0
}

type Signature struct {
	R []byte `protobuf:"bytes,1,opt,name=r,proto3" json:"r,omitempty"`
	S []byte `protobuf:"bytes,2,opt,name=s,proto3" json:"s,omitempty"`
}

func (m *Signature) Reset()                    { *m = Signature{} }
func (m *Signature) String() string            { return proto.CompactTextString(m) }
func (*Signature) ProtoMessage()               {}
func (*Signature) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *Signature) GetR() []byte {
	if m != nil {
		return m.R
	}
	return nil
}

func (m *Signature) GetS() []byte {
	if m != nil {
		return m.S
	}
	return nil
}

type PeerState struct {
	Id          string     `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	EntryHashes [][]byte   `protobuf:"bytes,2,rep,name=entryHashes,proto3" json:"entryHashes,omitempty"`
	RootHash    []byte     `protobuf:"bytes,3,opt,name=rootHash,proto3" json:"rootHash,omitempty"`
	PrevHash    []byte     `protobuf:"bytes,4,opt,name=prevHash,proto3" json:"prevHash,omitempty"`
	Epoch       uint64     `protobuf:"varint,5,opt,name=epoch" json:"epoch,omitempty"`
	Signature   *Signature `protobuf:"bytes,6,opt,name=signature" json:"signature,omitempty"`
	// Only for worm interaction
	// TODO solve this differently
	HttpAddr string `protobuf:"bytes,7,opt,name=httpAddr" json:"httpAddr,omitempty"`
}

func (m *PeerState) Reset()                    { *m = PeerState{} }
func (m *PeerState) String() string            { return proto.CompactTextString(m) }
func (*PeerState) ProtoMessage()               {}
func (*PeerState) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *PeerState) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *PeerState) GetEntryHashes() [][]byte {
	if m != nil {
		return m.EntryHashes
	}
	return nil
}

func (m *PeerState) GetRootHash() []byte {
	if m != nil {
		return m.RootHash
	}
	return nil
}

func (m *PeerState) GetPrevHash() []byte {
	if m != nil {
		return m.PrevHash
	}
	return nil
}

func (m *PeerState) GetEpoch() uint64 {
	if m != nil {
		return m.Epoch
	}
	return 0
}

func (m *PeerState) GetSignature() *Signature {
	if m != nil {
		return m.Signature
	}
	return nil
}

func (m *PeerState) GetHttpAddr() string {
	if m != nil {
		return m.HttpAddr
	}
	return ""
}

type PeerInfo struct {
	Id       string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Epoch    uint64 `protobuf:"varint,2,opt,name=epoch" json:"epoch,omitempty"`
	RootHash []byte `protobuf:"bytes,3,opt,name=rootHash,proto3" json:"rootHash,omitempty"`
}

func (m *PeerInfo) Reset()                    { *m = PeerInfo{} }
func (m *PeerInfo) String() string            { return proto.CompactTextString(m) }
func (*PeerInfo) ProtoMessage()               {}
func (*PeerInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *PeerInfo) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *PeerInfo) GetEpoch() uint64 {
	if m != nil {
		return m.Epoch
	}
	return 0
}

func (m *PeerInfo) GetRootHash() []byte {
	if m != nil {
		return m.RootHash
	}
	return nil
}

type State struct {
	Peers          map[string]*PeerInfo `protobuf:"bytes,1,rep,name=peers" json:"peers,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	MissingEntries [][]byte             `protobuf:"bytes,2,rep,name=missingEntries,proto3" json:"missingEntries,omitempty"`
}

func (m *State) Reset()                    { *m = State{} }
func (m *State) String() string            { return proto.CompactTextString(m) }
func (*State) ProtoMessage()               {}
func (*State) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *State) GetPeers() map[string]*PeerInfo {
	if m != nil {
		return m.Peers
	}
	return nil
}

func (m *State) GetMissingEntries() [][]byte {
	if m != nil {
		return m.MissingEntries
	}
	return nil
}

type StateResponse struct {
	Peers   []*PeerState  `protobuf:"bytes,1,rep,name=peers" json:"peers,omitempty"`
	Entries []*BlockEntry `protobuf:"bytes,2,rep,name=entries" json:"entries,omitempty"`
}

func (m *StateResponse) Reset()                    { *m = StateResponse{} }
func (m *StateResponse) String() string            { return proto.CompactTextString(m) }
func (*StateResponse) ProtoMessage()               {}
func (*StateResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *StateResponse) GetPeers() []*PeerState {
	if m != nil {
		return m.Peers
	}
	return nil
}

func (m *StateResponse) GetEntries() []*BlockEntry {
	if m != nil {
		return m.Entries
	}
	return nil
}

type ChainState struct {
	Blocks map[uint64]*BlockHeader `protobuf:"bytes,1,rep,name=blocks" json:"blocks,omitempty" protobuf_key:"varint,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *ChainState) Reset()                    { *m = ChainState{} }
func (m *ChainState) String() string            { return proto.CompactTextString(m) }
func (*ChainState) ProtoMessage()               {}
func (*ChainState) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *ChainState) GetBlocks() map[uint64]*BlockHeader {
	if m != nil {
		return m.Blocks
	}
	return nil
}

type ChainContent struct {
	Blocks []*BlockContent `protobuf:"bytes,1,rep,name=blocks" json:"blocks,omitempty"`
}

func (m *ChainContent) Reset()                    { *m = ChainContent{} }
func (m *ChainContent) String() string            { return proto.CompactTextString(m) }
func (*ChainContent) ProtoMessage()               {}
func (*ChainContent) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *ChainContent) GetBlocks() []*BlockContent {
	if m != nil {
		return m.Blocks
	}
	return nil
}

type WormPayload struct {
	Blocks map[uint64]*BlockHeader `protobuf:"bytes,1,rep,name=blocks" json:"blocks,omitempty" protobuf_key:"varint,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *WormPayload) Reset()                    { *m = WormPayload{} }
func (m *WormPayload) String() string            { return proto.CompactTextString(m) }
func (*WormPayload) ProtoMessage()               {}
func (*WormPayload) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

func (m *WormPayload) GetBlocks() map[uint64]*BlockHeader {
	if m != nil {
		return m.Blocks
	}
	return nil
}

func init() {
	proto.RegisterType((*EntryHeader)(nil), "blockchain.EntryHeader")
	proto.RegisterType((*BlockEntry)(nil), "blockchain.BlockEntry")
	proto.RegisterType((*BlockHeader)(nil), "blockchain.BlockHeader")
	proto.RegisterType((*BlockContent)(nil), "blockchain.BlockContent")
	proto.RegisterType((*Signature)(nil), "blockchain.Signature")
	proto.RegisterType((*PeerState)(nil), "blockchain.PeerState")
	proto.RegisterType((*PeerInfo)(nil), "blockchain.PeerInfo")
	proto.RegisterType((*State)(nil), "blockchain.State")
	proto.RegisterType((*StateResponse)(nil), "blockchain.StateResponse")
	proto.RegisterType((*ChainState)(nil), "blockchain.ChainState")
	proto.RegisterType((*ChainContent)(nil), "blockchain.ChainContent")
	proto.RegisterType((*WormPayload)(nil), "blockchain.WormPayload")
}

func init() { proto.RegisterFile("chainstate.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 541 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xbc, 0x54, 0xdd, 0x6e, 0xd3, 0x30,
	0x14, 0x96, 0x9b, 0xa6, 0x5d, 0x4f, 0xca, 0x34, 0x59, 0x1b, 0x44, 0x15, 0x17, 0xc1, 0x48, 0x50,
	0x81, 0xa8, 0xa6, 0xec, 0x06, 0x95, 0x1b, 0x60, 0x42, 0x1a, 0x12, 0x9a, 0xa6, 0xec, 0x82, 0x6b,
	0xaf, 0x31, 0x4b, 0xb4, 0x36, 0x8e, 0x6c, 0x77, 0x52, 0x25, 0x1e, 0x05, 0xf1, 0x12, 0xbc, 0x0e,
	0x0f, 0x83, 0x6c, 0xe7, 0xc7, 0x49, 0xd1, 0x2e, 0xb8, 0xe0, 0x2e, 0xe7, 0xf8, 0x3b, 0xe7, 0xfb,
	0xbe, 0x63, 0x9f, 0xc0, 0xd1, 0x2a, 0xa3, 0x79, 0x21, 0x15, 0x55, 0x6c, 0x51, 0x0a, 0xae, 0x38,
	0x86, 0x9b, 0x35, 0x5f, 0xdd, 0x99, 0x34, 0x79, 0x06, 0xc1, 0xa7, 0x42, 0x89, 0xdd, 0x05, 0xa3,
	0x29, 0x13, 0x18, 0xc3, 0x30, 0xa3, 0x32, 0x0b, 0x51, 0x84, 0xe6, 0xd3, 0xc4, 0x7c, 0x93, 0x25,
	0xc0, 0x47, 0x5d, 0x60, 0x70, 0x7f, 0x43, 0xe0, 0x10, 0xc6, 0x2b, 0x5e, 0x28, 0x56, 0xa8, 0x70,
	0x60, 0xd2, 0x75, 0x48, 0x28, 0x04, 0xa6, 0xb6, 0x6a, 0x3f, 0x83, 0x03, 0xc1, 0xb9, 0xba, 0x68,
	0x1b, 0x34, 0xb1, 0x3e, 0x2b, 0x05, 0xbb, 0x37, 0x67, 0xb6, 0x4b, 0x13, 0xeb, 0x33, 0xa3, 0xf9,
	0x72, 0xbb, 0x09, 0xbd, 0x08, 0xcd, 0x87, 0x49, 0x13, 0x93, 0xef, 0x30, 0x35, 0x14, 0xe7, 0x96,
	0xf2, 0x9f, 0x39, 0x1c, 0x13, 0x5e, 0xe4, 0x39, 0x26, 0x3a, 0xec, 0xc3, 0x1e, 0xfb, 0x4b, 0x98,
	0x5c, 0xe7, 0xb7, 0x05, 0x55, 0x5b, 0xc1, 0xf0, 0x14, 0x90, 0xa8, 0x38, 0x91, 0xd0, 0x91, 0xac,
	0x58, 0x90, 0x24, 0xbf, 0x11, 0x4c, 0xae, 0x18, 0x13, 0xd7, 0xfa, 0x22, 0xf0, 0x21, 0x0c, 0xf2,
	0xd4, 0x40, 0x27, 0xc9, 0x20, 0x4f, 0x71, 0x04, 0x01, 0x33, 0xd7, 0x40, 0x65, 0xc6, 0x74, 0x95,
	0x16, 0xe0, 0xa6, 0x3a, 0xb6, 0xbc, 0x07, 0x6c, 0x0d, 0x7b, 0xb6, 0x8e, 0xc1, 0x67, 0x25, 0x5f,
	0x65, 0xa1, 0x6f, 0x94, 0xdb, 0x00, 0x9f, 0xc1, 0x44, 0xd6, 0xb2, 0xc3, 0x51, 0x84, 0xe6, 0x41,
	0x7c, 0xb2, 0x68, 0x9f, 0xc5, 0xa2, 0xf1, 0x94, 0xb4, 0x38, 0x4d, 0x93, 0x29, 0x55, 0x7e, 0x48,
	0x53, 0x11, 0x8e, 0x8d, 0xf4, 0x26, 0x26, 0x5f, 0xe0, 0x40, 0xbb, 0xfb, 0x5c, 0x7c, 0xe3, 0x7b,
	0xe6, 0x1a, 0x09, 0x03, 0x57, 0xc2, 0x03, 0x86, 0xc8, 0x2f, 0x04, 0xbe, 0x1d, 0x54, 0x0c, 0x7e,
	0xc9, 0x98, 0x90, 0x21, 0x8a, 0xbc, 0x79, 0x10, 0x3f, 0xed, 0x88, 0x34, 0x6f, 0x5a, 0xd3, 0x4a,
	0xf3, 0x36, 0x13, 0x0b, 0xc5, 0x2f, 0xe0, 0x70, 0x93, 0x4b, 0x99, 0x17, 0xb7, 0x3a, 0x9d, 0x37,
	0xf3, 0xec, 0x65, 0x67, 0x97, 0x00, 0x6d, 0x31, 0x3e, 0x02, 0xef, 0x8e, 0xed, 0x2a, 0xd9, 0xfa,
	0x13, 0xbf, 0x02, 0xff, 0x9e, 0xae, 0xb7, 0xcc, 0xe8, 0x0e, 0xe2, 0x63, 0x97, 0xbb, 0x36, 0x9b,
	0x58, 0xc8, 0x72, 0xf0, 0x16, 0x91, 0x02, 0x1e, 0x19, 0x49, 0x09, 0x93, 0x25, 0x2f, 0x24, 0xc3,
	0xaf, 0xbb, 0xe2, 0x4f, 0xfa, 0x0d, 0x2c, 0xba, 0x52, 0x7d, 0x0a, 0x63, 0xe6, 0xc8, 0x0d, 0xe2,
	0xc7, 0x2e, 0xbc, 0xdd, 0xc0, 0xa4, 0x86, 0x91, 0x1f, 0x08, 0xe0, 0x5c, 0x9f, 0xda, 0x51, 0x2d,
	0x61, 0x64, 0x0a, 0x6a, 0x3a, 0xe2, 0xd6, 0xb7, 0x38, 0xdb, 0xaa, 0x9a, 0x58, 0x55, 0x31, 0x4b,
	0xaa, 0x3d, 0xdd, 0x9f, 0xc5, 0xd0, 0xce, 0xe2, 0x4d, 0x77, 0x16, 0x4f, 0xf6, 0xb4, 0xd9, 0x0d,
	0x77, 0xc7, 0xf1, 0x1e, 0xa6, 0x86, 0xb5, 0x5e, 0xcc, 0xd3, 0x9e, 0xbe, 0x70, 0xaf, 0x47, 0x85,
	0xac, 0x55, 0x91, 0x9f, 0x08, 0x82, 0xaf, 0x5c, 0x6c, 0xae, 0xe8, 0x6e, 0xcd, 0x69, 0x8a, 0xdf,
	0xf5, 0x3a, 0x3c, 0x77, 0x3b, 0x38, 0xc0, 0xff, 0x65, 0xf1, 0x66, 0x64, 0x7e, 0xa8, 0x67, 0x7f,
	0x02, 0x00, 0x00, 0xff, 0xff, 0x35, 0xe8, 0x39, 0xab, 0x64, 0x05, 0x00, 0x00,
}
