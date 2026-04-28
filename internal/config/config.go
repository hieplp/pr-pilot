package config

// Config holds all runtime configuration for pr-pilot.
type Config struct {
	Provider        string
	Model           string
	Base            string // default base branch for `pr` subcommand
	AnthropicAPIKey string
	OpenAIAPIKey    string
	OllamaBaseURL   string
	MaxDiffBytes    int
}

// APIKey returns the API key for the active provider.
func (c *Config) APIKey() string {
	switch c.Provider {
	case "openai":
		return c.OpenAIAPIKey
	default:
		return c.AnthropicAPIKey
	}
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
