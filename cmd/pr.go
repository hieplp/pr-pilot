package cmd

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

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
	Example: `  pr-pilot pr                        # generate and review interactively
  pr-pilot pr --create               # generate, review, then open the PR on GitHub
  pr-pilot pr --push --create        # push branch first, then open the PR
  pr-pilot pr -y --create            # accept without review and open PR immediately`,
	RunE: runPR,
}

func init() {
	rootCmd.AddCommand(prCmd)
	prCmd.Flags().String("base", "", "Base branch to diff against (overrides config, default: main)")
	prCmd.Flags().BoolP("yes", "y", false, "Accept generated description without interactive review")
	prCmd.Flags().Bool("push", false, "Push the current branch to origin before creating the PR")
	prCmd.Flags().Bool("create", false, "Create the PR on GitHub using `gh` after accepting the description")
}

func runPR(cmd *cobra.Command, _ []string) error {
	yes, _ := cmd.Flags().GetBool("yes")
	doPush, _ := cmd.Flags().GetBool("push")
	doCreate, _ := cmd.Flags().GetBool("create")

	if doCreate {
		if _, err := exec.LookPath("gh"); err != nil {
			return errors.New("--create requires the GitHub CLI (gh) — install it from https://cli.github.com")
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

	p, err := provider.New(cfg.Provider, cfg.Model, cfg.APIKey())
	if err != nil {
		return err
	}

	system, user := prompt.PRPrompt(branch, base, diff, log)

	for {
		msg, err := p.Complete(cmd.Context(), system, user)
		if err != nil {
			return err
		}

		var body string

		if yes {
			fmt.Println(msg)
			body = msg
		} else {
			result, err := tui.Review(msg)
			if err != nil {
				return err
			}

			switch result.Action {
			case tui.ActionAccept, tui.ActionEdit:
				fmt.Println(result.Content)
				body = result.Content
			case tui.ActionRegenerate:
				continue
			case tui.ActionQuit:
				return nil
			}
		}

		if body == "" {
			return nil
		}

		if doPush {
			fmt.Printf("Pushing branch %q to origin…\n", strings.TrimSpace(branch))
			if err := git.Push(strings.TrimSpace(branch)); err != nil {
				return fmt.Errorf("push failed: %w", err)
			}
		}

		if doCreate {
			return createGitHubPR(body, base)
		}
		return nil
	}
}

func createGitHubPR(body, base string) error {
	title := prompt.PRTitle(body)
	out, err := exec.Command("gh", "pr", "create",
		"--title", title,
		"--body", body,
		"--base", base,
	).CombinedOutput()
	if err != nil {
		return fmt.Errorf("gh pr create failed: %s", strings.TrimSpace(string(out)))
	}
	fmt.Print(string(out))
	return nil
}
