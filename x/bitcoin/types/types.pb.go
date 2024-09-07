// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: goat/bitcoin/v1/types.proto

package types

import (
	fmt "fmt"
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

// WithdrawalIds
type WithdrawalIds struct {
	Id []uint64 `protobuf:"varint,1,rep,packed,name=id,proto3" json:"id,omitempty"`
}

func (m *WithdrawalIds) Reset()         { *m = WithdrawalIds{} }
func (m *WithdrawalIds) String() string { return proto.CompactTextString(m) }
func (*WithdrawalIds) ProtoMessage()    {}
func (*WithdrawalIds) Descriptor() ([]byte, []int) {
	return fileDescriptor_71f9ee0c92692d26, []int{0}
}
func (m *WithdrawalIds) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *WithdrawalIds) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_WithdrawalIds.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *WithdrawalIds) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WithdrawalIds.Merge(m, src)
}
func (m *WithdrawalIds) XXX_Size() int {
	return m.Size()
}
func (m *WithdrawalIds) XXX_DiscardUnknown() {
	xxx_messageInfo_WithdrawalIds.DiscardUnknown(m)
}

var xxx_messageInfo_WithdrawalIds proto.InternalMessageInfo

func (m *WithdrawalIds) GetId() []uint64 {
	if m != nil {
		return m.Id
	}
	return nil
}

// DepositReceipt
type DepositExecReceipt struct {
	Txid    []byte `protobuf:"bytes,1,opt,name=txid,proto3" json:"txid,omitempty"`
	Txout   uint32 `protobuf:"varint,2,opt,name=txout,proto3" json:"txout,omitempty"`
	Address []byte `protobuf:"bytes,3,opt,name=address,proto3" json:"address,omitempty"`
	Amount  uint64 `protobuf:"varint,4,opt,name=amount,proto3" json:"amount,omitempty"`
}

func (m *DepositExecReceipt) Reset()         { *m = DepositExecReceipt{} }
func (m *DepositExecReceipt) String() string { return proto.CompactTextString(m) }
func (*DepositExecReceipt) ProtoMessage()    {}
func (*DepositExecReceipt) Descriptor() ([]byte, []int) {
	return fileDescriptor_71f9ee0c92692d26, []int{1}
}
func (m *DepositExecReceipt) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *DepositExecReceipt) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_DepositExecReceipt.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *DepositExecReceipt) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DepositExecReceipt.Merge(m, src)
}
func (m *DepositExecReceipt) XXX_Size() int {
	return m.Size()
}
func (m *DepositExecReceipt) XXX_DiscardUnknown() {
	xxx_messageInfo_DepositExecReceipt.DiscardUnknown(m)
}

var xxx_messageInfo_DepositExecReceipt proto.InternalMessageInfo

func (m *DepositExecReceipt) GetTxid() []byte {
	if m != nil {
		return m.Txid
	}
	return nil
}

func (m *DepositExecReceipt) GetTxout() uint32 {
	if m != nil {
		return m.Txout
	}
	return 0
}

func (m *DepositExecReceipt) GetAddress() []byte {
	if m != nil {
		return m.Address
	}
	return nil
}

func (m *DepositExecReceipt) GetAmount() uint64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

