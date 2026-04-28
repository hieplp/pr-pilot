package config_test

import (
	"testing"

	"github.com/hieplp/pr-pilot/internal/config"
)

func TestOverride(t *testing.T) {
	cfg := &config.Config{Provider: "claude", Model: "claude-sonnet-4-6"}

	cfg.Override("openai", "gpt-4o")

	if cfg.Provider != "openai" {
		t.Errorf("Provider = %q, want %q", cfg.Provider, "openai")
	}
	if cfg.Model != "gpt-4o" {
		t.Errorf("Model = %q, want %q", cfg.Model, "gpt-4o")
	}
}

func TestOverrideEmptyPreservesExisting(t *testing.T) {
	cfg := &config.Config{Provider: "claude", Model: "claude-sonnet-4-6"}

	cfg.Override("", "")

	if cfg.Provider != "claude" {
		t.Errorf("Provider should not change, got %q", cfg.Provider)
	}
	if cfg.Model != "claude-sonnet-4-6" {
		t.Errorf("Model should not change, got %q", cfg.Model)
	}
}

func TestAPIKey(t *testing.T) {
	cfg := &config.Config{
		AnthropicAPIKey: "sk-ant-123",
		OpenAIAPIKey:    "sk-openai-456",
	}

	cfg.Provider = "claude"
	if got := cfg.APIKey(); got != "sk-ant-123" {
		t.Errorf("claude APIKey() = %q, want %q", got, "sk-ant-123")
	}

	cfg.Provider = "openai"
	if got := cfg.APIKey(); got != "sk-openai-456" {
		t.Errorf("openai APIKey() = %q, want %q", got, "sk-openai-456")
	}

	cfg.Provider = "ollama"
	if got := cfg.APIKey(); got != "" {
		t.Errorf("ollama APIKey() = %q, want empty string", got)
	}
}
