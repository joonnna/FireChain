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
	PeerState
	State
	StateResponse
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

type BlockContent struct {
	RootHash []byte   `protobuf:"bytes,1,opt,name=rootHash,proto3" json:"rootHash,omitempty"`
	PrevHash []byte   `protobuf:"bytes,2,opt,name=prevHash,proto3" json:"prevHash,omitempty"`
	Content  [][]byte `protobuf:"bytes,3,rep,name=content,proto3" json:"content,omitempty"`
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

type PeerState struct {
	Id          string   `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	EntryHashes [][]byte `protobuf:"bytes,2,rep,name=entryHashes,proto3" json:"entryHashes,omitempty"`
	RootHash    []byte   `protobuf:"bytes,3,opt,name=rootHash,proto3" json:"rootHash,omitempty"`
	PrevHash    []byte   `protobuf:"bytes,4,opt,name=prevHash,proto3" json:"prevHash,omitempty"`
	Epoch       uint64   `protobuf:"varint,5,opt,name=epoch" json:"epoch,omitempty"`
	Signature   []byte   `protobuf:"bytes,6,opt,name=signature,proto3" json:"signature,omitempty"`
	// Only for worm interaction
	// TODO solve this differently
	HttpAddr string `protobuf:"bytes,7,opt,name=httpAddr" json:"httpAddr,omitempty"`
}

func (m *PeerState) Reset()                    { *m = PeerState{} }
func (m *PeerState) String() string            { return proto.CompactTextString(m) }
func (*PeerState) ProtoMessage()               {}
func (*PeerState) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

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

func (m *PeerState) GetSignature() []byte {
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

type State struct {
	Peers          map[string]*PeerState `protobuf:"bytes,1,rep,name=peers" json:"peers,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	MissingEntries [][]byte              `protobuf:"bytes,2,rep,name=missingEntries,proto3" json:"missingEntries,omitempty"`
}

func (m *State) Reset()                    { *m = State{} }
func (m *State) String() string            { return proto.CompactTextString(m) }
func (*State) ProtoMessage()               {}
func (*State) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *State) GetPeers() map[string]*PeerState {
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
	Entries []*BlockEntry `protobuf:"bytes,1,rep,name=entries" json:"entries,omitempty"`
}

func (m *StateResponse) Reset()                    { *m = StateResponse{} }
func (m *StateResponse) String() string            { return proto.CompactTextString(m) }
func (*StateResponse) ProtoMessage()               {}
func (*StateResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *StateResponse) GetEntries() []*BlockEntry {
	if m != nil {
		return m.Entries
	}
	return nil
}

