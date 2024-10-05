// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: goat/locking/v1/params.proto

package types

import (
	cosmossdk_io_math "cosmossdk.io/math"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	_ "github.com/cosmos/cosmos-sdk/types/tx/amino"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	github_com_cosmos_gogoproto_types "github.com/cosmos/gogoproto/types"
	_ "google.golang.org/protobuf/types/known/durationpb"
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
	// the partial unlock duation
	UnlockDuration time.Duration `protobuf:"bytes,1,opt,name=unlock_duration,json=unlockDuration,proto3,stdduration" json:"unlock_duration"`
	// if the token amount is less then threshold, the validator will be the
	// inactive status the validator operator should wait for long then paritial
	// unlock duation
	ExitingDuration      time.Duration `protobuf:"bytes,2,opt,name=exiting_duration,json=exitingDuration,proto3,stdduration" json:"exiting_duration"`
	DowntimeJailDuration time.Duration `protobuf:"bytes,3,opt,name=downtime_jail_duration,json=downtimeJailDuration,proto3,stdduration" json:"downtime_jail_duration"`
	// max_validators is the maximum number of validators.
	MaxValidators           int64                       `protobuf:"varint,4,opt,name=max_validators,json=maxValidators,proto3" json:"max_validators,omitempty"`
	SignedBlocksWindow      int64                       `protobuf:"varint,5,opt,name=signed_blocks_window,json=signedBlocksWindow,proto3" json:"signed_blocks_window,omitempty"`
	MaxMissedPerWindow      int64                       `protobuf:"varint,6,opt,name=max_missed_per_window,json=maxMissedPerWindow,proto3" json:"max_missed_per_window,omitempty"`
	SlashFractionDoubleSign cosmossdk_io_math.LegacyDec `protobuf:"bytes,7,opt,name=slash_fraction_double_sign,json=slashFractionDoubleSign,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"slash_fraction_double_sign"`
	SlashFractionDowntime   cosmossdk_io_math.LegacyDec `protobuf:"bytes,8,opt,name=slash_fraction_downtime,json=slashFractionDowntime,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"slash_fraction_downtime"`
}

