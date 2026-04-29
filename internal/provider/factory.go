package provider

import "fmt"

// New returns a Provider for the given name, model, API key, and (Ollama-only) base URL.
func New(name, model, apiKey, ollamaBaseURL string) (Provider, error) {
	switch name {
	case "claude", "":
		if apiKey == "" {
			return nil, fmt.Errorf("claude provider requires ANTHROPIC_API_KEY or anthropic_api_key in config")
		}
		return NewClaude(apiKey, model), nil
	case "openai":
		if apiKey == "" {
			return nil, fmt.Errorf("openai provider requires OPENAI_API_KEY or openai_api_key in config")
		}
		return NewOpenAI(apiKey, model), nil
	case "ollama":
		return NewOllama(model, ollamaBaseURL), nil
	default:
		return nil, fmt.Errorf("unknown provider %q — supported: claude, openai, ollama", name)
	}
}
