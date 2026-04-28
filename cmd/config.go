package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hieplp/pr-pilot/internal/config"
	"github.com/hieplp/pr-pilot/internal/tui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure pr-pilot interactively",
	Long:  "Open a TUI form to set provider, model, and base branch. Saves to ~/.config/pr-pilot/config.toml.",
	RunE:  runConfig,
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func runConfig(_ *cobra.Command, _ []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	current := tui.ConfigFormResult{
		Provider:        cfg.Provider,
		Model:           cfg.Model,
		Base:            cfg.Base,
		AnthropicAPIKey: cfg.AnthropicAPIKey,
		OpenAIAPIKey:    cfg.OpenAIAPIKey,
	}

	result, submitted, err := tui.ConfigForm(current)
	if err != nil {
		return err
	}
	if !submitted {
		fmt.Println("Cancelled — no changes saved.")
		return nil
	}

	if err := saveConfig(result); err != nil {
		return err
	}

	fmt.Printf("Saved: provider=%s  model=%s  base=%s\n", result.Provider, result.Model, result.Base)
	return nil
}

func saveConfig(r tui.ConfigFormResult) error {
	dir := configDir()
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}

	v := viper.New()
	v.Set("provider", r.Provider)
	v.Set("model", r.Model)
	v.Set("base", r.Base)
	v.Set("anthropic_api_key", r.AnthropicAPIKey)
	v.Set("openai_api_key", r.OpenAIAPIKey)
	v.SetConfigType("toml")

	path := filepath.Join(dir, "config.toml")
	return v.WriteConfigAs(path)
}

func configDir() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "pr-pilot")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "pr-pilot")
}
