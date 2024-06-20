// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: hub-genesis/state.proto

package types

import (
	fmt "fmt"
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
	// accounts on the Hub to fund with some bootstrapping transfers
	GenesisAccounts   []GenesisAccount `protobuf:"bytes,3,rep,name=genesis_accounts,json=genesisAccounts,proto3" json:"genesis_accounts"`
	UnackedSeqNums    []uint64         `protobuf:"varint,4,rep,packed,name=unacked_seq_nums,json=unackedSeqNums,proto3" json:"unacked_seq_nums,omitempty"`
	NumUnackedSeqNums uint64           `protobuf:"varint,5,opt,name=num_unacked_seq_nums,json=numUnackedSeqNums,proto3" json:"num_unacked_seq_nums,omitempty"`
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

func (m *State) GetGenesisAccounts() []GenesisAccount {
	if m != nil {
		return m.GenesisAccounts
	}
	return nil
}

func (m *State) GetUnackedSeqNums() []uint64 {
	if m != nil {
		return m.UnackedSeqNums
	}
	return nil
}

func (m *State) GetNumUnackedSeqNums() uint64 {
	if m != nil {
		return m.NumUnackedSeqNums
	}
	return 0
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

func init() {
	proto.RegisterType((*State)(nil), "rollapp.hub_genesis.State")
	proto.RegisterType((*GenesisAccount)(nil), "rollapp.hub_genesis.GenesisAccount")
}

func init() { proto.RegisterFile("hub-genesis/state.proto", fileDescriptor_1ad65c2fe0d953ab) }

var fileDescriptor_1ad65c2fe0d953ab = []byte{
	// 358 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x51, 0x31, 0x6f, 0xe2, 0x30,
	0x14, 0x8e, 0x49, 0xe0, 0x38, 0x23, 0x71, 0xb9, 0x1c, 0xd2, 0xe5, 0x18, 0x72, 0x11, 0x5d, 0xb2,
	0x60, 0x0b, 0x3a, 0x54, 0x1d, 0x4b, 0x87, 0x4a, 0x0c, 0x1d, 0xa0, 0x5d, 0xba, 0x44, 0x4e, 0x62,
	0x85, 0x08, 0x6c, 0x87, 0x38, 0xae, 0xa0, 0xbf, 0xa2, 0x3f, 0x0b, 0x75, 0x62, 0xec, 0x54, 0x55,
	0xf0, 0x47, 0xaa, 0x90, 0x54, 0x2d, 0x6a, 0x27, 0xbf, 0xe7, 0xef, 0xfb, 0xde, 0x7b, 0xfa, 0x3e,
	0xf8, 0x77, 0xa6, 0x82, 0x7e, 0x4c, 0x39, 0x95, 0x89, 0xc4, 0x32, 0x27, 0x39, 0x45, 0x69, 0x26,
	0x72, 0x61, 0xfd, 0xc9, 0xc4, 0x62, 0x41, 0xd2, 0x14, 0xcd, 0x54, 0xe0, 0x57, 0x84, 0x6e, 0x27,
	0x16, 0xb1, 0x38, 0xe0, 0xb8, 0xa8, 0x4a, 0x6a, 0xd7, 0x09, 0x85, 0x64, 0x42, 0xe2, 0x80, 0x48,
	0x8a, 0xef, 0x07, 0x01, 0xcd, 0xc9, 0x00, 0x87, 0x22, 0xe1, 0x25, 0xde, 0x7b, 0x02, 0xb0, 0x3e,
	0x2d, 0x46, 0x5b, 0x37, 0xd0, 0xac, 0x46, 0xf9, 0x24, 0x0c, 0x85, 0xe2, 0xb9, 0xb4, 0x75, 0x57,
	0xf7, 0x5a, 0xc3, 0x13, 0xf4, 0xcd, 0x3e, 0x74, 0x55, 0xbe, 0x17, 0x25, 0x77, 0x64, 0x6c, 0x5e,
	0xfe, 0x6b, 0x93, 0x5f, 0xf1, 0xd1, 0xaf, 0xb4, 0x10, 0x34, 0x15, 0x27, 0xe1, 0x9c, 0x46, 0xbe,
	0xa4, 0x4b, 0x9f, 0x2b, 0x26, 0x6d, 0xc3, 0xd5, 0x3d, 0xa3, 0x12, 0xb4, 0x2b, 0x74, 0x4a, 0x97,
	0xd7, 0x8a, 0x49, 0x0b, 0xc3, 0x0e, 0x57, 0xcc, 0xff, 0xa2, 0xa9, 0xbb, 0xc0, 0x33, 0x26, 0xbf,
	0xb9, 0x62, 0xb7, 0x47, 0x82, 0xb1, 0xd1, 0x04, 0x66, 0x6d, 0x6c, 0x34, 0x6b, 0xa6, 0xde, 0x0b,
	0x61, 0xfb, 0xf8, 0x2a, 0xeb, 0x0c, 0x36, 0x08, 0x2b, 0x2a, 0x1b, 0xb8, 0xc0, 0x6b, 0x0d, 0xff,
	0xa1, 0xd2, 0x0f, 0x54, 0xf8, 0x81, 0x2a, 0x3f, 0xd0, 0xa5, 0x48, 0x78, 0x75, 0x4f, 0x45, 0xb7,
	0x6c, 0xf8, 0x83, 0x44, 0x51, 0x46, 0xa5, 0xb4, 0x6b, 0x2e, 0xf0, 0x7e, 0x4e, 0xde, 0xdb, 0xd1,
	0x74, 0xb3, 0x73, 0xc0, 0x76, 0xe7, 0x80, 0xd7, 0x9d, 0x03, 0x1e, 0xf7, 0x8e, 0xb6, 0xdd, 0x3b,
	0xda, 0xf3, 0xde, 0xd1, 0xee, 0xce, 0xe3, 0x24, 0x9f, 0xa9, 0x00, 0x85, 0x82, 0xe1, 0x68, 0xcd,
	0x28, 0x97, 0x89, 0xe0, 0xab, 0xf5, 0xc3, 0x47, 0xd3, 0xcf, 0xa2, 0x39, 0x5e, 0xe1, 0xcf, 0xb9,
	0xe6, 0xeb, 0x94, 0xca, 0xa0, 0x71, 0x48, 0xe3, 0xf4, 0x2d, 0x00, 0x00, 0xff, 0xff, 0xc8, 0x34,
	0x48, 0x78, 0xf3, 0x01, 0x00, 0x00,
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
	if m.NumUnackedSeqNums != 0 {
		i = encodeVarintState(dAtA, i, uint64(m.NumUnackedSeqNums))
		i--
		dAtA[i] = 0x28
	}
	if len(m.UnackedSeqNums) > 0 {
		dAtA2 := make([]byte, len(m.UnackedSeqNums)*10)
		var j1 int
		for _, num := range m.UnackedSeqNums {
			for num >= 1<<7 {
				dAtA2[j1] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j1++
			}
			dAtA2[j1] = uint8(num)
			j1++
		}
		i -= j1
		copy(dAtA[i:], dAtA2[:j1])
		i = encodeVarintState(dAtA, i, uint64(j1))
		i--
		dAtA[i] = 0x22
	}
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
			dAtA[i] = 0x1a
		}
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
	if len(m.GenesisAccounts) > 0 {
		for _, e := range m.GenesisAccounts {
			l = e.Size()
			n += 1 + l + sovState(uint64(l))
		}
	}
	if len(m.UnackedSeqNums) > 0 {
		l = 0
		for _, e := range m.UnackedSeqNums {
			l += sovState(uint64(e))
		}
		n += 1 + sovState(uint64(l)) + l
	}
	if m.NumUnackedSeqNums != 0 {
		n += 1 + sovState(uint64(m.NumUnackedSeqNums))
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
		case 3:
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
		case 4:
			if wireType == 0 {
				var v uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowState
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.UnackedSeqNums = append(m.UnackedSeqNums, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowState
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthState
				}
				postIndex := iNdEx + packedLen
				if postIndex < 0 {
					return ErrInvalidLengthState
				}
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				var count int
				for _, integer := range dAtA[iNdEx:postIndex] {
					if integer < 128 {
						count++
					}
				}
				elementCount = count
				if elementCount != 0 && len(m.UnackedSeqNums) == 0 {
					m.UnackedSeqNums = make([]uint64, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowState
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.UnackedSeqNums = append(m.UnackedSeqNums, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field UnackedSeqNums", wireType)
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field NumUnackedSeqNums", wireType)
			}
			m.NumUnackedSeqNums = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowState
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.NumUnackedSeqNums |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
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
