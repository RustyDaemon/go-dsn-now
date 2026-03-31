package data

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func BookmarksPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "go-dsn-now", "bookmarks.json")
}

func LoadBookmarks() map[string]bool {
	path := BookmarksPath()
	if path == "" {
		return make(map[string]bool)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return make(map[string]bool)
	}

	var bookmarks map[string]bool
	if err := json.Unmarshal(data, &bookmarks); err != nil {
		return make(map[string]bool)
	}
	return bookmarks
}

func SaveBookmarks(bookmarks map[string]bool) error {
	path := BookmarksPath()
	if path == "" {
		return nil
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(bookmarks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
