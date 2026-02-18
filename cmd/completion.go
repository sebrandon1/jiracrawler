package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion scripts",
	Long: `Generate shell completion scripts for jiracrawler.

To load completions:

Bash:
  $ source <(jiracrawler completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ jiracrawler completion bash > /etc/bash_completion.d/jiracrawler
  # macOS:
  $ jiracrawler completion bash > $(brew --prefix)/etc/bash_completion.d/jiracrawler

Zsh:
  $ jiracrawler completion zsh > "${fpath[1]}/_jiracrawler"

Fish:
  $ jiracrawler completion fish | source

  # To load completions for each session, execute once:
  $ jiracrawler completion fish > ~/.config/fish/completions/jiracrawler.fish

PowerShell:
  PS> jiracrawler completion powershell | Out-String | Invoke-Expression
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			_ = cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			_ = cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			_ = cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			_ = cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
