package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "jiracrawler",
	Short: "Jira issue crawler CLI",
}

func Execute() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(getCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.jiracrawler")
	_ = viper.ReadInConfig() // ignore error if config does not exist
}
