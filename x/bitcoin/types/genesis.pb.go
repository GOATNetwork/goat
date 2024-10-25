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
	Id         uint64     `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Withdrawal Withdrawal `protobuf:"bytes,2,opt,name=withdrawal,proto3" json:"withdrawal"`
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

func (m *WithdrawalGenesis) GetWithdrawal() Withdrawal {
	if m != nil {
		return m.Withdrawal
	}
	return Withdrawal{}
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

// ProcessingGenesis
type ProcessingGenesis struct {
	Id         uint64     `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Processing Processing `protobuf:"bytes,2,opt,name=processing,proto3" json:"processing"`
}

func (m *ProcessingGenesis) Reset()         { *m = ProcessingGenesis{} }
func (m *ProcessingGenesis) String() string { return proto.CompactTextString(m) }
func (*ProcessingGenesis) ProtoMessage()    {}
func (*ProcessingGenesis) Descriptor() ([]byte, []int) {
	return fileDescriptor_3450fcbc4045b0af, []int{2}
}
func (m *ProcessingGenesis) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ProcessingGenesis) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ProcessingGenesis.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ProcessingGenesis) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProcessingGenesis.Merge(m, src)
}
func (m *ProcessingGenesis) XXX_Size() int {
	return m.Size()
}
func (m *ProcessingGenesis) XXX_DiscardUnknown() {
	xxx_messageInfo_ProcessingGenesis.DiscardUnknown(m)
}

var xxx_messageInfo_ProcessingGenesis proto.InternalMessageInfo

func (m *ProcessingGenesis) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *ProcessingGenesis) GetProcessing() Processing {
	if m != nil {
		return m.Processing
	}
	return Processing{}
}

