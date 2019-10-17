// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Xuanwo/storage (interfaces: Storager,Servicer)

// Package mock is a generated GoMock package.
package mock

import (
	storage "github.com/Xuanwo/storage"
	iterator "github.com/Xuanwo/storage/pkg/iterator"
	types "github.com/Xuanwo/storage/types"
	gomock "github.com/golang/mock/gomock"
	io "io"
	reflect "reflect"
)

// MockStorager is a mock of Storager interface
type MockStorager struct {
	ctrl     *gomock.Controller
	recorder *MockStoragerMockRecorder
}

// MockStoragerMockRecorder is the mock recorder for MockStorager
type MockStoragerMockRecorder struct {
	mock *MockStorager
}

// NewMockStorager creates a new mock instance
func NewMockStorager(ctrl *gomock.Controller) *MockStorager {
	mock := &MockStorager{ctrl: ctrl}
	mock.recorder = &MockStoragerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStorager) EXPECT() *MockStoragerMockRecorder {
	return m.recorder
}

// AbortSegment mocks base method
func (m *MockStorager) AbortSegment(arg0 string, arg1 ...*types.Pair) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AbortSegment", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// AbortSegment indicates an expected call of AbortSegment
func (mr *MockStoragerMockRecorder) AbortSegment(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AbortSegment", reflect.TypeOf((*MockStorager)(nil).AbortSegment), varargs...)
}

// Capability mocks base method
func (m *MockStorager) Capability() types.Capability {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Capability")
	ret0, _ := ret[0].(types.Capability)
	return ret0
}

// Capability indicates an expected call of Capability
func (mr *MockStoragerMockRecorder) Capability() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Capability", reflect.TypeOf((*MockStorager)(nil).Capability))
}

// CompleteSegment mocks base method
func (m *MockStorager) CompleteSegment(arg0 string, arg1 ...*types.Pair) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CompleteSegment", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// CompleteSegment indicates an expected call of CompleteSegment
func (mr *MockStoragerMockRecorder) CompleteSegment(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CompleteSegment", reflect.TypeOf((*MockStorager)(nil).CompleteSegment), varargs...)
}

// Copy mocks base method
func (m *MockStorager) Copy(arg0, arg1 string, arg2 ...*types.Pair) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Copy", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Copy indicates an expected call of Copy
func (mr *MockStoragerMockRecorder) Copy(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Copy", reflect.TypeOf((*MockStorager)(nil).Copy), varargs...)
}

// CreateDir mocks base method
func (m *MockStorager) CreateDir(arg0 string, arg1 ...*types.Pair) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateDir", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateDir indicates an expected call of CreateDir
func (mr *MockStoragerMockRecorder) CreateDir(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateDir", reflect.TypeOf((*MockStorager)(nil).CreateDir), varargs...)
}

// Delete mocks base method
func (m *MockStorager) Delete(arg0 string, arg1 ...*types.Pair) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Delete", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockStoragerMockRecorder) Delete(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockStorager)(nil).Delete), varargs...)
}

// InitSegment mocks base method
func (m *MockStorager) InitSegment(arg0 string, arg1 ...*types.Pair) (string, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "InitSegment", varargs...)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InitSegment indicates an expected call of InitSegment
func (mr *MockStoragerMockRecorder) InitSegment(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InitSegment", reflect.TypeOf((*MockStorager)(nil).InitSegment), varargs...)
}

