package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// sentinel embedded in the hook so we can safely identify and remove it
const hookSentinel = "# managed by pr-pilot"

const hookScript = `#!/bin/sh
# managed by pr-pilot — do not remove this line
# Runs only for fresh commits; skips amend, merge, squash, fixup.
[ -z "$2" ] || exit 0
msg=$(pr-pilot commit --yes 2>/dev/null) || exit 0
[ -n "$msg" ] && printf '%s\n' "$msg" > "$1"
`

var hookCmd = &cobra.Command{
	Use:   "hook",
	Short: "Manage the prepare-commit-msg git hook",
}

var hookInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install pr-pilot as a prepare-commit-msg hook in the current repo",
	RunE:  runHookInstall,
}

var hookUninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove the pr-pilot prepare-commit-msg hook",
	RunE:  runHookUninstall,
}

var hookStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show whether the pr-pilot hook is installed in the current repo",
	RunE:  runHookStatus,
}

func init() {
	rootCmd.AddCommand(hookCmd)
	hookCmd.AddCommand(hookInstallCmd, hookUninstallCmd, hookStatusCmd)
}

func hookFilePath() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if info, err := os.Stat(filepath.Join(dir, ".git")); err == nil && info.IsDir() {
			return filepath.Join(dir, ".git", "hooks", "prepare-commit-msg"), nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("not inside a git repository")
		}
		dir = parent
	}
}

func runHookInstall(_ *cobra.Command, _ []string) error {
	path, err := hookFilePath()
	if err != nil {
		return err
	}

	if existing, err := os.ReadFile(path); err == nil {
		if !strings.Contains(string(existing), hookSentinel) {
			return fmt.Errorf("a prepare-commit-msg hook already exists at %s\n"+
				"Remove it manually before installing pr-pilot's hook.", path)
		}
		fmt.Println("Hook already installed — overwriting.")
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	if err := os.WriteFile(path, []byte(hookScript), 0755); err != nil {
		return err
	}

	fmt.Printf("Hook installed: %s\n", path)
	fmt.Println("pr-pilot will now pre-fill commit messages automatically.")
	return nil
}

func runHookStatus(_ *cobra.Command, _ []string) error {
	path, err := hookFilePath()
	if err != nil {
		return err
	}

	content, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("Not installed.")
			return nil
		}
		return err
	}

	if strings.Contains(string(content), hookSentinel) {
		fmt.Printf("Installed: %s\n", path)
	} else {
		fmt.Printf("A hook exists at %s but was not installed by pr-pilot.\n", path)
	}
	return nil
}

func runHookUninstall(_ *cobra.Command, _ []string) error {
	path, err := hookFilePath()
	if err != nil {
		return err
	}

	content, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("No hook found — nothing to remove.")
			return nil
		}
		return err
	}

	if !strings.Contains(string(content), hookSentinel) {
		return fmt.Errorf("the hook at %s was not installed by pr-pilot — remove it manually", path)
	}

	if err := os.Remove(path); err != nil {
		return err
	}

	fmt.Printf("Hook removed: %s\n", path)
	return nil
}
