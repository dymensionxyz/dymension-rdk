// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: rollappparams/params.proto

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

// rollapp params defined in genesis and updated via gov proposal
type Params struct {
	// data availability type (e.g. celestia) used in the rollapp
	Da string `protobuf:"bytes,1,opt,name=da,proto3" json:"da,omitempty"`
	// drs version
	DrsVersion uint32 `protobuf:"varint,2,opt,name=drs_version,json=drsVersion,proto3" json:"drs_version,omitempty"`
	// MinGasPrices is globally-specified minimum gas prices for transactions. These values
	// determine which denoms validators can use for accepting fees as well as minimum gas prices
	// for fees in each denom. Values from this list overwrite the validator-specified minimum
	// gas prices if greater. If the list is empty, then validators can accept any denom they specify.
	//
	// For example:
	//
	//  Global:    [10adym 1stake  5uatom] <- Validator could only accept fees in these denoms.
	//  Validator: [1adym  10stake        1uosmo]
	//  Final:     [10adym 10stake]
	//
	// After merging, the validator would only be able to accept fees greater than 10adym or 10stake.
	// If a validator attempted to accept a fee of 6uatom or 2uosmo, the transaction would be rejected.
	//
	// Possible cases:
	//
	//  | Global    | Validator | Result                       |
	//  |-----------|-----------|------------------------------|
	//  | empty     | empty     | all txs are accepted         |
	//  | empty     | non-empty | validator values             |
	//  | non-empty | empty     | global values                |
	//  | non-empty | non-empty | intersect(global, validator) |
	MinGasPrices github_com_cosmos_cosmos_sdk_types.DecCoins `protobuf:"bytes,3,rep,name=min_gas_prices,json=minGasPrices,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.DecCoins" json:"min_gas_prices"`
}

