// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sequencers/sequencers.proto

package types

import (
	fmt "fmt"
	types1 "github.com/cosmos/cosmos-sdk/x/staking/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	_ "github.com/gogo/protobuf/types"
	github_com_gogo_protobuf_types "github.com/gogo/protobuf/types"
	io "io"
	math "math"
	math_bits "math/bits"
	time "time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// Params defines the parameters for the module.
type Params struct {
	// UnbondingTime is the time duration of unbonding.
	UnbondingTime time.Duration `protobuf:"bytes,1,opt,name=unbonding_time,json=unbondingTime,proto3,stdduration" json:"unbonding_time" yaml:"unbonding_time"`
	// HistoricalEntries is the number of historical entries to persist.
	HistoricalEntries uint32 `protobuf:"varint,2,opt,name=historical_entries,json=historicalEntries,proto3" json:"historical_entries,omitempty" yaml:"historical_entries"`
}

func (m *Params) Reset()      { *m = Params{} }
func (*Params) ProtoMessage() {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_632cf56d50d4d0f1, []int{0}
}
func (m *Params) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Params) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Params.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Params) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Params.Merge(m, src)
}
func (m *Params) XXX_Size() int {
	return m.Size()
}
func (m *Params) XXX_DiscardUnknown() {
	xxx_messageInfo_Params.DiscardUnknown(m)
}

var xxx_messageInfo_Params proto.InternalMessageInfo

func (m *Params) GetUnbondingTime() time.Duration {
	if m != nil {
		return m.UnbondingTime
	}
	return 0
}

func (m *Params) GetHistoricalEntries() uint32 {
	if m != nil {
		return m.HistoricalEntries
	}
	return 0
}

type Sequencer struct {
	// Validator is a convenient storage for e.g operator address and consensus pub key
	Validator *types1.Validator `protobuf:"bytes,1,opt,name=validator,proto3" json:"validator,omitempty"`
	// RewardAddr is the sdk acc address where the sequencer has opted to receive rewards. Empty if not set.
	RewardAddr string `protobuf:"bytes,2,opt,name=reward_addr,json=rewardAddr,proto3" json:"reward_addr,omitempty"`
	// Relayers is an array of the whitelisted relayer addresses. Addresses are bech32-encoded strings.
	Relayers []string `protobuf:"bytes,3,rep,name=relayers,proto3" json:"relayers,omitempty"`
}

func (m *Sequencer) Reset()         { *m = Sequencer{} }
func (m *Sequencer) String() string { return proto.CompactTextString(m) }
func (*Sequencer) ProtoMessage()    {}
func (*Sequencer) Descriptor() ([]byte, []int) {
	return fileDescriptor_632cf56d50d4d0f1, []int{1}
}
func (m *Sequencer) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Sequencer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Sequencer.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Sequencer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Sequencer.Merge(m, src)
}
func (m *Sequencer) XXX_Size() int {
	return m.Size()
}
func (m *Sequencer) XXX_DiscardUnknown() {
	xxx_messageInfo_Sequencer.DiscardUnknown(m)
}

var xxx_messageInfo_Sequencer proto.InternalMessageInfo

func (m *Sequencer) GetValidator() *types1.Validator {
	if m != nil {
		return m.Validator
	}
	return nil
}

func (m *Sequencer) GetRewardAddr() string {
	if m != nil {
		return m.RewardAddr
	}
	return ""
}

func (m *Sequencer) GetRelayers() []string {
	if m != nil {
		return m.Relayers
	}
	return nil
}

// WhitelistedRelayers is used for storing the whitelisted relater list in the state
type WhitelistedRelayers struct {
	// Relayers is an array of the whitelisted relayer addresses. Addresses are bech32-encoded strings.
	Relayers []string `protobuf:"bytes,1,rep,name=relayers,proto3" json:"relayers,omitempty"`
}

