package configs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig_SlackLoginByGithub(t *testing.T) {
	config := testingConfig()

	login, err := config.SlackLoginByGithub("github-login-1")
	assert.Equal(t, "slack-login-1", login)
	assert.Nil(t, err)

	login2, err2 := config.SlackLoginByGithub("not-found")
	assert.Equal(t, "not-found", login2)
	assert.Equal(t, errGithubLoginNotFound, err2)
}

func TestConfig_GithubLoginByGithub(t *testing.T) {
	config := testingConfig()

	login, err := config.GithubLoginBySlack("slack-login-1")
	assert.Equal(t, "github-login-1", login)
	assert.Nil(t, err)

	login2, err2 := config.GithubLoginBySlack("not-found")
	assert.Equal(t, "not-found", login2)
	assert.Equal(t, errSlackLoginNotFound, err2)
}

func TestConfig_ValidateConfig(t *testing.T) {
	config := NewConfig("", "", nil, "test-jiraProject", "brand/testing-repo")
	err := config.ValidateConfig()
	assert.Equal(t, errTokenMissing, err)

	config2 := NewConfig("123token", "", nil, "test-jiraProject", "brand/testing-repo")
	err2 := config2.ValidateConfig()
	assert.Equal(t, errWebhookUrlMissing, err2)

	config3 := NewConfig("123token", "webhook-url", nil, "", "brand/testing-repo")
	err3 := config3.ValidateConfig()
	assert.Equal(t, errJiraProjectMissing, err3)

	config4 := NewConfig("123token", "webhook-url", nil, "test-jiraProject", "")
	err4 := config4.ValidateConfig()
	assert.Equal(t, errGithubRepoMissing, err4)

	configValid := NewConfig("123token", "webhook-url", nil, "test-jiraProject", "brand/testing-repo")
	errNil := configValid.ValidateConfig()
	assert.Nil(t, errNil)
}

func TestConfig_Getters(t *testing.T) {
	config := testingConfig()

	reviewers := []Reviewer{
		{
			GhLogin: "github-login-1",
			SlLogin: "slack-login-1",
		},
		{
			GhLogin: "github-login-2",
			SlLogin: "slack-login-2",
		},
	}
	assert.Equal(t, "token", config.GithubToken())
	assert.Equal(t, "webhook-url", config.SlackWebhookUrl())
	assert.Equal(t, "test-jiraProject", config.JiraProject())
	assert.Equal(t, "brand/testing-repo", config.GithubRepo())
	assert.Equal(t, reviewers, config.Reviewers())
}

func testingConfig() ConfigInterface {
	return NewConfig(
		"token",
		"webhook-url",
		[]Reviewer{
			{
				GhLogin: "github-login-1",
				SlLogin: "slack-login-1",
			},
			{
				GhLogin: "github-login-2",
				SlLogin: "slack-login-2",
			},
		},
		"test-jiraProject",
		"brand/testing-repo",
	)
}
