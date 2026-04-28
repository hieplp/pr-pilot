package provider

import "context"

// StubProvider is a test double that returns a fixed response without calling any API.
type StubProvider struct {
	Response string
	Calls    int
	name     string
}

// NewStub returns a StubProvider that always returns response.
func NewStub(name, response string) *StubProvider {
	return &StubProvider{name: name, Response: response}
}

func (s *StubProvider) Name() string { return s.name }

func (s *StubProvider) Complete(_ context.Context, _, _ string) (string, error) {
	s.Calls++
	return s.Response, nil
}
