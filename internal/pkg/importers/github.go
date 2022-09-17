package importers

import (
	"encoding/json"
	"fmt"
	debugPkg "github.com/pawelgarbarz/github-notify/internal/pkg/debug"
	"github.com/pawelgarbarz/github-notify/internal/pkg/models"
	"time"
)

const pullRequestUrlTemplate = "https://api.github.com/repos/%s/pulls?state=open&sort=created&direction=desc"
const commitUrlTemplate = "https://api.github.com/repos/%s/commits?since=%s&per_page=500&page=1&author=%s&sha=%s"
const pullsUrlTemplate = "https://api.github.com/repos/%s/commits/%s/pulls"

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

func (g *github) Pulls(repositoryUrl string, sha string) (*models.PullRequestCollection, error) {
	url := fmt.Sprintf(pullsUrlTemplate, repositoryUrl, sha)

	if g.debug.Level() > debugPkg.NoDebug {
		fmt.Printf("CommitDetails URL: %s \n",
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

func (g *github) Commits(repositoryUrl string, branch string, author string) (*models.CommitCollection, error) {
	sinceStr := time.Now().AddDate(0, 0, -14).UTC().Format("2006-01-02")
	url := fmt.Sprintf(commitUrlTemplate, repositoryUrl, sinceStr, author, branch)

	if g.debug.Level() > debugPkg.NoDebug {
		fmt.Printf("CommitDetails URL: %s \n",
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

	var data []models.Commit
	err = json.Unmarshal(responseBody, &data)
	if err != nil {
		return nil, err
	}

	return models.NewCommitCollection(data), nil
}
