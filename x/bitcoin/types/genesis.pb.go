// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: goat/bitcoin/v1/genesis.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/types/tx/amino"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	types "github.com/goatnetwork/goat/x/relayer/types"
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

// WithdrawalGeneis
type WithdrawalGenesis struct {
	Id         uint64      `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Withdrawal *Withdrawal `protobuf:"bytes,2,opt,name=withdrawal,proto3" json:"withdrawal,omitempty"`
}

func (m *WithdrawalGenesis) Reset()         { *m = WithdrawalGenesis{} }
func (m *WithdrawalGenesis) String() string { return proto.CompactTextString(m) }
func (*WithdrawalGenesis) ProtoMessage()    {}
func (*WithdrawalGenesis) Descriptor() ([]byte, []int) {
	return fileDescriptor_3450fcbc4045b0af, []int{0}
}
func (m *WithdrawalGenesis) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *WithdrawalGenesis) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_WithdrawalGenesis.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *WithdrawalGenesis) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WithdrawalGenesis.Merge(m, src)
}
func (m *WithdrawalGenesis) XXX_Size() int {
	return m.Size()
}
func (m *WithdrawalGenesis) XXX_DiscardUnknown() {
	xxx_messageInfo_WithdrawalGenesis.DiscardUnknown(m)
}

var xxx_messageInfo_WithdrawalGenesis proto.InternalMessageInfo

func (m *WithdrawalGenesis) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *WithdrawalGenesis) GetWithdrawal() *Withdrawal {
	if m != nil {
		return m.Withdrawal
	}
	return nil
}

// DepositGenesis
type DepositGenesis struct {
	Txid   []byte `protobuf:"bytes,1,opt,name=txid,proto3" json:"txid,omitempty"`
	Txout  uint32 `protobuf:"varint,2,opt,name=txout,proto3" json:"txout,omitempty"`
	Amount uint64 `protobuf:"varint,3,opt,name=amount,proto3" json:"amount,omitempty"`
}

func (m *DepositGenesis) Reset()         { *m = DepositGenesis{} }
func (m *DepositGenesis) String() string { return proto.CompactTextString(m) }
func (*DepositGenesis) ProtoMessage()    {}
func (*DepositGenesis) Descriptor() ([]byte, []int) {
	return fileDescriptor_3450fcbc4045b0af, []int{1}
}
func (m *DepositGenesis) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *DepositGenesis) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_DepositGenesis.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *DepositGenesis) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DepositGenesis.Merge(m, src)
}
func (m *DepositGenesis) XXX_Size() int {
	return m.Size()
}
func (m *DepositGenesis) XXX_DiscardUnknown() {
	xxx_messageInfo_DepositGenesis.DiscardUnknown(m)
}

var xxx_messageInfo_DepositGenesis proto.InternalMessageInfo

func (m *DepositGenesis) GetTxid() []byte {
	if m != nil {
		return m.Txid
	}
	return nil
}

func (m *DepositGenesis) GetTxout() uint32 {
	if m != nil {
		return m.Txout
	}
	return 0
}

func (m *DepositGenesis) GetAmount() uint64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

// GenesisState defines the bitcoin module's genesis state.
type GenesisState struct {
	// params defines all the parameters of the module.
	Params      Params               `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
	BlockTip    uint64               `protobuf:"varint,2,opt,name=block_tip,json=blockTip,proto3" json:"block_tip,omitempty"`
	BlockHashes map[uint64][]byte    `protobuf:"bytes,3,rep,name=block_hashes,json=blockHashes,proto3" json:"block_hashes,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	EthTxNonce  uint64               `protobuf:"varint,4,opt,name=eth_tx_nonce,json=ethTxNonce,proto3" json:"eth_tx_nonce,omitempty"`
	Pubkey      *types.PublicKey     `protobuf:"bytes,5,opt,name=pubkey,proto3" json:"pubkey,omitempty"`
	Queue       *ExecuableQueue      `protobuf:"bytes,6,opt,name=queue,proto3" json:"queue,omitempty"`
	Deposits    []*DepositGenesis    `protobuf:"bytes,7,rep,name=deposits,proto3" json:"deposits,omitempty"`
	Withdrawals []*WithdrawalGenesis `protobuf:"bytes,8,rep,name=withdrawals,proto3" json:"withdrawals,omitempty"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_3450fcbc4045b0af, []int{2}
}
func (m *GenesisState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisState.Merge(m, src)
}
func (m *GenesisState) XXX_Size() int {
	return m.Size()
}
func (m *GenesisState) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisState.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisState proto.InternalMessageInfo

