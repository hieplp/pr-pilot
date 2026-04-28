package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pr-pilot",
	Short: "Generate commit messages and PR descriptions using LLMs",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("provider", "claude", "LLM provider: claude, openai, ollama")
	rootCmd.PersistentFlags().String("model", "", "Model to use (defaults to provider's recommended model)")
}
