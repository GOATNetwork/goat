// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: goat/bitcoin/v1/bitcoin.proto

package types

import (
	fmt "fmt"
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

// WithdrawalStatus is the status of a withdrawal.
type WithdrawalStatus int32

const (
	// WITHDRAWAL_STATUS_UNSPECIFIED defines an invalid status.
	WITHDRAWAL_STATUS_UNSPECIFIED WithdrawalStatus = 0
	// WITHDRAWAL_STATUS_PENDING
	WITHDRAWAL_STATUS_PENDING WithdrawalStatus = 1
	// WITHDRAWAL_STATUS_PROCESSING
	WITHDRAWAL_STATUS_PROCESSING WithdrawalStatus = 2
	// WITHDRAWAL_STATUS_CANCELING
	WITHDRAWAL_STATUS_CANCELING WithdrawalStatus = 3
	// WITHDRAWAL_STATUS_CANCELED
	WITHDRAWAL_STATUS_CANCELED WithdrawalStatus = 4
	// WITHDRAWAL_STATUS_PAID
	WITHDRAWAL_STATUS_PAID WithdrawalStatus = 5
)

var WithdrawalStatus_name = map[int32]string{
	0: "WITHDRAWAL_STATUS_UNSPECIFIED",
	1: "WITHDRAWAL_STATUS_PENDING",
	2: "WITHDRAWAL_STATUS_PROCESSING",
	3: "WITHDRAWAL_STATUS_CANCELING",
	4: "WITHDRAWAL_STATUS_CANCELED",
	5: "WITHDRAWAL_STATUS_PAID",
}

var WithdrawalStatus_value = map[string]int32{
	"WITHDRAWAL_STATUS_UNSPECIFIED": 0,
	"WITHDRAWAL_STATUS_PENDING":     1,
	"WITHDRAWAL_STATUS_PROCESSING":  2,
	"WITHDRAWAL_STATUS_CANCELING":   3,
	"WITHDRAWAL_STATUS_CANCELED":    4,
	"WITHDRAWAL_STATUS_PAID":        5,
}

func (x WithdrawalStatus) String() string {
	return proto.EnumName(WithdrawalStatus_name, int32(x))
}

func (WithdrawalStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_c8a1a8d7fb9d8c11, []int{0}
}

// Deposit defines the deposit transaction and its proof
type Deposit struct {
	Version uint32 `protobuf:"varint,1,opt,name=version,proto3" json:"version,omitempty"`
	// block_number the block number that transaction located at
	BlockNumber uint64 `protobuf:"varint,2,opt,name=block_number,json=blockNumber,proto3" json:"block_number,omitempty"`
	// tx_index is the index of transaction in the block
	TxIndex uint32 `protobuf:"varint,3,opt,name=tx_index,json=txIndex,proto3" json:"tx_index,omitempty"`
	// tx is the raw transaction withtout witness
	NoWitnessTx []byte `protobuf:"bytes,4,opt,name=no_witness_tx,json=noWitnessTx,proto3" json:"no_witness_tx,omitempty"`
	OutputIndex uint32 `protobuf:"varint,5,opt,name=output_index,json=outputIndex,proto3" json:"output_index,omitempty"`
	// intermediate proof is the proof without the txid and merkel root
	IntermediateProof []byte `protobuf:"bytes,6,opt,name=intermediate_proof,json=intermediateProof,proto3" json:"intermediate_proof,omitempty"`
	// evm_address is the user wallet address in goat evm
	EvmAddress []byte `protobuf:"bytes,7,opt,name=evm_address,json=evmAddress,proto3" json:"evm_address,omitempty"`
	// relayer_pubkey is the compressed secp256k1 public key which managed by the
	// relayer group
	RelayerPubkey *types.PublicKey `protobuf:"bytes,8,opt,name=relayer_pubkey,json=relayerPubkey,proto3" json:"relayer_pubkey,omitempty"`
}

func (m *Deposit) Reset()         { *m = Deposit{} }
func (m *Deposit) String() string { return proto.CompactTextString(m) }
func (*Deposit) ProtoMessage()    {}
func (*Deposit) Descriptor() ([]byte, []int) {
	return fileDescriptor_c8a1a8d7fb9d8c11, []int{0}
}
func (m *Deposit) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Deposit) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Deposit.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Deposit) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Deposit.Merge(m, src)
}
func (m *Deposit) XXX_Size() int {
	return m.Size()
}
func (m *Deposit) XXX_DiscardUnknown() {
	xxx_messageInfo_Deposit.DiscardUnknown(m)
}

