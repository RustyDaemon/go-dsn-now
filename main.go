package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/RustyDaemon/go-dsn-now/internal/config"
	"github.com/RustyDaemon/go-dsn-now/internal/data"
	"github.com/RustyDaemon/go-dsn-now/internal/model"
	"github.com/RustyDaemon/go-dsn-now/internal/tui"
	"github.com/RustyDaemon/go-dsn-now/internal/tui/style"
)

func main() {
	cfg := config.Load()

	appSettings := data.LoadSettings()
	if appSettings != nil {
		if appSettings.RefreshIntervalSeconds > 0 {
			interval := time.Duration(appSettings.RefreshIntervalSeconds) * time.Second
			if interval < 10*time.Second {
				interval = 10 * time.Second
			}
			cfg.RefreshInterval = interval
		}
	} else {
		appSettings = &data.Settings{
			RefreshIntervalSeconds: int(cfg.RefreshInterval.Seconds()),
		}
	}

	if appSettings.Theme != "" {
		style.SetThemeByName(appSettings.Theme)
	}

	httpClient := &http.Client{
		Timeout: cfg.HTTPTimeout,
	}

	appData := model.NewAppData()
	appData.Bookmarks = data.LoadBookmarks()

	m := tui.NewModel(cfg, appSettings, httpClient, appData)

	p := tea.NewProgram(m,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
