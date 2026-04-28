package provider

import "fmt"

// New returns a Provider for the given name, model, API key, and (Ollama-only) base URL.
func New(name, model, apiKey, ollamaBaseURL string) (Provider, error) {
	switch name {
	case "claude", "":
		return NewClaude(apiKey, model), nil
	case "openai":
		return NewOpenAI(apiKey, model), nil
	case "ollama":
		return NewOllama(model, ollamaBaseURL), nil
	default:
		return nil, fmt.Errorf("unknown provider %q — supported: claude, openai, ollama", name)
	}
}
