package provider

import "fmt"

// New returns a Provider for the given name and API key.
func New(name, model, apiKey string) (Provider, error) {
	switch name {
	case "claude", "":
		return NewClaude(apiKey, model), nil
	case "openai":
		return NewOpenAI(apiKey, model), nil
	case "ollama":
		return NewOllama(model), nil
	default:
		return nil, fmt.Errorf("unknown provider %q — supported: claude, openai, ollama", name)
	}
}
