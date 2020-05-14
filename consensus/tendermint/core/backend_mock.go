// Code generated by MockGen. DO NOT EDIT.
// Source: consensus/tendermint/core/core_backend.go

// Package core is a generated GoMock package.
package core

import (
	context "context"
	common "github.com/clearmatics/autonity/common"
	committee "github.com/clearmatics/autonity/consensus/tendermint/committee"
	types "github.com/clearmatics/autonity/core/types"
	event "github.com/clearmatics/autonity/event"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	time "time"
)

// MockBackend is a mock of Backend interface
type MockBackend struct {
	ctrl     *gomock.Controller
	recorder *MockBackendMockRecorder
}

// MockBackendMockRecorder is the mock recorder for MockBackend
type MockBackendMockRecorder struct {
	mock *MockBackend
}

// NewMockBackend creates a new mock instance
func NewMockBackend(ctrl *gomock.Controller) *MockBackend {
	mock := &MockBackend{ctrl: ctrl}
	mock.recorder = &MockBackendMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBackend) EXPECT() *MockBackendMockRecorder {
	return m.recorder
}

// Address mocks base method
func (m *MockBackend) Address() common.Address {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Address")
	ret0, _ := ret[0].(common.Address)
	return ret0
}

// Address indicates an expected call of Address
func (mr *MockBackendMockRecorder) Address() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Address", reflect.TypeOf((*MockBackend)(nil).Address))
}

// AddSeal mocks base method
func (m *MockBackend) AddSeal(block *types.Block) (*types.Block, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddSeal", block)
	ret0, _ := ret[0].(*types.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddSeal indicates an expected call of AddSeal
func (mr *MockBackendMockRecorder) AddSeal(block interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddSeal", reflect.TypeOf((*MockBackend)(nil).AddSeal), block)
}

// AskSync mocks base method
func (m *MockBackend) AskSync(set committee.Set) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AskSync", set)
}

// AskSync indicates an expected call of AskSync
func (mr *MockBackendMockRecorder) AskSync(set interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AskSync", reflect.TypeOf((*MockBackend)(nil).AskSync), set)
}

// Broadcast mocks base method
func (m *MockBackend) Broadcast(ctx context.Context, valSet committee.Set, payload []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Broadcast", ctx, valSet, payload)
	ret0, _ := ret[0].(error)
	return ret0
}

// Broadcast indicates an expected call of Broadcast
func (mr *MockBackendMockRecorder) Broadcast(ctx, valSet, payload interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Broadcast", reflect.TypeOf((*MockBackend)(nil).Broadcast), ctx, valSet, payload)
}

// Commit mocks base method
func (m *MockBackend) Commit(proposalBlock *types.Block, round int64, seals [][]byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Commit", proposalBlock, round, seals)
	ret0, _ := ret[0].(error)
	return ret0
}

// Commit indicates an expected call of Commit
func (mr *MockBackendMockRecorder) Commit(proposalBlock, round, seals interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Commit", reflect.TypeOf((*MockBackend)(nil).Commit), proposalBlock, round, seals)
}

// Committee mocks base method
func (m *MockBackend) Committee(number uint64) (committee.Set, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Committee", number)
	ret0, _ := ret[0].(committee.Set)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Committee indicates an expected call of Committee
func (mr *MockBackendMockRecorder) Committee(number interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Committee", reflect.TypeOf((*MockBackend)(nil).Committee), number)
}

// GetContractABI mocks base method
func (m *MockBackend) GetContractABI() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetContractABI")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetContractABI indicates an expected call of GetContractABI
func (mr *MockBackendMockRecorder) GetContractABI() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContractABI", reflect.TypeOf((*MockBackend)(nil).GetContractABI))
}

// GetContractAddress mocks base method
func (m *MockBackend) GetContractAddress() common.Address {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetContractAddress")
	ret0, _ := ret[0].(common.Address)
	return ret0
}

// GetContractAddress indicates an expected call of GetContractAddress
func (mr *MockBackendMockRecorder) GetContractAddress() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContractAddress", reflect.TypeOf((*MockBackend)(nil).GetContractAddress))
}

// Gossip mocks base method
func (m *MockBackend) Gossip(ctx context.Context, valSet committee.Set, payload []byte) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Gossip", ctx, valSet, payload)
}

// Gossip indicates an expected call of Gossip
func (mr *MockBackendMockRecorder) Gossip(ctx, valSet, payload interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Gossip", reflect.TypeOf((*MockBackend)(nil).Gossip), ctx, valSet, payload)
}

// HandleUnhandledMsgs mocks base method
func (m *MockBackend) HandleUnhandledMsgs(ctx context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "HandleUnhandledMsgs", ctx)
}

// HandleUnhandledMsgs indicates an expected call of HandleUnhandledMsgs
func (mr *MockBackendMockRecorder) HandleUnhandledMsgs(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleUnhandledMsgs", reflect.TypeOf((*MockBackend)(nil).HandleUnhandledMsgs), ctx)
}

// LastCommittedProposal mocks base method
func (m *MockBackend) LastCommittedProposal() (*types.Block, common.Address) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LastCommittedProposal")
	ret0, _ := ret[0].(*types.Block)
	ret1, _ := ret[1].(common.Address)
	return ret0, ret1
}

// LastCommittedProposal indicates an expected call of LastCommittedProposal
func (mr *MockBackendMockRecorder) LastCommittedProposal() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastCommittedProposal", reflect.TypeOf((*MockBackend)(nil).LastCommittedProposal))
}