func (m *GenesisState) GetParams() Params {
	if m != nil {
		return m.Params
	}
	return Params{}
}

func (m *GenesisState) GetBlockTip() uint64 {
	if m != nil {
		return m.BlockTip
	}
	return 0
}

func (m *GenesisState) GetBlockHashes() map[uint64][]byte {
	if m != nil {
		return m.BlockHashes
	}
	return nil
}

func (m *GenesisState) GetEthTxNonce() uint64 {
	if m != nil {
		return m.EthTxNonce
	}
	return 0
}

func (m *GenesisState) GetPubkey() *types.PublicKey {
	if m != nil {
		return m.Pubkey
	}
	return nil
}

func (m *GenesisState) GetQueue() *ExecuableQueue {
	if m != nil {
		return m.Queue
	}
	return nil
}

func (m *GenesisState) GetDeposits() []*DepositGenesis {
	if m != nil {
		return m.Deposits
	}
	return nil
}

func (m *GenesisState) GetWithdrawals() []*WithdrawalGenesis {
	if m != nil {
		return m.Withdrawals
	}
	return nil
}

func init() {
	proto.RegisterType((*WithdrawalGenesis)(nil), "goat.bitcoin.v1.WithdrawalGenesis")
	proto.RegisterType((*DepositGenesis)(nil), "goat.bitcoin.v1.DepositGenesis")
	proto.RegisterType((*GenesisState)(nil), "goat.bitcoin.v1.GenesisState")
	proto.RegisterMapType((map[uint64][]byte)(nil), "goat.bitcoin.v1.GenesisState.BlockHashesEntry")
}

func init() { proto.RegisterFile("goat/bitcoin/v1/genesis.proto", fileDescriptor_3450fcbc4045b0af) }

