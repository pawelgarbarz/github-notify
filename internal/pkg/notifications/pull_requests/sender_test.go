package pull_requests

import (
	"fmt"
	debugPkg "github.com/pawelgarbarz/github-notify/internal/pkg/debug"
	"github.com/pawelgarbarz/github-notify/internal/pkg/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var expectedMsg = "<https://acme.testing|Pull request> pending for review\n" +
	">*<https://acme.testing|unit test PR>*\n" +
	">Created `0 days 0 hours 0 minutes` ago\n" +
	">Author: testing-dev"

func TestGetMessageWithoutReviewers(t *testing.T) {
	debug := debugPkg.NewDebug()

	notify := NewSender(senderClientMock(), configMock(), debug)

	assert.Equal(t, expectedMsg, notify.PullRequestMessage(prWithoutReviewers()))
}

func TestGetMessageWithReviewers(t *testing.T) {
	debug := debugPkg.NewDebug()
	notify := NewSender(senderClientMock(), configMock(), debug)

	expectedMsgWithReviewers := expectedMsg + "\n>Reviewers: <@first-sender>, second"

	pr := prWithoutReviewers()
	pr.RequestedReviewers = []models.User{
		{Login: "first"},
		{Login: "second"},
	}

	assert.Equal(t, expectedMsgWithReviewers, notify.PullRequestMessage(pr))
}

func TestGetMessageWithSummary(t *testing.T) {
	debug := debugPkg.NewDebug()

	notify := NewSender(senderClientMock(), configMock(), debug)

	pr := prWithoutReviewers()
	pr.Body = "FirstLine\r\nSecondLine"

	expectedMsg := "<https://acme.testing|Pull request> pending for review\n" +
		">*<https://acme.testing|unit test PR>*\n" +
		">Created `0 days 0 hours 0 minutes` ago\n" +
		">Author: testing-dev\n" +
		"> FirstLine\n" +
		"> SecondLine"

	assert.Equal(t, expectedMsg, notify.PullRequestMessage(pr))
}

func TestSendWithoutError(t *testing.T) {
	debug := debugPkg.NewDebug()

	notify := NewSender(senderClientMock(), configMock(), debug)

	result := notify.Send("test msg")

	assert.Nil(t, result)
}

func TestSendWithError(t *testing.T) {
	debug := debugPkg.NewDebug()
	_ = debug.SetLevel(debugPkg.Detailed)

	notify := NewSender(senderClientMock(), configMock(), debug)

	result := notify.Send("errorThrown")

	assert.Equal(t, fmt.Errorf("pr send error: %s", errHttp), result)
}

func prWithoutReviewers() models.PullRequest {
	return models.PullRequest{
		HTMLURL:   "https://acme.testing",
		Title:     "unit test PR",
		CreatedAt: time.Now(),
		User: models.User{
			Login: "testing-dev",
		},
	}
}
