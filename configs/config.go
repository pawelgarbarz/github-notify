package configs

import (
	"errors"
)

type Config struct {
	Token         string     `mapstructure:"token"`
	WebhookURL    string     `mapstructure:"webhook-url"`
	ReviewersList []Reviewer `mapstructure:"reviewers"`
}

func NewConfig(token string, webhookURL string, reviewersList []Reviewer) Config {
	return Config{Token: token, WebhookURL: webhookURL, ReviewersList: reviewersList}
}

type Reviewer struct {
	GhLogin string `mapstructure:"github"`
	SlLogin string `mapstructure:"slack"`
}

var tokenMissingError = errors.New("`token` configuration must be set")
var webhookUrlMissingError = errors.New("`webhook-url` configuration must be set")

func (c Config) ValidateConfig() error {
	if c.Token == "" {
		return tokenMissingError
	}

	if c.WebhookURL == "" {
		return webhookUrlMissingError
	}

	return nil
}

func (c Config) GithubToken() string {
	return c.Token
}

func (c Config) Reviewers() []Reviewer {
	return c.ReviewersList
}

var githubLoginNotFoundErr = errors.New("github login not found")

func (c Config) SlackLoginByGithub(githubLogin string) (string, error) {
	for _, reviewer := range c.ReviewersList {
		if githubLogin == reviewer.GithubLogin() {
			return reviewer.SlackLogin(), nil
		}
	}

	return githubLogin, githubLoginNotFoundErr
}

func (c Config) GithubLoginBySlack(slackLogin string) (string, error) {
	for _, reviewer := range c.ReviewersList {
		if slackLogin == reviewer.SlackLogin() {
			return reviewer.GithubLogin(), nil
		}
	}

	return slackLogin, slackLoginNotFoundErr
}

var slackLoginNotFoundErr = errors.New("slack login not found")

func (c Config) SlackWebhookUrl() string {
	return c.WebhookURL
}

func (r Reviewer) GithubLogin() string {
	return r.GhLogin
}

func (r Reviewer) SlackLogin() string {
	return r.SlLogin
}
