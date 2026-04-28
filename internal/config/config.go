package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Provider string
	Model    string
	Base     string // default base branch for `pr` subcommand
}

// Load reads config in precedence order:
//
//	CLI flags (caller applies via Override)
//	> env vars (PR_PILOT_*)
//	> .pr-pilot.toml in cwd (project config)
//	> ~/.config/pr-pilot/config.toml (global config)
func Load() (*Config, error) {
	viper.SetDefault("provider", "claude")
	viper.SetDefault("model", "")
	viper.SetDefault("base", "main")

	viper.SetEnvPrefix("PR_PILOT")
	viper.AutomaticEnv()

	// Global config (lowest file priority).
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(configDir())
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	// Project config overlaid on top of global.
	pv := viper.New()
	pv.SetConfigFile(".pr-pilot.toml")
	if err := pv.ReadInConfig(); err == nil {
		_ = viper.MergeConfigMap(pv.AllSettings())
	}

	return &Config{
		Provider: viper.GetString("provider"),
		Model:    viper.GetString("model"),
		Base:     viper.GetString("base"),
	}, nil
}

// Override merges non-empty flag values on top of the loaded config.
func (c *Config) Override(provider, model string) {
	if provider != "" {
		c.Provider = provider
	}
	if model != "" {
		c.Model = model
	}
}

func configDir() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "pr-pilot")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "pr-pilot")
}
