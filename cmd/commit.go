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
	"time"
)

var commitCacheTTL = time.Hour * 24 * 14

// commitCmd represents the pr command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "get commits from github",
	Long: `Filter github commits and find those which titles consists JIRA project name.
	Than send message to slack chanel with details about pending pull requests.`,
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

		cache := initCache()

		for _, reviewer := range getConfig().Reviewers() {
			commits, err := github.Commits(getConfig().GithubRepo(), getConfig().GithubBranch(), reviewer.GithubLogin())
			if err != nil {
				log.Fatal(err)

				return
			}

			for _, commit := range commits.List {
				cacheExists, _ := cache.Exists(cacheCommitKey(getConfig().GithubRepo(), commit.Sha))
				if cacheExists {
					continue
				}

				pullCommits, err := github.Pulls(getConfig().GithubRepo(), commit.Sha)
				if err != nil {
					log.Fatal(err)

					return
				}

				if len(pullCommits.Data) > 0 { //means that this commit comes from pull request branch
					// skip commits form PR's

					value := fmt.Sprintf("SaveAt: %s, merge commit", time.Now().UTC().String())
					err := cache.Save(cacheCommitKey(getConfig().GithubRepo(), commit.Sha), value, commitCacheTTL)
					if err != nil {
						log.Printf("[Warning] Cache save error: %s", err.Error())
					}

					continue
				}

				sendCommit(commit, sender, cache)
			}
		}
	},
}

func sendCommit(commitData models.Commit, sender notifications.NotificationSender, cache cache.Cache) {
	msg := sender.CommitMessage(commitData, getConfig().GithubBranch())
	if err := sender.Send(msg); err != nil {
		log.Fatal(err)
	}

	value := fmt.Sprintf("SentAt: %s, Medium: Slack", time.Now().UTC().String())
	err := cache.Save(cacheCommitKey(getConfig().GithubRepo(), commitData.Sha), value, commitCacheTTL)
	if err != nil {
		log.Printf("[Warning] Cache save error: %s", err.Error())
	}
}

func cacheCommitKey(repo string, sha string) string {
	return fmt.Sprintf("commits/%s/%s", repo, sha)
}

func init() {
	rootCmd.AddCommand(commitCmd)
}