func (m *WhitelistedRelayers) Reset()         { *m = WhitelistedRelayers{} }
func (m *WhitelistedRelayers) String() string { return proto.CompactTextString(m) }
func (*WhitelistedRelayers) ProtoMessage()    {}
func (*WhitelistedRelayers) Descriptor() ([]byte, []int) {
	return fileDescriptor_632cf56d50d4d0f1, []int{2}
}
func (m *WhitelistedRelayers) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *WhitelistedRelayers) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_WhitelistedRelayers.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *WhitelistedRelayers) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WhitelistedRelayers.Merge(m, src)
}
func (m *WhitelistedRelayers) XXX_Size() int {
	return m.Size()
}
func (m *WhitelistedRelayers) XXX_DiscardUnknown() {
	xxx_messageInfo_WhitelistedRelayers.DiscardUnknown(m)
}

var xxx_messageInfo_WhitelistedRelayers proto.InternalMessageInfo

func (m *WhitelistedRelayers) GetRelayers() []string {
	if m != nil {
		return m.Relayers
	}
	return nil
}

func init() {
	proto.RegisterType((*Params)(nil), "rollapp.sequencers.types.Params")
	proto.RegisterType((*Sequencer)(nil), "rollapp.sequencers.types.Sequencer")
	proto.RegisterType((*WhitelistedRelayers)(nil), "rollapp.sequencers.types.WhitelistedRelayers")
}

func init() { proto.RegisterFile("sequencers/sequencers.proto", fileDescriptor_632cf56d50d4d0f1) }

