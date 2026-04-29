package config

import "github.com/spf13/viper"

// Load reads config in precedence order:
//
//	CLI flags (caller applies via Override)
//	> env vars (PR_PILOT_*)
//	> .pr-pilot.toml in cwd (project config)
//	> ~/.config/pr-pilot/config.toml (global config)
func Load() (*Config, error) {
	v := viper.New()
	v.SetDefault("provider", "claude")
	v.SetDefault("model", "")
	v.SetDefault("base", "main")
	v.SetDefault("anthropic_api_key", "")
	v.SetDefault("openai_api_key", "")
	v.SetDefault("ollama_base_url", "http://localhost:11434/v1")
	v.SetDefault("max_diff_bytes", 80_000)

	// Global config (lowest file priority).
	v.SetConfigName("config")
	v.SetConfigType("toml")
	v.AddConfigPath(ConfigDir())
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	// Project config overlays global config but remains below environment variables.
	pv := viper.New()
	pv.SetConfigFile(".pr-pilot.toml")
	if err := pv.ReadInConfig(); err == nil {
		if err := v.MergeConfigMap(pv.AllSettings()); err != nil {
			return nil, err
		}
	}

	v.SetEnvPrefix("PR_PILOT")
	v.AutomaticEnv()
	// Also bind bare env vars so existing ANTHROPIC_API_KEY / OPENAI_API_KEY still work.
	_ = v.BindEnv("anthropic_api_key", "PR_PILOT_ANTHROPIC_API_KEY", "ANTHROPIC_API_KEY")
	_ = v.BindEnv("openai_api_key", "PR_PILOT_OPENAI_API_KEY", "OPENAI_API_KEY")

	return &Config{
		Provider:        v.GetString("provider"),
		Model:           v.GetString("model"),
		Base:            v.GetString("base"),
		AnthropicAPIKey: v.GetString("anthropic_api_key"),
		OpenAIAPIKey:    v.GetString("openai_api_key"),
		OllamaBaseURL:   v.GetString("ollama_base_url"),
		MaxDiffBytes:    v.GetInt("max_diff_bytes"),
	}, nil
}
