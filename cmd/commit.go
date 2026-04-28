package cmd

import (
	"fmt"

	"github.com/hieplp/pr-pilot/internal/config"
	"github.com/hieplp/pr-pilot/internal/git"
	"github.com/hieplp/pr-pilot/internal/prompt"
	"github.com/hieplp/pr-pilot/internal/provider"
	"github.com/hieplp/pr-pilot/internal/tui"
	"github.com/spf13/cobra"
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Generate a commit message from staged changes",
	RunE:  runCommit,
}

func init() {
	rootCmd.AddCommand(commitCmd)
	commitCmd.Flags().BoolP("yes", "y", false, "Accept generated message without interactive review")
	commitCmd.Flags().BoolP("commit", "c", false, "Run git commit with the accepted message")
}

func runCommit(cmd *cobra.Command, _ []string) error {
	yes, _ := cmd.Flags().GetBool("yes")
	doCommit, _ := cmd.Flags().GetBool("commit")

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	var providerFlag, modelFlag string
	if cmd.Flags().Changed("provider") {
		providerFlag, _ = cmd.Flags().GetString("provider")
	}
	if cmd.Flags().Changed("model") {
		modelFlag, _ = cmd.Flags().GetString("model")
	}
	cfg.Override(providerFlag, modelFlag)

	diff, err := git.StagedDiff()
	if err != nil {
		return err
	}

	p, err := provider.New(cfg.Provider, cfg.Model, cfg.APIKey())
	if err != nil {
		return err
	}

	system, user := prompt.CommitPrompt(diff)

	for {
		msg, err := p.Complete(cmd.Context(), system, user)
		if err != nil {
			return err
		}

		if yes {
			fmt.Println(msg)
			if doCommit {
				return git.Commit(msg)
			}
			return nil
		}

		result, err := tui.Review(msg)
		if err != nil {
			return err
		}

		switch result.Action {
		case tui.ActionAccept, tui.ActionEdit:
			fmt.Println(result.Content)
			if doCommit {
				return git.Commit(result.Content)
			}
			return nil
		case tui.ActionRegenerate:
			continue
		case tui.ActionQuit:
			return nil
		}
	}
}
