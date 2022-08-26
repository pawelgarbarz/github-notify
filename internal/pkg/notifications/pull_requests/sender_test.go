package pull_requests

import (
	"fmt"
	debugPkg "github.com/pawelgarbarz/github-notify/internal/pkg/debug"
	"github.com/pawelgarbarz/github-notify/internal/pkg/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var expectedPullRequestMsg = "<https://acme.testing|Pull request> pending for review\n" +
	">*<https://acme.testing|unit test PR>*\n" +
	">Created `0 days 0 hours 0 minutes` ago\n" +
	">Author: testing-dev"

var expectedCommitMsg = "A new <https://commit.com|commit> have been pushed to test-branch branch\n" +
	">Message: test-commit\n" +
	">Created `1 days 0 hours 0 minutes` ago\n" +
	">Author: pawelgarbarz"

func TestGetMessageWithoutReviewers(t *testing.T) {
	debug := debugPkg.NewDebug()

	notify := NewSender(senderClientMock(), configMock(), debug)

	assert.Equal(t, expectedPullRequestMsg, notify.PullRequestMessage(prWithoutReviewers()))
}

func TestGetMessageWithReviewers(t *testing.T) {
	debug := debugPkg.NewDebug()
	notify := NewSender(senderClientMock(), configMock(), debug)

	expectedMsgWithReviewers := expectedPullRequestMsg + "\n>Reviewers: <@first-sender>, second"

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

func commitModel() models.Commit {
	return models.Commit{
		HTMLURL: "https://commit.com",
		Commit: models.CommitDetails{
			Message: "test-commit",
			Author: models.UserShort{
				Date: time.Now().AddDate(0, 0, -1),
			},
		},
		Author: models.User{
			Login: "pawelgarbarz",
		},
	}
}

func TestGetCommitMessage(t *testing.T) {
	debug := debugPkg.NewDebug()

	notify := NewSender(senderClientMock(), configMock(), debug)

	assert.Equal(t, expectedCommitMsg, notify.CommitMessage(commitModel(), "test-branch"))
}
