// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sequencers/gov_permission.proto

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

type GrantPermissionsProposal struct {
	Title              string               `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	Description        string               `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	AddressPermissions []AddressPermissions `protobuf:"bytes,3,rep,name=address_permissions,json=addressPermissions,proto3" json:"address_permissions"`
}

func (m *GrantPermissionsProposal) Reset()      { *m = GrantPermissionsProposal{} }
func (*GrantPermissionsProposal) ProtoMessage() {}
func (*GrantPermissionsProposal) Descriptor() ([]byte, []int) {
	return fileDescriptor_d40c11f5a74a5529, []int{0}
}
func (m *GrantPermissionsProposal) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GrantPermissionsProposal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GrantPermissionsProposal.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GrantPermissionsProposal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GrantPermissionsProposal.Merge(m, src)
}
func (m *GrantPermissionsProposal) XXX_Size() int {
	return m.Size()
}
func (m *GrantPermissionsProposal) XXX_DiscardUnknown() {
	xxx_messageInfo_GrantPermissionsProposal.DiscardUnknown(m)
}

var xxx_messageInfo_GrantPermissionsProposal proto.InternalMessageInfo

type RevokePermissionsProposal struct {
	Title              string               `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	Description        string               `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	AddressPermissions []AddressPermissions `protobuf:"bytes,3,rep,name=address_permissions,json=addressPermissions,proto3" json:"address_permissions"`
}

func (m *RevokePermissionsProposal) Reset()      { *m = RevokePermissionsProposal{} }
func (*RevokePermissionsProposal) ProtoMessage() {}
func (*RevokePermissionsProposal) Descriptor() ([]byte, []int) {
	return fileDescriptor_d40c11f5a74a5529, []int{1}
}
func (m *RevokePermissionsProposal) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *RevokePermissionsProposal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_RevokePermissionsProposal.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *RevokePermissionsProposal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RevokePermissionsProposal.Merge(m, src)
}
func (m *RevokePermissionsProposal) XXX_Size() int {
	return m.Size()
}
func (m *RevokePermissionsProposal) XXX_DiscardUnknown() {
	xxx_messageInfo_RevokePermissionsProposal.DiscardUnknown(m)
}

var xxx_messageInfo_RevokePermissionsProposal proto.InternalMessageInfo

func init() {
	proto.RegisterType((*GrantPermissionsProposal)(nil), "rollapp.sequencers.types.GrantPermissionsProposal")
	proto.RegisterType((*RevokePermissionsProposal)(nil), "rollapp.sequencers.types.RevokePermissionsProposal")
}

func init() { proto.RegisterFile("sequencers/gov_permission.proto", fileDescriptor_d40c11f5a74a5529) }

