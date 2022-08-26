package importers

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/mock"
	"time"
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

var errHttp = errors.New("clients error")

func clientMock(expected []byte) githubClient {
	client := &githubClientMockInterface{}

	ok := fmt.Sprintf(pullRequestUrlTemplate, "ok")
	err := fmt.Sprintf(pullRequestUrlTemplate, "error")

	sinceStr := time.Now().AddDate(0, 0, -14).UTC().Format("2006-01-02")
	okCommit := fmt.Sprintf(commitUrlTemplate, "ok", sinceStr, "pawelgarbarz", "main")
	errCommit := fmt.Sprintf(commitUrlTemplate, "error", sinceStr, "pawelgarbarz", "main")

	client.On("Get", ok).Return(expected, nil)
	client.On("Get", err).Return(nil, errHttp.Error())

	client.On("Get", okCommit).Return(expected, nil)
	client.On("Get", errCommit).Return(nil, errHttp.Error())

	return client
}