// WithdrawalExecReceipt
type WithdrawalExecReceipt struct {
	Id      uint64             `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Receipt *WithdrawalReceipt `protobuf:"bytes,2,opt,name=receipt,proto3" json:"receipt,omitempty"`
}

func (m *WithdrawalExecReceipt) Reset()         { *m = WithdrawalExecReceipt{} }
func (m *WithdrawalExecReceipt) String() string { return proto.CompactTextString(m) }
func (*WithdrawalExecReceipt) ProtoMessage()    {}
func (*WithdrawalExecReceipt) Descriptor() ([]byte, []int) {
	return fileDescriptor_71f9ee0c92692d26, []int{2}
}
func (m *WithdrawalExecReceipt) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *WithdrawalExecReceipt) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_WithdrawalExecReceipt.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *WithdrawalExecReceipt) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WithdrawalExecReceipt.Merge(m, src)
}
func (m *WithdrawalExecReceipt) XXX_Size() int {
	return m.Size()
}
func (m *WithdrawalExecReceipt) XXX_DiscardUnknown() {
	xxx_messageInfo_WithdrawalExecReceipt.DiscardUnknown(m)
}

var xxx_messageInfo_WithdrawalExecReceipt proto.InternalMessageInfo

func (m *WithdrawalExecReceipt) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *WithdrawalExecReceipt) GetReceipt() *WithdrawalReceipt {
	if m != nil {
		return m.Receipt
	}
	return nil
}

// ExecuableQueue
type ExecuableQueue struct {
	BlockNumber         uint64                   `protobuf:"varint,1,opt,name=block_number,json=blockNumber,proto3" json:"block_number,omitempty"`
	Deposits            []*DepositExecReceipt    `protobuf:"bytes,2,rep,name=deposits,proto3" json:"deposits,omitempty"`
	PaidWithdrawals     []*WithdrawalExecReceipt `protobuf:"bytes,3,rep,name=paid_withdrawals,json=paidWithdrawals,proto3" json:"paid_withdrawals,omitempty"`
	RejectedWithdrawals []uint64                 `protobuf:"varint,4,rep,packed,name=rejected_withdrawals,json=rejectedWithdrawals,proto3" json:"rejected_withdrawals,omitempty"`
}

func (m *ExecuableQueue) Reset()         { *m = ExecuableQueue{} }
func (m *ExecuableQueue) String() string { return proto.CompactTextString(m) }
func (*ExecuableQueue) ProtoMessage()    {}
func (*ExecuableQueue) Descriptor() ([]byte, []int) {
	return fileDescriptor_71f9ee0c92692d26, []int{3}
}
func (m *ExecuableQueue) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ExecuableQueue) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ExecuableQueue.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ExecuableQueue) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExecuableQueue.Merge(m, src)
}
func (m *ExecuableQueue) XXX_Size() int {
	return m.Size()
}
func (m *ExecuableQueue) XXX_DiscardUnknown() {
	xxx_messageInfo_ExecuableQueue.DiscardUnknown(m)
}

var xxx_messageInfo_ExecuableQueue proto.InternalMessageInfo

func (m *ExecuableQueue) GetBlockNumber() uint64 {
	if m != nil {
		return m.BlockNumber
	}
	return 0
}

func (m *ExecuableQueue) GetDeposits() []*DepositExecReceipt {
	if m != nil {
		return m.Deposits
	}
	return nil
}

func (m *ExecuableQueue) GetPaidWithdrawals() []*WithdrawalExecReceipt {
	if m != nil {
		return m.PaidWithdrawals
	}
	return nil
}

func (m *ExecuableQueue) GetRejectedWithdrawals() []uint64 {
	if m != nil {
		return m.RejectedWithdrawals
	}
	return nil
}

func init() {
	proto.RegisterType((*WithdrawalIds)(nil), "goat.bitcoin.v1.WithdrawalIds")
	proto.RegisterType((*DepositExecReceipt)(nil), "goat.bitcoin.v1.DepositExecReceipt")
	proto.RegisterType((*WithdrawalExecReceipt)(nil), "goat.bitcoin.v1.WithdrawalExecReceipt")
	proto.RegisterType((*ExecuableQueue)(nil), "goat.bitcoin.v1.ExecuableQueue")
}

func init() { proto.RegisterFile("goat/bitcoin/v1/types.proto", fileDescriptor_71f9ee0c92692d26) }

var fileDescriptor_71f9ee0c92692d26 = []byte{
	// 394 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x52, 0x4d, 0x4b, 0xf3, 0x40,
	0x10, 0xee, 0xa6, 0x79, 0xdb, 0x97, 0xed, 0x97, 0xac, 0x55, 0x82, 0x62, 0x8c, 0x11, 0x24, 0x20,
	0x24, 0xb4, 0x5e, 0x05, 0x41, 0xec, 0xc1, 0x8b, 0xd0, 0x5c, 0x04, 0x2f, 0x25, 0xc9, 0x2e, 0xed,
	0xda, 0x36, 0x1b, 0x92, 0x4d, 0x1b, 0xff, 0x85, 0x3f, 0xcb, 0x63, 0x8f, 0x1e, 0xa5, 0xfd, 0x17,
	0x9e, 0x24, 0x9b, 0xa4, 0xad, 0x2d, 0x78, 0xdb, 0x99, 0x67, 0x9e, 0x8f, 0x4c, 0x06, 0x9e, 0x0e,
	0x99, 0xc3, 0x2d, 0x97, 0x72, 0x8f, 0x51, 0xdf, 0x9a, 0x75, 0x2c, 0xfe, 0x16, 0x90, 0xc8, 0x0c,
	0x42, 0xc6, 0x19, 0x6a, 0xa5, 0xa0, 0x99, 0x83, 0xe6, 0xac, 0x73, 0x72, 0xb6, 0x3b, 0x5d, 0x60,
	0x62, 0x5e, 0x3f, 0x87, 0x8d, 0x67, 0xca, 0x47, 0x38, 0x74, 0xe6, 0xce, 0xe4, 0x11, 0x47, 0xa8,
	0x09, 0x25, 0x8a, 0x15, 0xa0, 0x95, 0x0d, 0xd9, 0x96, 0x28, 0xd6, 0x03, 0x88, 0x1e, 0x48, 0xc0,
	0x22, 0xca, 0x7b, 0x09, 0xf1, 0x6c, 0xe2, 0x11, 0x1a, 0x70, 0x84, 0xa0, 0xcc, 0x13, 0x31, 0x07,
	0x8c, 0xba, 0x2d, 0xde, 0xa8, 0x0d, 0xff, 0xf1, 0x84, 0xc5, 0x5c, 0x91, 0x34, 0x60, 0x34, 0xec,
	0xac, 0x40, 0x0a, 0xac, 0x3a, 0x18, 0x87, 0x24, 0x8a, 0x94, 0xb2, 0x18, 0x2e, 0x4a, 0x74, 0x0c,
	0x2b, 0xce, 0x94, 0xc5, 0x3e, 0x57, 0x64, 0x0d, 0x18, 0xb2, 0x9d, 0x57, 0x3a, 0x81, 0x47, 0x9b,
	0x48, 0xdb, 0xa6, 0x45, 0x34, 0x90, 0x45, 0x43, 0xb7, 0xb0, 0x1a, 0x66, 0x90, 0xb0, 0xac, 0x75,
	0x75, 0x73, 0xe7, 0xeb, 0xcd, 0x8d, 0x50, 0x2e, 0x62, 0x17, 0x14, 0xfd, 0x1b, 0xc0, 0x66, 0xaa,
	0x1e, 0x3b, 0xee, 0x84, 0xf4, 0x63, 0x12, 0x13, 0x74, 0x01, 0xeb, 0xee, 0x84, 0x79, 0xe3, 0x81,
	0x1f, 0x4f, 0x5d, 0x12, 0xe6, 0x56, 0x35, 0xd1, 0x7b, 0x12, 0x2d, 0x74, 0x07, 0xff, 0xe3, 0x6c,
	0x1d, 0x91, 0x22, 0x69, 0x65, 0xa3, 0xd6, 0xbd, 0xdc, 0x33, 0xdd, 0xdf, 0x97, 0xbd, 0x26, 0xa1,
	0x3e, 0x3c, 0x08, 0x1c, 0x8a, 0x07, 0xf3, 0x75, 0xb2, 0x74, 0x31, 0xa9, 0xd0, 0xd5, 0x1f, 0xe9,
	0xb7, 0xb5, 0x5a, 0x29, 0x7f, 0x03, 0x45, 0xa8, 0x03, 0xdb, 0x21, 0x79, 0x25, 0x1e, 0x27, 0xbf,
	0x65, 0x65, 0xf1, 0x13, 0x0f, 0x0b, 0x6c, 0x8b, 0x72, 0xdf, 0xfb, 0x58, 0xaa, 0x60, 0xb1, 0x54,
	0xc1, 0xd7, 0x52, 0x05, 0xef, 0x2b, 0xb5, 0xb4, 0x58, 0xa9, 0xa5, 0xcf, 0x95, 0x5a, 0x7a, 0xb9,
	0x1e, 0x52, 0x3e, 0x8a, 0x5d, 0xd3, 0x63, 0x53, 0x2b, 0xcd, 0xe3, 0x13, 0x3e, 0x67, 0xe1, 0x58,
	0xbc, 0xad, 0x64, 0x7d, 0x48, 0xe2, 0xe6, 0xdc, 0x8a, 0x38, 0xa2, 0x9b, 0x9f, 0x00, 0x00, 0x00,
	0xff, 0xff, 0xbc, 0xa4, 0x24, 0x23, 0x93, 0x02, 0x00, 0x00,
}

func (m *WithdrawalIds) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *WithdrawalIds) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *WithdrawalIds) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Id) > 0 {
		dAtA2 := make([]byte, len(m.Id)*10)
		var j1 int
		for _, num := range m.Id {
			for num >= 1<<7 {
				dAtA2[j1] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j1++
			}
			dAtA2[j1] = uint8(num)
			j1++
		}
		i -= j1
		copy(dAtA[i:], dAtA2[:j1])
		i = encodeVarintTypes(dAtA, i, uint64(j1))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *DepositExecReceipt) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DepositExecReceipt) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *DepositExecReceipt) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Amount != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.Amount))
		i--
		dAtA[i] = 0x20
	}
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0x1a
	}
	if m.Txout != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.Txout))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Txid) > 0 {
		i -= len(m.Txid)
		copy(dAtA[i:], m.Txid)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Txid)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *WithdrawalExecReceipt) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *WithdrawalExecReceipt) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *WithdrawalExecReceipt) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Receipt != nil {
		{
			size, err := m.Receipt.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTypes(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.Id != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.Id))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *ExecuableQueue) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ExecuableQueue) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ExecuableQueue) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.RejectedWithdrawals) > 0 {
		dAtA5 := make([]byte, len(m.RejectedWithdrawals)*10)
		var j4 int
		for _, num := range m.RejectedWithdrawals {
			for num >= 1<<7 {
				dAtA5[j4] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j4++
			}
			dAtA5[j4] = uint8(num)
			j4++
		}
		i -= j4
		copy(dAtA[i:], dAtA5[:j4])
		i = encodeVarintTypes(dAtA, i, uint64(j4))
		i--
		dAtA[i] = 0x22
	}
	if len(m.PaidWithdrawals) > 0 {
		for iNdEx := len(m.PaidWithdrawals) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.PaidWithdrawals[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTypes(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.Deposits) > 0 {
		for iNdEx := len(m.Deposits) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Deposits[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTypes(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if m.BlockNumber != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.BlockNumber))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintTypes(dAtA []byte, offset int, v uint64) int {
	offset -= sovTypes(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *WithdrawalIds) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Id) > 0 {
		l = 0
		for _, e := range m.Id {
			l += sovTypes(uint64(e))
		}
		n += 1 + sovTypes(uint64(l)) + l
	}
	return n
}

func (m *DepositExecReceipt) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Txid)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	if m.Txout != 0 {
		n += 1 + sovTypes(uint64(m.Txout))
	}
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	if m.Amount != 0 {
		n += 1 + sovTypes(uint64(m.Amount))
	}
	return n
}

func (m *WithdrawalExecReceipt) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Id != 0 {
		n += 1 + sovTypes(uint64(m.Id))
	}
	if m.Receipt != nil {
		l = m.Receipt.Size()
		n += 1 + l + sovTypes(uint64(l))
	}
	return n
}

func (m *ExecuableQueue) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.BlockNumber != 0 {
		n += 1 + sovTypes(uint64(m.BlockNumber))
	}
	if len(m.Deposits) > 0 {
		for _, e := range m.Deposits {
			l = e.Size()
			n += 1 + l + sovTypes(uint64(l))
		}
	}
	if len(m.PaidWithdrawals) > 0 {
		for _, e := range m.PaidWithdrawals {
			l = e.Size()
			n += 1 + l + sovTypes(uint64(l))
		}
	}
	if len(m.RejectedWithdrawals) > 0 {
		l = 0
		for _, e := range m.RejectedWithdrawals {
			l += sovTypes(uint64(e))
		}
		n += 1 + sovTypes(uint64(l)) + l
	}
	return n
}

func sovTypes(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTypes(x uint64) (n int) {
	return sovTypes(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *WithdrawalIds) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
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
			return fmt.Errorf("proto: WithdrawalIds: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: WithdrawalIds: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType == 0 {
				var v uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowTypes
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.Id = append(m.Id, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowTypes
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthTypes
				}
				postIndex := iNdEx + packedLen
				if postIndex < 0 {
					return ErrInvalidLengthTypes
				}
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				var count int
				for _, integer := range dAtA[iNdEx:postIndex] {
					if integer < 128 {
						count++
					}
				}
				elementCount = count
				if elementCount != 0 && len(m.Id) == 0 {
					m.Id = make([]uint64, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowTypes
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.Id = append(m.Id, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTypes
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
func (m *DepositExecReceipt) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
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
			return fmt.Errorf("proto: DepositExecReceipt: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DepositExecReceipt: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Txid", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Txid = append(m.Txid[:0], dAtA[iNdEx:postIndex]...)
			if m.Txid == nil {
				m.Txid = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Txout", wireType)
			}
			m.Txout = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Txout |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = append(m.Address[:0], dAtA[iNdEx:postIndex]...)
			if m.Address == nil {
				m.Address = []byte{}
			}
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			m.Amount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Amount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTypes
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
func (m *WithdrawalExecReceipt) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
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
			return fmt.Errorf("proto: WithdrawalExecReceipt: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: WithdrawalExecReceipt: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Id |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Receipt", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Receipt == nil {
				m.Receipt = &WithdrawalReceipt{}
			}
			if err := m.Receipt.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTypes
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
func (m *ExecuableQueue) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
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
			return fmt.Errorf("proto: ExecuableQueue: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ExecuableQueue: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field BlockNumber", wireType)
			}
			m.BlockNumber = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.BlockNumber |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Deposits", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Deposits = append(m.Deposits, &DepositExecReceipt{})
			if err := m.Deposits[len(m.Deposits)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PaidWithdrawals", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PaidWithdrawals = append(m.PaidWithdrawals, &WithdrawalExecReceipt{})
			if err := m.PaidWithdrawals[len(m.PaidWithdrawals)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType == 0 {
				var v uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowTypes
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.RejectedWithdrawals = append(m.RejectedWithdrawals, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowTypes
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthTypes
				}
				postIndex := iNdEx + packedLen
				if postIndex < 0 {
					return ErrInvalidLengthTypes
				}
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				var count int
				for _, integer := range dAtA[iNdEx:postIndex] {
					if integer < 128 {
						count++
					}
				}
				elementCount = count
				if elementCount != 0 && len(m.RejectedWithdrawals) == 0 {
					m.RejectedWithdrawals = make([]uint64, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowTypes
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.RejectedWithdrawals = append(m.RejectedWithdrawals, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field RejectedWithdrawals", wireType)
			}
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTypes
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
func skipTypes(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTypes
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
					return 0, ErrIntOverflowTypes
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
					return 0, ErrIntOverflowTypes
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
				return 0, ErrInvalidLengthTypes
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTypes
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTypes
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTypes        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTypes          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTypes = fmt.Errorf("proto: unexpected end of group")
)