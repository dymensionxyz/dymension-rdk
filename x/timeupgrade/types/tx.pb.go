// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: timeupgrade/tx.proto

package types

import (
	context "context"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/codec/types"
	_ "github.com/cosmos/cosmos-sdk/types/msgservice"
	types "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	_ "github.com/gogo/protobuf/gogoproto"
	grpc1 "github.com/gogo/protobuf/grpc"
	proto "github.com/gogo/protobuf/proto"
	types1 "github.com/gogo/protobuf/types"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type MsgSoftwareUpgrade struct {
	OriginalUpgrade *types.MsgSoftwareUpgrade `protobuf:"bytes,1,opt,name=original_upgrade,json=originalUpgrade,proto3" json:"original_upgrade,omitempty"`
	UpgradeTime     *types1.Timestamp         `protobuf:"bytes,2,opt,name=upgrade_time,json=upgradeTime,proto3" json:"upgrade_time,omitempty"`
}

func (m *MsgSoftwareUpgrade) Reset()         { *m = MsgSoftwareUpgrade{} }
func (m *MsgSoftwareUpgrade) String() string { return proto.CompactTextString(m) }
func (*MsgSoftwareUpgrade) ProtoMessage()    {}
func (*MsgSoftwareUpgrade) Descriptor() ([]byte, []int) {
	return fileDescriptor_309e066b5c0e7455, []int{0}
}
func (m *MsgSoftwareUpgrade) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgSoftwareUpgrade) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgSoftwareUpgrade.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgSoftwareUpgrade) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgSoftwareUpgrade.Merge(m, src)
}
func (m *MsgSoftwareUpgrade) XXX_Size() int {
	return m.Size()
}
func (m *MsgSoftwareUpgrade) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgSoftwareUpgrade.DiscardUnknown(m)
}

var xxx_messageInfo_MsgSoftwareUpgrade proto.InternalMessageInfo

func (m *MsgSoftwareUpgrade) GetOriginalUpgrade() *types.MsgSoftwareUpgrade {
	if m != nil {
		return m.OriginalUpgrade
	}
	return nil
}

func (m *MsgSoftwareUpgrade) GetUpgradeTime() *types1.Timestamp {
	if m != nil {
		return m.UpgradeTime
	}
	return nil
}

type MsgSoftwareUpgradeResponse struct {
}

func (m *MsgSoftwareUpgradeResponse) Reset()         { *m = MsgSoftwareUpgradeResponse{} }
func (m *MsgSoftwareUpgradeResponse) String() string { return proto.CompactTextString(m) }
func (*MsgSoftwareUpgradeResponse) ProtoMessage()    {}
func (*MsgSoftwareUpgradeResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_309e066b5c0e7455, []int{1}
}
func (m *MsgSoftwareUpgradeResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgSoftwareUpgradeResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgSoftwareUpgradeResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgSoftwareUpgradeResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgSoftwareUpgradeResponse.Merge(m, src)
}
func (m *MsgSoftwareUpgradeResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgSoftwareUpgradeResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgSoftwareUpgradeResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgSoftwareUpgradeResponse proto.InternalMessageInfo

type MsgCancelUpgrade struct {
}

func (m *MsgCancelUpgrade) Reset()         { *m = MsgCancelUpgrade{} }
func (m *MsgCancelUpgrade) String() string { return proto.CompactTextString(m) }
func (*MsgCancelUpgrade) ProtoMessage()    {}
func (*MsgCancelUpgrade) Descriptor() ([]byte, []int) {
	return fileDescriptor_309e066b5c0e7455, []int{2}
}
func (m *MsgCancelUpgrade) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgCancelUpgrade) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgCancelUpgrade.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgCancelUpgrade) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgCancelUpgrade.Merge(m, src)
}
func (m *MsgCancelUpgrade) XXX_Size() int {
	return m.Size()
}
func (m *MsgCancelUpgrade) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgCancelUpgrade.DiscardUnknown(m)
}

var xxx_messageInfo_MsgCancelUpgrade proto.InternalMessageInfo

type MsgCancelUpgradeResponse struct {
}

func (m *MsgCancelUpgradeResponse) Reset()         { *m = MsgCancelUpgradeResponse{} }
func (m *MsgCancelUpgradeResponse) String() string { return proto.CompactTextString(m) }
func (*MsgCancelUpgradeResponse) ProtoMessage()    {}
func (*MsgCancelUpgradeResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_309e066b5c0e7455, []int{3}
}
func (m *MsgCancelUpgradeResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgCancelUpgradeResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgCancelUpgradeResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgCancelUpgradeResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgCancelUpgradeResponse.Merge(m, src)
}
func (m *MsgCancelUpgradeResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgCancelUpgradeResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgCancelUpgradeResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgCancelUpgradeResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*MsgSoftwareUpgrade)(nil), "rollapp.timeupgrade.types.MsgSoftwareUpgrade")
	proto.RegisterType((*MsgSoftwareUpgradeResponse)(nil), "rollapp.timeupgrade.types.MsgSoftwareUpgradeResponse")
	proto.RegisterType((*MsgCancelUpgrade)(nil), "rollapp.timeupgrade.types.MsgCancelUpgrade")
	proto.RegisterType((*MsgCancelUpgradeResponse)(nil), "rollapp.timeupgrade.types.MsgCancelUpgradeResponse")
}

