// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/maxsnegir/url-shortener/internal/storage (interfaces: ShortenerStorage)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	storage "github.com/maxsnegir/url-shortener/internal/storage"
)

// MockShortenerStorage is a mock of ShortenerStorage interface.
type MockShortenerStorage struct {
	ctrl     *gomock.Controller
	recorder *MockShortenerStorageMockRecorder
}

// MockShortenerStorageMockRecorder is the mock recorder for MockShortenerStorage.
type MockShortenerStorageMockRecorder struct {
	mock *MockShortenerStorage
}

// NewMockShortenerStorage creates a new mock instance.
func NewMockShortenerStorage(ctrl *gomock.Controller) *MockShortenerStorage {
	mock := &MockShortenerStorage{ctrl: ctrl}
	mock.recorder = &MockShortenerStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockShortenerStorage) EXPECT() *MockShortenerStorageMockRecorder {
	return m.recorder
}

// GetOriginalURL mocks base method.
func (m *MockShortenerStorage) GetOriginalURL(arg0 context.Context, arg1 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOriginalURL", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOriginalURL indicates an expected call of GetOriginalURL.
func (mr *MockShortenerStorageMockRecorder) GetOriginalURL(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOriginalURL", reflect.TypeOf((*MockShortenerStorage)(nil).GetOriginalURL), arg0, arg1)
}

// GetUserURLs mocks base method.
func (m *MockShortenerStorage) GetUserURLs(arg0 context.Context, arg1 string) ([]storage.URLData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserURLs", arg0, arg1)
	ret0, _ := ret[0].([]storage.URLData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserURLs indicates an expected call of GetUserURLs.
func (mr *MockShortenerStorageMockRecorder) GetUserURLs(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserURLs", reflect.TypeOf((*MockShortenerStorage)(nil).GetUserURLs), arg0, arg1)
}

// Ping mocks base method.
func (m *MockShortenerStorage) Ping(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockShortenerStorageMockRecorder) Ping(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockShortenerStorage)(nil).Ping), arg0)
}

// SaveData mocks base method.
func (m *MockShortenerStorage) SaveData(arg0 context.Context, arg1 string, arg2 storage.URLData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveData", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveData indicates an expected call of SaveData.
func (mr *MockShortenerStorageMockRecorder) SaveData(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveData", reflect.TypeOf((*MockShortenerStorage)(nil).SaveData), arg0, arg1, arg2)
}

// SaveDataBatch mocks base method.
func (m *MockShortenerStorage) SaveDataBatch(arg0 context.Context, arg1 string, arg2 []storage.URLData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveDataBatch", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveDataBatch indicates an expected call of SaveDataBatch.
func (mr *MockShortenerStorageMockRecorder) SaveDataBatch(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveDataBatch", reflect.TypeOf((*MockShortenerStorage)(nil).SaveDataBatch), arg0, arg1, arg2)
}

// SetShortURL mocks base method.
func (m *MockShortenerStorage) SetShortURL(arg0 storage.URLData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetShortURL", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetShortURL indicates an expected call of SetShortURL.
func (mr *MockShortenerStorageMockRecorder) SetShortURL(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetShortURL", reflect.TypeOf((*MockShortenerStorage)(nil).SetShortURL), arg0)
}

// SetUserURL mocks base method.
func (m *MockShortenerStorage) SetUserURL(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetUserURL", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetUserURL indicates an expected call of SetUserURL.
func (mr *MockShortenerStorageMockRecorder) SetUserURL(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetUserURL", reflect.TypeOf((*MockShortenerStorage)(nil).SetUserURL), arg0, arg1)
}

// Shutdown mocks base method.
func (m *MockShortenerStorage) Shutdown(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Shutdown", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Shutdown indicates an expected call of Shutdown.
func (mr *MockShortenerStorageMockRecorder) Shutdown(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Shutdown", reflect.TypeOf((*MockShortenerStorage)(nil).Shutdown), arg0)
}
