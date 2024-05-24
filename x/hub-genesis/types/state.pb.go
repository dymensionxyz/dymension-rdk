// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: hub-genesis/state.proto

package types

import (
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// State holds the state of the genesis event
type State struct {
	// were the initial genesis tokens escrowed successfully?
	GenesisTokensWereEscrowed bool `protobuf:"varint,1,opt,name=genesis_tokens_were_escrowed,json=genesisTokensWereEscrowed,proto3" json:"genesis_tokens_were_escrowed,omitempty"` // Deprecated: Do not use.
	// genesis_tokens is the list of tokens that are expected to be locked on genesis event
	GenesisTokens   github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,2,rep,name=genesis_tokens,json=genesisTokens,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"genesis_tokens"` // Deprecated: Do not use.
	GenesisAccounts []GenesisAccount                         `protobuf:"bytes,9,rep,name=genesis_accounts,json=genesisAccounts,proto3" json:"genesis_accounts"`
	Metadatas       []TokenMetadata                          `protobuf:"bytes,8,rep,name=metadatas,proto3" json:"metadatas"`
}

func (m *State) Reset()         { *m = State{} }
func (m *State) String() string { return proto.CompactTextString(m) }
func (*State) ProtoMessage()    {}
func (*State) Descriptor() ([]byte, []int) {
	return fileDescriptor_1ad65c2fe0d953ab, []int{0}
}
func (m *State) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *State) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_State.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *State) XXX_Merge(src proto.Message) {
	xxx_messageInfo_State.Merge(m, src)
}
func (m *State) XXX_Size() int {
	return m.Size()
}
func (m *State) XXX_DiscardUnknown() {
	xxx_messageInfo_State.DiscardUnknown(m)
}

var xxx_messageInfo_State proto.InternalMessageInfo

// Deprecated: Do not use.
func (m *State) GetGenesisTokensWereEscrowed() bool {
	if m != nil {
		return m.GenesisTokensWereEscrowed
	}
	return false
}

// Deprecated: Do not use.
func (m *State) GetGenesisTokens() github_com_cosmos_cosmos_sdk_types.Coins {
	if m != nil {
		return m.GenesisTokens
	}
	return nil
}

func (m *State) GetGenesisAccounts() []GenesisAccount {
	if m != nil {
		return m.GenesisAccounts
	}
	return nil
}

func (m *State) GetMetadatas() []TokenMetadata {
	if m != nil {
		return m.Metadatas
	}
	return nil
}

// GenesisAccount is a struct for the genesis account for the rollapp
type GenesisAccount struct {
	// amount of coins to be sent to the genesis address
	Amount types.Coin `protobuf:"bytes,1,opt,name=amount,proto3" json:"amount"`
	// address is a bech-32 address of the genesis account
	Address string `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
}

func (m *GenesisAccount) Reset()         { *m = GenesisAccount{} }
func (m *GenesisAccount) String() string { return proto.CompactTextString(m) }
func (*GenesisAccount) ProtoMessage()    {}
func (*GenesisAccount) Descriptor() ([]byte, []int) {
	return fileDescriptor_1ad65c2fe0d953ab, []int{1}
}
func (m *GenesisAccount) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisAccount) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisAccount.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisAccount) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisAccount.Merge(m, src)
}
func (m *GenesisAccount) XXX_Size() int {
	return m.Size()
}
func (m *GenesisAccount) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisAccount.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisAccount proto.InternalMessageInfo

func (m *GenesisAccount) GetAmount() types.Coin {
	if m != nil {
		return m.Amount
	}
	return types.Coin{}
}

func (m *GenesisAccount) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

// DenomUnit represents a struct that describes a given
// denomination unit of the basic token.
type DenomUnit struct {
	// denom represents the string name of the given denom unit (e.g uatom).
	Denom string `protobuf:"bytes,1,opt,name=denom,proto3" json:"denom,omitempty"`
	// exponent represents power of 10 exponent that one must
	// raise the base_denom to in order to equal the given DenomUnit's denom
	// 1 denom = 10^exponent base_denom
	// (e.g. with a base_denom of uatom, one can create a DenomUnit of 'atom' with
	// exponent = 6, thus: 1 atom = 10^6 uatom).
	Exponent uint32 `protobuf:"varint,2,opt,name=exponent,proto3" json:"exponent,omitempty"`
	// aliases is a list of string aliases for the given denom
	Aliases []string `protobuf:"bytes,3,rep,name=aliases,proto3" json:"aliases,omitempty"`
}

