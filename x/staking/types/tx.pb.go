// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: governors/tx.proto

package types

import (
	context "context"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/codec/types"
	_ "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/cosmos-sdk/types/msgservice"
	types "github.com/cosmos/cosmos-sdk/x/staking/types"
	grpc1 "github.com/gogo/protobuf/grpc"
	proto "github.com/gogo/protobuf/proto"
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

// MsgCreateValidator defines a SDK message for creating a new validator.
type MsgCreateValidator struct {
	Value *types.MsgCreateValidator `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (m *MsgCreateValidator) Reset()         { *m = MsgCreateValidator{} }
func (m *MsgCreateValidator) String() string { return proto.CompactTextString(m) }
func (*MsgCreateValidator) ProtoMessage()    {}
func (*MsgCreateValidator) Descriptor() ([]byte, []int) {
	return fileDescriptor_a43bdaa658212415, []int{0}
}
func (m *MsgCreateValidator) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgCreateValidator) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgCreateValidator.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgCreateValidator) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgCreateValidator.Merge(m, src)
}
func (m *MsgCreateValidator) XXX_Size() int {
	return m.Size()
}
func (m *MsgCreateValidator) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgCreateValidator.DiscardUnknown(m)
}

var xxx_messageInfo_MsgCreateValidator proto.InternalMessageInfo

func (m *MsgCreateValidator) GetValue() *types.MsgCreateValidator {
	if m != nil {
		return m.Value
	}
	return nil
}

// MsgCreateValidatorResponse defines the Msg/CreateValidator response type.
type MsgCreateValidatorResponse struct {
	Value *types.MsgCreateValidatorResponse `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (m *MsgCreateValidatorResponse) Reset()         { *m = MsgCreateValidatorResponse{} }
func (m *MsgCreateValidatorResponse) String() string { return proto.CompactTextString(m) }
func (*MsgCreateValidatorResponse) ProtoMessage()    {}
func (*MsgCreateValidatorResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_a43bdaa658212415, []int{1}
}
func (m *MsgCreateValidatorResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgCreateValidatorResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgCreateValidatorResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgCreateValidatorResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgCreateValidatorResponse.Merge(m, src)
}
func (m *MsgCreateValidatorResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgCreateValidatorResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgCreateValidatorResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgCreateValidatorResponse proto.InternalMessageInfo

func (m *MsgCreateValidatorResponse) GetValue() *types.MsgCreateValidatorResponse {
	if m != nil {
		return m.Value
	}
	return nil
}

// MsgDelegate defines a SDK message for performing a delegation of coins
// from a delegator to a validator.
type MsgDelegate struct {
	Value *types.MsgDelegate `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (m *MsgDelegate) Reset()         { *m = MsgDelegate{} }
func (m *MsgDelegate) String() string { return proto.CompactTextString(m) }
func (*MsgDelegate) ProtoMessage()    {}
func (*MsgDelegate) Descriptor() ([]byte, []int) {
	return fileDescriptor_a43bdaa658212415, []int{2}
}
func (m *MsgDelegate) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgDelegate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgDelegate.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgDelegate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgDelegate.Merge(m, src)
}
func (m *MsgDelegate) XXX_Size() int {
	return m.Size()
}
func (m *MsgDelegate) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgDelegate.DiscardUnknown(m)
}

var xxx_messageInfo_MsgDelegate proto.InternalMessageInfo

func (m *MsgDelegate) GetValue() *types.MsgDelegate {
	if m != nil {
		return m.Value
	}
	return nil
}

// MsgDelegateResponse defines the Msg/Delegate response type.
type MsgDelegateResponse struct {
	Value *types.MsgDelegateResponse `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (m *MsgDelegateResponse) Reset()         { *m = MsgDelegateResponse{} }
func (m *MsgDelegateResponse) String() string { return proto.CompactTextString(m) }
func (*MsgDelegateResponse) ProtoMessage()    {}
func (*MsgDelegateResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_a43bdaa658212415, []int{3}
}
func (m *MsgDelegateResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgDelegateResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgDelegateResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgDelegateResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgDelegateResponse.Merge(m, src)
}
func (m *MsgDelegateResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgDelegateResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgDelegateResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgDelegateResponse proto.InternalMessageInfo

func (m *MsgDelegateResponse) GetValue() *types.MsgDelegateResponse {
	if m != nil {
		return m.Value
	}
	return nil
}

func init() {
	proto.RegisterType((*MsgCreateValidator)(nil), "rollapp.staking.MsgCreateValidator")
	proto.RegisterType((*MsgCreateValidatorResponse)(nil), "rollapp.staking.MsgCreateValidatorResponse")
	proto.RegisterType((*MsgDelegate)(nil), "rollapp.staking.MsgDelegate")
	proto.RegisterType((*MsgDelegateResponse)(nil), "rollapp.staking.MsgDelegateResponse")
}

func init() { proto.RegisterFile("governors/tx.proto", fileDescriptor_a43bdaa658212415) }

var fileDescriptor_a43bdaa658212415 = []byte{
	// 367 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x92, 0x41, 0x4b, 0xc3, 0x40,
	0x10, 0x85, 0x1b, 0x44, 0x0f, 0x5b, 0x44, 0x58, 0x05, 0x35, 0x48, 0x94, 0xd6, 0x83, 0x58, 0xcc,
	0xda, 0x88, 0x07, 0x6f, 0x6a, 0x15, 0x7a, 0x29, 0x42, 0x85, 0x22, 0xde, 0x36, 0xed, 0x76, 0x1b,
	0x9a, 0x64, 0x42, 0x76, 0x1b, 0x1a, 0x7f, 0x85, 0x3f, 0xca, 0x83, 0xc7, 0x1e, 0x3d, 0x4a, 0xfb,
	0x47, 0xc4, 0x66, 0xb7, 0x9a, 0x5a, 0x4a, 0x3d, 0x85, 0xcd, 0x7b, 0xef, 0x9b, 0x37, 0x30, 0x08,
	0x73, 0x48, 0x58, 0x1c, 0x42, 0x2c, 0x88, 0x1c, 0xda, 0x51, 0x0c, 0x12, 0xf0, 0x56, 0x0c, 0xbe,
	0x4f, 0xa3, 0xc8, 0x16, 0x92, 0xf6, 0xbd, 0x90, 0x9b, 0xbb, 0x6d, 0x10, 0x01, 0x08, 0x12, 0x08,
	0x4e, 0x92, 0xea, 0xf7, 0x27, 0x73, 0x9a, 0xfb, 0x1c, 0x80, 0xfb, 0x8c, 0x4c, 0x5f, 0xee, 0xa0,
	0x4b, 0x68, 0x98, 0x2a, 0xc9, 0x52, 0x19, 0x97, 0x0a, 0x46, 0x92, 0xaa, 0xcb, 0x24, 0xad, 0x92,
	0x36, 0x78, 0xa1, 0xd2, 0x0f, 0x95, 0xae, 0x66, 0xcc, 0x2c, 0xba, 0x45, 0xa9, 0x85, 0x70, 0x43,
	0xf0, 0x5a, 0xcc, 0xa8, 0x64, 0x2d, 0xea, 0x7b, 0x1d, 0x2a, 0x21, 0xc6, 0xd7, 0x68, 0x3d, 0xa1,
	0xfe, 0x80, 0xed, 0x19, 0x47, 0xc6, 0x49, 0xd1, 0x39, 0xb5, 0x33, 0x8c, 0xae, 0x6a, 0x2b, 0x8c,
	0xfd, 0x37, 0xda, 0xcc, 0x82, 0xa5, 0x2e, 0x32, 0x17, 0x88, 0x4c, 0x44, 0x10, 0x0a, 0x86, 0xeb,
	0x79, 0xbe, 0xf3, 0x0f, 0xbe, 0x42, 0xe8, 0x39, 0x75, 0x54, 0x6c, 0x08, 0x7e, 0xc7, 0x7c, 0xc6,
	0xa9, 0x64, 0xf8, 0x2a, 0x0f, 0x2e, 0x2f, 0x01, 0xeb, 0x8c, 0x26, 0x3d, 0xa1, 0xed, 0xdf, 0x7f,
	0x75, 0xd5, 0x9b, 0x3c, 0xb1, 0xb2, 0x0a, 0x31, 0xdf, 0xd1, 0x79, 0x33, 0xd0, 0x5a, 0x43, 0x70,
	0xdc, 0x43, 0x3b, 0x73, 0xdb, 0xdc, 0x37, 0x6b, 0xce, 0x39, 0x2e, 0xdb, 0x73, 0xa7, 0xb0, 0x60,
	0x6f, 0xb3, 0xb2, 0x82, 0x69, 0x56, 0xfa, 0x11, 0x6d, 0xea, 0x32, 0xd9, 0x88, 0x83, 0x45, 0x69,
	0x6d, 0x31, 0x8f, 0x97, 0xa9, 0x1a, 0x7a, 0xfb, 0xf0, 0x3e, 0xb6, 0x8c, 0xd1, 0xd8, 0x32, 0x3e,
	0xc7, 0x96, 0xf1, 0x3a, 0xb1, 0x0a, 0xa3, 0x89, 0x55, 0xf8, 0x98, 0x58, 0x85, 0xe7, 0x4b, 0xee,
	0xc9, 0xde, 0xc0, 0xb5, 0xdb, 0x10, 0x90, 0x4e, 0x1a, 0xb0, 0x50, 0x78, 0x10, 0x0e, 0xd3, 0x97,
	0x9f, 0xc7, 0x59, 0xdc, 0xe9, 0x93, 0xe1, 0xec, 0x0c, 0x65, 0x1a, 0x31, 0xe1, 0x6e, 0x4c, 0x4f,
	0xf0, 0xe2, 0x2b, 0x00, 0x00, 0xff, 0xff, 0x5d, 0x5f, 0x54, 0x27, 0x1e, 0x03, 0x00, 0x00,
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
	// CreateValidator defines a method for creating a new validator.
	CreateValidatorERC20(ctx context.Context, in *MsgCreateValidator, opts ...grpc.CallOption) (*MsgCreateValidatorResponse, error)
	// Delegate defines a method for performing a delegation of coins
	// from a delegator to a validator.
	DelegateERC20(ctx context.Context, in *MsgDelegate, opts ...grpc.CallOption) (*MsgDelegateResponse, error)
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) CreateValidatorERC20(ctx context.Context, in *MsgCreateValidator, opts ...grpc.CallOption) (*MsgCreateValidatorResponse, error) {
	out := new(MsgCreateValidatorResponse)
	err := c.cc.Invoke(ctx, "/rollapp.staking.Msg/CreateValidatorERC20", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) DelegateERC20(ctx context.Context, in *MsgDelegate, opts ...grpc.CallOption) (*MsgDelegateResponse, error) {
	out := new(MsgDelegateResponse)
	err := c.cc.Invoke(ctx, "/rollapp.staking.Msg/DelegateERC20", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	// CreateValidator defines a method for creating a new validator.
	CreateValidatorERC20(context.Context, *MsgCreateValidator) (*MsgCreateValidatorResponse, error)
	// Delegate defines a method for performing a delegation of coins
	// from a delegator to a validator.
	DelegateERC20(context.Context, *MsgDelegate) (*MsgDelegateResponse, error)
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (*UnimplementedMsgServer) CreateValidatorERC20(ctx context.Context, req *MsgCreateValidator) (*MsgCreateValidatorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateValidatorERC20 not implemented")
}
func (*UnimplementedMsgServer) DelegateERC20(ctx context.Context, req *MsgDelegate) (*MsgDelegateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DelegateERC20 not implemented")
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_CreateValidatorERC20_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCreateValidator)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CreateValidatorERC20(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rollapp.staking.Msg/CreateValidatorERC20",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CreateValidatorERC20(ctx, req.(*MsgCreateValidator))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_DelegateERC20_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgDelegate)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).DelegateERC20(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rollapp.staking.Msg/DelegateERC20",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).DelegateERC20(ctx, req.(*MsgDelegate))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "rollapp.staking.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateValidatorERC20",
			Handler:    _Msg_CreateValidatorERC20_Handler,
		},
		{
			MethodName: "DelegateERC20",
			Handler:    _Msg_DelegateERC20_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "governors/tx.proto",
}

func (m *MsgCreateValidator) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgCreateValidator) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgCreateValidator) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Value != nil {
		{
			size, err := m.Value.MarshalToSizedBuffer(dAtA[:i])
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

func (m *MsgCreateValidatorResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgCreateValidatorResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgCreateValidatorResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Value != nil {
		{
			size, err := m.Value.MarshalToSizedBuffer(dAtA[:i])
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

func (m *MsgDelegate) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgDelegate) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgDelegate) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Value != nil {
		{
			size, err := m.Value.MarshalToSizedBuffer(dAtA[:i])
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

func (m *MsgDelegateResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgDelegateResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgDelegateResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Value != nil {
		{
			size, err := m.Value.MarshalToSizedBuffer(dAtA[:i])
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
func (m *MsgCreateValidator) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Value != nil {
		l = m.Value.Size()
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgCreateValidatorResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Value != nil {
		l = m.Value.Size()
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgDelegate) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Value != nil {
		l = m.Value.Size()
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgDelegateResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Value != nil {
		l = m.Value.Size()
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func sovTx(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTx(x uint64) (n int) {
	return sovTx(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MsgCreateValidator) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgCreateValidator: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgCreateValidator: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Value", wireType)
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
			if m.Value == nil {
				m.Value = &types.MsgCreateValidator{}
			}
			if err := m.Value.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *MsgCreateValidatorResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgCreateValidatorResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgCreateValidatorResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Value", wireType)
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
			if m.Value == nil {
				m.Value = &types.MsgCreateValidatorResponse{}
			}
			if err := m.Value.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *MsgDelegate) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgDelegate: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgDelegate: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Value", wireType)
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
			if m.Value == nil {
				m.Value = &types.MsgDelegate{}
			}
			if err := m.Value.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *MsgDelegateResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgDelegateResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgDelegateResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Value", wireType)
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
			if m.Value == nil {
				m.Value = &types.MsgDelegateResponse{}
			}
			if err := m.Value.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
