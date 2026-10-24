package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/RustyDaemon/go-dsn-now/internal/tui/components"
	"github.com/RustyDaemon/go-dsn-now/internal/tui/style"
)

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == "ctrl+c" {
		return m, tea.Quit
	}
	if m.ui.searchEditing {
		switch msg.String() {
		case "enter":
			m.ui.searchEditing = false
			m.search.Blur()
			return m, nil
		case "esc":
			m.search.SetValue(m.ui.searchBefore)
			m.ui.searchEditing = false
			m.search.Blur()
			m.refreshCompactTable()
			return m, nil
		default:
			before := m.search.Value()
			var cmd tea.Cmd
			m.search, cmd = m.search.Update(msg)
			if m.search.Value() != before {
				m.refreshCompactTable()
			}
			return m, cmd
		}
	}
	if m.ui.activeModal != components.ModalNone {
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "esc":
			m.ui.activeModal = components.ModalNone
			return m, nil
		default:
			cmd := m.modal.Update(msg)
			return m, cmd
		}
	}

	if msg.String() == "q" {
		return m, tea.Quit
	}

	if !m.ui.ready {
		if msg.String() == "r" && !m.ui.loading {
			m.ui.loading = true
			m.appData.LastError = ""
			load := loadDSNConfig(m.httpClient, m.cfg)
			if m.ui.configLoaded {
				load = loadDSNData(m.httpClient, m.cfg)
			}
			return m, tea.Batch(m.spinner.Tick, load)
		}
		return m, nil
	}

	switch msg.String() {
	case "/":
		m.ui.searchBefore = m.search.Value()
		m.ui.searchEditing = true
		m.ui.compactView = true
		return m, m.search.Focus()
	case "a":
		m.ui.activeOnly = !m.ui.activeOnly
		m.refreshCompactTable()
	case "f":
		m.ui.bookmarkedOnly = !m.ui.bookmarkedOnly
		m.refreshCompactTable()
	case "esc":
		m.clearFilters()
	case "enter":
		if m.useCompactView() && m.compactTable.RowCount() > 0 {
			if m.ui.width < 80 || m.ui.height < 20 {
				m.ui.statusMessage = "Details need a terminal of at least 80×20"
				return m, clearStatusMessage(2 * time.Second)
			}
			m.ui.compactView = false
			m.refreshDishList()
		}
	case "s":
		m.cycleStation()
	case "t":
		m.cycleTarget()
	case "u":
		m.cycleUpSignal()
	case "d":
		m.cycleDownSignal()
	case "up", "k":
		if m.useCompactView() {
			cmd := m.compactTable.Update(msg)
			m.selectCompactDish()
			m.refreshDishList()
			return m, cmd
		}
		m.navigateDish(-1)
	case "down", "j":
		if m.useCompactView() {
			cmd := m.compactTable.Update(msg)
			m.selectCompactDish()
			m.refreshDishList()
			return m, cmd
		}
		m.navigateDish(1)
	case "b":
		m.toggleBookmark()
	case "c":
		m.ui.compactView = !m.ui.compactView
		m.refreshCompactTable()
	case "S":
		if m.useCompactView() {
			m.ui.compactSort.cycle()
			m.refreshCompactTable()
		}
	case "T":
		themeName := style.CycleTheme()
		m.prefs.setTheme(themeName)
		m.compactTable.RefreshStyles()
		m.ui.statusMessage = "Theme: " + themeName
		return m, clearStatusMessage(2 * time.Second)
	case "U":
		m.prefs.cycleDistanceUnit()
	case "y":
		return m, copyToClipboard(m.getVisibleContent())
	case "+", "=":
		m.cfg.RefreshInterval += 5 * time.Second
		m.prefs.setRefreshInterval(m.cfg.RefreshInterval)
	case "-":
		m.cfg.RefreshInterval -= 5 * time.Second
		if m.cfg.RefreshInterval < 10*time.Second {
			m.cfg.RefreshInterval = 10 * time.Second
		}
		m.prefs.setRefreshInterval(m.cfg.RefreshInterval)
	case "?":
		m.openHelpModal()
	}

	return m, nil
}

func (m Model) handleMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	if !m.ui.ready {
		return m, nil
	}

	isWheel := msg.Button == tea.MouseButtonWheelUp || msg.Button == tea.MouseButtonWheelDown
	if m.ui.activeModal != components.ModalNone {
		if isWheel {
			cmd := m.modal.Update(msg)
			return m, cmd
		}
		return m, nil
	}

	if m.useCompactView() {
		if isWheel {
			key := tea.KeyMsg{Type: tea.KeyDown}
			if msg.Button == tea.MouseButtonWheelUp {
				key.Type = tea.KeyUp
			}
			cmd := m.compactTable.Update(key)
			m.selectCompactDish()
			m.refreshDishList()
			return m, cmd
		}
		return m, nil
	}

	leftWidth := m.leftPanelWidth()
	switch {
	case isWheel:
		cmd := m.dishList.Update(msg)
		m.selectVisibleDish(m.dishList.Index())
		return m, cmd
	case msg.Button == tea.MouseButtonLeft && msg.Action == tea.MouseActionPress:
		if _, ok := m.selectedStation(); !ok || msg.X >= leftWidth {
			return m, nil
		}

		dishIdx := m.dishList.VisibleOffset() + msg.Y - 2 // filter bar and panel border
		m.selectVisibleDish(dishIdx)
	}

	return m, nil
}

func (m *Model) openHelpModal() {
	m.ui.activeModal = components.ModalHelp
	m.modal.SetContent("Help", components.BuildHelpContent(m.cfg.AppVersion, m.cfg.AppGithubURL))
	m.modal.SetSize(m.ui.width, m.ui.height-3)
}

func (m *Model) toggleBookmark() {
	dish, ok := m.selectedDish()
	if !ok {
		return
	}

	m.prefs.toggleBookmark(m.appData.Bookmarks, dish.Name)
	m.refreshCompactTable()
}