func (m *DenomUnit) Reset()         { *m = DenomUnit{} }
func (m *DenomUnit) String() string { return proto.CompactTextString(m) }
func (*DenomUnit) ProtoMessage()    {}
func (*DenomUnit) Descriptor() ([]byte, []int) {
	return fileDescriptor_1ad65c2fe0d953ab, []int{2}
}
func (m *DenomUnit) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *DenomUnit) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_DenomUnit.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *DenomUnit) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DenomUnit.Merge(m, src)
}
func (m *DenomUnit) XXX_Size() int {
	return m.Size()
}
func (m *DenomUnit) XXX_DiscardUnknown() {
	xxx_messageInfo_DenomUnit.DiscardUnknown(m)
}

var xxx_messageInfo_DenomUnit proto.InternalMessageInfo

func (m *DenomUnit) GetDenom() string {
	if m != nil {
		return m.Denom
	}
	return ""
}

func (m *DenomUnit) GetExponent() uint32 {
	if m != nil {
		return m.Exponent
	}
	return 0
}

func (m *DenomUnit) GetAliases() []string {
	if m != nil {
		return m.Aliases
	}
	return nil
}

// Metadata represents a struct that describes
// a basic token.
type TokenMetadata struct {
	Description string `protobuf:"bytes,1,opt,name=description,proto3" json:"description,omitempty"`
	// denom_units represents the list of DenomUnit's for a given coin
	DenomUnits []*DenomUnit `protobuf:"bytes,2,rep,name=denom_units,json=denomUnits,proto3" json:"denom_units,omitempty"`
	// base represents the base denom (should be the DenomUnit with exponent = 0).
	Base string `protobuf:"bytes,3,opt,name=base,proto3" json:"base,omitempty"`
	// display indicates the suggested denom that should be
	// displayed in clients.
	Display string `protobuf:"bytes,4,opt,name=display,proto3" json:"display,omitempty"`
	// name defines the name of the token (eg: Cosmos Atom)
	//
	// Since: cosmos-sdk 0.43
	Name string `protobuf:"bytes,5,opt,name=name,proto3" json:"name,omitempty"`
	// symbol is the token symbol usually shown on exchanges (eg: ATOM). This can
	// be the same as the display.
	//
	// Since: cosmos-sdk 0.43
	Symbol string `protobuf:"bytes,6,opt,name=symbol,proto3" json:"symbol,omitempty"`
	// URI to a document (on or off-chain) that contains additional information. Optional.
	//
	// Since: cosmos-sdk 0.46
	URI string `protobuf:"bytes,7,opt,name=uri,proto3" json:"uri,omitempty"`
	// URIHash is a sha256 hash of a document pointed by URI. It's used to verify that
	// the document didn't change. Optional.
	//
	// Since: cosmos-sdk 0.46
	URIHash string `protobuf:"bytes,8,opt,name=uri_hash,json=uriHash,proto3" json:"uri_hash,omitempty"`
}

func (m *TokenMetadata) Reset()         { *m = TokenMetadata{} }
func (m *TokenMetadata) String() string { return proto.CompactTextString(m) }
func (*TokenMetadata) ProtoMessage()    {}
func (*TokenMetadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_1ad65c2fe0d953ab, []int{3}
}
func (m *TokenMetadata) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TokenMetadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TokenMetadata.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TokenMetadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TokenMetadata.Merge(m, src)
}
func (m *TokenMetadata) XXX_Size() int {
	return m.Size()
}
func (m *TokenMetadata) XXX_DiscardUnknown() {
	xxx_messageInfo_TokenMetadata.DiscardUnknown(m)
}

var xxx_messageInfo_TokenMetadata proto.InternalMessageInfo

func (m *TokenMetadata) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *TokenMetadata) GetDenomUnits() []*DenomUnit {
	if m != nil {
		return m.DenomUnits
	}
	return nil
}

