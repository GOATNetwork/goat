// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: goat/bitcoin/v1/params.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/types/tx/amino"
	proto "github.com/cosmos/gogoproto/proto"
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

// Params defines the parameters for the module.
type Params struct {
	NetworkName string `protobuf:"bytes,1,opt,name=network_name,json=networkName,proto3" json:"network_name,omitempty"`
	// A block has the number should be considered as a finalized block
	ConfirmationNumber uint64 `protobuf:"varint,2,opt,name=confirmation_number,json=confirmationNumber,proto3" json:"confirmation_number,omitempty"`
	// min deposit amount in satoshi
	MinDepositAmount   uint64 `protobuf:"varint,3,opt,name=min_deposit_amount,json=minDepositAmount,proto3" json:"min_deposit_amount,omitempty"`
	DepositMagicPrefix []byte `protobuf:"bytes,4,opt,name=deposit_magic_prefix,json=depositMagicPrefix,proto3" json:"deposit_magic_prefix,omitempty"`
	DepositTaxRate     uint64 `protobuf:"varint,5,opt,name=deposit_tax_rate,json=depositTaxRate,proto3" json:"deposit_tax_rate,omitempty"`
	// max deposit tax in satoshi
	MaxDepositTax uint64 `protobuf:"varint,6,opt,name=max_deposit_tax,json=maxDepositTax,proto3" json:"max_deposit_tax,omitempty"`
}

