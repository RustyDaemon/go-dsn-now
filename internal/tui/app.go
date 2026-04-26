package tui

import (
	"net/http"

	"github.com/charmbracelet/bubbles/spinner"
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
}

func NewModel(cfg *config.Config, settings *data.Settings, client *http.Client, appData *model.AppData) Model {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(style.ColorPrimary)

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
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		loadDSNConfig(m.httpClient, m.cfg),
		tickClock(),
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
		return m, nil
	}
}
