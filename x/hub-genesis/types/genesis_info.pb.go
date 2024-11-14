// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: hub-genesis/genesis_info.proto

package types

import (
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/cosmos-sdk/x/bank/types"
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

// The genesis info of the rollapp, that is passed to the hub for validation.
// it's populated on the InitGenesis of the rollapp
type GenesisInfo struct {
	// checksum used to verify integrity of the genesis file. currently unused
	GenesisChecksum string `protobuf:"bytes,1,opt,name=genesis_checksum,json=genesisChecksum,proto3" json:"genesis_checksum,omitempty"`
	// unique bech32 prefix
	Bech32Prefix string `protobuf:"bytes,2,opt,name=bech32_prefix,json=bech32Prefix,proto3" json:"bech32_prefix,omitempty"`
	// native_denom is the base denom for the native token
	NativeDenom *DenomMetadata `protobuf:"bytes,3,opt,name=native_denom,json=nativeDenom,proto3" json:"native_denom,omitempty"`
	// initial_supply is the initial supply of the native token
	InitialSupply github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,4,opt,name=initial_supply,json=initialSupply,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"initial_supply"`
	// accounts on the Hub to fund with some bootstrapping transfers
	GenesisAccounts []GenesisAccount `protobuf:"bytes,5,rep,name=genesis_accounts,json=genesisAccounts,proto3" json:"genesis_accounts"`
}

func (m *GenesisInfo) Reset()         { *m = GenesisInfo{} }
func (m *GenesisInfo) String() string { return proto.CompactTextString(m) }
func (*GenesisInfo) ProtoMessage()    {}
func (*GenesisInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_db3d5bf55e08315f, []int{0}
}
func (m *GenesisInfo) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisInfo.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisInfo.Merge(m, src)
}
func (m *GenesisInfo) XXX_Size() int {
	return m.Size()
}
func (m *GenesisInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisInfo.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisInfo proto.InternalMessageInfo

func (m *GenesisInfo) GetGenesisChecksum() string {
	if m != nil {
		return m.GenesisChecksum
	}
	return ""
}

func (m *GenesisInfo) GetBech32Prefix() string {
	if m != nil {
		return m.Bech32Prefix
	}
	return ""
}

func (m *GenesisInfo) GetNativeDenom() *DenomMetadata {
	if m != nil {
		return m.NativeDenom
	}
	return nil
}

func (m *GenesisInfo) GetGenesisAccounts() []GenesisAccount {
	if m != nil {
		return m.GenesisAccounts
	}
	return nil
}

type DenomMetadata struct {
	Display  string `protobuf:"bytes,1,opt,name=display,proto3" json:"display,omitempty"`
	Base     string `protobuf:"bytes,2,opt,name=base,proto3" json:"base,omitempty"`
	Exponent uint32 `protobuf:"varint,3,opt,name=exponent,proto3" json:"exponent,omitempty"`
}

func (m *DenomMetadata) Reset()         { *m = DenomMetadata{} }
func (m *DenomMetadata) String() string { return proto.CompactTextString(m) }
func (*DenomMetadata) ProtoMessage()    {}
func (*DenomMetadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_db3d5bf55e08315f, []int{1}
}
func (m *DenomMetadata) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *DenomMetadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_DenomMetadata.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *DenomMetadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DenomMetadata.Merge(m, src)
}
func (m *DenomMetadata) XXX_Size() int {
	return m.Size()
}
func (m *DenomMetadata) XXX_DiscardUnknown() {
	xxx_messageInfo_DenomMetadata.DiscardUnknown(m)
}

var xxx_messageInfo_DenomMetadata proto.InternalMessageInfo

func (m *DenomMetadata) GetDisplay() string {
	if m != nil {
		return m.Display
	}
	return ""
}

func (m *DenomMetadata) GetBase() string {
	if m != nil {
		return m.Base
	}
	return ""
}

func (m *DenomMetadata) GetExponent() uint32 {
	if m != nil {
		return m.Exponent
	}
	return 0
}

// GenesisAccount is a struct for the genesis account for the rollapp
type GenesisAccount struct {
	// amount of coins to be sent to the genesis address
	Amount github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,1,opt,name=amount,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"amount"`
	// address is a bech-32 address of the genesis account
	Address string `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
}

func (m *GenesisAccount) Reset()         { *m = GenesisAccount{} }
func (m *GenesisAccount) String() string { return proto.CompactTextString(m) }
func (*GenesisAccount) ProtoMessage()    {}
func (*GenesisAccount) Descriptor() ([]byte, []int) {
	return fileDescriptor_db3d5bf55e08315f, []int{2}
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

func (m *GenesisAccount) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func init() {
	proto.RegisterType((*GenesisInfo)(nil), "rollapp.hub_genesis.GenesisInfo")
	proto.RegisterType((*DenomMetadata)(nil), "rollapp.hub_genesis.DenomMetadata")
	proto.RegisterType((*GenesisAccount)(nil), "rollapp.hub_genesis.GenesisAccount")
}

func init() { proto.RegisterFile("hub-genesis/genesis_info.proto", fileDescriptor_db3d5bf55e08315f) }

var fileDescriptor_db3d5bf55e08315f = []byte{
	// 456 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x52, 0x4d, 0x6f, 0xd3, 0x40,
	0x10, 0x8d, 0xdb, 0x50, 0x60, 0xd3, 0x14, 0xb4, 0x20, 0x61, 0xe5, 0xe0, 0x46, 0xa9, 0x84, 0xc2,
	0x21, 0xb6, 0x9a, 0x9e, 0x38, 0x12, 0xbe, 0xd4, 0x03, 0x12, 0x4a, 0xe1, 0x00, 0x17, 0x6b, 0x6d,
	0x4f, 0xe2, 0x55, 0xec, 0xdd, 0x95, 0x67, 0x5d, 0xc5, 0xfc, 0x0a, 0x7e, 0x56, 0x8f, 0xe5, 0x86,
	0x38, 0x54, 0x28, 0xf9, 0x23, 0xc8, 0xeb, 0x4d, 0x9b, 0x4a, 0xb9, 0x70, 0xda, 0x7d, 0x6f, 0x66,
	0xe7, 0xbd, 0x79, 0x5a, 0xe2, 0xa5, 0x65, 0x34, 0x9a, 0x83, 0x00, 0xe4, 0x18, 0xd8, 0x33, 0xe4,
	0x62, 0x26, 0x7d, 0x55, 0x48, 0x2d, 0xe9, 0xb3, 0x42, 0x66, 0x19, 0x53, 0xca, 0x4f, 0xcb, 0x28,
	0xb4, 0xf5, 0xde, 0xf3, 0xb9, 0x9c, 0x4b, 0x53, 0x0f, 0xea, 0x5b, 0xd3, 0xda, 0xf3, 0x62, 0x89,
	0xb9, 0xc4, 0x20, 0x62, 0x62, 0x11, 0x5c, 0x9e, 0x46, 0xa0, 0xd9, 0xa9, 0x01, 0xb6, 0xfe, 0x62,
	0x5b, 0x0a, 0x35, 0xd3, 0xd0, 0x14, 0x06, 0xbf, 0xf6, 0x48, 0xe7, 0x63, 0xc3, 0x9f, 0x8b, 0x99,
	0xa4, 0xaf, 0xc8, 0xd3, 0x8d, 0x93, 0x38, 0x85, 0x78, 0x81, 0x65, 0xee, 0x3a, 0x7d, 0x67, 0xf8,
	0x78, 0xfa, 0xc4, 0xf2, 0x6f, 0x2d, 0x4d, 0x4f, 0x48, 0x37, 0x82, 0x38, 0x3d, 0x1b, 0x87, 0xaa,
	0x80, 0x19, 0x5f, 0xba, 0x7b, 0xa6, 0xef, 0xb0, 0x21, 0x3f, 0x1b, 0x8e, 0xbe, 0x27, 0x87, 0x82,
	0x69, 0x7e, 0x09, 0x61, 0x02, 0x42, 0xe6, 0xee, 0x7e, 0xdf, 0x19, 0x76, 0xc6, 0x03, 0x7f, 0xc7,
	0x6a, 0xfe, 0xbb, 0xba, 0xe3, 0x13, 0x68, 0x96, 0x30, 0xcd, 0xa6, 0x9d, 0xe6, 0x9d, 0x21, 0xe9,
	0x57, 0x72, 0xc4, 0x05, 0xd7, 0x9c, 0x65, 0x21, 0x96, 0x4a, 0x65, 0x95, 0xdb, 0xae, 0xc5, 0x26,
	0xfe, 0xd5, 0xcd, 0x71, 0xeb, 0xcf, 0xcd, 0xf1, 0xcb, 0x39, 0xd7, 0x69, 0x19, 0xf9, 0xb1, 0xcc,
	0x03, 0x1b, 0x45, 0x73, 0x8c, 0x30, 0x59, 0x04, 0xba, 0x52, 0x80, 0xfe, 0xb9, 0xd0, 0xd3, 0xae,
	0x9d, 0x72, 0x61, 0x86, 0xd0, 0x2f, 0x77, 0xdb, 0xb2, 0x38, 0x96, 0xa5, 0xd0, 0xe8, 0x3e, 0xe8,
	0xef, 0x0f, 0x3b, 0xe3, 0x93, 0x9d, 0x0e, 0x6d, 0x52, 0x6f, 0x9a, 0xde, 0x49, 0xbb, 0x56, 0xbf,
	0x0d, 0xc6, 0xb2, 0x38, 0xf8, 0x46, 0xba, 0xf7, 0x56, 0xa1, 0x2e, 0x79, 0x98, 0x70, 0x54, 0x19,
	0xab, 0x6c, 0x96, 0x1b, 0x48, 0x29, 0x69, 0x47, 0x0c, 0xc1, 0x46, 0x67, 0xee, 0xb4, 0x47, 0x1e,
	0xc1, 0x52, 0x49, 0x01, 0x42, 0x9b, 0xb8, 0xba, 0xd3, 0x5b, 0x3c, 0x28, 0xc8, 0xd1, 0x7d, 0x0f,
	0xf4, 0x03, 0x39, 0x60, 0x79, 0x7d, 0x6b, 0x46, 0xff, 0x77, 0x22, 0xf6, 0x75, 0xed, 0x91, 0x25,
	0x49, 0x01, 0x88, 0xd6, 0xcc, 0x06, 0x4e, 0x2e, 0xae, 0x56, 0x9e, 0x73, 0xbd, 0xf2, 0x9c, 0xbf,
	0x2b, 0xcf, 0xf9, 0xb9, 0xf6, 0x5a, 0xd7, 0x6b, 0xaf, 0xf5, 0x7b, 0xed, 0xb5, 0xbe, 0xbf, 0xde,
	0xd2, 0x48, 0xaa, 0x1c, 0x04, 0x72, 0x29, 0x96, 0xd5, 0x8f, 0x3b, 0x30, 0x2a, 0x92, 0x45, 0xb0,
	0x0c, 0xb6, 0x7f, 0x9f, 0x91, 0x8e, 0x0e, 0xcc, 0xf7, 0x3b, 0xfb, 0x17, 0x00, 0x00, 0xff, 0xff,
	0x8f, 0x72, 0xef, 0x1e, 0x04, 0x03, 0x00, 0x00,
}

func (m *GenesisInfo) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisInfo) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisInfo) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
				i = encodeVarintGenesisInfo(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x2a
		}
	}
	{
		size := m.InitialSupply.Size()
		i -= size
		if _, err := m.InitialSupply.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGenesisInfo(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	if m.NativeDenom != nil {
		{
			size, err := m.NativeDenom.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintGenesisInfo(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Bech32Prefix) > 0 {
		i -= len(m.Bech32Prefix)
		copy(dAtA[i:], m.Bech32Prefix)
		i = encodeVarintGenesisInfo(dAtA, i, uint64(len(m.Bech32Prefix)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.GenesisChecksum) > 0 {
		i -= len(m.GenesisChecksum)
		copy(dAtA[i:], m.GenesisChecksum)
		i = encodeVarintGenesisInfo(dAtA, i, uint64(len(m.GenesisChecksum)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *DenomMetadata) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DenomMetadata) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *DenomMetadata) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Exponent != 0 {
		i = encodeVarintGenesisInfo(dAtA, i, uint64(m.Exponent))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Base) > 0 {
		i -= len(m.Base)
		copy(dAtA[i:], m.Base)
		i = encodeVarintGenesisInfo(dAtA, i, uint64(len(m.Base)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Display) > 0 {
		i -= len(m.Display)
		copy(dAtA[i:], m.Display)
		i = encodeVarintGenesisInfo(dAtA, i, uint64(len(m.Display)))
		i--
		dAtA[i] = 0xa
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
		i = encodeVarintGenesisInfo(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0x12
	}
	{
		size := m.Amount.Size()
		i -= size
		if _, err := m.Amount.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGenesisInfo(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintGenesisInfo(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenesisInfo(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GenesisInfo) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.GenesisChecksum)
	if l > 0 {
		n += 1 + l + sovGenesisInfo(uint64(l))
	}
	l = len(m.Bech32Prefix)
	if l > 0 {
		n += 1 + l + sovGenesisInfo(uint64(l))
	}
	if m.NativeDenom != nil {
		l = m.NativeDenom.Size()
		n += 1 + l + sovGenesisInfo(uint64(l))
	}
	l = m.InitialSupply.Size()
	n += 1 + l + sovGenesisInfo(uint64(l))
	if len(m.GenesisAccounts) > 0 {
		for _, e := range m.GenesisAccounts {
			l = e.Size()
			n += 1 + l + sovGenesisInfo(uint64(l))
		}
	}
	return n
}

func (m *DenomMetadata) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Display)
	if l > 0 {
		n += 1 + l + sovGenesisInfo(uint64(l))
	}
	l = len(m.Base)
	if l > 0 {
		n += 1 + l + sovGenesisInfo(uint64(l))
	}
	if m.Exponent != 0 {
		n += 1 + sovGenesisInfo(uint64(m.Exponent))
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
	n += 1 + l + sovGenesisInfo(uint64(l))
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovGenesisInfo(uint64(l))
	}
	return n
}

func sovGenesisInfo(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenesisInfo(x uint64) (n int) {
	return sovGenesisInfo(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GenesisInfo) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesisInfo
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
			return fmt.Errorf("proto: GenesisInfo: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisInfo: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GenesisChecksum", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesisInfo
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
				return ErrInvalidLengthGenesisInfo
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesisInfo
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.GenesisChecksum = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Bech32Prefix", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesisInfo
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
				return ErrInvalidLengthGenesisInfo
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesisInfo
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Bech32Prefix = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NativeDenom", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesisInfo
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
				return ErrInvalidLengthGenesisInfo
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesisInfo
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.NativeDenom == nil {
				m.NativeDenom = &DenomMetadata{}
			}
			if err := m.NativeDenom.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InitialSupply", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesisInfo
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
				return ErrInvalidLengthGenesisInfo
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesisInfo
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.InitialSupply.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GenesisAccounts", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesisInfo
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
				return ErrInvalidLengthGenesisInfo
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesisInfo
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
			skippy, err := skipGenesisInfo(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesisInfo
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
func (m *DenomMetadata) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesisInfo
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
			return fmt.Errorf("proto: DenomMetadata: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DenomMetadata: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Display", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesisInfo
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
				return ErrInvalidLengthGenesisInfo
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesisInfo
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Display = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Base", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesisInfo
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
				return ErrInvalidLengthGenesisInfo
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesisInfo
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Base = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Exponent", wireType)
			}
			m.Exponent = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesisInfo
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
		default:
			iNdEx = preIndex
			skippy, err := skipGenesisInfo(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesisInfo
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
				return ErrIntOverflowGenesisInfo
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
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesisInfo
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
				return ErrInvalidLengthGenesisInfo
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesisInfo
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
					return ErrIntOverflowGenesisInfo
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
				return ErrInvalidLengthGenesisInfo
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesisInfo
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesisInfo(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesisInfo
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
func skipGenesisInfo(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenesisInfo
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
					return 0, ErrIntOverflowGenesisInfo
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
					return 0, ErrIntOverflowGenesisInfo
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
				return 0, ErrInvalidLengthGenesisInfo
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenesisInfo
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenesisInfo
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenesisInfo        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenesisInfo          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenesisInfo = fmt.Errorf("proto: unexpected end of group")
)