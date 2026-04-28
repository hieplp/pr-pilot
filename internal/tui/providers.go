package tui

const customModelSentinel = "__custom__"

var providerDefaults = map[string]string{
	"claude": "claude-sonnet-4-6",
	"openai": "gpt-4o",
	"ollama": "llama3.2",
}

var providerModels = map[string][]string{
	"claude": {
		"claude-opus-4-7",
		"claude-sonnet-4-6",
		"claude-haiku-4-5-20251001",
	},
	"openai": {
		"gpt-4o",
		"gpt-4o-mini",
		"gpt-4-turbo",
		"gpt-3.5-turbo",
	},
	"ollama": {
		"llama3.2",
		"llama3.1",
		"mistral",
		"codellama",
		"phi3",
	},
}
