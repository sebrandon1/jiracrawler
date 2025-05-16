package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFile string

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

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	configFile = ".jiracrawler-config.yaml"
	viper.SetConfigType("yaml")
	viper.SetConfigFile(configFile)

	viper.AutomaticEnv()
	viper.SetEnvPrefix("JIRACRAWLER")

	// If the config file is not found, create it
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		err = viper.WriteConfig()
		if err != nil {
			panic(err)
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config file: ", err)
	}
}