func (m *Params) Reset()         { *m = Params{} }
func (m *Params) String() string { return proto.CompactTextString(m) }
func (*Params) ProtoMessage()    {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_96fcb82db75b9cd1, []int{0}
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

func (m *Params) GetNetworkName() string {
	if m != nil {
		return m.NetworkName
	}
	return ""
}

func (m *Params) GetConfirmationNumber() uint64 {
	if m != nil {
		return m.ConfirmationNumber
	}
	return 0
}

func (m *Params) GetMinDepositAmount() uint64 {
	if m != nil {
		return m.MinDepositAmount
	}
	return 0
}

func (m *Params) GetDepositMagicPrefix() []byte {
	if m != nil {
		return m.DepositMagicPrefix
	}
	return nil
}

func (m *Params) GetDepositTaxRate() uint64 {
	if m != nil {
		return m.DepositTaxRate
	}
	return 0
}

func (m *Params) GetMaxDepositTax() uint64 {
	if m != nil {
		return m.MaxDepositTax
	}
	return 0
}

func init() {
	proto.RegisterType((*Params)(nil), "goat.bitcoin.v1.Params")
}

func init() { proto.RegisterFile("goat/bitcoin/v1/params.proto", fileDescriptor_96fcb82db75b9cd1) }

var fileDescriptor_96fcb82db75b9cd1 = []byte{
	// 330 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x90, 0xb1, 0x4e, 0x32, 0x41,
	0x14, 0x85, 0x59, 0x7e, 0x7e, 0x12, 0x47, 0x14, 0x1c, 0x35, 0xd9, 0x10, 0xb3, 0x41, 0x0b, 0xb3,
	0x51, 0xb3, 0x23, 0xb1, 0xb3, 0xd3, 0x60, 0x29, 0x21, 0x1b, 0x2b, 0x9b, 0xc9, 0xdd, 0x75, 0xc0,
	0x89, 0xb9, 0x33, 0x9b, 0x61, 0xc0, 0xf5, 0x15, 0xac, 0x7c, 0x0e, 0x2b, 0x1f, 0xc3, 0x92, 0xd2,
	0xd2, 0x40, 0xe1, 0x6b, 0x18, 0x86, 0x85, 0x10, 0x9b, 0x5b, 0x9c, 0xef, 0x3b, 0xb7, 0x38, 0xe4,
	0x60, 0xa0, 0xc1, 0xb2, 0x44, 0xda, 0x54, 0x4b, 0xc5, 0xc6, 0x6d, 0x96, 0x81, 0x01, 0x1c, 0x46,
	0x99, 0xd1, 0x56, 0xd3, 0xfa, 0x9c, 0x46, 0x05, 0x8d, 0xc6, 0xed, 0xe6, 0x0e, 0xa0, 0x54, 0x9a,
	0xb9, 0xbb, 0x70, 0x8e, 0xde, 0xcb, 0xa4, 0xda, 0x73, 0x25, 0x7a, 0x48, 0x6a, 0x4a, 0xd8, 0x67,
	0x6d, 0x9e, 0xb8, 0x02, 0x14, 0xbe, 0xd7, 0xf2, 0xc2, 0x8d, 0x78, 0xb3, 0xc8, 0xba, 0x80, 0x82,
	0x32, 0xb2, 0x9b, 0x6a, 0xd5, 0x97, 0x06, 0xc1, 0x4a, 0xad, 0xb8, 0x1a, 0x61, 0x22, 0x8c, 0x5f,
	0x6e, 0x79, 0x61, 0x25, 0xa6, 0xeb, 0xa8, 0xeb, 0x08, 0x3d, 0x23, 0x14, 0xa5, 0xe2, 0x0f, 0x22,
	0xd3, 0x43, 0x69, 0x39, 0xa0, 0x1e, 0x29, 0xeb, 0xff, 0x73, 0x7e, 0x03, 0xa5, 0xea, 0x2c, 0xc0,
	0x95, 0xcb, 0xe9, 0x39, 0xd9, 0x5b, 0x9a, 0x08, 0x03, 0x99, 0xf2, 0xcc, 0x88, 0xbe, 0xcc, 0xfd,
	0x4a, 0xcb, 0x0b, 0x6b, 0x31, 0x2d, 0xd8, 0xed, 0x1c, 0xf5, 0x1c, 0xa1, 0x21, 0x69, 0x2c, 0x1b,
	0x16, 0x72, 0x6e, 0xc0, 0x0a, 0xff, 0xbf, 0xfb, 0xbe, 0x5d, 0xe4, 0x77, 0x90, 0xc7, 0x60, 0x05,
	0x3d, 0x26, 0x75, 0x84, 0x9c, 0xaf, 0xd9, 0x7e, 0xd5, 0x89, 0x5b, 0x08, 0x79, 0x67, 0xe5, 0x5e,
	0x36, 0x5f, 0x7f, 0x3e, 0x4e, 0xf6, 0xdd, 0xae, 0xf9, 0x6a, 0xd9, 0xc5, 0x42, 0xd7, 0x37, 0x9f,
	0xd3, 0xc0, 0x9b, 0x4c, 0x03, 0xef, 0x7b, 0x1a, 0x78, 0x6f, 0xb3, 0xa0, 0x34, 0x99, 0x05, 0xa5,
	0xaf, 0x59, 0x50, 0xba, 0x3f, 0x1d, 0x48, 0xfb, 0x38, 0x4a, 0xa2, 0x54, 0x23, 0x9b, 0x77, 0x8b,
	0xd1, 0xd8, 0x9f, 0x3f, 0xf6, 0x25, 0x13, 0xc3, 0xa4, 0xea, 0xa6, 0xbf, 0xf8, 0x0d, 0x00, 0x00,
	0xff, 0xff, 0x9f, 0x68, 0xd2, 0x56, 0xbe, 0x01, 0x00, 0x00,
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
	if m.MaxDepositTax != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.MaxDepositTax))
		i--
		dAtA[i] = 0x30
	}
	if m.DepositTaxRate != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.DepositTaxRate))
		i--
		dAtA[i] = 0x28
	}
	if len(m.DepositMagicPrefix) > 0 {
		i -= len(m.DepositMagicPrefix)
		copy(dAtA[i:], m.DepositMagicPrefix)
		i = encodeVarintParams(dAtA, i, uint64(len(m.DepositMagicPrefix)))
		i--
		dAtA[i] = 0x22
	}
	if m.MinDepositAmount != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.MinDepositAmount))
		i--
		dAtA[i] = 0x18
	}
	if m.ConfirmationNumber != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.ConfirmationNumber))
		i--
		dAtA[i] = 0x10
	}
	if len(m.NetworkName) > 0 {
		i -= len(m.NetworkName)
		copy(dAtA[i:], m.NetworkName)
		i = encodeVarintParams(dAtA, i, uint64(len(m.NetworkName)))
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
	l = len(m.NetworkName)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	if m.ConfirmationNumber != 0 {
		n += 1 + sovParams(uint64(m.ConfirmationNumber))
	}
	if m.MinDepositAmount != 0 {
		n += 1 + sovParams(uint64(m.MinDepositAmount))
	}
	l = len(m.DepositMagicPrefix)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	if m.DepositTaxRate != 0 {
		n += 1 + sovParams(uint64(m.DepositTaxRate))
	}
	if m.MaxDepositTax != 0 {
		n += 1 + sovParams(uint64(m.MaxDepositTax))
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
				return fmt.Errorf("proto: wrong wireType = %d for field NetworkName", wireType)
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
			m.NetworkName = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ConfirmationNumber", wireType)
			}
			m.ConfirmationNumber = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ConfirmationNumber |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinDepositAmount", wireType)
			}
			m.MinDepositAmount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MinDepositAmount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DepositMagicPrefix", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DepositMagicPrefix = append(m.DepositMagicPrefix[:0], dAtA[iNdEx:postIndex]...)
			if m.DepositMagicPrefix == nil {
				m.DepositMagicPrefix = []byte{}
			}
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field DepositTaxRate", wireType)
			}
			m.DepositTaxRate = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.DepositTaxRate |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxDepositTax", wireType)
			}
			m.MaxDepositTax = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MaxDepositTax |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
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
