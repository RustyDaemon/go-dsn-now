package tui

import (
	"net/http"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/RustyDaemon/go-dsn-now/internal/config"
	"github.com/RustyDaemon/go-dsn-now/internal/data"
	"github.com/RustyDaemon/go-dsn-now/internal/model"
	"github.com/RustyDaemon/go-dsn-now/internal/tui/components"
	"github.com/RustyDaemon/go-dsn-now/internal/tui/style"
)

type Model struct {
	appData    *model.AppData
	cfg        *config.Config
	prefs      *preferences
	httpClient *http.Client

	ui uiState

	dishList     components.DishList
	stationBar   components.StationBar
	statusBar    components.StatusBar
	compactTable components.CompactTable
	modal        components.Modal
	spinner      spinner.Model
	search       textinput.Model
}

func NewModel(cfg *config.Config, settings *data.Settings, client *http.Client, appData *model.AppData) Model {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(style.ColorPrimary)
	search := textinput.New()
	search.Prompt = "Search: "
	search.Placeholder = "station, dish, or spacecraft"
	search.CharLimit = 128
	search.Cursor.SetMode(cursor.CursorStatic)

	return Model{
		cfg:          cfg,
		prefs:        newPreferences(settings),
		httpClient:   client,
		appData:      appData,
		ui:           newUIState(),
		dishList:     components.NewDishList(),
		stationBar:   components.NewStationBar(),
		statusBar:    components.NewStatusBar(),
		compactTable: components.NewCompactTable(),
		modal:        components.NewModal(),
		spinner:      sp,
		search:       search,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		loadDSNConfig(m.httpClient, m.cfg),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if updated, cmd, handled := m.handleDataMessage(msg); handled {
		return updated, cmd
	}

	switch msg := msg.(type) {
	case tea.MouseMsg:
		return m.handleMouse(msg)
	case tea.KeyMsg:
		return m.handleKey(msg)
	default:
		if m.ui.searchEditing {
			before := m.search.Value()
			var cmd tea.Cmd
			m.search, cmd = m.search.Update(msg)
			if m.search.Value() != before {
				m.refreshCompactTable()
			}
			return m, cmd
		}
		return m, nil
	}
}
