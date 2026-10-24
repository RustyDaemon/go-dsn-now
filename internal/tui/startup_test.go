package tui

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/RustyDaemon/go-dsn-now/internal/config"
	"github.com/RustyDaemon/go-dsn-now/internal/model"
)

type startupTransport struct {
	failPath string
	failed   bool
	paths    []string
}

func (s *startupTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	s.paths = append(s.paths, req.URL.Path)
	code, body := 200, `<dsn><dish name="dss14"/></dsn>`
	if req.URL.Path == "/config" {
		body = `<dsn><sites><site name="gdscc" friendlyName="Goldstone"><dish name="dss14" friendlyName="DSS 14"/></site></sites></dsn>`
	}
	if !s.failed && req.URL.Path == s.failPath {
		code, body, s.failed = 503, "<html/>", true
	}
	return &http.Response{StatusCode: code, Status: fmt.Sprintf("%d %s", code, http.StatusText(code)), Body: io.NopCloser(strings.NewReader(body))}, nil
}

func TestStartupErrorCanRetryFailedStage(t *testing.T) {
	for _, path := range []string{"/config", "/data"} {
		t.Run(path, func(t *testing.T) {
			transport := &startupTransport{failPath: path}
			cfg := config.Load()
			cfg.DSNConfigURL, cfg.DSNDataURL = "https://test.invalid/config", "https://test.invalid/data?r=%v"
			client := &http.Client{Transport: transport}
			m := NewModel(cfg, nil, client, model.NewAppData())
			m.ui.width, m.ui.height = 100, 30
			configMsg := loadDSNConfig(client, cfg)()
			m, cmd, _ := m.handleDataMessage(configMsg)
			if path == "/data" {
				m, cmd, _ = m.handleDataMessage(cmd())
			}
			if cmd != nil || m.ui.ready || m.ui.loading {
				t.Fatal("startup failure should wait for retry")
			}
			view := m.View()
			if !strings.Contains(view, "HTTP 503") || !strings.Contains(view, "r retry") {
				t.Fatal("startup error and retry hint are not visible")
			}
			m, cmd = press(m, "r")
			if cmd == nil || !m.ui.loading || m.appData.LastError != "" {
				t.Fatal("retry did not restart loading")
			}
			var duplicate tea.Cmd
			m, duplicate = press(m, "r")
			if duplicate != nil {
				t.Fatal("repeated retry started a concurrent request")
			}
			batch, ok := cmd().(tea.BatchMsg)
			if !ok {
				t.Fatal("retry should restart spinner and request")
			}
			// Execute only the request, leaving spinner timers to Bubble Tea.
			m, cmd, _ = m.handleDataMessage(batch[len(batch)-1]())
			if path == "/config" {
				m, cmd, _ = m.handleDataMessage(cmd())
			}
			if !m.ui.ready || m.ui.loading || m.appData.LastError != "" || cmd == nil {
				t.Fatal("retry did not recover to live monitoring")
			}
			if len(transport.paths) != 3 || len(m.appData.FullData.Stations) != 1 {
				t.Fatalf("unexpected requests or duplicate stations: %v", transport.paths)
			}
			m, cmd, _ = m.handleDataMessage(TickRefreshMsg{})
			if cmd == nil || !m.ui.loading {
				t.Fatal("refresh did not start a request")
			}
			_, duplicate, _ = m.handleDataMessage(TickRefreshMsg{})
			if duplicate != nil {
				t.Fatal("refresh started overlapping requests")
			}
		})
	}
}

func TestQuitWorksDuringLoadingErrorAndSearch(t *testing.T) {
	m := testModel()
	for _, ready := range []bool{true, false} {
		m.ui.ready = ready
		_, cmd := press(m, "q")
		if cmd == nil {
			t.Fatal("quit unavailable")
		}
		if _, ok := cmd().(tea.QuitMsg); !ok {
			t.Fatal("q did not quit")
		}
	}
	m.ui.ready = true
	m, _ = press(m, "/")
	_, cmd := press(m, "ctrl+c")
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Fatal("Ctrl+C did not quit search")
	}
}
