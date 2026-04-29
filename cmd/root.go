package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pr-pilot",
	Short: "Generate commit messages, PR descriptions, and changelogs using LLMs",
	Long: `pr-pilot uses LLMs (Claude, OpenAI, Ollama) to generate commit messages,
pull-request descriptions, and changelog entries from your git history.`,
	Example: `  pr-pilot commit                    # generate a commit message for staged changes
  pr-pilot pr                        # generate a PR description for the current branch
  pr-pilot changelog                 # generate a changelog from latest tag to HEAD
  pr-pilot config                    # open interactive config TUI
  pr-pilot config show               # print current config
  pr-pilot config model              # switch active model interactively
  pr-pilot hook install              # install the prepare-commit-msg git hook`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("provider", "", "LLM provider: claude, openai, ollama")
	rootCmd.PersistentFlags().String("model", "", "Model to use (defaults to provider's recommended model)")
}
