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

var prCmd = &cobra.Command{
	Use:   "pr",
	Short: "Generate a PR description from the current branch diff",
	RunE:  runPR,
}

func init() {
	rootCmd.AddCommand(prCmd)
	prCmd.Flags().String("base", "", "Base branch to diff against (overrides config, default: main)")
	prCmd.Flags().BoolP("yes", "y", false, "Accept generated description without interactive review")
}

func runPR(cmd *cobra.Command, _ []string) error {
	yes, _ := cmd.Flags().GetBool("yes")

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	providerFlag, _ := cmd.Flags().GetString("provider")
	modelFlag, _ := cmd.Flags().GetString("model")
	cfg.Override(providerFlag, modelFlag)

	// --base flag wins over config; config wins over hardcoded "main".
	base, _ := cmd.Flags().GetString("base")
	if base == "" {
		base = cfg.Base
	}

	diff, err := git.BranchDiff(base)
	if err != nil {
		return err
	}
	log, err := git.CommitLog(base)
	if err != nil {
		return err
	}
	branch, err := git.CurrentBranch()
	if err != nil {
		return err
	}

	p, err := provider.New(cfg.Provider, cfg.Model)
	if err != nil {
		return err
	}

	promptStr := prompt.PRPrompt(branch, base, diff, log)

	for {
		msg, err := p.Complete(cmd.Context(), promptStr)
		if err != nil {
			return err
		}

		if yes {
			fmt.Println(msg)
			return nil
		}

		result, err := tui.Review(msg)
		if err != nil {
			return err
		}

		switch result.Action {
		case tui.ActionAccept, tui.ActionEdit:
			fmt.Println(result.Content)
			return nil
		case tui.ActionRegenerate:
			continue
		case tui.ActionQuit:
			return nil
		}
	}
}