// GenesisState defines the bitcoin module's genesis state.
type GenesisState struct {
	// params defines all the parameters of the module.
	Params          Params              `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
	BlockTip        uint64              `protobuf:"varint,2,opt,name=block_tip,json=blockTip,proto3" json:"block_tip,omitempty"`
	BlockHashes     [][]byte            `protobuf:"bytes,3,rep,name=block_hashes,json=blockHashes,proto3" json:"block_hashes,omitempty"`
	EthTxNonce      uint64              `protobuf:"varint,4,opt,name=eth_tx_nonce,json=ethTxNonce,proto3" json:"eth_tx_nonce,omitempty"`
	EthTxQueue      EthTxQueue          `protobuf:"bytes,5,opt,name=eth_tx_queue,json=ethTxQueue,proto3" json:"eth_tx_queue"`
	Pubkey          *types.PublicKey    `protobuf:"bytes,6,opt,name=pubkey,proto3" json:"pubkey,omitempty"`
	Deposits        []DepositGenesis    `protobuf:"bytes,7,rep,name=deposits,proto3" json:"deposits"`
	Withdrawals     []WithdrawalGenesis `protobuf:"bytes,8,rep,name=withdrawals,proto3" json:"withdrawals"`
	Processing      []ProcessingGenesis `protobuf:"bytes,9,rep,name=processing,proto3" json:"processing"`
	LatestProcessId uint64              `protobuf:"varint,10,opt,name=latest_process_id,json=latestProcessId,proto3" json:"latest_process_id,omitempty"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_3450fcbc4045b0af, []int{3}
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

func (m *GenesisState) GetBlockHashes() [][]byte {
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

func (m *GenesisState) GetEthTxQueue() EthTxQueue {
	if m != nil {
		return m.EthTxQueue
	}
	return EthTxQueue{}
}

func (m *GenesisState) GetPubkey() *types.PublicKey {
	if m != nil {
		return m.Pubkey
	}
	return nil
}

func (m *GenesisState) GetDeposits() []DepositGenesis {
	if m != nil {
		return m.Deposits
	}
	return nil
}

func (m *GenesisState) GetWithdrawals() []WithdrawalGenesis {
	if m != nil {
		return m.Withdrawals
	}
	return nil
}

func (m *GenesisState) GetProcessing() []ProcessingGenesis {
	if m != nil {
		return m.Processing
	}
	return nil
}

func (m *GenesisState) GetLatestProcessId() uint64 {
	if m != nil {
		return m.LatestProcessId
	}
	return 0
}

func init() {
	proto.RegisterType((*WithdrawalGenesis)(nil), "goat.bitcoin.v1.WithdrawalGenesis")
	proto.RegisterType((*DepositGenesis)(nil), "goat.bitcoin.v1.DepositGenesis")
	proto.RegisterType((*ProcessingGenesis)(nil), "goat.bitcoin.v1.ProcessingGenesis")
	proto.RegisterType((*GenesisState)(nil), "goat.bitcoin.v1.GenesisState")
}

func init() { proto.RegisterFile("goat/bitcoin/v1/genesis.proto", fileDescriptor_3450fcbc4045b0af) }

var fileDescriptor_3450fcbc4045b0af = []byte{
	// 564 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x53, 0xcd, 0x6e, 0xda, 0x4c,
	0x14, 0xc5, 0x40, 0xf8, 0x60, 0xe0, 0x4b, 0xc4, 0x28, 0x6a, 0x2d, 0x68, 0x1d, 0xca, 0x0a, 0xa5,
	0x92, 0xad, 0xd0, 0x5d, 0x97, 0x51, 0x93, 0xa6, 0xaa, 0xfa, 0xe7, 0x46, 0xaa, 0xd4, 0x0d, 0x1a,
	0xdb, 0x23, 0x7b, 0x04, 0x78, 0x5c, 0xcf, 0x75, 0x80, 0xb7, 0xe8, 0x63, 0x74, 0xd9, 0x27, 0xe8,
	0x3a, 0xcb, 0x2c, 0xbb, 0xaa, 0x2a, 0x58, 0xf4, 0x35, 0x2a, 0x8f, 0x07, 0xf3, 0x57, 0xb2, 0xb1,
	0xae, 0xef, 0xb9, 0xf7, 0x9c, 0x39, 0x9a, 0x33, 0xe8, 0xb1, 0xcf, 0x09, 0x58, 0x0e, 0x03, 0x97,
	0xb3, 0xd0, 0xba, 0x39, 0xb3, 0x7c, 0x1a, 0x52, 0xc1, 0x84, 0x19, 0xc5, 0x1c, 0x38, 0x3e, 0x4a,
	0x61, 0x53, 0xc1, 0xe6, 0xcd, 0x59, 0xab, 0x49, 0xc6, 0x2c, 0xe4, 0x96, 0xfc, 0x66, 0x33, 0xad,
	0x1d, 0x8a, 0xe5, 0x78, 0x06, 0x3f, 0xda, 0x86, 0x23, 0x12, 0x93, 0xb1, 0x12, 0x68, 0xb5, 0xb7,
	0x51, 0x98, 0x45, 0x54, 0x6c, 0x30, 0xc7, 0x74, 0x44, 0x66, 0x34, 0x4e, 0x41, 0x55, 0x2a, 0xf8,
	0xd8, 0xe7, 0x3e, 0x97, 0xa5, 0x95, 0x56, 0x59, 0xb7, 0x3b, 0x44, 0xcd, 0x4f, 0x0c, 0x02, 0x2f,
	0x26, 0x13, 0x32, 0x7a, 0x99, 0xb9, 0xc1, 0x87, 0xa8, 0xc8, 0x3c, 0x5d, 0xeb, 0x68, 0xbd, 0xb2,
	0x5d, 0x64, 0x1e, 0xbe, 0x44, 0x68, 0x92, 0x0f, 0xe9, 0xc5, 0x8e, 0xd6, 0xab, 0xf7, 0xdb, 0xe6,
	0x96, 0x59, 0x73, 0xc5, 0x73, 0x5e, 0xbb, 0xfd, 0x75, 0x52, 0xf8, 0xf6, 0xe7, 0xfb, 0xa9, 0x66,
	0xaf, 0x6d, 0x76, 0x6d, 0x74, 0xf8, 0x82, 0x46, 0x5c, 0x30, 0x58, 0x2a, 0x61, 0x54, 0x86, 0xa9,
	0xd2, 0x6a, 0xd8, 0xb2, 0xc6, 0xc7, 0xe8, 0x00, 0xa6, 0x3c, 0x01, 0x29, 0xf4, 0xbf, 0x9d, 0xfd,
	0xe0, 0x07, 0xa8, 0x42, 0xc6, 0x3c, 0x09, 0x41, 0x2f, 0xc9, 0x73, 0xa9, 0xbf, 0xd4, 0xc0, 0xfb,
	0x98, 0xbb, 0x54, 0x08, 0x16, 0xfa, 0xf7, 0x18, 0x88, 0xf2, 0xa1, 0xbd, 0x06, 0x56, 0x3c, 0x1b,
	0x06, 0x56, 0x9b, 0xdd, 0x1f, 0x65, 0xd4, 0x50, 0x1a, 0x1f, 0x81, 0x00, 0xc5, 0xcf, 0x51, 0x25,
	0xbb, 0x20, 0x29, 0x56, 0xef, 0x3f, 0xdc, 0x25, 0x95, 0xf0, 0x3a, 0xa1, 0xda, 0xc0, 0x6d, 0x54,
	0x73, 0x46, 0xdc, 0x1d, 0x0e, 0x80, 0x45, 0xf2, 0x4c, 0x65, 0xbb, 0x2a, 0x1b, 0xd7, 0x2c, 0xc2,
	0x4f, 0x50, 0x23, 0x03, 0x03, 0x22, 0x02, 0x2a, 0xf4, 0x52, 0xa7, 0xd4, 0x6b, 0xd8, 0x75, 0xd9,
	0xbb, 0x92, 0x2d, 0xdc, 0x41, 0x0d, 0x0a, 0xc1, 0x00, 0xa6, 0x83, 0x90, 0x87, 0x2e, 0xd5, 0xcb,
	0x92, 0x02, 0x51, 0x08, 0xae, 0xa7, 0x6f, 0xd3, 0x0e, 0xbe, 0xca, 0x27, 0xbe, 0x24, 0x34, 0xa1,
	0xfa, 0xc1, 0x1e, 0xe3, 0x17, 0xe9, 0xca, 0x87, 0x74, 0x64, 0xc3, 0x38, 0xcd, 0xdb, 0xb8, 0x8f,
	0x2a, 0x51, 0xe2, 0x0c, 0xe9, 0x4c, 0xaf, 0x48, 0x8e, 0x56, 0xc6, 0xb1, 0x4c, 0x58, 0xea, 0x33,
	0x71, 0x46, 0xcc, 0x7d, 0x4d, 0x67, 0xb6, 0x9a, 0xc4, 0x97, 0xa8, 0xea, 0x65, 0xb7, 0x2d, 0xf4,
	0xff, 0x3a, 0xa5, 0x5e, 0xbd, 0x7f, 0xb2, 0xa3, 0xbc, 0x19, 0x87, 0x75, 0xf5, 0x7c, 0x17, 0xbf,
	0x43, 0xf5, 0x55, 0x86, 0x84, 0x5e, 0x95, 0x54, 0xdd, 0x7b, 0xe2, 0xf7, 0x0f, 0xb6, 0x75, 0x06,
	0xfc, 0x66, 0x23, 0x0d, 0xb5, 0x3d, 0x7c, 0x3b, 0xa9, 0xda, 0x13, 0x0a, 0x7c, 0x8a, 0x9a, 0x23,
	0x02, 0x54, 0xc0, 0x40, 0x35, 0x07, 0xcc, 0xd3, 0x91, 0xbc, 0x8c, 0xa3, 0x0c, 0x50, 0x54, 0xaf,
	0xbc, 0xf3, 0x8b, 0xdb, 0xb9, 0xa1, 0xdd, 0xcd, 0x0d, 0xed, 0xf7, 0xdc, 0xd0, 0xbe, 0x2e, 0x8c,
	0xc2, 0xdd, 0xc2, 0x28, 0xfc, 0x5c, 0x18, 0x85, 0xcf, 0x4f, 0x7d, 0x06, 0x41, 0xe2, 0x98, 0x2e,
	0x1f, 0x5b, 0xe9, 0x51, 0x42, 0x0a, 0x13, 0x1e, 0x0f, 0x65, 0x6d, 0x4d, 0xf3, 0x37, 0x2f, 0x1f,
	0xbc, 0x53, 0x91, 0x8f, 0xf7, 0xd9, 0xdf, 0x00, 0x00, 0x00, 0xff, 0xff, 0xce, 0xd2, 0xf8, 0x82,
	0x90, 0x04, 0x00, 0x00,
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

func (m *ProcessingGenesis) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ProcessingGenesis) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ProcessingGenesis) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Processing.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if m.Id != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.Id))
		i--
		dAtA[i] = 0x8
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
	if m.LatestProcessId != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.LatestProcessId))
		i--
		dAtA[i] = 0x50
	}
	if len(m.Processing) > 0 {
		for iNdEx := len(m.Processing) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Processing[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x4a
		}
	}
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
		dAtA[i] = 0x32
	}
	{
		size, err := m.EthTxQueue.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	if m.EthTxNonce != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.EthTxNonce))
		i--
		dAtA[i] = 0x20
	}
	if len(m.BlockHashes) > 0 {
		for iNdEx := len(m.BlockHashes) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.BlockHashes[iNdEx])
			copy(dAtA[i:], m.BlockHashes[iNdEx])
			i = encodeVarintGenesis(dAtA, i, uint64(len(m.BlockHashes[iNdEx])))
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
	l = m.Withdrawal.Size()
	n += 1 + l + sovGenesis(uint64(l))
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

