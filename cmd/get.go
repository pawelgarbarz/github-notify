package cmd

import (
	"github.com/pawelgarbarz/gh-pr-listener/internal/pkg/debug"
	"github.com/pawelgarbarz/gh-pr-listener/internal/pkg/http"
	"github.com/pawelgarbarz/gh-pr-listener/internal/pkg/importers"
	"github.com/pawelgarbarz/gh-pr-listener/internal/pkg/models"
	notifications "github.com/pawelgarbarz/gh-pr-listener/internal/pkg/notifications/pull_requests"
	"github.com/spf13/cobra"
	"log"
	"strings"
)

var projectCode string
var repoUrl string
var debugLevel int
var notifyList map[int]models.PullRequest

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get pull requests from github",
	Long: `FilterByCodeAndReviewers github pull requests and find those which titles consists JIRA project name.
	Than send message to slack channel with details about pending pull requests.`,
	Run: func(cmd *cobra.Command, args []string) {
		debug := debug.NewDebug()
		if err := debug.SetLevel(debugLevel); err != nil {
			log.Fatal(err)

			return
		}

		githubClient := http.NewGithubClient(getConfig().Token)
		github := importers.NewGithub(githubClient, debug)

		slackClient := http.NewSlackClient(getConfig().SlackWebhookUrl())
		sender := notifications.NewSender(slackClient, getConfig(), debug)

		pullRequests, err := github.PullRequests(repoUrl)
		if err != nil {
			log.Fatal(err)

			return
		}

		notifyList = make(map[int]models.PullRequest)
		filterByCode(pullRequests, notifyList)
		filterByReviewers(pullRequests, notifyList)

		send(notifyList, sender)
	},
}

func send(prToNotifyList map[int]models.PullRequest, sender notifications.NotificationSender) {
	for _, pullRequest := range prToNotifyList {
		msg := sender.PullRequestMessage(pullRequest)
		if err := sender.Send(msg); err != nil {
			log.Fatal(err)
		}
	}
}

func filterByCode(pullRequests *models.PullRequestCollection, output map[int]models.PullRequest) {
	for _, pullRequest := range pullRequests.Data {
		if strings.Contains(strings.ToLower(pullRequest.Title), strings.ToLower(projectCode)) {
			if _, exists := output[pullRequest.Number]; !exists {
				output[pullRequest.Number] = pullRequest
			}
		}
	}
}

func filterByReviewers(pullRequests *models.PullRequestCollection, output map[int]models.PullRequest) {
	for _, pullRequest := range pullRequests.Data {
		for _, reviewer := range getConfig().Reviewers() {
			if pullRequest.User.Login == reviewer.GithubLogin() {
				if _, exists := output[pullRequest.Number]; !exists {
					output[pullRequest.Number] = pullRequest
				}
			}
		}
	}
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().StringVarP(&projectCode, "projectCode", "p", "JIRA", "jira project code")
	getCmd.Flags().StringVarP(&repoUrl, "repoUrl", "u", "pawelgarbarz/github-notify", "repository brand/name")
	getCmd.Flags().IntVarP(&debugLevel, "debugLevel", "d", 0, "debug level 0..3")
}
