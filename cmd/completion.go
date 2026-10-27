package cmd

import (
	"github.com/sebrandon1/go-toolbox/cobrautil"
)

var completionCmd = cobrautil.NewCompletionCommand(rootCmd, "jiracrawler")

func init() {
	completionCmd.DisableFlagsInUseLine = true
	rootCmd.AddCommand(completionCmd)
}
