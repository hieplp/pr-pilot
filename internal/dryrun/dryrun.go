package dryrun

import (
	"fmt"
	"math"
	"strings"
)

type Report struct {
	Provider      string
	Model         string
	InputTokens   int
	OutputTokens  int
	InputCostUSD  float64
	OutputCostUSD float64
	KnownPricing  bool
}

type pricing struct {
	inputPerMTok  float64
	outputPerMTok float64
}

var prices = map[string]pricing{
	"claude-3-5-haiku-latest":  {0.80, 4.00},
	"claude-3-5-sonnet-latest": {3.00, 15.00},
	"claude-sonnet-4-6":        {3.00, 15.00},
	"gpt-4o":                   {2.50, 10.00},
	"gpt-4o-mini":              {0.15, 0.60},
	"gpt-4.1":                  {2.00, 8.00},
	"gpt-4.1-mini":             {0.40, 1.60},
	"gpt-4.1-nano":             {0.10, 0.40},
}

// Estimate returns a rough dry-run report without calling the provider API.
func Estimate(provider, model, system, user string, maxOutputTokens int) Report {
	if model == "" {
		model = defaultModel(provider)
	}
	inputTokens := estimateTokens(system) + estimateTokens(user)
	r := Report{Provider: provider, Model: model, InputTokens: inputTokens, OutputTokens: maxOutputTokens}
	if p, ok := prices[strings.ToLower(model)]; ok {
		r.KnownPricing = true
		r.InputCostUSD = float64(inputTokens) / 1_000_000 * p.inputPerMTok
		r.OutputCostUSD = float64(maxOutputTokens) / 1_000_000 * p.outputPerMTok
	}
	return r
}

func (r Report) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "Dry run: no API call made\n")
	fmt.Fprintf(&b, "Provider: %s\n", r.Provider)
	fmt.Fprintf(&b, "Model: %s\n", r.Model)
	fmt.Fprintf(&b, "Estimated input tokens: %d\n", r.InputTokens)
	fmt.Fprintf(&b, "Assumed max output tokens: %d\n", r.OutputTokens)
	if r.KnownPricing {
		total := r.InputCostUSD + r.OutputCostUSD
		fmt.Fprintf(&b, "Estimated max cost: $%.6f (input $%.6f + output $%.6f)\n", total, r.InputCostUSD, r.OutputCostUSD)
	} else {
		fmt.Fprintf(&b, "Estimated max cost: unavailable (pricing unknown for this model)\n")
	}
	return b.String()
}

func estimateTokens(s string) int {
	if s == "" {
		return 0
	}
	return int(math.Ceil(float64(len([]rune(s))) / 4.0))
}

func defaultModel(provider string) string {
	switch provider {
	case "openai":
		return "gpt-4o"
	case "ollama":
		return "llama3"
	default:
		return "claude-sonnet-4-6"
	}
}
