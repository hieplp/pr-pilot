package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Generate a commit message from staged changes",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("commit message generation — coming soon")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
}
