package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func SetConfigValue(key, value string) {
	viper.Set(key, value)
	err := viper.WriteConfig()
	if err != nil {
		panic(err)
	}
}

func GetConfigValue(key string) string {
	return viper.GetString(key)
}

func UpdateConfigValue(key, value string) {
	viper.Set(key, value)
	err := viper.WriteConfig()
	if err != nil {
		panic(err)
	}
}

var defaultJiraURL = "https://issues.redhat.com"

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Set or get configuration values",
}

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set a key value pair to the configuration",
	Run: func(cmd *cobra.Command, args []string) {
		token, _ := cmd.Flags().GetString("token")
		url, _ := cmd.Flags().GetString("url")
		user, _ := cmd.Flags().GetString("user")

		if user != "" {
			UpdateConfigValue("jira_user", user)
			fmt.Println("Set Jira user in config")
		}

		if token != "" {
			UpdateConfigValue("apikey", token)
			fmt.Println("Set Jira API token in config")
		}
		if url == "" {
			url = defaultJiraURL
		}
		if url != "" {

			if url != defaultJiraURL {
				fmt.Printf("Set Jira custom URL in config: %s", url)
			}
			UpdateConfigValue("jira_url", url)
		}
	},
}

var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "View the configuration",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Current configuration:")
		settings := viper.AllSettings()
		for k, v := range settings {
			fmt.Printf("%s: %v\n", k, v)
		}
	},
}

func init() {
	setCmd.PersistentFlags().StringP("user", "s", "", "The Jira user email to set in the configuration.")
	setCmd.PersistentFlags().StringP("token", "t", "", "The Jira API token to set in the configuration.")
	setCmd.PersistentFlags().StringP("url", "u", "", "The Jira URL to set in the configuration.")

	configCmd.AddCommand(setCmd)
	configCmd.AddCommand(viewCmd)

	// Add config to root command in root.go
}
