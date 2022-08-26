package cmd

import (
	"fmt"
	"github.com/pawelgarbarz/github-notify/internal/pkg/clients"
	"github.com/pawelgarbarz/github-notify/internal/pkg/debug"
	"github.com/pawelgarbarz/github-notify/internal/pkg/importers"
	"github.com/spf13/cobra"
	"log"
	"time"
)

// cacheClearAllCmd represents the pr command
var cacheClearAllCmd = &cobra.Command{
	Use:   "all",
	Short: "clear ALL local cache",
	Long:  `Clear ALL local persistent caches user for example to not spam channels.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !getConfig().CacheEnabled() {
			log.Println("[ERROR] Cache is not enabled! Stop execution of command")

			return
		}

		cache := initCache()

		err := cache.ClearAll()
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Cache cleared!")
	},
}

// cacheClearOutdatedCmd represents the pr command
var cacheClearOutdatedCmd = &cobra.Command{
	Use:   "outdated",
	Short: "clear outdated local cache",
	Long:  `Clear outdated local persistent caches user for example to not spam channels.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !getConfig().CacheEnabled() {
			log.Println("[ERROR] Cache is not enabled! Stop execution of command")

			return
		}

		cache := initCache()

		err := cache.ClearOutdated()
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Cache cleared")
	},
}

// cacheCmd represents the pr command
var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "clear local cache",
	Long:  `Clear local persistent caches user for example to not spam channels.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !getConfig().CacheEnabled() {
			log.Println("[ERROR] Cache is not enabled! Stop execution of command")

			return
		}

		err := cmd.Help()
		if err != nil {
			log.Printf("[Warning] Print help error: %s", err.Error())
		}
	},
}

// cacheClearCmd represents the pr command
var cacheClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "clear local cache",
	Long:  `Clear local persistent caches user for example to not spam channels.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !getConfig().CacheEnabled() {
			log.Println("[ERROR] Cache is not enabled! Stop execution of command")

			return
		}

		err := cmd.Help()
		if err != nil {
			log.Printf("[Warning] Print help error: %s", err.Error())
		}
	},
}

// cacheCommitCmd represents the pr command
var cacheCommitCmd = &cobra.Command{
	Use:   "commit",
	Short: "commit cache actions",
	Long:  `Commit cache actions.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !getConfig().CacheEnabled() {
			log.Println("[ERROR] Cache is not enabled! Stop execution of command")

			return
		}

		err := cmd.Help()
		if err != nil {
			log.Printf("[Warning] Print help error: %s", err.Error())
		}
	},
}

// cacheCommitWarmupCmd represents the pr command
var cacheCommitWarmupCmd = &cobra.Command{
	Use:   "warmup",
	Short: "warm-up commit cache from github",
	Long:  `Filter github commits and cache those which we are interested in.`,
	Run: func(cmd *cobra.Command, args []string) {
		debug := debug.NewDebug()
		if err := debug.SetLevel(debugLevel); err != nil {
			log.Fatal(err)

			return
		}

		githubClient := clients.NewGithubClient(getConfig().GithubToken())
		github := importers.NewGithub(githubClient, debug)

		cache := initCache()

		for _, reviewer := range getConfig().Reviewers() {
			commits, err := github.Commits(getConfig().GithubRepo(), getConfig().GithubBranch(), reviewer.GithubLogin())
			if err != nil {
				log.Fatal(err)

				return
			}

			for _, commit := range commits.List {
				value := fmt.Sprintf("SentAt: %s, Medium: none, Source: Cache CMD", time.Now().UTC().String())
				err := cache.Save(cacheCommitKey(getConfig().GithubRepo(), commit.Sha), value, commitCacheTTL)
				if err != nil {
					log.Printf("[Warning] Cache save error: %s", err.Error())
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(cacheCmd)

	cacheCmd.AddCommand(cacheClearCmd)
	cacheCmd.AddCommand(cacheCommitCmd)

	cacheClearCmd.AddCommand(cacheClearAllCmd)
	cacheClearCmd.AddCommand(cacheClearOutdatedCmd)

	cacheCommitCmd.AddCommand(cacheCommitWarmupCmd)
}
