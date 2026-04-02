package tui

import (
	"context"
	"net/http"
	"time"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/RustyDaemon/go-dsn-now/internal/config"
	"github.com/RustyDaemon/go-dsn-now/internal/data"
	"github.com/RustyDaemon/go-dsn-now/internal/model/response"
)

func loadDSNConfig(client *http.Client, cfg *config.Config) tea.Cmd {
	return func() tea.Msg {
		ch := make(chan response.DSNConfig, 1)
		ce := make(chan error, 1)
		ctx, cancel := context.WithTimeout(context.Background(), cfg.HTTPTimeout)
		defer cancel()

		go data.LoadDSNConfig(ctx, client, cfg, ch, ce)

		select {
		case c := <-ch:
			return DSNConfigLoadedMsg{Config: c}
		case err := <-ce:
			return DSNConfigErrorMsg{Err: err}
		case <-ctx.Done():
			return DSNConfigErrorMsg{Err: ctx.Err()}
		}
	}
}

func loadDSNData(client *http.Client, cfg *config.Config) tea.Cmd {
	return func() tea.Msg {
		ch := make(chan response.DSN, 1)
		ce := make(chan error, 1)
		ctx, cancel := context.WithTimeout(context.Background(), cfg.HTTPTimeout)
		defer cancel()

		go data.LoadDSNData(ctx, client, cfg, ch, ce)

		select {
		case d := <-ch:
			return DSNDataLoadedMsg{Data: d}
		case err := <-ce:
			return DSNDataErrorMsg{Err: err}
		case <-ctx.Done():
			return DSNDataErrorMsg{Err: ctx.Err()}
		}
	}
}

func tickRefresh(interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(t time.Time) tea.Msg {
		return TickRefreshMsg(t)
	})
}

func tickClock() tea.Cmd {
	return tea.Every(time.Second, func(t time.Time) tea.Msg {
		return TickClockMsg(t)
	})
}

func copyToClipboard(text string) tea.Cmd {
	return func() tea.Msg {
		err := clipboard.WriteAll(text)
		return CopyResultMsg{Err: err}
	}
}

func clearStatusMessage(after time.Duration) tea.Cmd {
	return tea.Tick(after, func(t time.Time) tea.Msg {
		return StatusMessageExpiredMsg{}
	})
}
