// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	domain "middleman-capstone/domain"

	mock "github.com/stretchr/testify/mock"
)

// ProductData is an autogenerated mock type for the ProductData type
type ProductData struct {
	mock.Mock
}

// CreateProductData provides a mock function with given fields: newProduct
func (_m *ProductData) CreateProductData(newProduct domain.Product) domain.Product {
	ret := _m.Called(newProduct)

	var r0 domain.Product
	if rf, ok := ret.Get(0).(func(domain.Product) domain.Product); ok {
		r0 = rf(newProduct)
	} else {
		r0 = ret.Get(0).(domain.Product)
	}

	return r0
}

// DeleteProductData provides a mock function with given fields: productid
func (_m *ProductData) DeleteProductData(productid int) (int, error) {
	ret := _m.Called(productid)

	var r0 int
	if rf, ok := ret.Get(0).(func(int) int); ok {
		r0 = rf(productid)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(productid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllProductData provides a mock function with given fields: limit, offset
func (_m *ProductData) GetAllProductData(limit int, offset int) ([]domain.Product, error) {
	ret := _m.Called(limit, offset)

	var r0 []domain.Product
	if rf, ok := ret.Get(0).(func(int, int) []domain.Product); ok {
		r0 = rf(limit, offset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Product)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, int) error); ok {
		r1 = rf(limit, offset)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SearchRestoData provides a mock function with given fields: search
func (_m *ProductData) SearchRestoData(search string) ([]domain.Product, error) {
	ret := _m.Called(search)

	var r0 []domain.Product
	if rf, ok := ret.Get(0).(func(string) []domain.Product); ok {
		r0 = rf(search)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Product)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(search)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateProductData provides a mock function with given fields: data, idProduct
func (_m *ProductData) UpdateProductData(data map[string]interface{}, idProduct int) (int, error) {
	ret := _m.Called(data, idProduct)

	var r0 int
	if rf, ok := ret.Get(0).(func(map[string]interface{}, int) int); ok {
		r0 = rf(data, idProduct)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(map[string]interface{}, int) error); ok {
		r1 = rf(data, idProduct)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewProductData interface {
	mock.TestingT
	Cleanup(func())
}

// NewProductData creates a new instance of ProductData. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewProductData(t mockConstructorTestingTNewProductData) *ProductData {
	mock := &ProductData{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
