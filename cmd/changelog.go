package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/atotto/clipboard"
	"github.com/hieplp/pr-pilot/internal/config"
	"github.com/hieplp/pr-pilot/internal/dryrun"
	"github.com/hieplp/pr-pilot/internal/git"
	"github.com/hieplp/pr-pilot/internal/prompt"
	"github.com/hieplp/pr-pilot/internal/provider"
	"github.com/hieplp/pr-pilot/internal/tui"
	"github.com/spf13/cobra"
)

var changelogCmd = &cobra.Command{
	Use:   "changelog",
	Short: "Generate a changelog entry from a git revision range",
	Example: `  pr-pilot changelog                         # generate from latest tag to HEAD
  pr-pilot changelog --from v1.2.0 --to HEAD
  pr-pilot changelog --from main --yes`,
	RunE: runChangelog,
}

func init() {
	rootCmd.AddCommand(changelogCmd)
	changelogCmd.Flags().String("from", "", "Start revision, tag, or branch (default: latest tag)")
	changelogCmd.Flags().String("to", "HEAD", "End revision, tag, or branch")
	changelogCmd.Flags().BoolP("yes", "y", false, "Accept generated changelog without interactive review")
	changelogCmd.Flags().Bool("dry-run", false, "Estimate prompt tokens and cost without calling the provider")
}

func runChangelog(cmd *cobra.Command, _ []string) error {
	yes, _ := cmd.Flags().GetBool("yes")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	from, _ := cmd.Flags().GetString("from")
	to, _ := cmd.Flags().GetString("to")
	if from == "" {
		var err error
		from, err = git.LatestTag()
		if err != nil {
			return err
		}
	}

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

	diff, err := git.RangeDiff(from, to)
	if err != nil {
		return err
	}
	diff = git.Truncate(diff, cfg.MaxDiffBytes)

	log, err := git.RangeCommitLog(from, to)
	if err != nil {
		return err
	}

	p, err := provider.New(cfg.Provider, cfg.Model, cfg.APIKey(), cfg.OllamaBaseURL)
	if err != nil {
		return err
	}

	system, user := prompt.ChangelogPrompt(from, to, diff, log)
	if dryRun {
		fmt.Print(dryrun.Estimate(cfg.Provider, cfg.Model, system, user, 1024).String())
		return nil
	}

	for {
		msg, err := tui.Spin("Generating changelog…", func() (string, error) {
			ctx, cancel := context.WithTimeout(cmd.Context(), 60*time.Second)
			defer cancel()
			return p.Complete(ctx, system, user)
		})
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
		case tui.ActionCopy:
			if err := clipboard.WriteAll(result.Content); err != nil {
				fmt.Fprintf(os.Stderr, "clipboard: %v\n", err)
			} else {
				fmt.Println("Copied to clipboard.")
			}
			return nil
		case tui.ActionRegenerate:
			continue
		case tui.ActionQuit:
			return nil
		}
	}
}
