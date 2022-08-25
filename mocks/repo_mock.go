// Code generated by MockGen. DO NOT EDIT.
// Source: ./storage/repo.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	product "github.com/NeverlandMJ/arzon-market/pkg/product"
	store "github.com/NeverlandMJ/arzon-market/pkg/store"
	user "github.com/NeverlandMJ/arzon-market/pkg/user"
	gomock "github.com/golang/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// AddCard mocks base method.
func (m *MockRepository) AddCard(ctx context.Context, c user.Card) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCard", ctx, c)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddCard indicates an expected call of AddCard.
func (mr *MockRepositoryMockRecorder) AddCard(ctx, c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCard", reflect.TypeOf((*MockRepository)(nil).AddCard), ctx, c)
}

// AddProduct mocks base method.
func (m *MockRepository) AddProduct(ctx context.Context, p product.Product) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddProduct", ctx, p)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddProduct indicates an expected call of AddProduct.
func (mr *MockRepositoryMockRecorder) AddProduct(ctx, p interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddProduct", reflect.TypeOf((*MockRepository)(nil).AddProduct), ctx, p)
}

// AddProducts mocks base method.
func (m *MockRepository) AddProducts(ctx context.Context, ps []product.Product) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddProducts", ctx, ps)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddProducts indicates an expected call of AddProducts.
func (mr *MockRepositoryMockRecorder) AddProducts(ctx, ps interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddProducts", reflect.TypeOf((*MockRepository)(nil).AddProducts), ctx, ps)
}

// AddUser mocks base method.
func (m *MockRepository) AddUser(ctx context.Context, u user.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUser", ctx, u)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddUser indicates an expected call of AddUser.
func (mr *MockRepositoryMockRecorder) AddUser(ctx, u interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUser", reflect.TypeOf((*MockRepository)(nil).AddUser), ctx, u)
}

// GetBalance mocks base method.
func (m *MockRepository) GetBalance(ownerID string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBalance", ownerID)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBalance indicates an expected call of GetBalance.
func (mr *MockRepositoryMockRecorder) GetBalance(ownerID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBalance", reflect.TypeOf((*MockRepository)(nil).GetBalance), ownerID)
}

// GetProduct mocks base method.
func (m *MockRepository) GetProduct(ctx context.Context, name string) (product.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProduct", ctx, name)
	ret0, _ := ret[0].(product.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProduct indicates an expected call of GetProduct.
func (mr *MockRepositoryMockRecorder) GetProduct(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProduct", reflect.TypeOf((*MockRepository)(nil).GetProduct), ctx, name)
}

// GetUser mocks base method.
func (m *MockRepository) GetUser(ctx context.Context, email, pw string) (user.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", ctx, email, pw)
	ret0, _ := ret[0].(user.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockRepositoryMockRecorder) GetUser(ctx, email, pw interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockRepository)(nil).GetUser), ctx, email, pw)
}

// GetUserByID mocks base method.
func (m *MockRepository) GetUserByID(id string) (user.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByID", id)
	ret0, _ := ret[0].(user.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByID indicates an expected call of GetUserByID.
func (mr *MockRepositoryMockRecorder) GetUserByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByID", reflect.TypeOf((*MockRepository)(nil).GetUserByID), id)
}

// ListProducts mocks base method.
func (m *MockRepository) ListProducts(ctx context.Context) ([]product.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListProducts", ctx)
	ret0, _ := ret[0].([]product.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListProducts indicates an expected call of ListProducts.
func (mr *MockRepositoryMockRecorder) ListProducts(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListProducts", reflect.TypeOf((*MockRepository)(nil).ListProducts), ctx)
}

// ListUsers mocks base method.
func (m *MockRepository) ListUsers(ctx context.Context) ([]user.UserCard, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListUsers", ctx)
	ret0, _ := ret[0].([]user.UserCard)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListUsers indicates an expected call of ListUsers.
func (mr *MockRepositoryMockRecorder) ListUsers(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListUsers", reflect.TypeOf((*MockRepository)(nil).ListUsers), ctx)
}

// SellProduct mocks base method.
func (m *MockRepository) SellProduct(ctx context.Context, sale store.Sales, product product.Product) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SellProduct", ctx, sale, product)
	ret0, _ := ret[0].(error)
	return ret0
}

// SellProduct indicates an expected call of SellProduct.
func (mr *MockRepositoryMockRecorder) SellProduct(ctx, sale, product interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SellProduct", reflect.TypeOf((*MockRepository)(nil).SellProduct), ctx, sale, product)
}