func (m *TokenMetadata) GetBase() string {
	if m != nil {
		return m.Base
	}
	return ""
}

func (m *TokenMetadata) GetDisplay() string {
	if m != nil {
		return m.Display
	}
	return ""
}

func (m *TokenMetadata) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *TokenMetadata) GetSymbol() string {
	if m != nil {
		return m.Symbol
	}
	return ""
}

func (m *TokenMetadata) GetURI() string {
	if m != nil {
		return m.URI
	}
	return ""
}

func (m *TokenMetadata) GetURIHash() string {
	if m != nil {
		return m.URIHash
	}
	return ""
}

func init() {
	proto.RegisterType((*State)(nil), "rollapp.hub_genesis.State")
	proto.RegisterType((*GenesisAccount)(nil), "rollapp.hub_genesis.GenesisAccount")
	proto.RegisterType((*DenomUnit)(nil), "rollapp.hub_genesis.DenomUnit")
	proto.RegisterType((*TokenMetadata)(nil), "rollapp.hub_genesis.TokenMetadata")
}

func init() { proto.RegisterFile("hub-genesis/state.proto", fileDescriptor_1ad65c2fe0d953ab) }

var fileDescriptor_1ad65c2fe0d953ab = []byte{
	// 583 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x53, 0xbd, 0x6e, 0xdb, 0x3c,
	0x14, 0xb5, 0xac, 0xc4, 0x3f, 0x0c, 0x9c, 0xef, 0x03, 0x1b, 0xb4, 0xb4, 0x51, 0xc8, 0x86, 0x0b,
	0x14, 0x5e, 0x2c, 0x21, 0xe9, 0x50, 0x74, 0x2a, 0xea, 0xf4, 0x2f, 0x43, 0x17, 0x25, 0x46, 0x80,
	0x2e, 0x02, 0x25, 0x11, 0x16, 0x61, 0x8b, 0x14, 0x44, 0xaa, 0xb1, 0xfb, 0x12, 0xed, 0x73, 0xf4,
	0x49, 0x32, 0x66, 0xec, 0xe4, 0x06, 0xf6, 0x8b, 0x14, 0xa4, 0xe8, 0xc4, 0x06, 0x8c, 0x4e, 0xbc,
	0xe7, 0xf2, 0xdc, 0x73, 0x0f, 0x2f, 0x2e, 0xc1, 0xb3, 0xa4, 0x08, 0x87, 0x13, 0xc2, 0x88, 0xa0,
	0xc2, 0x13, 0x12, 0x4b, 0xe2, 0x66, 0x39, 0x97, 0x1c, 0x3e, 0xc9, 0xf9, 0x6c, 0x86, 0xb3, 0xcc,
	0x4d, 0x8a, 0x30, 0x30, 0x84, 0xce, 0xc9, 0x84, 0x4f, 0xb8, 0xbe, 0xf7, 0x54, 0x54, 0x52, 0x3b,
	0x4e, 0xc4, 0x45, 0xca, 0x85, 0x17, 0x62, 0x41, 0xbc, 0x6f, 0xa7, 0x21, 0x91, 0xf8, 0xd4, 0x8b,
	0x38, 0x65, 0xe5, 0x7d, 0xff, 0xbe, 0x0a, 0x0e, 0x2f, 0x95, 0x34, 0x3c, 0x07, 0xcf, 0x8d, 0x54,
	0x20, 0xf9, 0x94, 0x30, 0x11, 0xdc, 0x90, 0x9c, 0x04, 0x44, 0x44, 0x39, 0xbf, 0x21, 0x31, 0xb2,
	0x7a, 0xd6, 0xa0, 0x31, 0xaa, 0x22, 0xcb, 0x6f, 0x1b, 0xde, 0x95, 0xa6, 0x5d, 0x93, 0x9c, 0x7c,
	0x30, 0x24, 0x28, 0xc1, 0xf1, 0xae, 0x08, 0xaa, 0xf6, 0xec, 0xc1, 0xd1, 0x59, 0xdb, 0x2d, 0x7d,
	0xb8, 0xca, 0x87, 0x6b, 0x7c, 0xb8, 0xe7, 0x9c, 0xb2, 0xd1, 0xd9, 0xed, 0xb2, 0x5b, 0xf9, 0xf5,
	0xa7, 0x3b, 0x98, 0x50, 0x99, 0x14, 0xa1, 0x1b, 0xf1, 0xd4, 0x33, 0xa6, 0xcb, 0x63, 0x28, 0xe2,
	0xa9, 0x27, 0x17, 0x19, 0x11, 0xba, 0x40, 0x20, 0xcb, 0x6f, 0xed, 0x38, 0x80, 0x57, 0xe0, 0xff,
	0x4d, 0x57, 0x1c, 0x45, 0xbc, 0x60, 0x52, 0xa0, 0xa6, 0xee, 0xfb, 0xc2, 0xdd, 0x33, 0x2a, 0xf7,
	0x53, 0x79, 0xbe, 0x2b, 0xb9, 0xa3, 0x03, 0xe5, 0xc0, 0xff, 0x6f, 0xb2, 0x93, 0x15, 0xf0, 0x23,
	0x68, 0xa6, 0x44, 0xe2, 0x18, 0x4b, 0x2c, 0x50, 0x43, 0xcb, 0xf5, 0xf7, 0xca, 0x69, 0x17, 0x5f,
	0x0c, 0xd5, 0xa8, 0x3d, 0x96, 0xf6, 0x23, 0x70, 0xbc, 0xdb, 0x10, 0xbe, 0x06, 0x35, 0x9c, 0xaa,
	0x48, 0x0f, 0xf5, 0x9f, 0xd3, 0x29, 0xd5, 0x0c, 0x1d, 0x22, 0x50, 0xc7, 0x71, 0x9c, 0x13, 0xa1,
	0xe6, 0x6a, 0x0d, 0x9a, 0xfe, 0x06, 0xf6, 0xaf, 0x41, 0xf3, 0x3d, 0x61, 0x3c, 0x1d, 0x33, 0x2a,
	0xe1, 0x09, 0x38, 0x8c, 0x15, 0xd0, 0xf2, 0x4d, 0xbf, 0x04, 0xb0, 0x03, 0x1a, 0x64, 0x9e, 0x71,
	0x46, 0x98, 0xd4, 0xd5, 0x2d, 0xff, 0x01, 0x6b, 0xe1, 0x19, 0xc5, 0x82, 0x08, 0x64, 0xf7, 0x6c,
	0x2d, 0x5c, 0xc2, 0xfe, 0x8f, 0x2a, 0x68, 0xed, 0x3c, 0x10, 0xf6, 0xc0, 0x51, 0xac, 0xb6, 0x82,
	0x66, 0x92, 0x72, 0x66, 0x7a, 0x6c, 0xa7, 0xe0, 0x5b, 0xc5, 0x60, 0x3c, 0x0d, 0x0a, 0x46, 0xe5,
	0x66, 0x05, 0x9c, 0xbd, 0xb3, 0x7b, 0x30, 0xed, 0x83, 0x78, 0x13, 0x0a, 0x08, 0xc1, 0x81, 0x1a,
	0x05, 0xb2, 0xb5, 0xb6, 0x8e, 0x95, 0xc5, 0x98, 0x8a, 0x6c, 0x86, 0x17, 0xe8, 0xa0, 0x7c, 0xbb,
	0x81, 0x8a, 0xcd, 0x70, 0x4a, 0xd0, 0x61, 0xc9, 0x56, 0x31, 0x7c, 0x0a, 0x6a, 0x62, 0x91, 0x86,
	0x7c, 0x86, 0x6a, 0x3a, 0x6b, 0x10, 0x6c, 0x03, 0xbb, 0xc8, 0x29, 0xaa, 0xab, 0xe4, 0xa8, 0xbe,
	0x5a, 0x76, 0xed, 0xb1, 0x7f, 0xe1, 0xab, 0x1c, 0x7c, 0x09, 0x1a, 0x45, 0x4e, 0x83, 0x04, 0x8b,
	0x04, 0x35, 0xf4, 0xfd, 0xd1, 0x6a, 0xd9, 0xad, 0x8f, 0xfd, 0x8b, 0xcf, 0x58, 0x24, 0x7e, 0xbd,
	0xc8, 0xa9, 0x0a, 0x46, 0x97, 0xb7, 0x2b, 0xc7, 0xba, 0x5b, 0x39, 0xd6, 0xfd, 0xca, 0xb1, 0x7e,
	0xae, 0x9d, 0xca, 0xdd, 0xda, 0xa9, 0xfc, 0x5e, 0x3b, 0x95, 0xaf, 0x6f, 0xb6, 0x56, 0x38, 0x5e,
	0xa4, 0x84, 0x09, 0xca, 0xd9, 0x7c, 0xf1, 0xfd, 0x11, 0x0c, 0xf3, 0x78, 0xea, 0xcd, 0xbd, 0xed,
	0x8f, 0xad, 0x37, 0x3b, 0xac, 0xe9, 0xef, 0xf8, 0xea, 0x6f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x13,
	0x93, 0x98, 0x5a, 0xf4, 0x03, 0x00, 0x00,
}

