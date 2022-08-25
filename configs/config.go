package configs

import (
	"errors"
	"time"
)

var errTokenMissing = errors.New("`token` configuration must be set")
var errWebhookUrlMissing = errors.New("`webhook-url` configuration must be set")
var errJiraProjectMissing = errors.New("`jira-project` configuration must be set")
var errGithubRepoMissing = errors.New("`github-repo` configuration must be set")
var errGithubLoginNotFound = errors.New("github login not found")
var errSlackLoginNotFound = errors.New("slack login not found")

type ConfigInterface interface {
	ValidateConfig() error
	GithubLoginBySlack(slackLogin string) (string, error)
	SlackLoginByGithub(githubLogin string) (string, error)
	Reviewers() []Reviewer
	GithubToken() string
	SlackWebhookUrl() string
	JiraProject() string
	GithubRepo() string
	CacheEnabled() bool
	CacheTTL() time.Duration
}

type Config struct {
	token         string
	webhookURL    string
	reviewersList []Reviewer
	jiraProject   string
	githubRepo    string
	cacheEnabled  bool
	cacheTTL      int
}

type Reviewer struct {
	GhLogin string
	SlLogin string
}

func NewConfig(token string, webhookURL string, reviewersList []Reviewer, jiraProject string, githubRepo string, cacheEnabled bool, cacheTTL int) ConfigInterface {
	return &Config{token: token, webhookURL: webhookURL, reviewersList: reviewersList, jiraProject: jiraProject, githubRepo: githubRepo, cacheEnabled: cacheEnabled, cacheTTL: cacheTTL}
}

func (c Config) ValidateConfig() error {
	if c.token == "" {
		return errTokenMissing
	}

	if c.webhookURL == "" {
		return errWebhookUrlMissing
	}

	if c.jiraProject == "" {
		return errJiraProjectMissing
	}

	if c.githubRepo == "" {
		return errGithubRepoMissing
	}

	return nil
}

func (c Config) GithubToken() string {
	return c.token
}

func (c Config) Reviewers() []Reviewer {
	return c.reviewersList
}

func (c Config) SlackLoginByGithub(githubLogin string) (string, error) {
	for _, reviewer := range c.reviewersList {
		if githubLogin == reviewer.GithubLogin() {
			return reviewer.SlackLogin(), nil
		}
	}

	return githubLogin, errGithubLoginNotFound
}

func (c Config) GithubLoginBySlack(slackLogin string) (string, error) {
	for _, reviewer := range c.reviewersList {
		if slackLogin == reviewer.SlackLogin() {
			return reviewer.GithubLogin(), nil
		}
	}

	return slackLogin, errSlackLoginNotFound
}

func (c Config) SlackWebhookUrl() string {
	return c.webhookURL
}

func (c Config) JiraProject() string {
	return c.jiraProject
}

func (c Config) GithubRepo() string {
	return c.githubRepo
}

func (r Reviewer) GithubLogin() string {
	return r.GhLogin
}

func (r Reviewer) SlackLogin() string {
	return r.SlLogin
}

func (c Config) CacheEnabled() bool {
	return c.cacheEnabled
}

func (c Config) CacheTTL() time.Duration {
	if c.cacheTTL > 0 {
		return time.Second * time.Duration(c.cacheTTL)
	}

	return 0
}