var fileDescriptor_d40c11f5a74a5529 = []byte{
	// 297 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x2f, 0x4e, 0x2d, 0x2c,
	0x4d, 0xcd, 0x4b, 0x4e, 0x2d, 0x2a, 0xd6, 0x4f, 0xcf, 0x2f, 0x8b, 0x2f, 0x48, 0x2d, 0xca, 0xcd,
	0x2c, 0x2e, 0xce, 0xcc, 0xcf, 0xd3, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x92, 0x28, 0xca, 0xcf,
	0xc9, 0x49, 0x2c, 0x28, 0xd0, 0x43, 0x28, 0xd4, 0x2b, 0xa9, 0x2c, 0x48, 0x2d, 0x96, 0x12, 0x49,
	0xcf, 0x4f, 0xcf, 0x07, 0x2b, 0xd2, 0x07, 0xb1, 0x20, 0xea, 0xa5, 0xa4, 0x91, 0x0c, 0x44, 0x37,
	0x4c, 0xe9, 0x28, 0x23, 0x97, 0x84, 0x7b, 0x51, 0x62, 0x5e, 0x49, 0x00, 0x5c, 0xa6, 0x38, 0xa0,
	0x28, 0xbf, 0x20, 0xbf, 0x38, 0x31, 0x47, 0x48, 0x84, 0x8b, 0xb5, 0x24, 0xb3, 0x24, 0x27, 0x55,
	0x82, 0x51, 0x81, 0x51, 0x83, 0x33, 0x08, 0xc2, 0x11, 0x52, 0xe0, 0xe2, 0x4e, 0x49, 0x2d, 0x4e,
	0x2e, 0xca, 0x2c, 0x28, 0xc9, 0xcc, 0xcf, 0x93, 0x60, 0x02, 0xcb, 0x21, 0x0b, 0x09, 0x25, 0x73,
	0x09, 0x27, 0xa6, 0xa4, 0x14, 0xa5, 0x16, 0x17, 0x23, 0xb9, 0xbe, 0x58, 0x82, 0x59, 0x81, 0x59,
	0x83, 0xdb, 0x48, 0x47, 0x0f, 0x97, 0xfb, 0xf5, 0x1c, 0x21, 0x9a, 0x90, 0x9c, 0xe2, 0xc4, 0x72,
	0xe2, 0x9e, 0x3c, 0x43, 0x90, 0x50, 0x22, 0x86, 0x8c, 0x15, 0x4f, 0xc7, 0x02, 0x79, 0x86, 0x19,
	0x0b, 0xe4, 0x19, 0x5e, 0x2c, 0x90, 0x67, 0x54, 0x3a, 0xc6, 0xc8, 0x25, 0x19, 0x94, 0x5a, 0x96,
	0x9f, 0x9d, 0x3a, 0xb4, 0x3d, 0xe2, 0x14, 0x74, 0xe2, 0x91, 0x1c, 0xe3, 0x85, 0x47, 0x72, 0x8c,
	0x0f, 0x1e, 0xc9, 0x31, 0x4e, 0x78, 0x2c, 0xc7, 0x70, 0xe1, 0xb1, 0x1c, 0xc3, 0x8d, 0xc7, 0x72,
	0x0c, 0x51, 0x16, 0xe9, 0x99, 0x25, 0x19, 0xa5, 0x49, 0x7a, 0xc9, 0xf9, 0xb9, 0xfa, 0x29, 0x95,
	0xb9, 0xa9, 0x79, 0x20, 0xed, 0x15, 0x95, 0x55, 0x08, 0x8e, 0x6e, 0x51, 0x4a, 0xb6, 0x7e, 0x85,
	0x3e, 0x52, 0x7c, 0x83, 0x9d, 0x93, 0xc4, 0x06, 0x8e, 0x6b, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff,
	0xff, 0x87, 0x43, 0xb4, 0x54, 0x5b, 0x02, 0x00, 0x00,
}

