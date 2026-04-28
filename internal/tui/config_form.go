package tui

import "github.com/charmbracelet/huh"

// ConfigFormResult holds values collected from the config TUI form.
type ConfigFormResult struct {
	Provider        string
	Model           string
	Base            string
	AnthropicAPIKey string
	OpenAIAPIKey    string
}

// ConfigForm opens an interactive form pre-filled with current values.
// Returns (result, submitted, error). submitted is false when the user cancels.
func ConfigForm(current ConfigFormResult) (ConfigFormResult, bool, error) {
	result := current

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Provider").
				Description("LLM provider to use").
				Options(
					huh.NewOption("Claude (Anthropic)", "claude"),
					huh.NewOption("OpenAI", "openai"),
					huh.NewOption("Ollama (local)", "ollama"),
				).
				Value(&result.Provider),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Anthropic API key").
				Description("Stored in ~/.config/pr-pilot/config.toml").
				EchoMode(huh.EchoModePassword).
				Value(&result.AnthropicAPIKey),
		).WithHideFunc(func() bool { return result.Provider != "claude" }),
		huh.NewGroup(
			huh.NewInput().
				Title("OpenAI API key").
				Description("Stored in ~/.config/pr-pilot/config.toml").
				EchoMode(huh.EchoModePassword).
				Value(&result.OpenAIAPIKey),
		).WithHideFunc(func() bool { return result.Provider != "openai" }),
		huh.NewGroup(
			huh.NewInput().
				Title("Model").
				Description("Leave blank to use the provider default").
				Value(&result.Model),
			huh.NewInput().
				Title("Base branch").
				Description("Default base for PR diff (e.g. main, master)").
				Value(&result.Base),
		),
	)

	err := form.Run()
	if err == huh.ErrUserAborted {
		return ConfigFormResult{}, false, nil
	}
	if err != nil {
		return ConfigFormResult{}, false, err
	}

	if result.Model == "" {
		if def, ok := providerDefaults[result.Provider]; ok {
			result.Model = def
		}
	}

	return result, true, nil
}
