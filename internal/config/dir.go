package config

import (
	"os"
	"path/filepath"
)

// ConfigDir returns the path to the pr-pilot config directory.
func ConfigDir() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "pr-pilot")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "pr-pilot")
}
