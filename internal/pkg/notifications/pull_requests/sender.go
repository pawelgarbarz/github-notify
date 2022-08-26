package pull_requests

import (
	"fmt"
	debugPkg "github.com/pawelgarbarz/github-notify/internal/pkg/debug"
	"github.com/pawelgarbarz/github-notify/internal/pkg/models"
	"math"
	"regexp"
	"strings"
	"time"
)

type NotificationSender interface {
	Send(msg string) error
	PullRequestMessage(pullRequest models.PullRequest) string
	CommitMessage(pullRequest models.Commit, branch string) string
}

type sender struct {
	client  notificationClient
	debug   debug
	configs config
}

func NewSender(client notificationClient, config config, debug debug) *sender {
	return &sender{
		client:  client,
		configs: config,
		debug:   debug,
	}
}

func (s *sender) Send(msg string) error {
	if s.debug.Level() >= debugPkg.Detailed {
		fmt.Println(msg)
	}

	err := s.client.Send(msg)
	if err != nil {
		return fmt.Errorf("pr send error: %s", err)
	}

	return nil
}

func (s *sender) PullRequestMessage(pullRequest models.PullRequest) string {

	text := fmt.Sprintf(
		"<%s|Pull request> pending for review\n>*<%s|%s>*\n>Created `%s` ago\n>Author: %s",
		pullRequest.HTMLURL,
		pullRequest.HTMLURL,
		pullRequest.Title,
		s.calculateAge(pullRequest.CreatedAt),
		s.slackUsername(pullRequest.User.Login),
	)

	if pullRequest.Body != "" {
		re := regexp.MustCompile(`\r?\n`)
		summary := re.ReplaceAllString(pullRequest.Body, "\n> ")

		text = text + fmt.Sprintf(
			"\n> %s",
			summary,
		)
	}

	if len(pullRequest.RequestedReviewers) > 0 {
		text = text + fmt.Sprintf(
			"\n>Reviewers: %s",
			s.reviewers(pullRequest),
		)
	}

	return text
}

func (s *sender) CommitMessage(commitData models.Commit, branch string) string {
	return fmt.Sprintf(
		"A new <%s|commit> have been pushed to %s branch\n>Message: %s\n>Created `%s` ago\n>Author: %s",
		commitData.HTMLURL,
		branch,
		commitData.Commit.Message,
		s.calculateAge(commitData.Commit.Author.Date),
		s.slackUsername(commitData.Author.Login),
	)
}

func (s *sender) calculateAge(createdAt time.Time) string {
	now := time.Now()
	diff := now.Sub(createdAt)

	days := math.Floor(diff.Hours() / 24)
	hours := math.Floor(diff.Hours() - (days * 24))
	minutes := math.Floor(diff.Minutes() - (days * 24 * 60) - (hours * 60))

	return fmt.Sprintf("%v days %v hours %v minutes", days, hours, minutes)
}

func (s *sender) reviewers(pullRequest models.PullRequest) string {
	var slackReviewers []string

	for _, reviewer := range pullRequest.RequestedReviewers {
		slackReviewers = append(slackReviewers, s.slackUsername(reviewer.Login))
	}

	return strings.Join(slackReviewers, ", ")
}

func (s *sender) slackUsername(githubUsername string) string {
	if slackUsername, err := s.configs.SlackLoginByGithub(githubUsername); err == nil {
		return fmt.Sprintf("<@%s>", slackUsername)
	}

	return githubUsername
}