var xxx_messageInfo_Deposit proto.InternalMessageInfo

func (m *Deposit) GetVersion() uint32 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *Deposit) GetBlockNumber() uint64 {
	if m != nil {
		return m.BlockNumber
	}
	return 0
}

func (m *Deposit) GetTxIndex() uint32 {
	if m != nil {
		return m.TxIndex
	}
	return 0
}

func (m *Deposit) GetNoWitnessTx() []byte {
	if m != nil {
		return m.NoWitnessTx
	}
	return nil
}

func (m *Deposit) GetOutputIndex() uint32 {
	if m != nil {
		return m.OutputIndex
	}
	return 0
}

func (m *Deposit) GetIntermediateProof() []byte {
	if m != nil {
		return m.IntermediateProof
	}
	return nil
}

func (m *Deposit) GetEvmAddress() []byte {
	if m != nil {
		return m.EvmAddress
	}
	return nil
}

func (m *Deposit) GetRelayerPubkey() *types.PublicKey {
	if m != nil {
		return m.RelayerPubkey
	}
	return nil
}

// WithdrawalReceipt
type WithdrawalReceipt struct {
	Txid   []byte `protobuf:"bytes,1,opt,name=txid,proto3" json:"txid,omitempty"`
	Txout  uint32 `protobuf:"varint,2,opt,name=txout,proto3" json:"txout,omitempty"`
	Amount uint64 `protobuf:"varint,3,opt,name=amount,proto3" json:"amount,omitempty"`
}

func (m *WithdrawalReceipt) Reset()         { *m = WithdrawalReceipt{} }
func (m *WithdrawalReceipt) String() string { return proto.CompactTextString(m) }
func (*WithdrawalReceipt) ProtoMessage()    {}
func (*WithdrawalReceipt) Descriptor() ([]byte, []int) {
	return fileDescriptor_c8a1a8d7fb9d8c11, []int{1}
}
func (m *WithdrawalReceipt) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *WithdrawalReceipt) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_WithdrawalReceipt.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *WithdrawalReceipt) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WithdrawalReceipt.Merge(m, src)
}
func (m *WithdrawalReceipt) XXX_Size() int {
	return m.Size()
}
func (m *WithdrawalReceipt) XXX_DiscardUnknown() {
	xxx_messageInfo_WithdrawalReceipt.DiscardUnknown(m)
}

var xxx_messageInfo_WithdrawalReceipt proto.InternalMessageInfo

func (m *WithdrawalReceipt) GetTxid() []byte {
	if m != nil {
		return m.Txid
	}
	return nil
}

func (m *WithdrawalReceipt) GetTxout() uint32 {
	if m != nil {
		return m.Txout
	}
	return 0
}

func (m *WithdrawalReceipt) GetAmount() uint64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

