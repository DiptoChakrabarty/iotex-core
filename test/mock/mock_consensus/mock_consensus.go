// Code generated by MockGen. DO NOT EDIT.
// Source: ./consensus/consensus.go

// Package mock_consensus is a generated GoMock package.
package mock_consensus

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	block "github.com/iotexproject/iotex-core/blockchain/block"
	scheme "github.com/iotexproject/iotex-core/consensus/scheme"
	iotextypes "github.com/iotexproject/iotex-proto/golang/iotextypes"
)

// MockConsensus is a mock of Consensus interface.
type MockConsensus struct {
	ctrl     *gomock.Controller
	recorder *MockConsensusMockRecorder
}

// MockConsensusMockRecorder is the mock recorder for MockConsensus.
type MockConsensusMockRecorder struct {
	mock *MockConsensus
}

// NewMockConsensus creates a new mock instance.
func NewMockConsensus(ctrl *gomock.Controller) *MockConsensus {
	mock := &MockConsensus{ctrl: ctrl}
	mock.recorder = &MockConsensusMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConsensus) EXPECT() *MockConsensusMockRecorder {
	return m.recorder
}

// Activate mocks base method.
func (m *MockConsensus) Activate(arg0 bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Activate", arg0)
}

// Activate indicates an expected call of Activate.
func (mr *MockConsensusMockRecorder) Activate(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Activate", reflect.TypeOf((*MockConsensus)(nil).Activate), arg0)
}

// Active mocks base method.
func (m *MockConsensus) Active() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Active")
	ret0, _ := ret[0].(bool)
	return ret0
}

// Active indicates an expected call of Active.
func (mr *MockConsensusMockRecorder) Active() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Active", reflect.TypeOf((*MockConsensus)(nil).Active))
}

// Calibrate mocks base method.
func (m *MockConsensus) Calibrate(arg0 uint64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Calibrate", arg0)
}

// Calibrate indicates an expected call of Calibrate.
func (mr *MockConsensusMockRecorder) Calibrate(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Calibrate", reflect.TypeOf((*MockConsensus)(nil).Calibrate), arg0)
}

// HandleConsensusMsg mocks base method.
func (m *MockConsensus) HandleConsensusMsg(arg0 *iotextypes.ConsensusMessage) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandleConsensusMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// HandleConsensusMsg indicates an expected call of HandleConsensusMsg.
func (mr *MockConsensusMockRecorder) HandleConsensusMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleConsensusMsg", reflect.TypeOf((*MockConsensus)(nil).HandleConsensusMsg), arg0)
}

// Metrics mocks base method.
func (m *MockConsensus) Metrics() (scheme.ConsensusMetrics, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Metrics")
	ret0, _ := ret[0].(scheme.ConsensusMetrics)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Metrics indicates an expected call of Metrics.
func (mr *MockConsensusMockRecorder) Metrics() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Metrics", reflect.TypeOf((*MockConsensus)(nil).Metrics))
}

// Start mocks base method.
func (m *MockConsensus) Start(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start.
func (mr *MockConsensusMockRecorder) Start(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockConsensus)(nil).Start), arg0)
}

// Stop mocks base method.
func (m *MockConsensus) Stop(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stop", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Stop indicates an expected call of Stop.
func (mr *MockConsensusMockRecorder) Stop(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockConsensus)(nil).Stop), arg0)
}

// ValidateBlockFooter mocks base method.
func (m *MockConsensus) ValidateBlockFooter(arg0 *block.Block) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateBlockFooter", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateBlockFooter indicates an expected call of ValidateBlockFooter.
func (mr *MockConsensusMockRecorder) ValidateBlockFooter(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateBlockFooter", reflect.TypeOf((*MockConsensus)(nil).ValidateBlockFooter), arg0)
}
