package dryrun

import (
	"strings"
	"testing"
)

func TestEstimateKnownPricing(t *testing.T) {
	r := Estimate("openai", "gpt-4o", "system", "user", 1024)

	if r.InputTokens == 0 {
		t.Fatal("InputTokens should be estimated")
	}
	if !r.KnownPricing {
		t.Fatal("KnownPricing should be true for gpt-4o")
	}
	if r.InputCostUSD <= 0 || r.OutputCostUSD <= 0 {
		t.Fatalf("costs should be positive, got input=%f output=%f", r.InputCostUSD, r.OutputCostUSD)
	}
}

func TestEstimateDefaultModel(t *testing.T) {
	r := Estimate("claude", "", "system", "user", 1024)

	if r.Model != "claude-sonnet-4-6" {
		t.Fatalf("Model = %q, want default claude model", r.Model)
	}
}

func TestReportStringUnknownPricing(t *testing.T) {
	r := Estimate("ollama", "custom-local", "system", "user", 1024)
	out := r.String()

	if !strings.Contains(out, "no API call made") {
		t.Fatal("report should say no API call was made")
	}
	if !strings.Contains(out, "pricing unknown") {
		t.Fatal("report should mention unknown pricing")
	}
}