var fileDescriptor_632cf56d50d4d0f1 = []byte{
	// 426 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x52, 0x3d, 0x6f, 0xd4, 0x40,
	0x10, 0xf5, 0x12, 0x29, 0xca, 0x6d, 0x14, 0x24, 0x0c, 0x48, 0x97, 0x43, 0xd8, 0x17, 0x8b, 0xe2,
	0x1a, 0x76, 0x75, 0xd0, 0xa0, 0x34, 0x88, 0x13, 0x74, 0x14, 0xc8, 0x20, 0x90, 0x68, 0x4e, 0x6b,
	0xef, 0xe0, 0x5b, 0x65, 0xbd, 0x6b, 0x76, 0xd7, 0x21, 0xe6, 0x1f, 0xd0, 0x51, 0xa6, 0xcc, 0x9f,
	0x41, 0x4a, 0x99, 0x92, 0xea, 0x40, 0x77, 0x0d, 0x75, 0x7e, 0x01, 0xc2, 0x1f, 0xf1, 0x9d, 0xe8,
	0xe6, 0xbd, 0x79, 0xcf, 0xe3, 0x79, 0x3b, 0xf8, 0x81, 0x85, 0xcf, 0x25, 0xa8, 0x14, 0x8c, 0xa5,
	0x7d, 0x49, 0x0a, 0xa3, 0x9d, 0xf6, 0x87, 0x46, 0x4b, 0xc9, 0x8a, 0x82, 0x6c, 0x74, 0x5c, 0x55,
	0x80, 0x1d, 0xdd, 0xcb, 0x74, 0xa6, 0x6b, 0x11, 0xfd, 0x57, 0x35, 0xfa, 0x51, 0x90, 0x69, 0x9d,
	0x49, 0xa0, 0x35, 0x4a, 0xca, 0x4f, 0x94, 0x97, 0x86, 0x39, 0xa1, 0x55, 0xdb, 0x7f, 0x94, 0x6a,
	0x9b, 0x6b, 0x4b, 0xad, 0x63, 0x27, 0x42, 0x65, 0xf4, 0x74, 0x9a, 0x80, 0x63, 0xd3, 0x0e, 0x37,
	0xaa, 0xe8, 0x07, 0xc2, 0xbb, 0x6f, 0x98, 0x61, 0xb9, 0xf5, 0x53, 0x7c, 0xbb, 0x54, 0x89, 0x56,
	0x5c, 0xa8, 0x6c, 0xee, 0x44, 0x0e, 0x43, 0x34, 0x46, 0x93, 0xfd, 0x27, 0x87, 0xa4, 0x99, 0x44,
	0xba, 0x49, 0xe4, 0x65, 0x3b, 0x69, 0x76, 0x74, 0xb9, 0x0c, 0xbd, 0xeb, 0x65, 0x78, 0xbf, 0x62,
	0xb9, 0x3c, 0x8e, 0xb6, 0xed, 0xd1, 0xf9, 0xaf, 0x10, 0xc5, 0x07, 0x37, 0xe4, 0x3b, 0x91, 0x83,
	0xff, 0x1a, 0xfb, 0x0b, 0x61, 0x9d, 0x36, 0x22, 0x65, 0x72, 0x0e, 0xca, 0x19, 0x01, 0x76, 0x78,
	0x6b, 0x8c, 0x26, 0x07, 0xb3, 0x87, 0xd7, 0xcb, 0xf0, 0xb0, 0xf9, 0xd2, 0xff, 0x9a, 0x28, 0xbe,
	0xd3, 0x93, 0xaf, 0x1a, 0xee, 0x78, 0xef, 0xfc, 0x22, 0xf4, 0xfe, 0x5c, 0x84, 0x28, 0xfa, 0x86,
	0xf0, 0xe0, 0x6d, 0x17, 0x9c, 0xff, 0x1c, 0x0f, 0x4e, 0x99, 0x14, 0x9c, 0x39, 0x6d, 0xda, 0x2d,
	0x8e, 0x48, 0x93, 0x07, 0xe9, 0xf6, 0x6f, 0xf3, 0x20, 0xef, 0x3b, 0x61, 0xdc, 0x7b, 0xfc, 0x10,
	0xef, 0x1b, 0xf8, 0xc2, 0x0c, 0x9f, 0x33, 0xce, 0x4d, 0xfd, 0x7f, 0x83, 0x18, 0x37, 0xd4, 0x0b,
	0xce, 0x8d, 0x3f, 0xc2, 0x7b, 0x06, 0x24, 0xab, 0xc0, 0xd8, 0xe1, 0xce, 0x78, 0x67, 0x32, 0x88,
	0x6f, 0x70, 0x34, 0xc5, 0x77, 0x3f, 0x2c, 0x84, 0x03, 0x29, 0xac, 0x03, 0x1e, 0xb7, 0xf4, 0x96,
	0x05, 0x6d, 0x5b, 0x66, 0xf1, 0xe5, 0x2a, 0x40, 0x57, 0xab, 0x00, 0xfd, 0x5e, 0x05, 0xe8, 0xfb,
	0x3a, 0xf0, 0xae, 0xd6, 0x81, 0xf7, 0x73, 0x1d, 0x78, 0x1f, 0x9f, 0x65, 0xc2, 0x2d, 0xca, 0x84,
	0xa4, 0x3a, 0xa7, 0xbc, 0xca, 0x41, 0x59, 0xa1, 0xd5, 0x59, 0xf5, 0xb5, 0x07, 0x8f, 0x0d, 0x3f,
	0xa1, 0x67, 0x1b, 0x07, 0x45, 0xeb, 0xb3, 0x49, 0x76, 0xeb, 0xf7, 0x7a, 0xfa, 0x37, 0x00, 0x00,
	0xff, 0xff, 0x49, 0x03, 0xc8, 0x30, 0x76, 0x02, 0x00, 0x00,
}