func (m *State) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *State) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *State) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.GenesisAccounts) > 0 {
		for iNdEx := len(m.GenesisAccounts) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.GenesisAccounts[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintState(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x4a
		}
	}
	if len(m.Metadatas) > 0 {
		for iNdEx := len(m.Metadatas) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Metadatas[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintState(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x42
		}
	}
	if len(m.GenesisTokens) > 0 {
		for iNdEx := len(m.GenesisTokens) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.GenesisTokens[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintState(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if m.GenesisTokensWereEscrowed {
		i--
		if m.GenesisTokensWereEscrowed {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *GenesisAccount) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisAccount) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisAccount) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintState(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0x12
	}
	{
		size, err := m.Amount.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintState(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *DenomUnit) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DenomUnit) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *DenomUnit) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Aliases) > 0 {
		for iNdEx := len(m.Aliases) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Aliases[iNdEx])
			copy(dAtA[i:], m.Aliases[iNdEx])
			i = encodeVarintState(dAtA, i, uint64(len(m.Aliases[iNdEx])))
			i--
			dAtA[i] = 0x1a
		}
	}
	if m.Exponent != 0 {
		i = encodeVarintState(dAtA, i, uint64(m.Exponent))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Denom) > 0 {
		i -= len(m.Denom)
		copy(dAtA[i:], m.Denom)
		i = encodeVarintState(dAtA, i, uint64(len(m.Denom)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *TokenMetadata) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TokenMetadata) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TokenMetadata) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.URIHash) > 0 {
		i -= len(m.URIHash)
		copy(dAtA[i:], m.URIHash)
		i = encodeVarintState(dAtA, i, uint64(len(m.URIHash)))
		i--
		dAtA[i] = 0x42
	}
	if len(m.URI) > 0 {
		i -= len(m.URI)
		copy(dAtA[i:], m.URI)
		i = encodeVarintState(dAtA, i, uint64(len(m.URI)))
		i--
		dAtA[i] = 0x3a
	}
	if len(m.Symbol) > 0 {
		i -= len(m.Symbol)
		copy(dAtA[i:], m.Symbol)
		i = encodeVarintState(dAtA, i, uint64(len(m.Symbol)))
		i--
		dAtA[i] = 0x32
	}
	if len(m.Name) > 0 {
		i -= len(m.Name)
		copy(dAtA[i:], m.Name)
		i = encodeVarintState(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.Display) > 0 {
		i -= len(m.Display)
		copy(dAtA[i:], m.Display)
		i = encodeVarintState(dAtA, i, uint64(len(m.Display)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Base) > 0 {
		i -= len(m.Base)
		copy(dAtA[i:], m.Base)
		i = encodeVarintState(dAtA, i, uint64(len(m.Base)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.DenomUnits) > 0 {
		for iNdEx := len(m.DenomUnits) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.DenomUnits[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintState(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.Description) > 0 {
		i -= len(m.Description)
		copy(dAtA[i:], m.Description)
		i = encodeVarintState(dAtA, i, uint64(len(m.Description)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintState(dAtA []byte, offset int, v uint64) int {
	offset -= sovState(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *State) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.GenesisTokensWereEscrowed {
		n += 2
	}
	if len(m.GenesisTokens) > 0 {
		for _, e := range m.GenesisTokens {
			l = e.Size()
			n += 1 + l + sovState(uint64(l))
		}
	}
	if len(m.Metadatas) > 0 {
		for _, e := range m.Metadatas {
			l = e.Size()
			n += 1 + l + sovState(uint64(l))
		}
	}
	if len(m.GenesisAccounts) > 0 {
		for _, e := range m.GenesisAccounts {
			l = e.Size()
			n += 1 + l + sovState(uint64(l))
		}
	}
	return n
}

func (m *GenesisAccount) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Amount.Size()
	n += 1 + l + sovState(uint64(l))
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovState(uint64(l))
	}
	return n
}

func (m *DenomUnit) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Denom)
	if l > 0 {
		n += 1 + l + sovState(uint64(l))
	}
	if m.Exponent != 0 {
		n += 1 + sovState(uint64(m.Exponent))
	}
	if len(m.Aliases) > 0 {
		for _, s := range m.Aliases {
			l = len(s)
			n += 1 + l + sovState(uint64(l))
		}
	}
	return n
}

func (m *TokenMetadata) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Description)
	if l > 0 {
		n += 1 + l + sovState(uint64(l))
	}
	if len(m.DenomUnits) > 0 {
		for _, e := range m.DenomUnits {
			l = e.Size()
			n += 1 + l + sovState(uint64(l))
		}
	}
	l = len(m.Base)
	if l > 0 {
		n += 1 + l + sovState(uint64(l))
	}
	l = len(m.Display)
	if l > 0 {
		n += 1 + l + sovState(uint64(l))
	}
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovState(uint64(l))
	}
	l = len(m.Symbol)
	if l > 0 {
		n += 1 + l + sovState(uint64(l))
	}
	l = len(m.URI)
	if l > 0 {
		n += 1 + l + sovState(uint64(l))
	}
	l = len(m.URIHash)
	if l > 0 {
		n += 1 + l + sovState(uint64(l))
	}
	return n
}

func sovState(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozState(x uint64) (n int) {
	return sovState(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *State) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowState
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: State: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: State: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field GenesisTokensWereEscrowed", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowState
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.GenesisTokensWereEscrowed = bool(v != 0)
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GenesisTokens", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowState
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthState
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthState
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.GenesisTokens = append(m.GenesisTokens, types.Coin{})
			if err := m.GenesisTokens[len(m.GenesisTokens)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Metadatas", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowState
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthState
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthState
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Metadatas = append(m.Metadatas, TokenMetadata{})
			if err := m.Metadatas[len(m.Metadatas)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GenesisAccounts", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowState
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthState
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthState
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.GenesisAccounts = append(m.GenesisAccounts, GenesisAccount{})
			if err := m.GenesisAccounts[len(m.GenesisAccounts)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipState(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthState
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *GenesisAccount) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowState
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: GenesisAccount: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisAccount: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowState
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthState
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthState
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Amount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowState
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthState
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthState
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipState(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthState
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *DenomUnit) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowState
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: DenomUnit: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DenomUnit: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Denom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowState
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthState
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthState
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Denom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Exponent", wireType)
			}
			m.Exponent = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowState
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Exponent |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Aliases", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowState
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthState
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthState
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Aliases = append(m.Aliases, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipState(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthState
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *TokenMetadata) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowState
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: TokenMetadata: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TokenMetadata: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Description", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowState
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthState
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthState
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Description = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DenomUnits", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowState
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthState
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthState
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DenomUnits = append(m.DenomUnits, &DenomUnit{})
			if err := m.DenomUnits[len(m.DenomUnits)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Base", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowState
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthState
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthState
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Base = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Display", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowState
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthState
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthState
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Display = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowState
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthState
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthState
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Symbol", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowState
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthState
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthState
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Symbol = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field URI", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowState
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthState
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthState
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.URI = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field URIHash", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowState
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthState
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthState
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.URIHash = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipState(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthState
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipState(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowState
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowState
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowState
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthState
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupState
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthState
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthState        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowState          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupState = fmt.Errorf("proto: unexpected end of group")
)
