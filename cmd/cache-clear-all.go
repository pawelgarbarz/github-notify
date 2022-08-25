package cmd

import (
	"github.com/spf13/cobra"
	"log"
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

func init() {
	rootCmd.AddCommand(cacheCmd)
	cacheCmd.AddCommand(cacheClearCmd)
	cacheClearCmd.AddCommand(cacheClearAllCmd)
	cacheClearCmd.AddCommand(cacheClearOutdatedCmd)
}
