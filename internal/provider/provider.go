package provider

import "context"

type Provider interface {
	Complete(ctx context.Context, prompt string) (string, error)
	Name() string
}
