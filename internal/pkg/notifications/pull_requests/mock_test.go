package pull_requests

import (
	"errors"
	"github.com/stretchr/testify/mock"
)

type configMockInterface struct {
	mock.Mock
}

func (_m *configMockInterface) SlackLoginByGithub(ghUsername string) (string, error) {
	ret := _m.Called(ghUsername)

	//error passing mock
	if ret.Get(1) != nil {
		return "", errors.New(ret.Get(1).(string))
	}

	return ret.Get(0).(string), nil
}

type senderClientMockInterface struct {
	mock.Mock
}

func (_m *senderClientMockInterface) Send(msg string) error {
	ret := _m.Called(msg)

	//error passing mock
	if ret.Get(0) != nil {
		return errors.New(ret.Get(0).(string))
	}

	return nil
}

func configMock() config {
	config := &configMockInterface{}

	config.On("SlackLoginByGithub", "first").Return("first-sender", nil).Once()
	config.On("SlackLoginByGithub", mock.Anything).Return("", "login not found")

	return config
}

var httpError = errors.New("http error")

func senderClientMock() notificationClient {
	slackClient := &senderClientMockInterface{}

	slackClient.On("Send", "errorThrown").Return(httpError.Error())
	slackClient.On("Send", mock.Anything).Return(nil)

	return slackClient
}