var fileDescriptor_3450fcbc4045b0af = []byte{
	// 551 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x53, 0xcd, 0x6a, 0xdb, 0x4e,
	0x10, 0xb7, 0x6c, 0xc7, 0x7f, 0x67, 0xec, 0x7f, 0x9a, 0x2c, 0xa1, 0x15, 0x76, 0xab, 0x18, 0x9f,
	0x4c, 0x0b, 0x12, 0x71, 0x29, 0x94, 0x04, 0x7a, 0x30, 0x31, 0x2d, 0x14, 0x4a, 0xb3, 0x0d, 0x14,
	0x7a, 0x71, 0x57, 0xf2, 0x62, 0x2d, 0x96, 0xb5, 0xaa, 0xb4, 0xb2, 0xe5, 0x57, 0xe8, 0xa9, 0x8f,
	0xd1, 0x63, 0x1f, 0x23, 0xc7, 0x1c, 0x7b, 0x2a, 0xc5, 0x3e, 0xf4, 0x35, 0xca, 0xee, 0xca, 0x1f,
	0xb1, 0xa1, 0x17, 0x31, 0x3b, 0xbf, 0x8f, 0x99, 0xd1, 0xec, 0xc2, 0x93, 0x11, 0x27, 0xc2, 0x71,
	0x99, 0xf0, 0x38, 0x0b, 0x9d, 0xe9, 0xb9, 0x33, 0xa2, 0x21, 0x4d, 0x58, 0x62, 0x47, 0x31, 0x17,
	0x1c, 0x3d, 0x90, 0xb0, 0x9d, 0xc3, 0xf6, 0xf4, 0xbc, 0x71, 0x42, 0x26, 0x2c, 0xe4, 0x8e, 0xfa,
	0x6a, 0x4e, 0x63, 0xcf, 0x62, 0x45, 0xd7, 0xf0, 0xe3, 0x5d, 0x38, 0x22, 0x31, 0x99, 0xe4, 0x05,
	0x1a, 0xcd, 0x5d, 0x54, 0xcc, 0x23, 0x9a, 0xdc, 0x73, 0x8e, 0x69, 0x40, 0xe6, 0x34, 0x96, 0x60,
	0x1e, 0xe6, 0xf0, 0xe9, 0x88, 0x8f, 0xb8, 0x0a, 0x1d, 0x19, 0xe9, 0x6c, 0xfb, 0x33, 0x9c, 0x7c,
	0x64, 0xc2, 0x1f, 0xc6, 0x64, 0x46, 0x82, 0xd7, 0x7a, 0x1a, 0x74, 0x04, 0x45, 0x36, 0x34, 0x8d,
	0x96, 0xd1, 0x29, 0xe3, 0x22, 0x1b, 0xa2, 0x4b, 0x80, 0xd9, 0x9a, 0x64, 0x16, 0x5b, 0x46, 0xa7,
	0xd6, 0x6d, 0xda, 0x3b, 0xc3, 0xda, 0x1b, 0x1f, 0xbc, 0x45, 0x6f, 0x63, 0x38, 0xba, 0xa2, 0x11,
	0x4f, 0x98, 0x58, 0xd9, 0x23, 0x28, 0x8b, 0x2c, 0x2f, 0x50, 0xc7, 0x2a, 0x46, 0xa7, 0x70, 0x20,
	0x32, 0x9e, 0x0a, 0xe5, 0xfe, 0x3f, 0xd6, 0x07, 0xf4, 0x10, 0x2a, 0x64, 0xc2, 0xd3, 0x50, 0x98,
	0x25, 0xd5, 0x4c, 0x7e, 0x6a, 0x7f, 0x2d, 0x43, 0x3d, 0x77, 0xfb, 0x20, 0x88, 0xa0, 0xe8, 0x02,
	0x2a, 0xfa, 0x47, 0x29, 0xd3, 0x5a, 0xf7, 0xd1, 0x5e, 0x77, 0xef, 0x15, 0xdc, 0x3b, 0xbc, 0xfd,
	0x75, 0x56, 0xf8, 0xfe, 0xe7, 0xc7, 0x53, 0x03, 0xe7, 0x0a, 0xd4, 0x84, 0x43, 0x37, 0xe0, 0xde,
	0x78, 0x20, 0x58, 0xa4, 0xca, 0x97, 0x71, 0x55, 0x25, 0x6e, 0x58, 0x84, 0xae, 0xa1, 0xae, 0x41,
	0x9f, 0x24, 0x3e, 0x4d, 0xcc, 0x52, 0xab, 0xd4, 0xa9, 0x75, 0xed, 0x3d, 0xfb, 0xed, 0x6e, 0xec,
	0x9e, 0x54, 0xbc, 0x51, 0x82, 0x7e, 0x28, 0xe2, 0x39, 0xae, 0xb9, 0x9b, 0x0c, 0x6a, 0x41, 0x9d,
	0x0a, 0x7f, 0x20, 0xb2, 0x41, 0xc8, 0x43, 0x8f, 0x9a, 0x65, 0x55, 0x12, 0xa8, 0xf0, 0x6f, 0xb2,
	0x77, 0x32, 0x83, 0xba, 0x50, 0x89, 0x52, 0x77, 0x4c, 0xe7, 0xe6, 0x81, 0x9a, 0xa6, 0xa1, 0xcb,
	0xad, 0xf6, 0x29, 0xa7, 0x49, 0xdd, 0x80, 0x79, 0x6f, 0xe9, 0x1c, 0xe7, 0x4c, 0xf4, 0x02, 0x0e,
	0xbe, 0xa4, 0x34, 0xa5, 0x66, 0x45, 0x49, 0xce, 0xf6, 0x3a, 0xec, 0x67, 0xd4, 0x4b, 0x89, 0x1b,
	0xd0, 0x6b, 0x49, 0xc3, 0x9a, 0x8d, 0x2e, 0xa1, 0x3a, 0xd4, 0xdb, 0x49, 0xcc, 0xff, 0xd4, 0x6c,
	0xfb, 0xca, 0xfb, 0xeb, 0xc3, 0x6b, 0x01, 0xba, 0x82, 0xda, 0x66, 0xd1, 0x89, 0x59, 0x55, 0xfa,
	0xf6, 0x3f, 0x2e, 0xc6, 0xca, 0x62, 0x5b, 0xd6, 0x78, 0x05, 0xc7, 0xbb, 0x3f, 0x0c, 0x1d, 0x43,
	0x49, 0x8e, 0xaf, 0xaf, 0xa0, 0x0c, 0xe5, 0x05, 0x99, 0x92, 0x20, 0xa5, 0x6a, 0x43, 0x75, 0xac,
	0x0f, 0x17, 0xc5, 0x97, 0x46, 0xaf, 0x7f, 0xbb, 0xb0, 0x8c, 0xbb, 0x85, 0x65, 0xfc, 0x5e, 0x58,
	0xc6, 0xb7, 0xa5, 0x55, 0xb8, 0x5b, 0x5a, 0x85, 0x9f, 0x4b, 0xab, 0xf0, 0xe9, 0xd9, 0x88, 0x09,
	0x3f, 0x75, 0x6d, 0x8f, 0x4f, 0x1c, 0xd9, 0x54, 0x48, 0xc5, 0x8c, 0xc7, 0x63, 0x15, 0x3b, 0xd9,
	0xfa, 0x1d, 0xa9, 0x47, 0xe4, 0x56, 0xd4, 0x83, 0x78, 0xfe, 0x37, 0x00, 0x00, 0xff, 0xff, 0xb0,
	0xb3, 0x0f, 0xc2, 0xe4, 0x03, 0x00, 0x00,
}

