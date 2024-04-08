// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: hub-genesis/locked.proto

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

// locked holds the state of the genesis event
type Locked struct {
	// locked is the state of the genesis event
	Locked bool `protobuf:"varint,1,opt,name=locked,proto3" json:"locked,omitempty"`
	// tokens is the list of tokens that are expected to be locked on genesis
	Tokens github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,2,rep,name=tokens,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"tokens"`
}

func (m *Locked) Reset()         { *m = Locked{} }
func (m *Locked) String() string { return proto.CompactTextString(m) }
func (*Locked) ProtoMessage()    {}
func (*Locked) Descriptor() ([]byte, []int) {
	return fileDescriptor_d673da6e390f4b86, []int{0}
}
func (m *Locked) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Locked) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Locked.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Locked) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Locked.Merge(m, src)
}
func (m *Locked) XXX_Size() int {
	return m.Size()
}
func (m *Locked) XXX_DiscardUnknown() {
	xxx_messageInfo_Locked.DiscardUnknown(m)
}

var xxx_messageInfo_Locked proto.InternalMessageInfo

func (m *Locked) GetLocked() bool {
	if m != nil {
		return m.Locked
	}
	return false
}

func (m *Locked) GetTokens() github_com_cosmos_cosmos_sdk_types.Coins {
	if m != nil {
		return m.Tokens
	}
	return nil
}

func init() {
	proto.RegisterType((*Locked)(nil), "rollapp.hub_genesis.Locked")
}

func init() { proto.RegisterFile("hub-genesis/locked.proto", fileDescriptor_d673da6e390f4b86) }

var fileDescriptor_d673da6e390f4b86 = []byte{
	// 265 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x90, 0x31, 0x4f, 0xc3, 0x30,
	0x10, 0x85, 0x63, 0x90, 0x22, 0x14, 0xb6, 0x80, 0x50, 0xe8, 0xe0, 0x56, 0x4c, 0x59, 0x62, 0x53,
	0x98, 0x58, 0xcb, 0xca, 0x54, 0x36, 0x16, 0x14, 0x27, 0x56, 0x12, 0x25, 0xf1, 0x45, 0xb9, 0x04,
	0x35, 0xec, 0xec, 0xfc, 0x0e, 0x7e, 0x49, 0xc7, 0x8e, 0x4c, 0x80, 0x92, 0x3f, 0x82, 0x6a, 0x5b,
	0x6a, 0x27, 0xfb, 0xe9, 0x3d, 0x7f, 0xcf, 0x77, 0x5e, 0x90, 0xf7, 0x22, 0xca, 0xa4, 0x92, 0x58,
	0x20, 0xaf, 0x20, 0x29, 0x65, 0xca, 0x9a, 0x16, 0x3a, 0xf0, 0x2f, 0x5a, 0xa8, 0xaa, 0xb8, 0x69,
	0x58, 0xde, 0x8b, 0x57, 0x9b, 0x98, 0x5d, 0x66, 0x90, 0x81, 0xf6, 0xf9, 0xfe, 0x66, 0xa2, 0x33,
	0x9a, 0x00, 0xd6, 0x80, 0x5c, 0xc4, 0x28, 0xf9, 0xdb, 0x52, 0xc8, 0x2e, 0x5e, 0xf2, 0x04, 0x0a,
	0x65, 0xfc, 0x9b, 0x0f, 0xe2, 0xb9, 0x4f, 0x9a, 0xed, 0x5f, 0x79, 0xae, 0x69, 0x09, 0xc8, 0x82,
	0x84, 0x67, 0x6b, 0xab, 0xfc, 0xc4, 0x73, 0x3b, 0x28, 0xa5, 0xc2, 0xe0, 0x64, 0x71, 0x1a, 0x9e,
	0xdf, 0x5d, 0x33, 0xc3, 0x64, 0x7b, 0x26, 0xb3, 0x4c, 0xf6, 0x08, 0x85, 0x5a, 0xdd, 0x6e, 0x7f,
	0xe6, 0xce, 0xd7, 0xef, 0x3c, 0xcc, 0x8a, 0x2e, 0xef, 0x05, 0x4b, 0xa0, 0xe6, 0xf6, 0x03, 0xe6,
	0x88, 0x30, 0x2d, 0x79, 0x37, 0x34, 0x12, 0xf5, 0x03, 0x5c, 0x5b, 0xf4, 0xea, 0x79, 0x3b, 0x52,
	0xb2, 0x1b, 0x29, 0xf9, 0x1b, 0x29, 0xf9, 0x9c, 0xa8, 0xb3, 0x9b, 0xa8, 0xf3, 0x3d, 0x51, 0xe7,
	0xe5, 0xe1, 0x88, 0x95, 0x0e, 0xb5, 0x54, 0x58, 0x80, 0xda, 0x0c, 0xef, 0x07, 0x11, 0xb5, 0x69,
	0xc9, 0x37, 0xfc, 0x78, 0x5d, 0xba, 0x42, 0xb8, 0x7a, 0xc6, 0xfb, 0xff, 0x00, 0x00, 0x00, 0xff,
	0xff, 0xa1, 0xa6, 0x01, 0xf7, 0x4a, 0x01, 0x00, 0x00,
}

func (m *Locked) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Locked) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Locked) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Tokens) > 0 {
		for iNdEx := len(m.Tokens) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Tokens[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintLocked(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if m.Locked {
		i--
		if m.Locked {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintLocked(dAtA []byte, offset int, v uint64) int {
	offset -= sovLocked(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Locked) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Locked {
		n += 2
	}
	if len(m.Tokens) > 0 {
		for _, e := range m.Tokens {
			l = e.Size()
			n += 1 + l + sovLocked(uint64(l))
		}
	}
	return n
}

func sovLocked(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozLocked(x uint64) (n int) {
	return sovLocked(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Locked) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowLocked
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
			return fmt.Errorf("proto: Locked: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Locked: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Locked", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLocked
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
			m.Locked = bool(v != 0)
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Tokens", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLocked
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
				return ErrInvalidLengthLocked
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthLocked
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Tokens = append(m.Tokens, types.Coin{})
			if err := m.Tokens[len(m.Tokens)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipLocked(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthLocked
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
func skipLocked(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowLocked
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
					return 0, ErrIntOverflowLocked
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
					return 0, ErrIntOverflowLocked
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
				return 0, ErrInvalidLengthLocked
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupLocked
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthLocked
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthLocked        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowLocked          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupLocked = fmt.Errorf("proto: unexpected end of group")
)
