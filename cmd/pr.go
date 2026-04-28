package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var prCmd = &cobra.Command{
	Use:   "pr",
	Short: "Generate a PR description from the current branch diff",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("PR description generation — coming soon")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(prCmd)
	prCmd.Flags().String("base", "main", "Base branch to diff against")
}
