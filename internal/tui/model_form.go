package tui

import "github.com/charmbracelet/huh"

// ModelForm opens a quick model-switcher for the given provider.
// If the current model is not in the predefined list it is pre-filled as a
// custom value. Returns (model, submitted, error).
func ModelForm(provider, currentModel string) (string, bool, error) {
	models := providerModels[provider]

	selected := customModelSentinel
	for _, m := range models {
		if m == currentModel {
			selected = m
			break
		}
	}

	customModel := ""
	if selected == customModelSentinel {
		customModel = currentModel
	}

	opts := make([]huh.Option[string], 0, len(models)+1)
	for _, m := range models {
		opts = append(opts, huh.NewOption(m, m))
	}
	opts = append(opts, huh.NewOption("Custom...", customModelSentinel))

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Model").
				Description("Provider: "+provider).
				Options(opts...).
				Value(&selected),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Custom model name").
				Value(&customModel),
		).WithHideFunc(func() bool { return selected != customModelSentinel }),
	)

	err := form.Run()
	if err == huh.ErrUserAborted {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}

	if selected == customModelSentinel {
		return customModel, true, nil
	}
	return selected, true, nil
}