// GetProposer mocks base method
func (m *MockBackend) GetProposer(number uint64) common.Address {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProposer", number)
	ret0, _ := ret[0].(common.Address)
	return ret0
}

// GetProposer indicates an expected call of GetProposer
func (mr *MockBackendMockRecorder) GetProposer(number interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProposer", reflect.TypeOf((*MockBackend)(nil).GetProposer), number)
}

// GetProposerFromAC mocks base method
func (m *MockBackend) GetProposerFromAC(height uint64, round int64) common.Address {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProposerFromAC", height, round)
	ret0, _ := ret[0].(common.Address)
	return ret0
}

// GetProposerFromAC indicates an expected call of GetProposerFromAC
func (mr *MockBackendMockRecorder) GetProposerFromAC(height, round interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProposerFromAC", reflect.TypeOf((*MockBackend)(nil).GetProposerFromAC), height, round)
}

// HasBadProposal mocks base method
func (m *MockBackend) HasBadProposal(hash common.Hash) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasBadProposal", hash)
	ret0, _ := ret[0].(bool)
	return ret0
}

// HasBadProposal indicates an expected call of HasBadProposal
func (mr *MockBackendMockRecorder) HasBadProposal(hash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasBadProposal", reflect.TypeOf((*MockBackend)(nil).HasBadProposal), hash)
}

// Post mocks base method
func (m *MockBackend) Post(ev interface{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Post", ev)
}

// Post indicates an expected call of Post
func (mr *MockBackendMockRecorder) Post(ev interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Post", reflect.TypeOf((*MockBackend)(nil).Post), ev)
}

// SetProposedBlockHash mocks base method
func (m *MockBackend) SetProposedBlockHash(hash common.Hash) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetProposedBlockHash", hash)
}

// SetProposedBlockHash indicates an expected call of SetProposedBlockHash
func (mr *MockBackendMockRecorder) SetProposedBlockHash(hash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetProposedBlockHash", reflect.TypeOf((*MockBackend)(nil).SetProposedBlockHash), hash)
}

// Sign mocks base method
func (m *MockBackend) Sign(arg0 []byte) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sign", arg0)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Sign indicates an expected call of Sign
func (mr *MockBackendMockRecorder) Sign(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sign", reflect.TypeOf((*MockBackend)(nil).Sign), arg0)
}

// Subscribe mocks base method
func (m *MockBackend) Subscribe(types ...interface{}) *event.TypeMuxSubscription {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range types {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Subscribe", varargs...)
	ret0, _ := ret[0].(*event.TypeMuxSubscription)
	return ret0
}

// Subscribe indicates an expected call of Subscribe
func (mr *MockBackendMockRecorder) Subscribe(types ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockBackend)(nil).Subscribe), types...)
}

// SyncPeer mocks base method
func (m *MockBackend) SyncPeer(address common.Address) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SyncPeer", address)
}

// SyncPeer indicates an expected call of SyncPeer
func (mr *MockBackendMockRecorder) SyncPeer(address interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SyncPeer", reflect.TypeOf((*MockBackend)(nil).SyncPeer), address)
}

// VerifyProposal mocks base method
func (m *MockBackend) VerifyProposal(arg0 types.Block) (time.Duration, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyProposal", arg0)
	ret0, _ := ret[0].(time.Duration)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VerifyProposal indicates an expected call of VerifyProposal
func (mr *MockBackendMockRecorder) VerifyProposal(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyProposal", reflect.TypeOf((*MockBackend)(nil).VerifyProposal), arg0)
}

// WhiteList mocks base method
func (m *MockBackend) WhiteList() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WhiteList")
	ret0, _ := ret[0].([]string)
	return ret0
}

// WhiteList indicates an expected call of WhiteList
func (mr *MockBackendMockRecorder) WhiteList() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WhiteList", reflect.TypeOf((*MockBackend)(nil).WhiteList))
}

// MockTendermint is a mock of Tendermint interface
type MockTendermint struct {
	ctrl     *gomock.Controller
	recorder *MockTendermintMockRecorder
}

// MockTendermintMockRecorder is the mock recorder for MockTendermint
type MockTendermintMockRecorder struct {
	mock *MockTendermint
}

// NewMockTendermint creates a new mock instance
func NewMockTendermint(ctrl *gomock.Controller) *MockTendermint {
	mock := &MockTendermint{ctrl: ctrl}
	mock.recorder = &MockTendermintMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTendermint) EXPECT() *MockTendermintMockRecorder {
	return m.recorder
}

// Start mocks base method
func (m *MockTendermint) Start(ctx context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Start", ctx)
}

// Start indicates an expected call of Start
func (mr *MockTendermintMockRecorder) Start(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockTendermint)(nil).Start), ctx)
}

// Stop mocks base method
func (m *MockTendermint) Stop() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Stop")
}

// Stop indicates an expected call of Stop
func (mr *MockTendermintMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockTendermint)(nil).Stop))
}

// GetCurrentHeightMessages mocks base method
func (m *MockTendermint) GetCurrentHeightMessages() []*Message {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCurrentHeightMessages")
	ret0, _ := ret[0].([]*Message)
	return ret0
}

// GetCurrentHeightMessages indicates an expected call of GetCurrentHeightMessages
func (mr *MockTendermintMockRecorder) GetCurrentHeightMessages() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCurrentHeightMessages", reflect.TypeOf((*MockTendermint)(nil).GetCurrentHeightMessages))
}
