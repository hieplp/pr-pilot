package provider

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

const (
	defaultOpenAIModel = "gpt-4o"
	defaultOllamaModel = "llama3"
	ollamaBaseURL      = "http://localhost:11434/v1"
)

type OpenAIProvider struct {
	client *openai.Client
	model  string
	name   string
}

func NewOpenAI(apiKey, model string) *OpenAIProvider {
	if model == "" {
		model = defaultOpenAIModel
	}
	return &OpenAIProvider{
		client: openai.NewClient(apiKey),
		model:  model,
		name:   "openai",
	}
}

// NewOllama creates an OpenAI-compatible provider pointed at a local Ollama instance.
func NewOllama(model string) *OpenAIProvider {
	if model == "" {
		model = defaultOllamaModel
	}
	cfg := openai.DefaultConfig("ollama") // key unused by Ollama but required by client
	cfg.BaseURL = ollamaBaseURL
	return &OpenAIProvider{
		client: openai.NewClientWithConfig(cfg),
		model:  model,
		name:   "ollama",
	}
}

func (p *OpenAIProvider) Name() string { return p.name }

func (p *OpenAIProvider) Complete(ctx context.Context, prompt string) (string, error) {
	resp, err := p.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: p.model,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
	})
	if err != nil {
		return "", fmt.Errorf("%s: %w", p.name, err)
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("%s: empty response", p.name)
	}
	return resp.Choices[0].Message.Content, nil
}
