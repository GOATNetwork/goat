// Code generated by MockGen. DO NOT EDIT.
// Source: expected_keepers.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	types "github.com/cosmos/cosmos-sdk/types"
	types0 "github.com/goatnetwork/goat/x/relayer/types"
	gomock "github.com/golang/mock/gomock"
)

// MockRelayerKeeper is a mock of RelayerKeeper interface.
type MockRelayerKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockRelayerKeeperMockRecorder
}

// MockRelayerKeeperMockRecorder is the mock recorder for MockRelayerKeeper.
type MockRelayerKeeperMockRecorder struct {
	mock *MockRelayerKeeper
}

// NewMockRelayerKeeper creates a new mock instance.
func NewMockRelayerKeeper(ctrl *gomock.Controller) *MockRelayerKeeper {
	mock := &MockRelayerKeeper{ctrl: ctrl}
	mock.recorder = &MockRelayerKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRelayerKeeper) EXPECT() *MockRelayerKeeperMockRecorder {
	return m.recorder
}

// AddNewKey mocks base method.
func (m *MockRelayerKeeper) AddNewKey(ctx context.Context, raw []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddNewKey", ctx, raw)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddNewKey indicates an expected call of AddNewKey.
func (mr *MockRelayerKeeperMockRecorder) AddNewKey(ctx, raw interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddNewKey", reflect.TypeOf((*MockRelayerKeeper)(nil).AddNewKey), ctx, raw)
}

// HasPubkey mocks base method.
func (m *MockRelayerKeeper) HasPubkey(ctx context.Context, raw []byte) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasPubkey", ctx, raw)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HasPubkey indicates an expected call of HasPubkey.
func (mr *MockRelayerKeeperMockRecorder) HasPubkey(ctx, raw interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasPubkey", reflect.TypeOf((*MockRelayerKeeper)(nil).HasPubkey), ctx, raw)
}

// SetProposalSeq mocks base method.
func (m *MockRelayerKeeper) SetProposalSeq(ctx context.Context, seq uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetProposalSeq", ctx, seq)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetProposalSeq indicates an expected call of SetProposalSeq.
func (mr *MockRelayerKeeperMockRecorder) SetProposalSeq(ctx, seq interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetProposalSeq", reflect.TypeOf((*MockRelayerKeeper)(nil).SetProposalSeq), ctx, seq)
}

// UpdateRandao mocks base method.
func (m *MockRelayerKeeper) UpdateRandao(ctx context.Context, req types0.IVoteMsg) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateRandao", ctx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateRandao indicates an expected call of UpdateRandao.
func (mr *MockRelayerKeeperMockRecorder) UpdateRandao(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateRandao", reflect.TypeOf((*MockRelayerKeeper)(nil).UpdateRandao), ctx, req)
}

// VerifyNonProposal mocks base method.
func (m *MockRelayerKeeper) VerifyNonProposal(ctx context.Context, req types0.INonVoteMsg) (types0.IRelayer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyNonProposal", ctx, req)
	ret0, _ := ret[0].(types0.IRelayer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VerifyNonProposal indicates an expected call of VerifyNonProposal.
func (mr *MockRelayerKeeperMockRecorder) VerifyNonProposal(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyNonProposal", reflect.TypeOf((*MockRelayerKeeper)(nil).VerifyNonProposal), ctx, req)
}

// VerifyProposal mocks base method.
func (m *MockRelayerKeeper) VerifyProposal(ctx context.Context, req types0.IVoteMsg, verifyFn ...func([]byte) error) (uint64, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, req}
	for _, a := range verifyFn {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "VerifyProposal", varargs...)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VerifyProposal indicates an expected call of VerifyProposal.
func (mr *MockRelayerKeeperMockRecorder) VerifyProposal(ctx, req interface{}, verifyFn ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, req}, verifyFn...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyProposal", reflect.TypeOf((*MockRelayerKeeper)(nil).VerifyProposal), varargs...)
}

// MockAccountKeeper is a mock of AccountKeeper interface.
type MockAccountKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockAccountKeeperMockRecorder
}

// MockAccountKeeperMockRecorder is the mock recorder for MockAccountKeeper.
type MockAccountKeeperMockRecorder struct {
	mock *MockAccountKeeper
}

// NewMockAccountKeeper creates a new mock instance.
func NewMockAccountKeeper(ctrl *gomock.Controller) *MockAccountKeeper {
	mock := &MockAccountKeeper{ctrl: ctrl}
	mock.recorder = &MockAccountKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccountKeeper) EXPECT() *MockAccountKeeperMockRecorder {
	return m.recorder
}

// GetAccount mocks base method.
func (m *MockAccountKeeper) GetAccount(arg0 context.Context, arg1 types.AccAddress) types.AccountI {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccount", arg0, arg1)
	ret0, _ := ret[0].(types.AccountI)
	return ret0
}

// GetAccount indicates an expected call of GetAccount.
func (mr *MockAccountKeeperMockRecorder) GetAccount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccount", reflect.TypeOf((*MockAccountKeeper)(nil).GetAccount), arg0, arg1)
}

// MockBankKeeper is a mock of BankKeeper interface.
type MockBankKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockBankKeeperMockRecorder
}

// MockBankKeeperMockRecorder is the mock recorder for MockBankKeeper.
type MockBankKeeperMockRecorder struct {
	mock *MockBankKeeper
}

// NewMockBankKeeper creates a new mock instance.
func NewMockBankKeeper(ctrl *gomock.Controller) *MockBankKeeper {
	mock := &MockBankKeeper{ctrl: ctrl}
	mock.recorder = &MockBankKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBankKeeper) EXPECT() *MockBankKeeperMockRecorder {
	return m.recorder
}

// SpendableCoins mocks base method.
func (m *MockBankKeeper) SpendableCoins(arg0 context.Context, arg1 types.AccAddress) types.Coins {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SpendableCoins", arg0, arg1)
	ret0, _ := ret[0].(types.Coins)
	return ret0
}

// SpendableCoins indicates an expected call of SpendableCoins.
func (mr *MockBankKeeperMockRecorder) SpendableCoins(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SpendableCoins", reflect.TypeOf((*MockBankKeeper)(nil).SpendableCoins), arg0, arg1)
}

// MockParamSubspace is a mock of ParamSubspace interface.
type MockParamSubspace struct {
	ctrl     *gomock.Controller
	recorder *MockParamSubspaceMockRecorder
}

// MockParamSubspaceMockRecorder is the mock recorder for MockParamSubspace.
type MockParamSubspaceMockRecorder struct {
	mock *MockParamSubspace
}

// NewMockParamSubspace creates a new mock instance.
func NewMockParamSubspace(ctrl *gomock.Controller) *MockParamSubspace {
	mock := &MockParamSubspace{ctrl: ctrl}
	mock.recorder = &MockParamSubspaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockParamSubspace) EXPECT() *MockParamSubspaceMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockParamSubspace) Get(arg0 context.Context, arg1 []byte, arg2 interface{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Get", arg0, arg1, arg2)
}

// Get indicates an expected call of Get.
func (mr *MockParamSubspaceMockRecorder) Get(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockParamSubspace)(nil).Get), arg0, arg1, arg2)
}

// Set mocks base method.
func (m *MockParamSubspace) Set(arg0 context.Context, arg1 []byte, arg2 interface{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Set", arg0, arg1, arg2)
}

// Set indicates an expected call of Set.
func (mr *MockParamSubspaceMockRecorder) Set(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockParamSubspace)(nil).Set), arg0, arg1, arg2)
}