func (m *WithdrawalGenesis) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *WithdrawalGenesis) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *WithdrawalGenesis) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Withdrawal != nil {
		{
			size, err := m.Withdrawal.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintGenesis(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.Id != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.Id))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *DepositGenesis) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DepositGenesis) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *DepositGenesis) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Amount != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.Amount))
		i--
		dAtA[i] = 0x18
	}
	if m.Txout != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.Txout))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Txid) > 0 {
		i -= len(m.Txid)
		copy(dAtA[i:], m.Txid)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.Txid)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *GenesisState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Withdrawals) > 0 {
		for iNdEx := len(m.Withdrawals) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Withdrawals[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x42
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
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x3a
		}
	}
	if m.Queue != nil {
		{
			size, err := m.Queue.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintGenesis(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x32
	}
	if m.Pubkey != nil {
		{
			size, err := m.Pubkey.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintGenesis(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x2a
	}
	if m.EthTxNonce != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.EthTxNonce))
		i--
		dAtA[i] = 0x20
	}
	if len(m.BlockHashes) > 0 {
		for k := range m.BlockHashes {
			v := m.BlockHashes[k]
			baseI := i
			if len(v) > 0 {
				i -= len(v)
				copy(dAtA[i:], v)
				i = encodeVarintGenesis(dAtA, i, uint64(len(v)))
				i--
				dAtA[i] = 0x12
			}
			i = encodeVarintGenesis(dAtA, i, uint64(k))
			i--
			dAtA[i] = 0x8
			i = encodeVarintGenesis(dAtA, i, uint64(baseI-i))
			i--
			dAtA[i] = 0x1a
		}
	}
	if m.BlockTip != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.BlockTip))
		i--
		dAtA[i] = 0x10
	}
	{
		size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintGenesis(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenesis(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *WithdrawalGenesis) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Id != 0 {
		n += 1 + sovGenesis(uint64(m.Id))
	}
	if m.Withdrawal != nil {
		l = m.Withdrawal.Size()
		n += 1 + l + sovGenesis(uint64(l))
	}
	return n
}

func (m *DepositGenesis) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Txid)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	if m.Txout != 0 {
		n += 1 + sovGenesis(uint64(m.Txout))
	}
	if m.Amount != 0 {
		n += 1 + sovGenesis(uint64(m.Amount))
	}
	return n
}

