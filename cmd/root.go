package cmd

import (
	"fmt"
	"github.com/pawelgarbarz/github-notify/configs"
	"github.com/pawelgarbarz/github-notify/internal/pkg/cache"
	"github.com/pawelgarbarz/github-notify/internal/pkg/clients"
	"github.com/spf13/viper"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var jiraProjectOpt string
var githubRepoOpt string
var cfgFile string
var config configs.ConfigInterface
var debugLevel int

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "github-notify",
	Short: "Application to interact with github and send notifications to communicators",
	Long:  `Application to interact with github and send notifications to communicators`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.github-notify.yaml)")
	rootCmd.PersistentFlags().StringVarP(&jiraProjectOpt, "jira-project", "p", "", "jira project code")
	rootCmd.PersistentFlags().StringVarP(&githubRepoOpt, "github-repo", "u", "", "github repository brand/name")

	rootCmd.PersistentFlags().IntVarP(&debugLevel, "debugLevel", "d", 0, "debug level 0..3")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	viper.SetDefault("cache-enabled", true)
	viper.SetDefault("cache-ttl", 14400) // 4h = 60 (seconds) * 60 (minutes) * 4 (hours)

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		cwd, err := os.Getwd()
		cobra.CheckErr(err)

		// Search config in home directory with name ".github-notify.yaml".
		viper.AddConfigPath(home)
		viper.AddConfigPath(cwd) // optionally look for config in the working directory
		viper.SetConfigType("yaml")
		viper.SetConfigName(".github-notify")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		_, _ = fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else {
		_, _ = fmt.Fprintln(os.Stderr, "config file not fount, expected in:", viper.ConfigFileUsed())
	}

	var reviewers []configs.Reviewer
	err := viper.UnmarshalKey("reviewers", &reviewers)
	if err != nil {
		log.Fatal(err)
	}

	jiraProject := viper.GetString("jira-project")
	if jiraProjectOpt != "" {
		jiraProject = jiraProjectOpt
	}

	githubRepo := viper.GetString("github-repo")
	if githubRepoOpt != "" {
		githubRepo = githubRepoOpt
	}

	config = configs.NewConfig(
		viper.GetString("token"),
		viper.GetString("webhook-url"),
		reviewers,
		jiraProject,
		githubRepo,
		viper.GetBool("cache-enabled"),
		viper.GetInt("cache-ttl"),
	)

	if err := config.ValidateConfig(); err != nil {
		log.Fatalf(err.Error())
	}
}

func getConfig() configs.ConfigInterface {
	return config
}

func initCache() cache.Cache {
	db, err := clients.NewSQLiteClient()
	if err != nil {
		log.Fatalf("Cannot create SQL driver: %s", err.Error())
	}

	return cache.NewCache(db)
}