func init() { proto.RegisterFile("timeupgrade/tx.proto", fileDescriptor_309e066b5c0e7455) }

var fileDescriptor_309e066b5c0e7455 = []byte{
	// 375 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x92, 0xcf, 0x4e, 0xc2, 0x30,
	0x1c, 0xc7, 0x99, 0x26, 0x1e, 0x8a, 0x06, 0xb2, 0x90, 0x08, 0x8d, 0x99, 0x66, 0x27, 0xa3, 0xa1,
	0x0d, 0x10, 0x0f, 0x1e, 0xbc, 0xe8, 0x99, 0x0b, 0xc8, 0xc5, 0x0b, 0xe9, 0xa0, 0xd4, 0xc5, 0x6d,
	0xbf, 0xb9, 0x96, 0x3f, 0xf3, 0x29, 0x7c, 0x0e, 0x9f, 0xc4, 0x23, 0x47, 0x8f, 0x06, 0x9e, 0xc0,
	0x37, 0x30, 0xfb, 0x53, 0x15, 0x44, 0x13, 0x4e, 0x5b, 0xdb, 0xcf, 0xef, 0xfb, 0x27, 0x2d, 0xaa,
	0x28, 0xd7, 0xe7, 0xe3, 0x50, 0x44, 0x6c, 0xc8, 0xa9, 0x9a, 0x91, 0x30, 0x02, 0x05, 0x66, 0x2d,
	0x02, 0xcf, 0x63, 0x61, 0x48, 0x7e, 0x9c, 0x12, 0x15, 0x87, 0x5c, 0xe2, 0x8a, 0x00, 0x01, 0x29,
	0x45, 0x93, 0xbf, 0x6c, 0x00, 0x1f, 0x0f, 0x40, 0xfa, 0x20, 0xa9, 0x56, 0x9a, 0x34, 0x1c, 0xae,
	0x58, 0xe3, 0x4b, 0x11, 0x1f, 0xe6, 0x80, 0x2f, 0x05, 0x9d, 0x34, 0x92, 0x8f, 0x9e, 0x14, 0x00,
	0xc2, 0xe3, 0x34, 0x5d, 0x39, 0xe3, 0x11, 0x4d, 0x2c, 0xa5, 0x62, 0x7e, 0x98, 0x03, 0xb5, 0x75,
	0x80, 0x05, 0x71, 0x76, 0x64, 0xbf, 0x18, 0xc8, 0x6c, 0x4b, 0xd1, 0x85, 0x91, 0x9a, 0xb2, 0x88,
	0xf7, 0x32, 0x73, 0xb3, 0x87, 0xca, 0x10, 0xb9, 0xc2, 0x0d, 0x98, 0xd7, 0xcf, 0x03, 0x55, 0x8d,
	0x13, 0xe3, 0xb4, 0xd8, 0x3c, 0x23, 0x59, 0x0c, 0xa2, 0x3b, 0xe5, 0x39, 0xc9, 0x6f, 0x95, 0x4e,
	0x49, 0x6b, 0x68, 0xd9, 0x2b, 0xb4, 0x9f, 0x8f, 0xf5, 0x93, 0x8c, 0xd5, 0x9d, 0x54, 0x12, 0x93,
	0x2c, 0x1f, 0xd1, 0xf9, 0xc8, 0xad, 0x2e, 0xd0, 0x29, 0xe6, 0x7c, 0xb2, 0x63, 0x1f, 0x21, 0xbc,
	0xc1, 0x85, 0xcb, 0x10, 0x02, 0xc9, 0x6d, 0x13, 0x95, 0xdb, 0x52, 0xdc, 0xb0, 0x60, 0xc0, 0xb5,
	0xa1, 0x8d, 0x51, 0x75, 0x7d, 0x4f, 0xf3, 0xcd, 0x0f, 0x03, 0xed, 0xb6, 0xa5, 0x30, 0xa7, 0xa8,
	0xb4, 0x5e, 0xbf, 0x4e, 0xfe, 0xbc, 0xbd, 0x0d, 0x3d, 0xf1, 0xc5, 0x56, 0xb8, 0x0e, 0x60, 0x3e,
	0xa2, 0x83, 0x95, 0x64, 0xe6, 0xf9, 0xff, 0x3a, 0x2b, 0x30, 0x6e, 0x6d, 0x01, 0x6b, 0xcb, 0xeb,
	0xee, 0xeb, 0xc2, 0x32, 0xe6, 0x0b, 0xcb, 0x78, 0x5f, 0x58, 0xc6, 0xf3, 0xd2, 0x2a, 0xcc, 0x97,
	0x56, 0xe1, 0x6d, 0x69, 0x15, 0xee, 0x2e, 0x85, 0xab, 0xee, 0xc7, 0x0e, 0x19, 0x80, 0x4f, 0x87,
	0xb1, 0xcf, 0x03, 0xe9, 0x42, 0x30, 0x8b, 0x9f, 0xbe, 0x17, 0xf5, 0x68, 0xf8, 0x40, 0x67, 0x74,
	0xe5, 0xb5, 0x27, 0x6e, 0xce, 0x5e, 0x7a, 0x6f, 0xad, 0xcf, 0x00, 0x00, 0x00, 0xff, 0xff, 0xba,
	0xf0, 0x43, 0x7b, 0x09, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MsgClient interface {
	SoftwareUpgrade(ctx context.Context, in *MsgSoftwareUpgrade, opts ...grpc.CallOption) (*MsgSoftwareUpgradeResponse, error)
	CancelUpgrade(ctx context.Context, in *MsgCancelUpgrade, opts ...grpc.CallOption) (*MsgCancelUpgradeResponse, error)
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) SoftwareUpgrade(ctx context.Context, in *MsgSoftwareUpgrade, opts ...grpc.CallOption) (*MsgSoftwareUpgradeResponse, error) {
	out := new(MsgSoftwareUpgradeResponse)
	err := c.cc.Invoke(ctx, "/rollapp.timeupgrade.types.Msg/SoftwareUpgrade", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) CancelUpgrade(ctx context.Context, in *MsgCancelUpgrade, opts ...grpc.CallOption) (*MsgCancelUpgradeResponse, error) {
	out := new(MsgCancelUpgradeResponse)
	err := c.cc.Invoke(ctx, "/rollapp.timeupgrade.types.Msg/CancelUpgrade", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	SoftwareUpgrade(context.Context, *MsgSoftwareUpgrade) (*MsgSoftwareUpgradeResponse, error)
	CancelUpgrade(context.Context, *MsgCancelUpgrade) (*MsgCancelUpgradeResponse, error)
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (*UnimplementedMsgServer) SoftwareUpgrade(ctx context.Context, req *MsgSoftwareUpgrade) (*MsgSoftwareUpgradeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SoftwareUpgrade not implemented")
}
func (*UnimplementedMsgServer) CancelUpgrade(ctx context.Context, req *MsgCancelUpgrade) (*MsgCancelUpgradeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CancelUpgrade not implemented")
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_SoftwareUpgrade_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgSoftwareUpgrade)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).SoftwareUpgrade(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rollapp.timeupgrade.types.Msg/SoftwareUpgrade",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).SoftwareUpgrade(ctx, req.(*MsgSoftwareUpgrade))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_CancelUpgrade_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCancelUpgrade)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CancelUpgrade(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rollapp.timeupgrade.types.Msg/CancelUpgrade",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CancelUpgrade(ctx, req.(*MsgCancelUpgrade))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "rollapp.timeupgrade.types.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SoftwareUpgrade",
			Handler:    _Msg_SoftwareUpgrade_Handler,
		},
		{
			MethodName: "CancelUpgrade",
			Handler:    _Msg_CancelUpgrade_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "timeupgrade/tx.proto",
}

func (m *MsgSoftwareUpgrade) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgSoftwareUpgrade) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgSoftwareUpgrade) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.UpgradeTime != nil {
		{
			size, err := m.UpgradeTime.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTx(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.OriginalUpgrade != nil {
		{
			size, err := m.OriginalUpgrade.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTx(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgSoftwareUpgradeResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgSoftwareUpgradeResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgSoftwareUpgradeResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *MsgCancelUpgrade) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgCancelUpgrade) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgCancelUpgrade) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *MsgCancelUpgradeResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgCancelUpgradeResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgCancelUpgradeResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func encodeVarintTx(dAtA []byte, offset int, v uint64) int {
	offset -= sovTx(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *MsgSoftwareUpgrade) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.OriginalUpgrade != nil {
		l = m.OriginalUpgrade.Size()
		n += 1 + l + sovTx(uint64(l))
	}
	if m.UpgradeTime != nil {
		l = m.UpgradeTime.Size()
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgSoftwareUpgradeResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *MsgCancelUpgrade) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *MsgCancelUpgradeResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func sovTx(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTx(x uint64) (n int) {
	return sovTx(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MsgSoftwareUpgrade) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgSoftwareUpgrade: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgSoftwareUpgrade: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OriginalUpgrade", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.OriginalUpgrade == nil {
				m.OriginalUpgrade = &types.MsgSoftwareUpgrade{}
			}
			if err := m.OriginalUpgrade.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UpgradeTime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.UpgradeTime == nil {
				m.UpgradeTime = &types1.Timestamp{}
			}
			if err := m.UpgradeTime.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func (m *MsgSoftwareUpgradeResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgSoftwareUpgradeResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgSoftwareUpgradeResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func (m *MsgCancelUpgrade) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgCancelUpgrade: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgCancelUpgrade: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func (m *MsgCancelUpgradeResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgCancelUpgradeResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgCancelUpgradeResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func skipTx(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTx
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
					return 0, ErrIntOverflowTx
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
					return 0, ErrIntOverflowTx
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
				return 0, ErrInvalidLengthTx
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTx
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTx
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTx        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTx          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTx = fmt.Errorf("proto: unexpected end of group")
)
