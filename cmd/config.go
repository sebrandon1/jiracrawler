package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
}

var setCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a config value",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := args[1]
		viper.Set(key, value)
		configDir := os.ExpandEnv("$HOME/.jiracrawler")
		os.MkdirAll(configDir, 0700)
		viper.WriteConfigAs(configDir + "/config.yaml")
		fmt.Printf("Set %s in config\n", key)
	},
}

func init() {
	configCmd.AddCommand(setCmd)

	// Ensure configCmd and getCmd are initialized for root.go
}
