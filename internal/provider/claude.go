package provider

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

const defaultClaudeModel = "claude-sonnet-4-6"

type ClaudeProvider struct {
	client anthropic.Client
	model  string
}

func NewClaude(apiKey, model string) *ClaudeProvider {
	opts := []option.RequestOption{}
	if apiKey != "" {
		opts = append(opts, option.WithAPIKey(apiKey))
	}
	if model == "" {
		model = defaultClaudeModel
	}
	return &ClaudeProvider{
		client: anthropic.NewClient(opts...),
		model:  model,
	}
}

func (p *ClaudeProvider) Name() string { return "claude" }

func (p *ClaudeProvider) Complete(ctx context.Context, prompt string) (string, error) {
	msg, err := p.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.Model(p.model),
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
	})
	if err != nil {
		return "", fmt.Errorf("claude: %w", err)
	}
	if len(msg.Content) == 0 {
		return "", fmt.Errorf("claude: empty response")
	}
	return msg.Content[0].Text, nil
}
