package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "jiracrawler",
	Short: "Jira issue crawler CLI",
}

func Execute() error {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(getCmd)
	return rootCmd.Execute()
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.jiracrawler")
	_ = viper.ReadInConfig() // ignore error if config does not exist
}
