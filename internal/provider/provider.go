package provider

import "context"

type Provider interface {
	Complete(ctx context.Context, system, user string) (string, error)
	Name() string
}