func (this *Params) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Params)
	if !ok {
		that2, ok := that.(Params)
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
	if this.UnbondingTime != that1.UnbondingTime {
		return false
	}
	if this.HistoricalEntries != that1.HistoricalEntries {
		return false
	}
	return true
}
func (m *Params) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Params) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Params) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.HistoricalEntries != 0 {
		i = encodeVarintSequencers(dAtA, i, uint64(m.HistoricalEntries))
		i--
		dAtA[i] = 0x10
	}
	n1, err1 := github_com_gogo_protobuf_types.StdDurationMarshalTo(m.UnbondingTime, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdDuration(m.UnbondingTime):])
	if err1 != nil {
		return 0, err1
	}
	i -= n1
	i = encodeVarintSequencers(dAtA, i, uint64(n1))
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *Sequencer) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Sequencer) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Sequencer) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Relayers) > 0 {
		for iNdEx := len(m.Relayers) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Relayers[iNdEx])
			copy(dAtA[i:], m.Relayers[iNdEx])
			i = encodeVarintSequencers(dAtA, i, uint64(len(m.Relayers[iNdEx])))
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.RewardAddr) > 0 {
		i -= len(m.RewardAddr)
		copy(dAtA[i:], m.RewardAddr)
		i = encodeVarintSequencers(dAtA, i, uint64(len(m.RewardAddr)))
		i--
		dAtA[i] = 0x12
	}
	if m.Validator != nil {
		{
			size, err := m.Validator.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSequencers(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *WhitelistedRelayers) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *WhitelistedRelayers) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *WhitelistedRelayers) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Relayers) > 0 {
		for iNdEx := len(m.Relayers) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Relayers[iNdEx])
			copy(dAtA[i:], m.Relayers[iNdEx])
			i = encodeVarintSequencers(dAtA, i, uint64(len(m.Relayers[iNdEx])))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintSequencers(dAtA []byte, offset int, v uint64) int {
	offset -= sovSequencers(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = github_com_gogo_protobuf_types.SizeOfStdDuration(m.UnbondingTime)
	n += 1 + l + sovSequencers(uint64(l))
	if m.HistoricalEntries != 0 {
		n += 1 + sovSequencers(uint64(m.HistoricalEntries))
	}
	return n
}

func (m *Sequencer) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Validator != nil {
		l = m.Validator.Size()
		n += 1 + l + sovSequencers(uint64(l))
	}
	l = len(m.RewardAddr)
	if l > 0 {
		n += 1 + l + sovSequencers(uint64(l))
	}
	if len(m.Relayers) > 0 {
		for _, s := range m.Relayers {
			l = len(s)
			n += 1 + l + sovSequencers(uint64(l))
		}
	}
	return n
}

func (m *WhitelistedRelayers) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Relayers) > 0 {
		for _, s := range m.Relayers {
			l = len(s)
			n += 1 + l + sovSequencers(uint64(l))
		}
	}
	return n
}

func sovSequencers(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozSequencers(x uint64) (n int) {
	return sovSequencers(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSequencers
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
			return fmt.Errorf("proto: Params: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Params: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UnbondingTime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSequencers
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
				return ErrInvalidLengthSequencers
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSequencers
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdDurationUnmarshal(&m.UnbondingTime, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field HistoricalEntries", wireType)
			}
			m.HistoricalEntries = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSequencers
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.HistoricalEntries |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipSequencers(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSequencers
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
func (m *Sequencer) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSequencers
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
			return fmt.Errorf("proto: Sequencer: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Sequencer: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Validator", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSequencers
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
				return ErrInvalidLengthSequencers
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSequencers
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Validator == nil {
				m.Validator = &types1.Validator{}
			}
			if err := m.Validator.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RewardAddr", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSequencers
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
				return ErrInvalidLengthSequencers
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSequencers
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.RewardAddr = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Relayers", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSequencers
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
				return ErrInvalidLengthSequencers
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSequencers
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Relayers = append(m.Relayers, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSequencers(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSequencers
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
func (m *WhitelistedRelayers) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSequencers
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
			return fmt.Errorf("proto: WhitelistedRelayers: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: WhitelistedRelayers: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Relayers", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSequencers
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
				return ErrInvalidLengthSequencers
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSequencers
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Relayers = append(m.Relayers, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSequencers(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSequencers
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
func skipSequencers(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowSequencers
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
					return 0, ErrIntOverflowSequencers
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
					return 0, ErrIntOverflowSequencers
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
				return 0, ErrInvalidLengthSequencers
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupSequencers
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthSequencers
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthSequencers        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowSequencers          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupSequencers = fmt.Errorf("proto: unexpected end of group")
)