func (m *GenesisState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Params.Size()
	n += 1 + l + sovGenesis(uint64(l))
	if m.BlockTip != 0 {
		n += 1 + sovGenesis(uint64(m.BlockTip))
	}
	if len(m.BlockHashes) > 0 {
		for k, v := range m.BlockHashes {
			_ = k
			_ = v
			l = 0
			if len(v) > 0 {
				l = 1 + len(v) + sovGenesis(uint64(len(v)))
			}
			mapEntrySize := 1 + sovGenesis(uint64(k)) + l
			n += mapEntrySize + 1 + sovGenesis(uint64(mapEntrySize))
		}
	}
	if m.EthTxNonce != 0 {
		n += 1 + sovGenesis(uint64(m.EthTxNonce))
	}
	if m.Pubkey != nil {
		l = m.Pubkey.Size()
		n += 1 + l + sovGenesis(uint64(l))
	}
	if m.Queue != nil {
		l = m.Queue.Size()
		n += 1 + l + sovGenesis(uint64(l))
	}
	if len(m.Deposits) > 0 {
		for _, e := range m.Deposits {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.Withdrawals) > 0 {
		for _, e := range m.Withdrawals {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	return n
}

func sovGenesis(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenesis(x uint64) (n int) {
	return sovGenesis(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *WithdrawalGenesis) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: WithdrawalGenesis: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: WithdrawalGenesis: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return fmt.Errorf("proto: wrong wireType = %d for field Withdrawal", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Withdrawal == nil {
				m.Withdrawal = &Withdrawal{}
			}
			if err := m.Withdrawal.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func (m *DepositGenesis) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: DepositGenesis: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DepositGenesis: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Txid", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
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
					return ErrIntOverflowGenesis
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
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			m.Amount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func (m *GenesisState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: GenesisState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field BlockTip", wireType)
			}
			m.BlockTip = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.BlockTip |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BlockHashes", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.BlockHashes == nil {
				m.BlockHashes = make(map[uint64][]byte)
			}
			var mapkey uint64
			mapvalue := []byte{}
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowGenesis
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
				if fieldNum == 1 {
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowGenesis
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapkey |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
				} else if fieldNum == 2 {
					var mapbyteLen uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowGenesis
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapbyteLen |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intMapbyteLen := int(mapbyteLen)
					if intMapbyteLen < 0 {
						return ErrInvalidLengthGenesis
					}
					postbytesIndex := iNdEx + intMapbyteLen
					if postbytesIndex < 0 {
						return ErrInvalidLengthGenesis
					}
					if postbytesIndex > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = make([]byte, mapbyteLen)
					copy(mapvalue, dAtA[iNdEx:postbytesIndex])
					iNdEx = postbytesIndex
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipGenesis(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if (skippy < 0) || (iNdEx+skippy) < 0 {
						return ErrInvalidLengthGenesis
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.BlockHashes[mapkey] = mapvalue
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field EthTxNonce", wireType)
			}
			m.EthTxNonce = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.EthTxNonce |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pubkey", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Pubkey == nil {
				m.Pubkey = &types.PublicKey{}
			}
			if err := m.Pubkey.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Queue", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Queue == nil {
				m.Queue = &ExecuableQueue{}
			}
			if err := m.Queue.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Deposits", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Deposits = append(m.Deposits, &DepositGenesis{})
			if err := m.Deposits[len(m.Deposits)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Withdrawals", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Withdrawals = append(m.Withdrawals, &WithdrawalGenesis{})
			if err := m.Withdrawals[len(m.Withdrawals)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func skipGenesis(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
				return 0, ErrInvalidLengthGenesis
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenesis
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenesis
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenesis        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenesis          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenesis = fmt.Errorf("proto: unexpected end of group")
)
