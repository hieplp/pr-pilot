package provider

import (
	"fmt"
	"os"
)

// New returns a Provider for the given name, reading API keys from the environment.
func New(name, model string) (Provider, error) {
	switch name {
	case "claude", "":
		return NewClaude(os.Getenv("ANTHROPIC_API_KEY"), model), nil
	case "openai":
		return NewOpenAI(os.Getenv("OPENAI_API_KEY"), model), nil
	case "ollama":
		return NewOllama(model), nil
	default:
		return nil, fmt.Errorf("unknown provider %q — supported: claude, openai, ollama", name)
	}
}
