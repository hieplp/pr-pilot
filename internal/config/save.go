package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Save writes the config to ~/.config/pr-pilot/config.toml.
func Save(c *Config) error {
	dir := ConfigDir()
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}

	v := viper.New()
	v.Set("provider", c.Provider)
	v.Set("model", c.Model)
	v.Set("base", c.Base)
	v.Set("anthropic_api_key", c.AnthropicAPIKey)
	v.Set("openai_api_key", c.OpenAIAPIKey)
	v.Set("ollama_base_url", c.OllamaBaseURL)
	v.Set("max_diff_bytes", c.MaxDiffBytes)
	v.SetConfigType("toml")

	return v.WriteConfigAs(filepath.Join(dir, "config.toml"))
}
