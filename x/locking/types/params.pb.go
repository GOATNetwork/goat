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
	HalvingInterval         int64                       `protobuf:"varint,9,opt,name=halving_interval,json=halvingInterval,proto3" json:"halving_interval,omitempty"`
	InitialBlockReward      int64                       `protobuf:"varint,10,opt,name=initial_block_reward,json=initialBlockReward,proto3" json:"initial_block_reward,omitempty"`
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

func (m *Params) GetHalvingInterval() int64 {
	if m != nil {
		return m.HalvingInterval
	}
	return 0
}

func (m *Params) GetInitialBlockReward() int64 {
	if m != nil {
		return m.InitialBlockReward
	}
	return 0
}

func init() {
	proto.RegisterType((*Params)(nil), "goat.locking.v1.Params")
}

func init() { proto.RegisterFile("goat/locking/v1/params.proto", fileDescriptor_a94ad1e7519f5b55) }

var fileDescriptor_a94ad1e7519f5b55 = []byte{
	// 534 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x93, 0xcd, 0x6e, 0x13, 0x31,
	0x10, 0xc7, 0xb3, 0x94, 0x06, 0x30, 0xb4, 0x29, 0xab, 0x94, 0x6e, 0x03, 0xda, 0x44, 0x48, 0x48,
	0x01, 0xc4, 0x2e, 0x01, 0x89, 0x07, 0x88, 0x02, 0x12, 0x08, 0xa4, 0x92, 0x4a, 0x20, 0x71, 0xc0,
	0x72, 0x76, 0x5d, 0xc7, 0x64, 0x6d, 0x47, 0xb6, 0xf3, 0xd1, 0xb7, 0xe0, 0xc8, 0x23, 0x70, 0xe4,
	0xc0, 0x43, 0xf4, 0x58, 0x71, 0x42, 0x1c, 0x02, 0x4a, 0x0e, 0xbc, 0x06, 0xf2, 0xc7, 0xa6, 0x12,
	0x9c, 0x2a, 0x2e, 0xab, 0xf5, 0xfc, 0x67, 0x7e, 0xe3, 0xf9, 0x6b, 0x0c, 0x6e, 0x11, 0x81, 0x74,
	0x5a, 0x88, 0x6c, 0x44, 0x39, 0x49, 0xa7, 0x9d, 0x74, 0x8c, 0x24, 0x62, 0x2a, 0x19, 0x4b, 0xa1,
	0x45, 0x58, 0x33, 0x6a, 0xe2, 0xd5, 0x64, 0xda, 0x69, 0x5c, 0x47, 0x8c, 0x72, 0x91, 0xda, 0xaf,
	0xcb, 0x69, 0xec, 0x67, 0x42, 0x31, 0xa1, 0xa0, 0x3d, 0xa5, 0xee, 0xe0, 0xa5, 0x3a, 0x11, 0x44,
	0xb8, 0xb8, 0xf9, 0xf3, 0xd1, 0x98, 0x08, 0x41, 0x0a, 0x9c, 0xda, 0xd3, 0x60, 0x72, 0x94, 0xe6,
	0x13, 0x89, 0x34, 0x15, 0xdc, 0xe9, 0xb7, 0x17, 0x9b, 0xa0, 0x7a, 0x60, 0x6f, 0x11, 0xbe, 0x06,
	0xb5, 0x09, 0x37, 0xed, 0x61, 0x99, 0x13, 0x05, 0xad, 0xa0, 0x7d, 0xf5, 0xd1, 0x7e, 0xe2, 0x20,
	0x49, 0x09, 0x49, 0x7a, 0x3e, 0xa1, 0xbb, 0x75, 0xb2, 0x68, 0x56, 0x3e, 0xfd, 0x6c, 0x06, 0x9f,
	0x7f, 0x7f, 0xb9, 0x17, 0xf4, 0xb7, 0x1d, 0xa0, 0x94, 0xc3, 0x43, 0xb0, 0x83, 0xe7, 0x54, 0x53,
	0x4e, 0xce, 0x98, 0x17, 0xce, 0xc9, 0xac, 0x79, 0xc2, 0x1a, 0xfa, 0x1e, 0xdc, 0xc8, 0xc5, 0x8c,
	0x6b, 0xca, 0x30, 0xfc, 0x80, 0x68, 0x71, 0x86, 0xde, 0x38, 0x27, 0xba, 0x5e, 0x72, 0x5e, 0x20,
	0x5a, 0xac, 0xf9, 0x77, 0xc0, 0x36, 0x43, 0x73, 0x38, 0x45, 0x05, 0xcd, 0x91, 0x16, 0x52, 0x45,
	0x17, 0x5b, 0x41, 0x7b, 0xa3, 0xbf, 0xc5, 0xd0, 0xfc, 0xcd, 0x3a, 0x18, 0x3e, 0x04, 0x75, 0x45,
	0x09, 0xc7, 0x39, 0x1c, 0x98, 0x99, 0x15, 0x9c, 0x51, 0x9e, 0x8b, 0x59, 0xb4, 0x69, 0x93, 0x43,
	0xa7, 0x75, 0xad, 0xf4, 0xd6, 0x2a, 0x61, 0x07, 0xec, 0x1a, 0x30, 0xa3, 0x4a, 0xe1, 0x1c, 0x8e,
	0xb1, 0x2c, 0x4b, 0xaa, 0xae, 0x84, 0xa1, 0xf9, 0x2b, 0xab, 0x1d, 0x60, 0xe9, 0x4b, 0x14, 0x68,
	0xa8, 0x02, 0xa9, 0x21, 0x3c, 0x92, 0x28, 0x33, 0xb7, 0x83, 0xb9, 0x98, 0x0c, 0x0a, 0x0c, 0x0d,
	0x3e, 0xba, 0xd4, 0x0a, 0xda, 0xd7, 0xba, 0x4f, 0xcc, 0x50, 0x3f, 0x16, 0xcd, 0x9b, 0x6e, 0x1d,
	0x54, 0x3e, 0x4a, 0xa8, 0x48, 0x19, 0xd2, 0xc3, 0xe4, 0x25, 0x26, 0x28, 0x3b, 0xee, 0xe1, 0xec,
	0xdb, 0xd7, 0x07, 0xc0, 0x6f, 0x4b, 0x0f, 0x67, 0x6e, 0xfa, 0x3d, 0x4b, 0x7e, 0xe6, 0xc1, 0x3d,
	0xcb, 0x3d, 0xa4, 0x84, 0x87, 0x1c, 0xec, 0xfd, 0xd3, 0xd4, 0xf9, 0x14, 0x5d, 0xfe, 0xaf, 0x8e,
	0xbb, 0x7f, 0x75, 0x74, 0xd0, 0xf0, 0x2e, 0xd8, 0x19, 0xa2, 0x62, 0x6a, 0xb6, 0x84, 0x72, 0x8d,
	0xe5, 0x14, 0x15, 0xd1, 0x15, 0x6b, 0x49, 0xcd, 0xc7, 0x9f, 0xfb, 0xb0, 0x31, 0x9d, 0x72, 0xaa,
	0x29, 0x2a, 0x9c, 0xeb, 0x50, 0xe2, 0x19, 0x92, 0x79, 0x04, 0x9c, 0x83, 0x5e, 0xb3, 0xae, 0xf7,
	0xad, 0xd2, 0x7d, 0x7a, 0xb2, 0x8c, 0x83, 0xd3, 0x65, 0x1c, 0xfc, 0x5a, 0xc6, 0xc1, 0xc7, 0x55,
	0x5c, 0x39, 0x5d, 0xc5, 0x95, 0xef, 0xab, 0xb8, 0xf2, 0xee, 0x3e, 0xa1, 0x7a, 0x38, 0x19, 0x24,
	0x99, 0x60, 0xa9, 0x79, 0x7a, 0x1c, 0xeb, 0x99, 0x90, 0x23, 0xfb, 0x9f, 0xce, 0xd7, 0xcf, 0x54,
	0x1f, 0x8f, 0xb1, 0x1a, 0x54, 0xed, 0x32, 0x3d, 0xfe, 0x13, 0x00, 0x00, 0xff, 0xff, 0xa1, 0x19,
	0xad, 0x11, 0xc3, 0x03, 0x00, 0x00,
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
	if m.InitialBlockReward != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.InitialBlockReward))
		i--
		dAtA[i] = 0x50
	}
	if m.HalvingInterval != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.HalvingInterval))
		i--
		dAtA[i] = 0x48
	}
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
	if m.HalvingInterval != 0 {
		n += 1 + sovParams(uint64(m.HalvingInterval))
	}
	if m.InitialBlockReward != 0 {
		n += 1 + sovParams(uint64(m.InitialBlockReward))
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
		case 9:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field HalvingInterval", wireType)
			}
			m.HalvingInterval = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.HalvingInterval |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 10:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field InitialBlockReward", wireType)
			}
			m.InitialBlockReward = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.InitialBlockReward |= int64(b&0x7F) << shift
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
