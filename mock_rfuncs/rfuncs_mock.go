// Code generated by MockGen. DO NOT EDIT.
// Source: rfuncs/rfuncs.pb.go

// Package mock_rfuncs is a generated GoMock package.
package mock_rfuncs

import (
	gomock "github.com/golang/mock/gomock"
	rfuncs "github.com/yukithm/rfunc/rfuncs"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
	reflect "reflect"
)

// MockisClipboardContent_Content is a mock of isClipboardContent_Content interface
type MockisClipboardContent_Content struct {
	ctrl     *gomock.Controller
	recorder *MockisClipboardContent_ContentMockRecorder
}

// MockisClipboardContent_ContentMockRecorder is the mock recorder for MockisClipboardContent_Content
type MockisClipboardContent_ContentMockRecorder struct {
	mock *MockisClipboardContent_Content
}

// NewMockisClipboardContent_Content creates a new mock instance
func NewMockisClipboardContent_Content(ctrl *gomock.Controller) *MockisClipboardContent_Content {
	mock := &MockisClipboardContent_Content{ctrl: ctrl}
	mock.recorder = &MockisClipboardContent_ContentMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockisClipboardContent_Content) EXPECT() *MockisClipboardContent_ContentMockRecorder {
	return m.recorder
}

// isClipboardContent_Content mocks base method
func (m *MockisClipboardContent_Content) isClipboardContent_Content() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "isClipboardContent_Content")
}

// isClipboardContent_Content indicates an expected call of isClipboardContent_Content
func (mr *MockisClipboardContent_ContentMockRecorder) isClipboardContent_Content() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "isClipboardContent_Content", reflect.TypeOf((*MockisClipboardContent_Content)(nil).isClipboardContent_Content))
}

// MockRFuncsClient is a mock of RFuncsClient interface
type MockRFuncsClient struct {
	ctrl     *gomock.Controller
	recorder *MockRFuncsClientMockRecorder
}

// MockRFuncsClientMockRecorder is the mock recorder for MockRFuncsClient
type MockRFuncsClientMockRecorder struct {
	mock *MockRFuncsClient
}

// NewMockRFuncsClient creates a new mock instance
func NewMockRFuncsClient(ctrl *gomock.Controller) *MockRFuncsClient {
	mock := &MockRFuncsClient{ctrl: ctrl}
	mock.recorder = &MockRFuncsClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRFuncsClient) EXPECT() *MockRFuncsClientMockRecorder {
	return m.recorder
}

// Copy mocks base method
func (m *MockRFuncsClient) Copy(ctx context.Context, in *rfuncs.CopyRequest, opts ...grpc.CallOption) (*rfuncs.CopyReply, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Copy", varargs...)
	ret0, _ := ret[0].(*rfuncs.CopyReply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Copy indicates an expected call of Copy
func (mr *MockRFuncsClientMockRecorder) Copy(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Copy", reflect.TypeOf((*MockRFuncsClient)(nil).Copy), varargs...)
}

// Paste mocks base method
func (m *MockRFuncsClient) Paste(ctx context.Context, in *rfuncs.PasteRequest, opts ...grpc.CallOption) (*rfuncs.PasteReply, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Paste", varargs...)
	ret0, _ := ret[0].(*rfuncs.PasteReply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Paste indicates an expected call of Paste
func (mr *MockRFuncsClientMockRecorder) Paste(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Paste", reflect.TypeOf((*MockRFuncsClient)(nil).Paste), varargs...)
}

// OpenURL mocks base method
func (m *MockRFuncsClient) OpenURL(ctx context.Context, in *rfuncs.OpenURLRequest, opts ...grpc.CallOption) (*rfuncs.OpenURLReply, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "OpenURL", varargs...)
	ret0, _ := ret[0].(*rfuncs.OpenURLReply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OpenURL indicates an expected call of OpenURL
func (mr *MockRFuncsClientMockRecorder) OpenURL(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenURL", reflect.TypeOf((*MockRFuncsClient)(nil).OpenURL), varargs...)
}

// MockRFuncsServer is a mock of RFuncsServer interface
type MockRFuncsServer struct {
	ctrl     *gomock.Controller
	recorder *MockRFuncsServerMockRecorder
}

// MockRFuncsServerMockRecorder is the mock recorder for MockRFuncsServer
type MockRFuncsServerMockRecorder struct {
	mock *MockRFuncsServer
}

// NewMockRFuncsServer creates a new mock instance
func NewMockRFuncsServer(ctrl *gomock.Controller) *MockRFuncsServer {
	mock := &MockRFuncsServer{ctrl: ctrl}
	mock.recorder = &MockRFuncsServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRFuncsServer) EXPECT() *MockRFuncsServerMockRecorder {
	return m.recorder
}

// Copy mocks base method
func (m *MockRFuncsServer) Copy(arg0 context.Context, arg1 *rfuncs.CopyRequest) (*rfuncs.CopyReply, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Copy", arg0, arg1)
	ret0, _ := ret[0].(*rfuncs.CopyReply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Copy indicates an expected call of Copy
func (mr *MockRFuncsServerMockRecorder) Copy(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Copy", reflect.TypeOf((*MockRFuncsServer)(nil).Copy), arg0, arg1)
}

// Paste mocks base method
func (m *MockRFuncsServer) Paste(arg0 context.Context, arg1 *rfuncs.PasteRequest) (*rfuncs.PasteReply, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Paste", arg0, arg1)
	ret0, _ := ret[0].(*rfuncs.PasteReply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Paste indicates an expected call of Paste
func (mr *MockRFuncsServerMockRecorder) Paste(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Paste", reflect.TypeOf((*MockRFuncsServer)(nil).Paste), arg0, arg1)
}

// OpenURL mocks base method
func (m *MockRFuncsServer) OpenURL(arg0 context.Context, arg1 *rfuncs.OpenURLRequest) (*rfuncs.OpenURLReply, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenURL", arg0, arg1)
	ret0, _ := ret[0].(*rfuncs.OpenURLReply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OpenURL indicates an expected call of OpenURL
func (mr *MockRFuncsServerMockRecorder) OpenURL(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenURL", reflect.TypeOf((*MockRFuncsServer)(nil).OpenURL), arg0, arg1)
}
