package importers

import (
	debugPkg "github.com/pawelgarbarz/github-notify/internal/pkg/debug"
	"github.com/pawelgarbarz/github-notify/internal/pkg/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGithub_PullRequests_Get_Error(t *testing.T) {
	debug := debugPkg.NewDebug()
	_ = debug.SetLevel(debugPkg.SuperDetailed)
	prFetcher := NewGithub(clientMock(nil), debug)

	result, err := prFetcher.PullRequests("error")
	assert.Nil(t, result)
	assert.Equal(t, httpError, err)
}

func TestGithub_PullRequests_Fetch(t *testing.T) {
	debug := debugPkg.NewDebug()
	_ = debug.SetLevel(debugPkg.SuperDetailed)

	httpResult := `[{"html_url": "testUrl"},{"html_url": "testUrl-2"}]`
	byteResult := []byte(httpResult)

	prs := []models.PullRequest{
		{
			HTMLURL: "testUrl",
		},
		{
			HTMLURL: "testUrl-2",
		},
	}
	expected := models.NewPullRequestCollection(prs)

	prFetcher := NewGithub(clientMock(byteResult), debug)

	result, err := prFetcher.PullRequests("ok")
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}
func TestGithub_PullRequests_Unmarshall_Error(t *testing.T) {
	debug := debugPkg.NewDebug()
	_ = debug.SetLevel(debugPkg.SuperDetailed)

	httpResult := `[{"html_url..BROKEN...]`
	byteResult := []byte(httpResult)

	prFetcher := NewGithub(clientMock(byteResult), debug)

	result, err := prFetcher.PullRequests("ok")
	assert.Nil(t, result)
	assert.Equal(t, "unexpected end of JSON input", err.Error())
}
