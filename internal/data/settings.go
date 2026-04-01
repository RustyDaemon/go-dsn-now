package data

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Settings struct {
	RefreshIntervalSeconds int    `json:"refreshIntervalSeconds"`
	Theme                  string `json:"theme"`
	LastStation            string `json:"lastStation,omitempty"`
}

func SettingsPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "go-dsn-now", "settings.json")
}

func LoadSettings() *Settings {
	path := SettingsPath()
	if path == "" {
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	var settings Settings
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil
	}
	return &settings
}

func SaveSettings(settings *Settings) error {
	path := SettingsPath()
	if path == "" {
		return nil
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
