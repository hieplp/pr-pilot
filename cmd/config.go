package cmd

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/hieplp/pr-pilot/internal/config"
	"github.com/hieplp/pr-pilot/internal/tui"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure pr-pilot interactively",
	Long:  "Open a TUI form to set provider, model, base branch, and API keys.",
	Example: `  pr-pilot config             # open interactive TUI
  pr-pilot config show        # print current settings
  pr-pilot config model       # switch active model
  pr-pilot config set base develop
  pr-pilot config set max_diff_bytes 50000
  pr-pilot config init        # scaffold a .pr-pilot.toml in the current directory`,
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

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a single config key",
	Example: `  pr-pilot config set provider openai
  pr-pilot config set model gpt-4o
  pr-pilot config set base develop
  pr-pilot config set ollama_base_url http://myhost:11434/v1
  pr-pilot config set max_diff_bytes 50000`,
	Args: cobra.ExactArgs(2),
	RunE: runConfigSet,
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a .pr-pilot.toml project config in the current directory",
	RunE:  runConfigInit,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configModelCmd)
	configCmd.AddCommand(configSetCmd)
	configInitCmd.Flags().Bool("force", false, "Overwrite an existing .pr-pilot.toml")
	configCmd.AddCommand(configInitCmd)
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
		OllamaBaseURL:   cfg.OllamaBaseURL,
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
	cfg.OllamaBaseURL = result.OllamaBaseURL

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

	fmt.Printf("Provider        : %s\n", cfg.Provider)
	fmt.Printf("Model           : %s\n", cfg.Model)
	fmt.Printf("Base            : %s\n", cfg.Base)
	fmt.Printf("API Key         : %s\n", config.MaskKey(cfg.APIKey()))
	fmt.Printf("Ollama base URL : %s\n", cfg.OllamaBaseURL)
	fmt.Printf("Max diff bytes  : %d\n", cfg.MaxDiffBytes)
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

// allowedKeys maps config key names to the setter that applies them to *config.Config.
var allowedKeys = map[string]func(*config.Config, string) error{
	"provider": func(c *config.Config, v string) error {
		switch v {
		case "claude", "openai", "ollama":
			c.Provider = v
			return nil
		default:
			return fmt.Errorf("provider must be one of: claude, openai, ollama")
		}
	},
	"model": func(c *config.Config, v string) error { c.Model = v; return nil },
	"base":  func(c *config.Config, v string) error { c.Base = v; return nil },
	"anthropic_api_key": func(c *config.Config, v string) error {
		c.AnthropicAPIKey = v
		return nil
	},
	"openai_api_key": func(c *config.Config, v string) error {
		c.OpenAIAPIKey = v
		return nil
	},
	"ollama_base_url": func(c *config.Config, v string) error {
		c.OllamaBaseURL = v
		return nil
	},
	"max_diff_bytes": func(c *config.Config, v string) error {
		n, err := strconv.Atoi(v)
		if err != nil || n < 0 {
			return errors.New("max_diff_bytes must be a non-negative integer")
		}
		c.MaxDiffBytes = n
		return nil
	},
}

func runConfigSet(_ *cobra.Command, args []string) error {
	key, value := args[0], args[1]

	setter, ok := allowedKeys[key]
	if !ok {
		keys := make([]string, 0, len(allowedKeys))
		for k := range allowedKeys {
			keys = append(keys, k)
		}
		return fmt.Errorf("unknown key %q — allowed: provider, model, base, anthropic_api_key, openai_api_key, ollama_base_url, max_diff_bytes", key)
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	if err := setter(cfg, value); err != nil {
		return err
	}
	if err := config.Save(cfg); err != nil {
		return err
	}

	fmt.Printf("%s = %s\n", key, value)
	return nil
}

const projectConfigTemplate = `# pr-pilot project configuration — .pr-pilot.toml
# Settings here override ~/.config/pr-pilot/config.toml for this repository.

# LLM provider: claude, openai, ollama
# provider = "claude"

# Model name (leave empty for provider default)
# model = ""

# Default base branch for PR diffs
# base = "main"

# Anthropic API key (prefer the ANTHROPIC_API_KEY env var instead)
# anthropic_api_key = ""

# OpenAI API key (prefer the OPENAI_API_KEY env var instead)
# openai_api_key = ""

# Ollama base URL (only used when provider = "ollama")
# ollama_base_url = "http://localhost:11434/v1"

# Maximum diff size sent to the LLM in bytes (0 = unlimited)
# max_diff_bytes = 80000
`

func runConfigInit(cmd *cobra.Command, _ []string) error {
	force, _ := cmd.Flags().GetBool("force")
	const path = ".pr-pilot.toml"

	if _, err := os.Stat(path); err == nil && !force {
		return fmt.Errorf("%s already exists — use --force to overwrite", path)
	}

	if err := os.WriteFile(path, []byte(projectConfigTemplate), 0o644); err != nil {
		return err
	}
	fmt.Printf("Created %s\n", path)
	return nil
}