// Withdrawal
type Withdrawal struct {
	Address       string             `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	RequestAmount uint64             `protobuf:"varint,2,opt,name=request_amount,json=requestAmount,proto3" json:"request_amount,omitempty"`
	MaxTxPrice    uint64             `protobuf:"varint,3,opt,name=max_tx_price,json=maxTxPrice,proto3" json:"max_tx_price,omitempty"`
	Status        WithdrawalStatus   `protobuf:"varint,4,opt,name=status,proto3,enum=goat.bitcoin.v1.WithdrawalStatus" json:"status,omitempty"`
	Receipt       *WithdrawalReceipt `protobuf:"bytes,5,opt,name=receipt,proto3" json:"receipt,omitempty"`
}

func (m *Withdrawal) Reset()         { *m = Withdrawal{} }
func (m *Withdrawal) String() string { return proto.CompactTextString(m) }
func (*Withdrawal) ProtoMessage()    {}
func (*Withdrawal) Descriptor() ([]byte, []int) {
	return fileDescriptor_c8a1a8d7fb9d8c11, []int{2}
}
func (m *Withdrawal) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Withdrawal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Withdrawal.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Withdrawal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Withdrawal.Merge(m, src)
}
func (m *Withdrawal) XXX_Size() int {
	return m.Size()
}
func (m *Withdrawal) XXX_DiscardUnknown() {
	xxx_messageInfo_Withdrawal.DiscardUnknown(m)
}

var xxx_messageInfo_Withdrawal proto.InternalMessageInfo

func (m *Withdrawal) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *Withdrawal) GetRequestAmount() uint64 {
	if m != nil {
		return m.RequestAmount
	}
	return 0
}

func (m *Withdrawal) GetMaxTxPrice() uint64 {
	if m != nil {
		return m.MaxTxPrice
	}
	return 0
}

func (m *Withdrawal) GetStatus() WithdrawalStatus {
	if m != nil {
		return m.Status
	}
	return WITHDRAWAL_STATUS_UNSPECIFIED
}

func (m *Withdrawal) GetReceipt() *WithdrawalReceipt {
	if m != nil {
		return m.Receipt
	}
	return nil
}

func init() {
	proto.RegisterEnum("goat.bitcoin.v1.WithdrawalStatus", WithdrawalStatus_name, WithdrawalStatus_value)
	proto.RegisterType((*Deposit)(nil), "goat.bitcoin.v1.Deposit")
	proto.RegisterType((*WithdrawalReceipt)(nil), "goat.bitcoin.v1.WithdrawalReceipt")
	proto.RegisterType((*Withdrawal)(nil), "goat.bitcoin.v1.Withdrawal")
}

func init() { proto.RegisterFile("goat/bitcoin/v1/bitcoin.proto", fileDescriptor_c8a1a8d7fb9d8c11) }

var fileDescriptor_c8a1a8d7fb9d8c11 = []byte{
	// 630 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x53, 0x41, 0x6e, 0xd3, 0x40,
	0x14, 0x8d, 0xdb, 0x34, 0x29, 0x3f, 0x49, 0x49, 0x47, 0x55, 0xe5, 0x06, 0xea, 0xa6, 0x91, 0x90,
	0x22, 0x10, 0x8e, 0x5a, 0x56, 0x48, 0x6c, 0x4c, 0x62, 0x20, 0xa2, 0x0a, 0xd1, 0x24, 0x55, 0x24,
	0x36, 0x96, 0x9d, 0x0c, 0xe9, 0xa8, 0xb1, 0xc7, 0xd8, 0x63, 0xd7, 0xb9, 0x01, 0x4b, 0xee, 0xc0,
	0x31, 0xb8, 0x00, 0xcb, 0x2e, 0x61, 0x87, 0x9a, 0x0d, 0xc7, 0x40, 0x33, 0xb6, 0x29, 0x6a, 0x0b,
	0xbb, 0xf7, 0xdf, 0x7f, 0xff, 0x8d, 0xe7, 0x7d, 0x0f, 0xec, 0xcf, 0x99, 0xcd, 0x3b, 0x0e, 0xe5,
	0x53, 0x46, 0xbd, 0x4e, 0x7c, 0x94, 0x43, 0xdd, 0x0f, 0x18, 0x67, 0xe8, 0xbe, 0x68, 0xeb, 0x39,
	0x17, 0x1f, 0x35, 0x52, 0x7d, 0x40, 0x16, 0xf6, 0x92, 0x04, 0x42, 0x9f, 0xc1, 0x54, 0xdf, 0xd8,
	0x99, 0xb3, 0x39, 0x93, 0xb0, 0x23, 0x50, 0xca, 0xb6, 0xbe, 0xae, 0x41, 0xb9, 0x47, 0x7c, 0x16,
	0x52, 0x8e, 0x54, 0x28, 0xc7, 0x24, 0x08, 0x29, 0xf3, 0x54, 0xa5, 0xa9, 0xb4, 0x6b, 0x38, 0x2f,
	0xd1, 0x21, 0x54, 0x9d, 0x05, 0x9b, 0x9e, 0x5b, 0x5e, 0xe4, 0x3a, 0x24, 0x50, 0xd7, 0x9a, 0x4a,
	0xbb, 0x88, 0x2b, 0x92, 0x1b, 0x48, 0x0a, 0xed, 0xc1, 0x26, 0x4f, 0x2c, 0xea, 0xcd, 0x48, 0xa2,
	0xae, 0xa7, 0xd3, 0x3c, 0xe9, 0x8b, 0x12, 0xb5, 0xa0, 0xe6, 0x31, 0xeb, 0x82, 0x72, 0x8f, 0x84,
	0xa1, 0xc5, 0x13, 0xb5, 0xd8, 0x54, 0xda, 0x55, 0x5c, 0xf1, 0xd8, 0x24, 0xe5, 0xc6, 0x89, 0x38,
	0x81, 0x45, 0xdc, 0x8f, 0x78, 0x66, 0xb1, 0x21, 0x2d, 0x2a, 0x29, 0x97, 0xda, 0x3c, 0x05, 0x44,
	0x3d, 0x4e, 0x02, 0x97, 0xcc, 0xa8, 0xcd, 0x89, 0xe5, 0x07, 0x8c, 0x7d, 0x50, 0x4b, 0xd2, 0x6b,
	0xfb, 0xef, 0xce, 0x50, 0x34, 0xd0, 0x01, 0x54, 0x48, 0xec, 0x5a, 0xf6, 0x6c, 0x16, 0x90, 0x30,
	0x54, 0xcb, 0x52, 0x07, 0x24, 0x76, 0x8d, 0x94, 0x41, 0x06, 0x6c, 0x65, 0x09, 0x59, 0x7e, 0xe4,
	0x9c, 0x93, 0xa5, 0xba, 0xd9, 0x54, 0xda, 0x95, 0xe3, 0x86, 0x2e, 0x93, 0xcd, 0xd3, 0x8b, 0x8f,
	0xf4, 0x61, 0xe4, 0x2c, 0xe8, 0xf4, 0x2d, 0x59, 0xe2, 0x5a, 0xc6, 0x0e, 0xe5, 0x40, 0xeb, 0x14,
	0xb6, 0x27, 0x94, 0x9f, 0xcd, 0x02, 0xfb, 0xc2, 0x5e, 0x60, 0x32, 0x25, 0xd4, 0xe7, 0x08, 0x41,
	0x91, 0x27, 0x74, 0x26, 0x33, 0xac, 0x62, 0x89, 0xd1, 0x0e, 0x6c, 0xf0, 0x84, 0x45, 0x5c, 0x26,
	0x57, 0xc3, 0x69, 0x81, 0x76, 0xa1, 0x64, 0xbb, 0x2c, 0xf2, 0xb8, 0x4c, 0xac, 0x88, 0xb3, 0xaa,
	0xf5, 0x4b, 0x01, 0xb8, 0xf6, 0x15, 0x7b, 0xc9, 0x6f, 0x21, 0x3c, 0xef, 0xe1, 0xbc, 0x44, 0x8f,
	0xc4, 0x15, 0x3e, 0x46, 0x24, 0xe4, 0x56, 0x66, 0x94, 0x6e, 0xa6, 0x96, 0xb1, 0x86, 0x24, 0x51,
	0x13, 0xaa, 0xae, 0x9d, 0x58, 0x3c, 0xb1, 0xfc, 0x80, 0x4e, 0x49, 0x76, 0x1a, 0xb8, 0x76, 0x32,
	0x4e, 0x86, 0x82, 0x41, 0xcf, 0xa1, 0x14, 0x72, 0x9b, 0x47, 0xa1, 0xdc, 0xcd, 0xd6, 0xf1, 0xa1,
	0x7e, 0xe3, 0xef, 0xd2, 0xaf, 0xbf, 0x67, 0x24, 0x85, 0x38, 0x1b, 0x40, 0x2f, 0xa0, 0x1c, 0xa4,
	0x37, 0x97, 0x4b, 0xab, 0x1c, 0xb7, 0xfe, 0x33, 0x9b, 0x65, 0x84, 0xf3, 0x91, 0xc7, 0x3f, 0x14,
	0xa8, 0xdf, 0xb4, 0x46, 0x87, 0xb0, 0x3f, 0xe9, 0x8f, 0xdf, 0xf4, 0xb0, 0x31, 0x31, 0x4e, 0xac,
	0xd1, 0xd8, 0x18, 0x9f, 0x8e, 0xac, 0xd3, 0xc1, 0x68, 0x68, 0x76, 0xfb, 0xaf, 0xfa, 0x66, 0xaf,
	0x5e, 0x40, 0xfb, 0xb0, 0x77, 0x5b, 0x32, 0x34, 0x07, 0xbd, 0xfe, 0xe0, 0x75, 0x5d, 0x41, 0x4d,
	0x78, 0x78, 0x47, 0x1b, 0xbf, 0xeb, 0x9a, 0xa3, 0x91, 0x50, 0xac, 0xa1, 0x03, 0x78, 0x70, 0x5b,
	0xd1, 0x35, 0x06, 0x5d, 0xf3, 0x44, 0x08, 0xd6, 0x91, 0x06, 0x8d, 0x7f, 0x09, 0xcc, 0x5e, 0xbd,
	0x88, 0x1a, 0xb0, 0x7b, 0xc7, 0x11, 0x46, 0xbf, 0x57, 0xdf, 0x68, 0x14, 0x3f, 0x7d, 0xd1, 0x0a,
	0x2f, 0xcd, 0x6f, 0x57, 0x9a, 0x72, 0x79, 0xa5, 0x29, 0x3f, 0xaf, 0x34, 0xe5, 0xf3, 0x4a, 0x2b,
	0x5c, 0xae, 0xb4, 0xc2, 0xf7, 0x95, 0x56, 0x78, 0xff, 0x64, 0x4e, 0xf9, 0x59, 0xe4, 0xe8, 0x53,
	0xe6, 0x76, 0x44, 0x58, 0x1e, 0xe1, 0x17, 0x2c, 0x38, 0x97, 0xb8, 0x93, 0xfc, 0x79, 0xf3, 0x7c,
	0xe9, 0x93, 0xd0, 0x29, 0xc9, 0x97, 0xfa, 0xec, 0x77, 0x00, 0x00, 0x00, 0xff, 0xff, 0x6c, 0x3c,
	0xd2, 0xf9, 0x10, 0x04, 0x00, 0x00,
}

func (m *Deposit) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Deposit) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Deposit) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.RelayerPubkey != nil {
		{
			size, err := m.RelayerPubkey.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintBitcoin(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x42
	}
	if len(m.EvmAddress) > 0 {
		i -= len(m.EvmAddress)
		copy(dAtA[i:], m.EvmAddress)
		i = encodeVarintBitcoin(dAtA, i, uint64(len(m.EvmAddress)))
		i--
		dAtA[i] = 0x3a
	}
	if len(m.IntermediateProof) > 0 {
		i -= len(m.IntermediateProof)
		copy(dAtA[i:], m.IntermediateProof)
		i = encodeVarintBitcoin(dAtA, i, uint64(len(m.IntermediateProof)))
		i--
		dAtA[i] = 0x32
	}
	if m.OutputIndex != 0 {
		i = encodeVarintBitcoin(dAtA, i, uint64(m.OutputIndex))
		i--
		dAtA[i] = 0x28
	}
	if len(m.NoWitnessTx) > 0 {
		i -= len(m.NoWitnessTx)
		copy(dAtA[i:], m.NoWitnessTx)
		i = encodeVarintBitcoin(dAtA, i, uint64(len(m.NoWitnessTx)))
		i--
		dAtA[i] = 0x22
	}
	if m.TxIndex != 0 {
		i = encodeVarintBitcoin(dAtA, i, uint64(m.TxIndex))
		i--
		dAtA[i] = 0x18
	}
	if m.BlockNumber != 0 {
		i = encodeVarintBitcoin(dAtA, i, uint64(m.BlockNumber))
		i--
		dAtA[i] = 0x10
	}
	if m.Version != 0 {
		i = encodeVarintBitcoin(dAtA, i, uint64(m.Version))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *WithdrawalReceipt) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *WithdrawalReceipt) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *WithdrawalReceipt) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Amount != 0 {
		i = encodeVarintBitcoin(dAtA, i, uint64(m.Amount))
		i--
		dAtA[i] = 0x18
	}
	if m.Txout != 0 {
		i = encodeVarintBitcoin(dAtA, i, uint64(m.Txout))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Txid) > 0 {
		i -= len(m.Txid)
		copy(dAtA[i:], m.Txid)
		i = encodeVarintBitcoin(dAtA, i, uint64(len(m.Txid)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Withdrawal) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Withdrawal) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Withdrawal) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
			i = encodeVarintBitcoin(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x2a
	}
	if m.Status != 0 {
		i = encodeVarintBitcoin(dAtA, i, uint64(m.Status))
		i--
		dAtA[i] = 0x20
	}
	if m.MaxTxPrice != 0 {
		i = encodeVarintBitcoin(dAtA, i, uint64(m.MaxTxPrice))
		i--
		dAtA[i] = 0x18
	}
	if m.RequestAmount != 0 {
		i = encodeVarintBitcoin(dAtA, i, uint64(m.RequestAmount))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintBitcoin(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintBitcoin(dAtA []byte, offset int, v uint64) int {
	offset -= sovBitcoin(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Deposit) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Version != 0 {
		n += 1 + sovBitcoin(uint64(m.Version))
	}
	if m.BlockNumber != 0 {
		n += 1 + sovBitcoin(uint64(m.BlockNumber))
	}
	if m.TxIndex != 0 {
		n += 1 + sovBitcoin(uint64(m.TxIndex))
	}
	l = len(m.NoWitnessTx)
	if l > 0 {
		n += 1 + l + sovBitcoin(uint64(l))
	}
	if m.OutputIndex != 0 {
		n += 1 + sovBitcoin(uint64(m.OutputIndex))
	}
	l = len(m.IntermediateProof)
	if l > 0 {
		n += 1 + l + sovBitcoin(uint64(l))
	}
	l = len(m.EvmAddress)
	if l > 0 {
		n += 1 + l + sovBitcoin(uint64(l))
	}
	if m.RelayerPubkey != nil {
		l = m.RelayerPubkey.Size()
		n += 1 + l + sovBitcoin(uint64(l))
	}
	return n
}

func (m *WithdrawalReceipt) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Txid)
	if l > 0 {
		n += 1 + l + sovBitcoin(uint64(l))
	}
	if m.Txout != 0 {
		n += 1 + sovBitcoin(uint64(m.Txout))
	}
	if m.Amount != 0 {
		n += 1 + sovBitcoin(uint64(m.Amount))
	}
	return n
}

func (m *Withdrawal) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovBitcoin(uint64(l))
	}
	if m.RequestAmount != 0 {
		n += 1 + sovBitcoin(uint64(m.RequestAmount))
	}
	if m.MaxTxPrice != 0 {
		n += 1 + sovBitcoin(uint64(m.MaxTxPrice))
	}
	if m.Status != 0 {
		n += 1 + sovBitcoin(uint64(m.Status))
	}
	if m.Receipt != nil {
		l = m.Receipt.Size()
		n += 1 + l + sovBitcoin(uint64(l))
	}
	return n
}

func sovBitcoin(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozBitcoin(x uint64) (n int) {
	return sovBitcoin(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Deposit) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBitcoin
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
			return fmt.Errorf("proto: Deposit: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Deposit: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Version", wireType)
			}
			m.Version = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBitcoin
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Version |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field BlockNumber", wireType)
			}
			m.BlockNumber = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBitcoin
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
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TxIndex", wireType)
			}
			m.TxIndex = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBitcoin
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.TxIndex |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NoWitnessTx", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBitcoin
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
				return ErrInvalidLengthBitcoin
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthBitcoin
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.NoWitnessTx = append(m.NoWitnessTx[:0], dAtA[iNdEx:postIndex]...)
			if m.NoWitnessTx == nil {
				m.NoWitnessTx = []byte{}
			}
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field OutputIndex", wireType)
			}
			m.OutputIndex = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBitcoin
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.OutputIndex |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field IntermediateProof", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBitcoin
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
				return ErrInvalidLengthBitcoin
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthBitcoin
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.IntermediateProof = append(m.IntermediateProof[:0], dAtA[iNdEx:postIndex]...)
			if m.IntermediateProof == nil {
				m.IntermediateProof = []byte{}
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EvmAddress", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBitcoin
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
				return ErrInvalidLengthBitcoin
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthBitcoin
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.EvmAddress = append(m.EvmAddress[:0], dAtA[iNdEx:postIndex]...)
			if m.EvmAddress == nil {
				m.EvmAddress = []byte{}
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RelayerPubkey", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBitcoin
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
				return ErrInvalidLengthBitcoin
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthBitcoin
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.RelayerPubkey == nil {
				m.RelayerPubkey = &types.PublicKey{}
			}
			if err := m.RelayerPubkey.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipBitcoin(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthBitcoin
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
func (m *WithdrawalReceipt) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBitcoin
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
			return fmt.Errorf("proto: WithdrawalReceipt: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: WithdrawalReceipt: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Txid", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBitcoin
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
				return ErrInvalidLengthBitcoin
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthBitcoin
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
					return ErrIntOverflowBitcoin
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
					return ErrIntOverflowBitcoin
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
			skippy, err := skipBitcoin(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthBitcoin
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
func (m *Withdrawal) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBitcoin
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
			return fmt.Errorf("proto: Withdrawal: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Withdrawal: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBitcoin
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
				return ErrInvalidLengthBitcoin
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthBitcoin
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field RequestAmount", wireType)
			}
			m.RequestAmount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBitcoin
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.RequestAmount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxTxPrice", wireType)
			}
			m.MaxTxPrice = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBitcoin
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MaxTxPrice |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Status", wireType)
			}
			m.Status = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBitcoin
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Status |= WithdrawalStatus(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Receipt", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBitcoin
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
				return ErrInvalidLengthBitcoin
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthBitcoin
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
			skippy, err := skipBitcoin(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthBitcoin
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
func skipBitcoin(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowBitcoin
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
					return 0, ErrIntOverflowBitcoin
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
					return 0, ErrIntOverflowBitcoin
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
				return 0, ErrInvalidLengthBitcoin
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupBitcoin
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthBitcoin
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthBitcoin        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowBitcoin          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupBitcoin = fmt.Errorf("proto: unexpected end of group")
)
