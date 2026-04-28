package config

import (
	"github.com/spf13/viper"
)

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
	viper.SetDefault("anthropic_api_key", "")
	viper.SetDefault("openai_api_key", "")

	viper.SetEnvPrefix("PR_PILOT")
	viper.AutomaticEnv()
	// Also bind bare env vars so existing ANTHROPIC_API_KEY / OPENAI_API_KEY still work.
	viper.BindEnv("anthropic_api_key", "ANTHROPIC_API_KEY")
	viper.BindEnv("openai_api_key", "OPENAI_API_KEY")

	// Global config (lowest file priority).
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(ConfigDir())
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
		Provider:        viper.GetString("provider"),
		Model:           viper.GetString("model"),
		Base:            viper.GetString("base"),
		AnthropicAPIKey: viper.GetString("anthropic_api_key"),
		OpenAIAPIKey:    viper.GetString("openai_api_key"),
	}, nil
}
