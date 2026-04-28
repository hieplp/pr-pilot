package cmd

import (
	"fmt"

	"github.com/hieplp/pr-pilot/internal/config"
	"github.com/hieplp/pr-pilot/internal/tui"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure pr-pilot interactively",
	Long:  "Open a TUI form to set provider, model, base branch, and API keys.",
	Example: `  pr-pilot config          # open interactive TUI
  pr-pilot config show     # print current settings
  pr-pilot config model    # switch active model`,
	RunE: runConfig,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	RunE:  runConfigShow,
}

var configModelCmd = &cobra.Command{
	Use:   "model",
	Short: "Quickly switch the active model",
	RunE:  runConfigModel,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configModelCmd)
}

func runConfig(_ *cobra.Command, _ []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	result, submitted, err := tui.ConfigForm(tui.ConfigFormResult{
		Provider:        cfg.Provider,
		Model:           cfg.Model,
		Base:            cfg.Base,
		AnthropicAPIKey: cfg.AnthropicAPIKey,
		OpenAIAPIKey:    cfg.OpenAIAPIKey,
	})
	if err != nil {
		return err
	}
	if !submitted {
		fmt.Println("Cancelled — no changes saved.")
		return nil
	}

	cfg.Provider = result.Provider
	cfg.Model = result.Model
	cfg.Base = result.Base
	cfg.AnthropicAPIKey = result.AnthropicAPIKey
	cfg.OpenAIAPIKey = result.OpenAIAPIKey

	if err := config.Save(cfg); err != nil {
		return err
	}
	fmt.Printf("Saved: provider=%s  model=%s  base=%s\n", cfg.Provider, cfg.Model, cfg.Base)
	return nil
}

func runConfigShow(_ *cobra.Command, _ []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	fmt.Printf("Provider : %s\n", cfg.Provider)
	fmt.Printf("Model    : %s\n", cfg.Model)
	fmt.Printf("Base     : %s\n", cfg.Base)
	fmt.Printf("API Key  : %s\n", config.MaskKey(cfg.APIKey()))
	return nil
}

func runConfigModel(_ *cobra.Command, _ []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	model, submitted, err := tui.ModelForm(cfg.Provider, cfg.Model)
	if err != nil {
		return err
	}
	if !submitted {
		fmt.Println("Cancelled — no changes saved.")
		return nil
	}

	cfg.Model = model
	if err := config.Save(cfg); err != nil {
		return err
	}
	fmt.Printf("Model set to: %s\n", model)
	return nil
}
