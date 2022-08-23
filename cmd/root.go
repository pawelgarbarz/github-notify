package cmd

import (
	"fmt"
	"github.com/pawelgarbarz/github-notify/configs"
	"github.com/spf13/viper"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string
var config configs.Config

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

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".github-notify.yaml".
		viper.AddConfigPath(home)
		viper.AddConfigPath(".") // optionally look for config in the working directory
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

	config = configs.NewConfig(
		viper.GetString("token"),
		viper.GetString("webhook-url"),
		reviewers,
	)

	if err := config.ValidateConfig(); err != nil {
		log.Fatalf(err.Error())
	}
}

func getConfig() configs.Config {
	return config
}
