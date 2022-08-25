package cache

import (
	"errors"
	"github.com/stretchr/testify/mock"
)

type dbMockInterface struct {
	mock.Mock
}

func newDbMock() *dbMockInterface {
	return &dbMockInterface{}
}

func (_m *dbMockInterface) Exec(query string, args ...any) error {
	x := append([]interface{}{query}, args...)
	ret := _m.Called(x...)

	//error passing mock
	if ret.Get(0) != nil {
		return errors.New(ret.Get(0).(string))
	}

	return nil
}

func (_m *dbMockInterface) QueryRow(destination interface{}, query string, args ...any) error {
	x := append([]interface{}{destination, query}, args...)
	ret := _m.Called(x...)

	//error passing mock
	if rf, ok := ret.Get(0).(func(destination interface{}, query string, args ...any) error); ok {
		return rf(destination, query, args...)
	}

	return ret.Error(0)
}

func (_m *dbMockInterface) IsEmptyResult(err error) bool {
	ret := _m.Called(err)

	return ret.Get(0).(bool)
}
