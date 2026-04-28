package config_test

import (
	"testing"

	"github.com/hieplp/pr-pilot/internal/config"
	"github.com/spf13/viper"
)

func TestLoadDefaults(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir()) // isolate from real config file
	t.Cleanup(viper.Reset)

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	checks := []struct {
		name string
		got  any
		want any
	}{
		{"Provider", cfg.Provider, "claude"},
		{"Base", cfg.Base, "main"},
		{"MaxDiffBytes", cfg.MaxDiffBytes, 80_000},
		{"OllamaBaseURL", cfg.OllamaBaseURL, "http://localhost:11434/v1"},
	}
	for _, c := range checks {
		if c.got != c.want {
			t.Errorf("%s = %v, want %v", c.name, c.got, c.want)
		}
	}
}

func TestLoadEnvOverrides(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Cleanup(viper.Reset)

	t.Setenv("PR_PILOT_PROVIDER", "openai")
	t.Setenv("PR_PILOT_MODEL", "gpt-4o")
	t.Setenv("PR_PILOT_BASE", "develop")
	t.Setenv("ANTHROPIC_API_KEY", "sk-ant-test")
	t.Setenv("OPENAI_API_KEY", "sk-openai-test")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if cfg.Provider != "openai" {
		t.Errorf("Provider = %q, want %q", cfg.Provider, "openai")
	}
	if cfg.Model != "gpt-4o" {
		t.Errorf("Model = %q, want %q", cfg.Model, "gpt-4o")
	}
	if cfg.Base != "develop" {
		t.Errorf("Base = %q, want %q", cfg.Base, "develop")
	}
	if cfg.AnthropicAPIKey != "sk-ant-test" {
		t.Errorf("AnthropicAPIKey = %q, want %q", cfg.AnthropicAPIKey, "sk-ant-test")
	}
	if cfg.OpenAIAPIKey != "sk-openai-test" {
		t.Errorf("OpenAIAPIKey = %q, want %q", cfg.OpenAIAPIKey, "sk-openai-test")
	}
}

func TestLoadMaxDiffBytesEnv(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Cleanup(viper.Reset)

	t.Setenv("PR_PILOT_MAX_DIFF_BYTES", "50000")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if cfg.MaxDiffBytes != 50_000 {
		t.Errorf("MaxDiffBytes = %d, want %d", cfg.MaxDiffBytes, 50_000)
	}
}
