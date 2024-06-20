// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: hub/hub.proto

package types

import (
	fmt "fmt"
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

type RegisteredDenom_Status int32

const (
	RegisteredDenom_PENDING  RegisteredDenom_Status = 0
	RegisteredDenom_ACTIVE   RegisteredDenom_Status = 1
	RegisteredDenom_INACTIVE RegisteredDenom_Status = 2
)

var RegisteredDenom_Status_name = map[int32]string{
	0: "PENDING",
	1: "ACTIVE",
	2: "INACTIVE",
}

var RegisteredDenom_Status_value = map[string]int32{
	"PENDING":  0,
	"ACTIVE":   1,
	"INACTIVE": 2,
}

func (x RegisteredDenom_Status) String() string {
	return proto.EnumName(RegisteredDenom_Status_name, int32(x))
}

func (RegisteredDenom_Status) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_87629b1556de20e1, []int{1, 0}
}

// Hub is a proto message that represents the metadata of the Hub
type Hub struct {
	// id is the unique identifier of the Hub
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// channel_id is the unique identifier of the channel that the Hub is connected to
	ChannelId string `protobuf:"bytes,2,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
	// registeredDenoms is a list of registered denoms on the Hub
	RegisteredDenoms []*RegisteredDenom `protobuf:"bytes,3,rep,name=registered_denoms,json=registeredDenoms,proto3" json:"registered_denoms,omitempty"`
}

func (m *Hub) Reset()         { *m = Hub{} }
func (m *Hub) String() string { return proto.CompactTextString(m) }
func (*Hub) ProtoMessage()    {}
func (*Hub) Descriptor() ([]byte, []int) {
	return fileDescriptor_87629b1556de20e1, []int{0}
}
func (m *Hub) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Hub) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Hub.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Hub) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Hub.Merge(m, src)
}
func (m *Hub) XXX_Size() int {
	return m.Size()
}
func (m *Hub) XXX_DiscardUnknown() {
	xxx_messageInfo_Hub.DiscardUnknown(m)
}

var xxx_messageInfo_Hub proto.InternalMessageInfo

func (m *Hub) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Hub) GetChannelId() string {
	if m != nil {
		return m.ChannelId
	}
	return ""
}

func (m *Hub) GetRegisteredDenoms() []*RegisteredDenom {
	if m != nil {
		return m.RegisteredDenoms
	}
	return nil
}

type RegisteredDenom struct {
	// base is the base of the denom
	Base string `protobuf:"bytes,1,opt,name=base,proto3" json:"base,omitempty"`
	// status is the status of the denom registration in the Hub
	Status RegisteredDenom_Status `protobuf:"varint,2,opt,name=status,proto3,enum=rollapp.hub.RegisteredDenom_Status" json:"status,omitempty"`
}

func (m *RegisteredDenom) Reset()         { *m = RegisteredDenom{} }
func (m *RegisteredDenom) String() string { return proto.CompactTextString(m) }
func (*RegisteredDenom) ProtoMessage()    {}
func (*RegisteredDenom) Descriptor() ([]byte, []int) {
	return fileDescriptor_87629b1556de20e1, []int{1}
}
func (m *RegisteredDenom) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *RegisteredDenom) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_RegisteredDenom.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *RegisteredDenom) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegisteredDenom.Merge(m, src)
}
func (m *RegisteredDenom) XXX_Size() int {
	return m.Size()
}
func (m *RegisteredDenom) XXX_DiscardUnknown() {
	xxx_messageInfo_RegisteredDenom.DiscardUnknown(m)
}

var xxx_messageInfo_RegisteredDenom proto.InternalMessageInfo

func (m *RegisteredDenom) GetBase() string {
	if m != nil {
		return m.Base
	}
	return ""
}

func (m *RegisteredDenom) GetStatus() RegisteredDenom_Status {
	if m != nil {
		return m.Status
	}
	return RegisteredDenom_PENDING
}

func init() {
	proto.RegisterEnum("rollapp.hub.RegisteredDenom_Status", RegisteredDenom_Status_name, RegisteredDenom_Status_value)
	proto.RegisterType((*Hub)(nil), "rollapp.hub.Hub")
	proto.RegisterType((*RegisteredDenom)(nil), "rollapp.hub.RegisteredDenom")
}

func init() { proto.RegisterFile("hub/hub.proto", fileDescriptor_87629b1556de20e1) }

var fileDescriptor_87629b1556de20e1 = []byte{
	// 311 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xcd, 0x28, 0x4d, 0xd2,
	0xcf, 0x28, 0x4d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x2e, 0xca, 0xcf, 0xc9, 0x49,
	0x2c, 0x28, 0xd0, 0xcb, 0x28, 0x4d, 0x92, 0x12, 0x49, 0xcf, 0x4f, 0xcf, 0x07, 0x8b, 0xeb, 0x83,
	0x58, 0x10, 0x25, 0x4a, 0xf5, 0x5c, 0xcc, 0x1e, 0xa5, 0x49, 0x42, 0x7c, 0x5c, 0x4c, 0x99, 0x29,
	0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x9c, 0x41, 0x4c, 0x99, 0x29, 0x42, 0xb2, 0x5c, 0x5c, 0xc9, 0x19,
	0x89, 0x79, 0x79, 0xa9, 0x39, 0xf1, 0x99, 0x29, 0x12, 0x4c, 0x60, 0x71, 0x4e, 0xa8, 0x88, 0x67,
	0x8a, 0x90, 0x27, 0x97, 0x60, 0x51, 0x6a, 0x7a, 0x66, 0x71, 0x49, 0x6a, 0x51, 0x6a, 0x4a, 0x7c,
	0x4a, 0x6a, 0x5e, 0x7e, 0x6e, 0xb1, 0x04, 0xb3, 0x02, 0xb3, 0x06, 0xb7, 0x91, 0x8c, 0x1e, 0x92,
	0xa5, 0x7a, 0x41, 0x70, 0x55, 0x2e, 0x20, 0x45, 0x41, 0x02, 0x45, 0xa8, 0x02, 0xc5, 0x4a, 0x93,
	0x19, 0xb9, 0xf8, 0xd1, 0x54, 0x09, 0x09, 0x71, 0xb1, 0x24, 0x25, 0x16, 0xa7, 0x42, 0xdd, 0x03,
	0x66, 0x0b, 0x59, 0x73, 0xb1, 0x15, 0x97, 0x24, 0x96, 0x94, 0x16, 0x83, 0x5d, 0xc3, 0x67, 0xa4,
	0x8c, 0xcf, 0x1e, 0xbd, 0x60, 0xb0, 0xd2, 0x20, 0xa8, 0x16, 0x25, 0x7d, 0x2e, 0x36, 0x88, 0x88,
	0x10, 0x37, 0x17, 0x7b, 0x80, 0xab, 0x9f, 0x8b, 0xa7, 0x9f, 0xbb, 0x00, 0x83, 0x10, 0x17, 0x17,
	0x9b, 0xa3, 0x73, 0x88, 0x67, 0x98, 0xab, 0x00, 0xa3, 0x10, 0x0f, 0x17, 0x87, 0xa7, 0x1f, 0x94,
	0xc7, 0xe4, 0xe4, 0x7d, 0xe2, 0x91, 0x1c, 0xe3, 0x85, 0x47, 0x72, 0x8c, 0x0f, 0x1e, 0xc9, 0x31,
	0x4e, 0x78, 0x2c, 0xc7, 0x70, 0xe1, 0xb1, 0x1c, 0xc3, 0x8d, 0xc7, 0x72, 0x0c, 0x51, 0x86, 0xe9,
	0x99, 0x25, 0x20, 0x1b, 0x93, 0xf3, 0x73, 0xf5, 0x53, 0x2a, 0x73, 0x53, 0xf3, 0x8a, 0x33, 0xf3,
	0xf3, 0x2a, 0x2a, 0xab, 0x10, 0x1c, 0xdd, 0xa2, 0x94, 0x6c, 0xfd, 0x0a, 0x50, 0x34, 0xe8, 0x97,
	0x54, 0x16, 0xa4, 0x16, 0x27, 0xb1, 0x81, 0x83, 0xda, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0x56,
	0x90, 0xf3, 0x3a, 0x9e, 0x01, 0x00, 0x00,
}

func (m *Hub) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Hub) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Hub) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.RegisteredDenoms) > 0 {
		for iNdEx := len(m.RegisteredDenoms) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.RegisteredDenoms[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintHub(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.ChannelId) > 0 {
		i -= len(m.ChannelId)
		copy(dAtA[i:], m.ChannelId)
		i = encodeVarintHub(dAtA, i, uint64(len(m.ChannelId)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Id) > 0 {
		i -= len(m.Id)
		copy(dAtA[i:], m.Id)
		i = encodeVarintHub(dAtA, i, uint64(len(m.Id)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *RegisteredDenom) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RegisteredDenom) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *RegisteredDenom) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Status != 0 {
		i = encodeVarintHub(dAtA, i, uint64(m.Status))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Base) > 0 {
		i -= len(m.Base)
		copy(dAtA[i:], m.Base)
		i = encodeVarintHub(dAtA, i, uint64(len(m.Base)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintHub(dAtA []byte, offset int, v uint64) int {
	offset -= sovHub(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Hub) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Id)
	if l > 0 {
		n += 1 + l + sovHub(uint64(l))
	}
	l = len(m.ChannelId)
	if l > 0 {
		n += 1 + l + sovHub(uint64(l))
	}
	if len(m.RegisteredDenoms) > 0 {
		for _, e := range m.RegisteredDenoms {
			l = e.Size()
			n += 1 + l + sovHub(uint64(l))
		}
	}
	return n
}

func (m *RegisteredDenom) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Base)
	if l > 0 {
		n += 1 + l + sovHub(uint64(l))
	}
	if m.Status != 0 {
		n += 1 + sovHub(uint64(m.Status))
	}
	return n
}

func sovHub(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozHub(x uint64) (n int) {
	return sovHub(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Hub) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowHub
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
			return fmt.Errorf("proto: Hub: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Hub: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHub
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
				return ErrInvalidLengthHub
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthHub
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Id = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChannelId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHub
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
				return ErrInvalidLengthHub
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthHub
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ChannelId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RegisteredDenoms", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHub
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
				return ErrInvalidLengthHub
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthHub
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.RegisteredDenoms = append(m.RegisteredDenoms, &RegisteredDenom{})
			if err := m.RegisteredDenoms[len(m.RegisteredDenoms)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipHub(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthHub
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
func (m *RegisteredDenom) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowHub
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
			return fmt.Errorf("proto: RegisteredDenom: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RegisteredDenom: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Base", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHub
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
				return ErrInvalidLengthHub
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthHub
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Base = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Status", wireType)
			}
			m.Status = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHub
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Status |= RegisteredDenom_Status(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipHub(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthHub
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
func skipHub(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowHub
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
					return 0, ErrIntOverflowHub
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
					return 0, ErrIntOverflowHub
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
				return 0, ErrInvalidLengthHub
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupHub
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthHub
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthHub        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowHub          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupHub = fmt.Errorf("proto: unexpected end of group")
)