func (m *Params) Reset()         { *m = Params{} }
func (m *Params) String() string { return proto.CompactTextString(m) }
func (*Params) ProtoMessage()    {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_a94ad1e7519f5b55, []int{0}
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

func (m *Params) GetUnlockDuration() time.Duration {
	if m != nil {
		return m.UnlockDuration
	}
	return 0
}

func (m *Params) GetExitingDuration() time.Duration {
	if m != nil {
		return m.ExitingDuration
	}
	return 0
}

func (m *Params) GetDowntimeJailDuration() time.Duration {
	if m != nil {
		return m.DowntimeJailDuration
	}
	return 0
}

func (m *Params) GetMaxValidators() int64 {
	if m != nil {
		return m.MaxValidators
	}
	return 0
}

func (m *Params) GetSignedBlocksWindow() int64 {
	if m != nil {
		return m.SignedBlocksWindow
	}
	return 0
}

func (m *Params) GetMaxMissedPerWindow() int64 {
	if m != nil {
		return m.MaxMissedPerWindow
	}
	return 0
}

func init() {
	proto.RegisterType((*Params)(nil), "goat.locking.v1.Params")
}

func init() { proto.RegisterFile("goat/locking/v1/params.proto", fileDescriptor_a94ad1e7519f5b55) }

var fileDescriptor_a94ad1e7519f5b55 = []byte{
	// 491 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x93, 0xc1, 0x6e, 0x13, 0x3d,
	0x10, 0xc7, 0xb3, 0x5f, 0xbf, 0x06, 0x64, 0x68, 0x03, 0xab, 0x94, 0xa6, 0x01, 0x6d, 0x22, 0x24,
	0xa4, 0x08, 0xc4, 0x9a, 0x80, 0xc4, 0x03, 0x44, 0x81, 0x03, 0x02, 0xa9, 0xb4, 0x12, 0x48, 0x1c,
	0xb0, 0xbc, 0xbb, 0xae, 0x63, 0xb2, 0xf6, 0x44, 0x6b, 0x6f, 0xb2, 0x3d, 0xf2, 0x06, 0x1c, 0x79,
	0x04, 0x8e, 0x1c, 0x78, 0x88, 0x1e, 0x2b, 0x4e, 0x88, 0x43, 0x41, 0xc9, 0x81, 0xd7, 0x40, 0x6b,
	0xef, 0xa6, 0x12, 0x9c, 0x2a, 0x2e, 0x2b, 0xcf, 0xfc, 0x67, 0x7e, 0x7f, 0xcf, 0xda, 0x46, 0xb7,
	0x38, 0x50, 0x83, 0x53, 0x88, 0xa7, 0x42, 0x71, 0x3c, 0x1f, 0xe2, 0x19, 0xcd, 0xa8, 0xd4, 0xe1,
	0x2c, 0x03, 0x03, 0x7e, 0xab, 0x54, 0xc3, 0x4a, 0x0d, 0xe7, 0xc3, 0xee, 0x75, 0x2a, 0x85, 0x02,
	0x6c, 0xbf, 0xae, 0xa6, 0xbb, 0x17, 0x83, 0x96, 0xa0, 0x89, 0x8d, 0xb0, 0x0b, 0x2a, 0xa9, 0xcd,
	0x81, 0x83, 0xcb, 0x97, 0xab, 0x2a, 0x1b, 0x70, 0x00, 0x9e, 0x32, 0x6c, 0xa3, 0x28, 0x3f, 0xc2,
	0x49, 0x9e, 0x51, 0x23, 0x40, 0x39, 0xfd, 0xf6, 0xfb, 0x4d, 0xd4, 0xdc, 0xb7, 0xbb, 0xf0, 0x5f,
	0xa2, 0x56, 0xae, 0x4a, 0x7b, 0x52, 0xd7, 0x74, 0xbc, 0xbe, 0x37, 0xb8, 0xf2, 0x70, 0x2f, 0x74,
	0x90, 0xb0, 0x86, 0x84, 0xe3, 0xaa, 0x60, 0xb4, 0x75, 0x72, 0xd6, 0x6b, 0x7c, 0xfc, 0xd1, 0xf3,
	0x3e, 0xfd, 0xfa, 0x7c, 0xd7, 0x3b, 0xd8, 0x76, 0x80, 0x5a, 0xf6, 0x0f, 0xd1, 0x35, 0x56, 0x08,
	0x23, 0x14, 0x3f, 0x67, 0xfe, 0x77, 0x41, 0x66, 0xab, 0x22, 0xac, 0xa1, 0x6f, 0xd1, 0x8d, 0x04,
	0x16, 0xca, 0x08, 0xc9, 0xc8, 0x3b, 0x2a, 0xd2, 0x73, 0xf4, 0xc6, 0x05, 0xd1, 0xed, 0x9a, 0xf3,
	0x8c, 0x8a, 0x74, 0xcd, 0xbf, 0x83, 0xb6, 0x25, 0x2d, 0xc8, 0x9c, 0xa6, 0x22, 0xa1, 0x06, 0x32,
	0xdd, 0xf9, 0xbf, 0xef, 0x0d, 0x36, 0x0e, 0xb6, 0x24, 0x2d, 0x5e, 0xad, 0x93, 0xfe, 0x03, 0xd4,
	0xd6, 0x82, 0x2b, 0x96, 0x90, 0xa8, 0x9c, 0x59, 0x93, 0x85, 0x50, 0x09, 0x2c, 0x3a, 0x9b, 0xb6,
	0xd8, 0x77, 0xda, 0xc8, 0x4a, 0xaf, 0xad, 0xe2, 0x0f, 0xd1, 0x4e, 0x09, 0x96, 0x42, 0x6b, 0x96,
	0x90, 0x19, 0xcb, 0xea, 0x96, 0xa6, 0x6b, 0x91, 0xb4, 0x78, 0x61, 0xb5, 0x7d, 0x96, 0x55, 0x2d,
	0x1a, 0x75, 0x75, 0x4a, 0xf5, 0x84, 0x1c, 0x65, 0x34, 0x2e, 0x77, 0x47, 0x12, 0xc8, 0xa3, 0x94,
	0x91, 0x12, 0xdf, 0xb9, 0xd4, 0xf7, 0x06, 0x57, 0x47, 0x8f, 0xcb, 0xa1, 0xbe, 0x9f, 0xf5, 0x6e,
	0xba, 0xeb, 0xa0, 0x93, 0x69, 0x28, 0x00, 0x4b, 0x6a, 0x26, 0xe1, 0x73, 0xc6, 0x69, 0x7c, 0x3c,
	0x66, 0xf1, 0xd7, 0x2f, 0xf7, 0x51, 0x75, 0x5b, 0xc6, 0x2c, 0x76, 0xd3, 0xef, 0x5a, 0xf2, 0xd3,
	0x0a, 0x3c, 0xb6, 0xdc, 0x43, 0xc1, 0x95, 0xaf, 0xd0, 0xee, 0x5f, 0xa6, 0xee, 0x3f, 0x75, 0x2e,
	0xff, 0x93, 0xe3, 0xce, 0x1f, 0x8e, 0x0e, 0x3a, 0x7a, 0x72, 0xb2, 0x0c, 0xbc, 0xd3, 0x65, 0xe0,
	0xfd, 0x5c, 0x06, 0xde, 0x87, 0x55, 0xd0, 0x38, 0x5d, 0x05, 0x8d, 0x6f, 0xab, 0xa0, 0xf1, 0xe6,
	0x1e, 0x17, 0x66, 0x92, 0x47, 0x61, 0x0c, 0x12, 0x97, 0xaf, 0x43, 0x31, 0xb3, 0x80, 0x6c, 0x6a,
	0xd7, 0xb8, 0x58, 0xbf, 0x24, 0x73, 0x3c, 0x63, 0x3a, 0x6a, 0xda, 0xf3, 0x7e, 0xf4, 0x3b, 0x00,
	0x00, 0xff, 0xff, 0xc3, 0x87, 0xb6, 0x2a, 0x66, 0x03, 0x00, 0x00,
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
	{
		size := m.SlashFractionDowntime.Size()
		i -= size
		if _, err := m.SlashFractionDowntime.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x42
	{
		size := m.SlashFractionDoubleSign.Size()
		i -= size
		if _, err := m.SlashFractionDoubleSign.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x3a
	if m.MaxMissedPerWindow != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.MaxMissedPerWindow))
		i--
		dAtA[i] = 0x30
	}
	if m.SignedBlocksWindow != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.SignedBlocksWindow))
		i--
		dAtA[i] = 0x28
	}
	if m.MaxValidators != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.MaxValidators))
		i--
		dAtA[i] = 0x20
	}
	n1, err1 := github_com_cosmos_gogoproto_types.StdDurationMarshalTo(m.DowntimeJailDuration, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdDuration(m.DowntimeJailDuration):])
	if err1 != nil {
		return 0, err1
	}
	i -= n1
	i = encodeVarintParams(dAtA, i, uint64(n1))
	i--
	dAtA[i] = 0x1a
	n2, err2 := github_com_cosmos_gogoproto_types.StdDurationMarshalTo(m.ExitingDuration, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdDuration(m.ExitingDuration):])
	if err2 != nil {
		return 0, err2
	}
	i -= n2
	i = encodeVarintParams(dAtA, i, uint64(n2))
	i--
	dAtA[i] = 0x12
	n3, err3 := github_com_cosmos_gogoproto_types.StdDurationMarshalTo(m.UnlockDuration, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdDuration(m.UnlockDuration):])
	if err3 != nil {
		return 0, err3
	}
	i -= n3
	i = encodeVarintParams(dAtA, i, uint64(n3))
	i--
	dAtA[i] = 0xa
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
	l = github_com_cosmos_gogoproto_types.SizeOfStdDuration(m.UnlockDuration)
	n += 1 + l + sovParams(uint64(l))
	l = github_com_cosmos_gogoproto_types.SizeOfStdDuration(m.ExitingDuration)
	n += 1 + l + sovParams(uint64(l))
	l = github_com_cosmos_gogoproto_types.SizeOfStdDuration(m.DowntimeJailDuration)
	n += 1 + l + sovParams(uint64(l))
	if m.MaxValidators != 0 {
		n += 1 + sovParams(uint64(m.MaxValidators))
	}
	if m.SignedBlocksWindow != 0 {
		n += 1 + sovParams(uint64(m.SignedBlocksWindow))
	}
	if m.MaxMissedPerWindow != 0 {
		n += 1 + sovParams(uint64(m.MaxMissedPerWindow))
	}
	l = m.SlashFractionDoubleSign.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.SlashFractionDowntime.Size()
	n += 1 + l + sovParams(uint64(l))
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
				return fmt.Errorf("proto: wrong wireType = %d for field UnlockDuration", wireType)
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
			if err := github_com_cosmos_gogoproto_types.StdDurationUnmarshal(&m.UnlockDuration, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExitingDuration", wireType)
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
			if err := github_com_cosmos_gogoproto_types.StdDurationUnmarshal(&m.ExitingDuration, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DowntimeJailDuration", wireType)
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
			if err := github_com_cosmos_gogoproto_types.StdDurationUnmarshal(&m.DowntimeJailDuration, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxValidators", wireType)
			}
			m.MaxValidators = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MaxValidators |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SignedBlocksWindow", wireType)
			}
			m.SignedBlocksWindow = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SignedBlocksWindow |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxMissedPerWindow", wireType)
			}
			m.MaxMissedPerWindow = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MaxMissedPerWindow |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SlashFractionDoubleSign", wireType)
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
			if err := m.SlashFractionDoubleSign.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SlashFractionDowntime", wireType)
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
			if err := m.SlashFractionDowntime.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
