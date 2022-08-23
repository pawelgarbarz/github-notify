package importers

import (
	"encoding/json"
	"fmt"
	debugPkg "github.com/pawelgarbarz/github-notify/internal/pkg/debug"
	"github.com/pawelgarbarz/github-notify/internal/pkg/models"
)

const pullRequestUrlTemplate = "https://api.github.com/repos/%s/pulls?state=open&sort=created&direction=desc"

type github struct {
	githubClient githubClient
	debug        debug
}

func NewGithub(githubClient githubClient, debug debug) *github {
	return &github{githubClient, debug}
}

func (g *github) PullRequests(repositoryUrl string) (*models.PullRequestCollection, error) {
	url := fmt.Sprintf(pullRequestUrlTemplate, repositoryUrl)

	if g.debug.Level() > debugPkg.NoDebug {
		fmt.Printf("PR URL: %s \n",
			url,
		)
	}

	responseBody, err := g.githubClient.Get(url)
	if err != nil {
		return nil, err
	}

	if g.debug.Level() >= debugPkg.SuperDetailed {
		fmt.Printf("Response: %s \n", responseBody)
	}

	var data []models.PullRequest
	err = json.Unmarshal(responseBody, &data)
	if err != nil {
		return nil, err
	}

	return models.NewPullRequestCollection(data), nil
}