func (m *Params) Reset()         { *m = Params{} }
func (m *Params) String() string { return proto.CompactTextString(m) }
func (*Params) ProtoMessage()    {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_c7f89f13f4d953f6, []int{0}
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

func (m *Params) GetDa() string {
	if m != nil {
		return m.Da
	}
	return ""
}

func (m *Params) GetDrsVersion() uint32 {
	if m != nil {
		return m.DrsVersion
	}
	return 0
}

func (m *Params) GetMinGasPrices() github_com_cosmos_cosmos_sdk_types.DecCoins {
	if m != nil {
		return m.MinGasPrices
	}
	return nil
}

func init() {
	proto.RegisterType((*Params)(nil), "rollapp.params.types.Params")
}

func init() { proto.RegisterFile("rollappparams/params.proto", fileDescriptor_c7f89f13f4d953f6) }

var fileDescriptor_c7f89f13f4d953f6 = []byte{
	// 302 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x90, 0xcf, 0x4e, 0x02, 0x31,
	0x10, 0xc6, 0xb7, 0x90, 0x90, 0xb8, 0x28, 0x87, 0x0d, 0x07, 0x42, 0x4c, 0x21, 0x9e, 0x48, 0x0c,
	0x6d, 0x90, 0xa3, 0x37, 0x34, 0xf1, 0x4a, 0x48, 0xf4, 0xe0, 0x85, 0x74, 0xb7, 0xcd, 0xda, 0x40,
	0xff, 0xa4, 0xb3, 0x22, 0xf8, 0x14, 0x3e, 0x87, 0x47, 0x9f, 0x82, 0x23, 0x47, 0x4f, 0x6a, 0xe0,
	0x45, 0xcc, 0xb6, 0x6b, 0xd4, 0xd3, 0xb4, 0x33, 0xd3, 0xef, 0xf7, 0xf5, 0x8b, 0xbb, 0xce, 0x2c,
	0x97, 0xcc, 0x5a, 0xcb, 0x1c, 0x53, 0x40, 0x43, 0x21, 0xd6, 0x99, 0xc2, 0x24, 0xed, 0x6a, 0x46,
	0xaa, 0x6e, 0xb1, 0xb1, 0x02, 0xba, 0xed, 0xdc, 0xe4, 0xc6, 0x2f, 0xd0, 0xf2, 0x14, 0x76, 0xbb,
	0x38, 0x33, 0xa0, 0x0c, 0xd0, 0x94, 0x81, 0xa0, 0xab, 0x51, 0x2a, 0x0a, 0x36, 0xa2, 0x99, 0x91,
	0x3a, 0xcc, 0xcf, 0xde, 0x50, 0xdc, 0x98, 0x7a, 0x99, 0xa4, 0x15, 0xd7, 0x38, 0xeb, 0xa0, 0x3e,
	0x1a, 0x1c, 0xcd, 0x6a, 0x9c, 0x25, 0xbd, 0xb8, 0xc9, 0x1d, 0xcc, 0x57, 0xc2, 0x81, 0x34, 0xba,
	0x53, 0xeb, 0xa3, 0xc1, 0xc9, 0x2c, 0xe6, 0x0e, 0xee, 0x42, 0x27, 0x79, 0x8a, 0x5b, 0x4a, 0xea,
	0x79, 0xce, 0x60, 0x6e, 0x9d, 0xcc, 0x04, 0x74, 0xea, 0xfd, 0xfa, 0xa0, 0x79, 0x71, 0x4a, 0x02,
	0x94, 0x94, 0x50, 0x52, 0x41, 0xc9, 0xb5, 0xc8, 0xae, 0x8c, 0xd4, 0x93, 0xf1, 0xf6, 0xa3, 0x17,
	0xbd, 0x7e, 0xf6, 0xce, 0x73, 0x59, 0x3c, 0x3c, 0xa6, 0x24, 0x33, 0x8a, 0x56, 0x26, 0x43, 0x19,
	0x02, 0x5f, 0x50, 0xff, 0xa7, 0x9f, 0x37, 0x30, 0x3b, 0x56, 0x52, 0xdf, 0x30, 0x98, 0x7a, 0xcc,
	0xe4, 0x76, 0xbb, 0xc7, 0x68, 0xb7, 0xc7, 0xe8, 0x6b, 0x8f, 0xd1, 0xcb, 0x01, 0x47, 0xbb, 0x03,
	0x8e, 0xde, 0x0f, 0x38, 0xba, 0xbf, 0xfc, 0x23, 0xca, 0x37, 0x4a, 0xe8, 0xd2, 0xe8, 0x7a, 0xf3,
	0xfc, 0x7b, 0x19, 0x3a, 0xbe, 0xa0, 0x6b, 0xfa, 0x3f, 0x5e, 0x4f, 0x4b, 0x1b, 0x3e, 0x92, 0xf1,
	0x77, 0x00, 0x00, 0x00, 0xff, 0xff, 0x87, 0xcf, 0x9b, 0x4a, 0x7c, 0x01, 0x00, 0x00,
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
	if len(m.MinGasPrices) > 0 {
		for iNdEx := len(m.MinGasPrices) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.MinGasPrices[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintParams(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if m.DrsVersion != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.DrsVersion))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Da) > 0 {
		i -= len(m.Da)
		copy(dAtA[i:], m.Da)
		i = encodeVarintParams(dAtA, i, uint64(len(m.Da)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintParams(dAtA []byte, offset int, v uint64) int {
	offset -= sovParams(v)
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
	l = len(m.Da)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	if m.DrsVersion != 0 {
		n += 1 + sovParams(uint64(m.DrsVersion))
	}
	if len(m.MinGasPrices) > 0 {
		for _, e := range m.MinGasPrices {
			l = e.Size()
			n += 1 + l + sovParams(uint64(l))
		}
	}
	return n
}

func sovParams(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozParams(x uint64) (n int) {
	return sovParams(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowParams
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
				return fmt.Errorf("proto: wrong wireType = %d for field Da", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Da = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field DrsVersion", wireType)
			}
			m.DrsVersion = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.DrsVersion |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinGasPrices", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.MinGasPrices = append(m.MinGasPrices, types.DecCoin{})
			if err := m.MinGasPrices[len(m.MinGasPrices)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipParams(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthParams
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
func skipParams(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowParams
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
					return 0, ErrIntOverflowParams
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
					return 0, ErrIntOverflowParams
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
				return 0, ErrInvalidLengthParams
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupParams
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthParams
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthParams        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowParams          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupParams = fmt.Errorf("proto: unexpected end of group")
)
