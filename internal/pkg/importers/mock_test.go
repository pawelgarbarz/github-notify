package importers

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/mock"
)

type githubClientMockInterface struct {
	mock.Mock
}

func (_m *githubClientMockInterface) Get(url string) ([]byte, error) {
	ret := _m.Called(url)

	//error passing mock
	if ret.Get(1) != nil {
		return nil, errors.New(ret.Get(1).(string))
	}

	return ret.Get(0).([]byte), nil
}

var httpError = errors.New("http error")

func clientMock(expected []byte) githubClient {
	client := &githubClientMockInterface{}

	ok := fmt.Sprintf(pullRequestUrlTemplate, "ok")
	err := fmt.Sprintf(pullRequestUrlTemplate, "error")

	client.On("Get", ok).Return(expected, nil)
	client.On("Get", err).Return(nil, httpError.Error())

	return client
}