func (m *ProcessingGenesis) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Id != 0 {
		n += 1 + sovGenesis(uint64(m.Id))
	}
	l = m.Processing.Size()
	n += 1 + l + sovGenesis(uint64(l))
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
		for _, b := range m.BlockHashes {
			l = len(b)
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if m.EthTxNonce != 0 {
		n += 1 + sovGenesis(uint64(m.EthTxNonce))
	}
	l = m.EthTxQueue.Size()
	n += 1 + l + sovGenesis(uint64(l))
	if m.Pubkey != nil {
		l = m.Pubkey.Size()
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
	if len(m.Processing) > 0 {
		for _, e := range m.Processing {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if m.LatestProcessId != 0 {
		n += 1 + sovGenesis(uint64(m.LatestProcessId))
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
func (m *ProcessingGenesis) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: ProcessingGenesis: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ProcessingGenesis: illegal tag %d (wire type %d)", fieldNum, wire)
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
				return fmt.Errorf("proto: wrong wireType = %d for field Processing", wireType)
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
			if err := m.Processing.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
			m.BlockHashes = append(m.BlockHashes, make([]byte, postIndex-iNdEx))
			copy(m.BlockHashes[len(m.BlockHashes)-1], dAtA[iNdEx:postIndex])
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
				return fmt.Errorf("proto: wrong wireType = %d for field EthTxQueue", wireType)
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
			if err := m.EthTxQueue.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
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
			m.Deposits = append(m.Deposits, DepositGenesis{})
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
			m.Withdrawals = append(m.Withdrawals, WithdrawalGenesis{})
			if err := m.Withdrawals[len(m.Withdrawals)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Processing", wireType)
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
			m.Processing = append(m.Processing, ProcessingGenesis{})
			if err := m.Processing[len(m.Processing)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 10:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LatestProcessId", wireType)
			}
			m.LatestProcessId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.LatestProcessId |= uint64(b&0x7F) << shift
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