// IsPairAvailable mocks base method
func (m *MockStorager) IsPairAvailable(arg0, arg1 string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsPairAvailable", arg0, arg1)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsPairAvailable indicates an expected call of IsPairAvailable
func (mr *MockStoragerMockRecorder) IsPairAvailable(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsPairAvailable", reflect.TypeOf((*MockStorager)(nil).IsPairAvailable), arg0, arg1)
}

// ListDir mocks base method
func (m *MockStorager) ListDir(arg0 string, arg1 ...*types.Pair) iterator.ObjectIterator {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListDir", varargs...)
	ret0, _ := ret[0].(iterator.ObjectIterator)
	return ret0
}

// ListDir indicates an expected call of ListDir
func (mr *MockStoragerMockRecorder) ListDir(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListDir", reflect.TypeOf((*MockStorager)(nil).ListDir), varargs...)
}

// ListSegments mocks base method
func (m *MockStorager) ListSegments(arg0 string, arg1 ...*types.Pair) iterator.SegmentIterator {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListSegments", varargs...)
	ret0, _ := ret[0].(iterator.SegmentIterator)
	return ret0
}

// ListSegments indicates an expected call of ListSegments
func (mr *MockStoragerMockRecorder) ListSegments(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListSegments", reflect.TypeOf((*MockStorager)(nil).ListSegments), varargs...)
}

// Metadata mocks base method
func (m *MockStorager) Metadata() (types.Metadata, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Metadata")
	ret0, _ := ret[0].(types.Metadata)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Metadata indicates an expected call of Metadata
func (mr *MockStoragerMockRecorder) Metadata() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Metadata", reflect.TypeOf((*MockStorager)(nil).Metadata))
}

// Move mocks base method
func (m *MockStorager) Move(arg0, arg1 string, arg2 ...*types.Pair) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Move", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Move indicates an expected call of Move
func (mr *MockStoragerMockRecorder) Move(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Move", reflect.TypeOf((*MockStorager)(nil).Move), varargs...)
}

// Reach mocks base method
func (m *MockStorager) Reach(arg0 string, arg1 ...*types.Pair) (string, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Reach", varargs...)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Reach indicates an expected call of Reach
func (mr *MockStoragerMockRecorder) Reach(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reach", reflect.TypeOf((*MockStorager)(nil).Reach), varargs...)
}

// Read mocks base method
func (m *MockStorager) Read(arg0 string, arg1 ...*types.Pair) (io.ReadCloser, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Read", varargs...)
	ret0, _ := ret[0].(io.ReadCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read
func (mr *MockStoragerMockRecorder) Read(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockStorager)(nil).Read), varargs...)
}

// Stat mocks base method
func (m *MockStorager) Stat(arg0 string, arg1 ...*types.Pair) (*types.Object, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Stat", varargs...)
	ret0, _ := ret[0].(*types.Object)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Stat indicates an expected call of Stat
func (mr *MockStoragerMockRecorder) Stat(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stat", reflect.TypeOf((*MockStorager)(nil).Stat), varargs...)
}

// WriteFile mocks base method
func (m *MockStorager) WriteFile(arg0 string, arg1 int64, arg2 io.Reader, arg3 ...*types.Pair) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2}
	for _, a := range arg3 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "WriteFile", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteFile indicates an expected call of WriteFile
func (mr *MockStoragerMockRecorder) WriteFile(arg0, arg1, arg2 interface{}, arg3 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2}, arg3...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteFile", reflect.TypeOf((*MockStorager)(nil).WriteFile), varargs...)
}

// WriteSegment mocks base method
func (m *MockStorager) WriteSegment(arg0 string, arg1, arg2 int64, arg3 io.Reader, arg4 ...*types.Pair) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2, arg3}
	for _, a := range arg4 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "WriteSegment", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteSegment indicates an expected call of WriteSegment
func (mr *MockStoragerMockRecorder) WriteSegment(arg0, arg1, arg2, arg3 interface{}, arg4 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2, arg3}, arg4...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteSegment", reflect.TypeOf((*MockStorager)(nil).WriteSegment), varargs...)
}

// WriteStream mocks base method
func (m *MockStorager) WriteStream(arg0 string, arg1 io.Reader, arg2 ...*types.Pair) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "WriteStream", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteStream indicates an expected call of WriteStream
func (mr *MockStoragerMockRecorder) WriteStream(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteStream", reflect.TypeOf((*MockStorager)(nil).WriteStream), varargs...)
}

// MockServicer is a mock of Servicer interface
type MockServicer struct {
	ctrl     *gomock.Controller
	recorder *MockServicerMockRecorder
}

// MockServicerMockRecorder is the mock recorder for MockServicer
type MockServicerMockRecorder struct {
	mock *MockServicer
}

// NewMockServicer creates a new mock instance
func NewMockServicer(ctrl *gomock.Controller) *MockServicer {
	mock := &MockServicer{ctrl: ctrl}
	mock.recorder = &MockServicerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockServicer) EXPECT() *MockServicerMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockServicer) Create(arg0 string, arg1 ...*types.Pair) (storage.Storager, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Create", varargs...)
	ret0, _ := ret[0].(storage.Storager)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockServicerMockRecorder) Create(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockServicer)(nil).Create), varargs...)
}

// Delete mocks base method
func (m *MockServicer) Delete(arg0 string, arg1 ...*types.Pair) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Delete", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockServicerMockRecorder) Delete(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockServicer)(nil).Delete), varargs...)
}

// Get mocks base method
func (m *MockServicer) Get(arg0 string, arg1 ...*types.Pair) (storage.Storager, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Get", varargs...)
	ret0, _ := ret[0].(storage.Storager)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockServicerMockRecorder) Get(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockServicer)(nil).Get), varargs...)
}

// Init mocks base method
func (m *MockServicer) Init(arg0 ...*types.Pair) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Init", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Init indicates an expected call of Init
func (mr *MockServicerMockRecorder) Init(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Init", reflect.TypeOf((*MockServicer)(nil).Init), arg0...)
}

// List mocks base method
func (m *MockServicer) List(arg0 ...*types.Pair) ([]storage.Storager, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "List", varargs...)
	ret0, _ := ret[0].([]storage.Storager)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List
func (mr *MockServicerMockRecorder) List(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockServicer)(nil).List), arg0...)
}
