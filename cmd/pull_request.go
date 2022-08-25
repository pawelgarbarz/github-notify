package cmd

import (
	"fmt"
	"github.com/pawelgarbarz/github-notify/internal/pkg/cache"
	"github.com/pawelgarbarz/github-notify/internal/pkg/clients"
	"github.com/pawelgarbarz/github-notify/internal/pkg/debug"
	"github.com/pawelgarbarz/github-notify/internal/pkg/importers"
	"github.com/pawelgarbarz/github-notify/internal/pkg/models"
	notifications "github.com/pawelgarbarz/github-notify/internal/pkg/notifications/pull_requests"
	"github.com/spf13/cobra"
	"log"
	"strings"
	"time"
)

var notifyList map[int]models.PullRequest

// prCmd represents the pr command
var prCmd = &cobra.Command{
	Use:   "pull-request",
	Short: "get pull requests from github",
	Long: `FilterByCodeAndReviewers github pull requests and find those which titles consists JIRA project name.
	Than send message to slack channel with details about pending pull requests.`,
	Run: func(cmd *cobra.Command, args []string) {
		debug := debug.NewDebug()
		if err := debug.SetLevel(debugLevel); err != nil {
			log.Fatal(err)

			return
		}

		githubClient := clients.NewGithubClient(getConfig().GithubToken())
		github := importers.NewGithub(githubClient, debug)

		slackClient := clients.NewSlackClient(getConfig().SlackWebhookUrl())
		sender := notifications.NewSender(slackClient, getConfig(), debug)

		pullRequests, err := github.PullRequests(getConfig().GithubRepo())
		if err != nil {
			log.Fatal(err)

			return
		}

		notifyList = make(map[int]models.PullRequest)

		cache := initCache()
		filterByCode(pullRequests, notifyList, cache)
		filterByReviewers(pullRequests, notifyList, cache)

		send(notifyList, sender, cache)
	},
}

func send(prToNotifyList map[int]models.PullRequest, sender notifications.NotificationSender, cache cache.Cache) {
	for _, pullRequest := range prToNotifyList {
		msg := sender.PullRequestMessage(pullRequest)
		if err := sender.Send(msg); err != nil {
			log.Fatal(err)
		}

		value := fmt.Sprintf("SentAt: %s, Medium: Slack", time.Now().UTC().String())
		err := cache.Save(cacheKey(getConfig().GithubRepo(), pullRequest.Number), value, getConfig().CacheTTL())
		if err != nil {
			log.Printf("[Warning] Cache save error: %s", err.Error())
		}
	}
}

func filterByCode(pullRequests *models.PullRequestCollection, output map[int]models.PullRequest, cache cache.Cache) {
	for _, pullRequest := range pullRequests.Data {
		if strings.Contains(strings.ToLower(pullRequest.Title), strings.ToLower(getConfig().JiraProject())) {
			if _, exists := output[pullRequest.Number]; !exists {
				cacheExists, _ := cache.Exists(cacheKey(getConfig().GithubRepo(), pullRequest.Number))
				if cacheExists {
					continue
				}

				output[pullRequest.Number] = pullRequest
			}
		}
	}
}

func filterByReviewers(pullRequests *models.PullRequestCollection, output map[int]models.PullRequest, cache cache.Cache) {
	for _, pullRequest := range pullRequests.Data {
		for _, reviewer := range getConfig().Reviewers() {
			if pullRequest.User.Login == reviewer.GithubLogin() {
				if _, exists := output[pullRequest.Number]; !exists {
					cacheExists, _ := cache.Exists(cacheKey(getConfig().GithubRepo(), pullRequest.Number))
					if cacheExists {
						continue
					}

					output[pullRequest.Number] = pullRequest
				}
			}
		}
	}
}

func cacheKey(repo string, id int) string {
	return fmt.Sprintf("pull-request/%s/%d", repo, id)
}

func init() {
	rootCmd.AddCommand(prCmd)
}