func (this *GrantPermissionsProposal) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*GrantPermissionsProposal)
	if !ok {
		that2, ok := that.(GrantPermissionsProposal)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Title != that1.Title {
		return false
	}
	if this.Description != that1.Description {
		return false
	}
	if len(this.AddressPermissions) != len(that1.AddressPermissions) {
		return false
	}
	for i := range this.AddressPermissions {
		if !this.AddressPermissions[i].Equal(&that1.AddressPermissions[i]) {
			return false
		}
	}
	return true
}
func (this *RevokePermissionsProposal) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*RevokePermissionsProposal)
	if !ok {
		that2, ok := that.(RevokePermissionsProposal)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Title != that1.Title {
		return false
	}
	if this.Description != that1.Description {
		return false
	}
	if len(this.AddressPermissions) != len(that1.AddressPermissions) {
		return false
	}
	for i := range this.AddressPermissions {
		if !this.AddressPermissions[i].Equal(&that1.AddressPermissions[i]) {
			return false
		}
	}
	return true
}
func (m *GrantPermissionsProposal) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GrantPermissionsProposal) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GrantPermissionsProposal) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.AddressPermissions) > 0 {
		for iNdEx := len(m.AddressPermissions) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.AddressPermissions[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGovPermission(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.Description) > 0 {
		i -= len(m.Description)
		copy(dAtA[i:], m.Description)
		i = encodeVarintGovPermission(dAtA, i, uint64(len(m.Description)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Title) > 0 {
		i -= len(m.Title)
		copy(dAtA[i:], m.Title)
		i = encodeVarintGovPermission(dAtA, i, uint64(len(m.Title)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *RevokePermissionsProposal) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RevokePermissionsProposal) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *RevokePermissionsProposal) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.AddressPermissions) > 0 {
		for iNdEx := len(m.AddressPermissions) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.AddressPermissions[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGovPermission(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.Description) > 0 {
		i -= len(m.Description)
		copy(dAtA[i:], m.Description)
		i = encodeVarintGovPermission(dAtA, i, uint64(len(m.Description)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Title) > 0 {
		i -= len(m.Title)
		copy(dAtA[i:], m.Title)
		i = encodeVarintGovPermission(dAtA, i, uint64(len(m.Title)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintGovPermission(dAtA []byte, offset int, v uint64) int {
	offset -= sovGovPermission(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GrantPermissionsProposal) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Title)
	if l > 0 {
		n += 1 + l + sovGovPermission(uint64(l))
	}
	l = len(m.Description)
	if l > 0 {
		n += 1 + l + sovGovPermission(uint64(l))
	}
	if len(m.AddressPermissions) > 0 {
		for _, e := range m.AddressPermissions {
			l = e.Size()
			n += 1 + l + sovGovPermission(uint64(l))
		}
	}
	return n
}

func (m *RevokePermissionsProposal) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Title)
	if l > 0 {
		n += 1 + l + sovGovPermission(uint64(l))
	}
	l = len(m.Description)
	if l > 0 {
		n += 1 + l + sovGovPermission(uint64(l))
	}
	if len(m.AddressPermissions) > 0 {
		for _, e := range m.AddressPermissions {
			l = e.Size()
			n += 1 + l + sovGovPermission(uint64(l))
		}
	}
	return n
}

func sovGovPermission(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGovPermission(x uint64) (n int) {
	return sovGovPermission(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GrantPermissionsProposal) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGovPermission
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
			return fmt.Errorf("proto: GrantPermissionsProposal: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GrantPermissionsProposal: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Title", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGovPermission
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
				return ErrInvalidLengthGovPermission
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGovPermission
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Title = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Description", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGovPermission
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
				return ErrInvalidLengthGovPermission
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGovPermission
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Description = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AddressPermissions", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGovPermission
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
				return ErrInvalidLengthGovPermission
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGovPermission
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AddressPermissions = append(m.AddressPermissions, AddressPermissions{})
			if err := m.AddressPermissions[len(m.AddressPermissions)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGovPermission(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGovPermission
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
func (m *RevokePermissionsProposal) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGovPermission
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
			return fmt.Errorf("proto: RevokePermissionsProposal: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RevokePermissionsProposal: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Title", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGovPermission
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
				return ErrInvalidLengthGovPermission
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGovPermission
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Title = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Description", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGovPermission
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
				return ErrInvalidLengthGovPermission
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGovPermission
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Description = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AddressPermissions", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGovPermission
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
				return ErrInvalidLengthGovPermission
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGovPermission
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AddressPermissions = append(m.AddressPermissions, AddressPermissions{})
			if err := m.AddressPermissions[len(m.AddressPermissions)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGovPermission(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGovPermission
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
func skipGovPermission(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGovPermission
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
					return 0, ErrIntOverflowGovPermission
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
					return 0, ErrIntOverflowGovPermission
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
				return 0, ErrInvalidLengthGovPermission
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGovPermission
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGovPermission
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGovPermission        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGovPermission          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGovPermission = fmt.Errorf("proto: unexpected end of group")
)