type WormPayload struct {
	Blocks map[uint64]*BlockHeader `protobuf:"bytes,1,rep,name=blocks" json:"blocks,omitempty" protobuf_key:"varint,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *WormPayload) Reset()                    { *m = WormPayload{} }
func (m *WormPayload) String() string            { return proto.CompactTextString(m) }
func (*WormPayload) ProtoMessage()               {}
func (*WormPayload) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

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
	proto.RegisterType((*PeerState)(nil), "blockchain.PeerState")
	proto.RegisterType((*State)(nil), "blockchain.State")
	proto.RegisterType((*StateResponse)(nil), "blockchain.StateResponse")
	proto.RegisterType((*WormPayload)(nil), "blockchain.WormPayload")
}

func init() { proto.RegisterFile("chainstate.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 431 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x53, 0x4d, 0x6b, 0xdb, 0x40,
	0x14, 0x64, 0xf5, 0x61, 0xd7, 0x4f, 0x69, 0x08, 0x4b, 0x3f, 0x16, 0x93, 0x83, 0xaa, 0x42, 0x31,
	0x94, 0x8a, 0xe2, 0x5e, 0x4a, 0x7a, 0x4a, 0x4b, 0x20, 0xb7, 0x06, 0xf5, 0xd0, 0x6b, 0x37, 0xd2,
	0x23, 0x12, 0x71, 0xb4, 0x62, 0x77, 0x13, 0xf0, 0x9f, 0xe9, 0xaf, 0xe8, 0x5f, 0xe8, 0xff, 0x2a,
	0x7a, 0x2b, 0x4b, 0x1b, 0xbb, 0xf4, 0xd0, 0x9b, 0xde, 0xec, 0xcc, 0x78, 0xe6, 0xed, 0x1a, 0x4e,
	0xca, 0x5a, 0x36, 0xad, 0xb1, 0xd2, 0x62, 0xde, 0x69, 0x65, 0x15, 0x87, 0xeb, 0x8d, 0x2a, 0x6f,
	0x09, 0xce, 0x5e, 0x41, 0x72, 0xd1, 0x5a, 0xbd, 0xbd, 0x44, 0x59, 0xa1, 0xe6, 0x1c, 0xa2, 0x5a,
	0x9a, 0x5a, 0xb0, 0x94, 0xad, 0x8e, 0x0a, 0xfa, 0xce, 0xce, 0x00, 0x3e, 0xf7, 0x02, 0xe2, 0xfd,
	0x8d, 0xc1, 0x05, 0xcc, 0x4b, 0xd5, 0x5a, 0x6c, 0xad, 0x08, 0x08, 0xde, 0x8d, 0xd9, 0x05, 0x24,
	0xa4, 0x1d, 0xec, 0x97, 0xf0, 0x44, 0x2b, 0x65, 0x2f, 0x27, 0x83, 0x71, 0xee, 0xcf, 0x3a, 0x8d,
	0x0f, 0x74, 0xe6, 0x5c, 0xc6, 0x39, 0xfb, 0x01, 0x47, 0x64, 0xf3, 0xc5, 0xd9, 0xfe, 0xaf, 0x8f,
	0x1f, 0x34, 0x4c, 0x43, 0x3f, 0xe8, 0x6f, 0x06, 0x8b, 0x2b, 0x44, 0xfd, 0xad, 0xdf, 0x13, 0x3f,
	0x86, 0xa0, 0xa9, 0xc8, 0x79, 0x51, 0x04, 0x4d, 0xc5, 0x53, 0x48, 0x90, 0xb6, 0x24, 0x4d, 0x8d,
	0x46, 0x04, 0xa4, 0xf5, 0xa1, 0x47, 0x89, 0xc2, 0x7f, 0x24, 0x8a, 0xf6, 0x12, 0x3d, 0x83, 0x18,
	0x3b, 0x55, 0xd6, 0x22, 0x4e, 0xd9, 0x2a, 0x2a, 0xdc, 0xc0, 0x4f, 0x61, 0x61, 0x9a, 0x9b, 0x56,
	0xda, 0x7b, 0x8d, 0x62, 0x46, 0x92, 0x09, 0xe8, 0xfd, 0x6a, 0x6b, 0xbb, 0xf3, 0xaa, 0xd2, 0x62,
	0x4e, 0x19, 0xc7, 0x39, 0xfb, 0xc5, 0x20, 0x76, 0x1d, 0xd6, 0x10, 0x77, 0x88, 0xda, 0x08, 0x96,
	0x86, 0xab, 0x64, 0x7d, 0x9a, 0x4f, 0xb7, 0x9e, 0x13, 0x23, 0xef, 0xfb, 0x1a, 0xba, 0xd5, 0xc2,
	0x51, 0xf9, 0x1b, 0x38, 0xbe, 0x6b, 0x8c, 0x69, 0xda, 0x9b, 0x1e, 0x6e, 0xc6, 0xaa, 0x7b, 0xe8,
	0xf2, 0x2b, 0xc0, 0x24, 0xe6, 0x27, 0x10, 0xde, 0xe2, 0x76, 0x58, 0x57, 0xff, 0xc9, 0xdf, 0x42,
	0xfc, 0x20, 0x37, 0xf7, 0x48, 0x17, 0x90, 0xac, 0x9f, 0xfb, 0xbf, 0x3d, 0x6e, 0xb9, 0x70, 0x9c,
	0xb3, 0xe0, 0x23, 0xcb, 0xce, 0xe1, 0xa9, 0xc3, 0xd0, 0x74, 0xaa, 0x35, 0xc8, 0xdf, 0xc3, 0x1c,
	0x87, 0x08, 0x2e, 0xff, 0x0b, 0xdf, 0x63, 0x7a, 0x8f, 0xc5, 0x8e, 0x96, 0xfd, 0x64, 0x90, 0x7c,
	0x57, 0xfa, 0xee, 0x4a, 0x6e, 0x37, 0x4a, 0x56, 0xfc, 0x13, 0xcc, 0x48, 0xb1, 0x33, 0x78, 0xed,
	0x1b, 0x78, 0x44, 0x67, 0x36, 0xec, 0x61, 0x90, 0x2c, 0x8b, 0xe1, 0xdd, 0x1e, 0x36, 0x8c, 0x5c,
	0xc3, 0x77, 0x8f, 0x1b, 0xbe, 0x3c, 0x48, 0xe7, 0x5e, 0xbc, 0xd7, 0xf1, 0x7a, 0x46, 0xff, 0xbe,
	0x0f, 0x7f, 0x02, 0x00, 0x00, 0xff, 0xff, 0xc1, 0x2b, 0xaa, 0xf4, 0x91, 0x03, 0x00, 0x00,
}
