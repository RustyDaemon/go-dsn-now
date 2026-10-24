package tui

import (
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/RustyDaemon/go-dsn-now/internal/data"
	"github.com/RustyDaemon/go-dsn-now/internal/model"
)

func (m Model) handleDataMessage(msg tea.Msg) (Model, tea.Cmd, bool) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.ui.width = msg.Width
		m.ui.height = msg.Height
		m.resizeComponents()
		return m, nil, true
	case spinner.TickMsg:
		if !m.ui.ready && m.ui.loading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd, true
		}
		return m, nil, true
	case DSNConfigLoadedMsg:
		m.appData.FullData = model.FullData{}
		data.MapConfigToFullData(msg.Config, &m.appData.FullData)
		m.ui.configLoaded = true
		m.ui.loading = true
		m.appData.LastError = ""
		m.restoreInitialStationSelection()
		m.resizeComponents()
		return m, loadDSNData(m.httpClient, m.cfg), true
	case DSNConfigErrorMsg:
		m.ui.loading = false
		m.appData.LastError = "Could not load NASA configuration: " + msg.Err.Error()
		return m, nil, true
	case DSNDataLoadedMsg:
		m.ui.loading = false
		m.applyLoadedData(msg)
		return m, tickRefresh(m.cfg.RefreshInterval), true
	case DSNDataErrorMsg:
		m.ui.loading = false
		m.appData.LastError = "Could not load NASA live data: " + msg.Err.Error()
		m.appData.ConsecutiveErrors++
		if m.ui.ready {
			return m, tickRefresh(m.cfg.RefreshInterval), true
		}
		return m, nil, true
	case TickRefreshMsg:
		if !m.ui.loading {
			m.ui.loading = true
			return m, loadDSNData(m.httpClient, m.cfg), true
		}
		return m, nil, true
	case CopyResultMsg:
		if msg.Err != nil {
			m.ui.statusMessage = "Clipboard unavailable"
		} else {
			m.ui.statusMessage = "Copied to clipboard"
		}
		return m, clearStatusMessage(2 * time.Second), true
	case StatusMessageExpiredMsg:
		m.ui.statusMessage = ""
		return m, nil, true
	default:
		return m, nil, false
	}
}

func (m *Model) applyLoadedData(msg DSNDataLoadedMsg) {
	data.MapDataToFullData(msg.Data, &m.appData.FullData)
	m.appData.LastError = ""
	m.appData.ConsecutiveErrors = 0
	m.appData.LastUpdated = time.Now()

	m.appendSignalHistory()
	if len(m.appData.FullData.Stations) == 0 {
		return
	}

	m.ui.ready = true
	m.restoreInitialStationSelection()
	m.clampSelection()
	m.appData.DetectSignalChanges()
	m.refreshCompactTable()
}

func (m *Model) appendSignalHistory() {
	for _, station := range m.appData.FullData.Stations {
		for _, dish := range station.Dishes {
			var downVal float64
			for _, sig := range dish.DownSignals {
				if sig.IsActive {
					if v, err := strconv.ParseFloat(sig.DataRate, 64); err == nil && v >= 0 {
						downVal = v
						break
					}
				}
			}
			downKey := dish.Name + ":down"
			downHistory := append(m.appData.SignalHistory[downKey], downVal)
			if len(downHistory) > 20 {
				downHistory = downHistory[len(downHistory)-20:]
			}
			m.appData.SignalHistory[downKey] = downHistory

			var upVal float64
			for _, sig := range dish.UpSignals {
				if sig.IsActive {
					if v, err := strconv.ParseFloat(sig.Power, 64); err == nil && v >= 0 {
						upVal = v
						break
					}
				}
			}
			upKey := dish.Name + ":up"
			upHistory := append(m.appData.SignalHistory[upKey], upVal)
			if len(upHistory) > 20 {
				upHistory = upHistory[len(upHistory)-20:]
			}
			m.appData.SignalHistory[upKey] = upHistory
		}
	}
}
