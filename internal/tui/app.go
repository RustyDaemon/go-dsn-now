package tui

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

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
	// Dimensions
	width  int
	height int

	// Core data
	appData    *model.AppData
	cfg        *config.Config
	settings   *data.Settings
	httpClient *http.Client

	// UI state
	ready       bool
	activeModal components.ModalType

	// Sub-components
	dishList     components.DishList
	stationBar   components.StationBar
	statusBar    components.StatusBar
	compactTable components.CompactTable
	modal        components.Modal
	spinner      spinner.Model

	// Status message (transient)
	statusMessage string
}

func NewModel(cfg *config.Config, settings *data.Settings, client *http.Client, appData *model.AppData) Model {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(style.ColorPrimary)

	return Model{
		cfg:          cfg,
		settings:     settings,
		httpClient:   client,
		appData:      appData,
		dishList:     components.NewDishList(),
		stationBar:   components.NewStationBar(),
		statusBar:    components.NewStatusBar(),
		compactTable: components.NewCompactTable(),
		modal:        components.NewModal(),
		spinner:      sp,
		activeModal:  components.ModalNone,
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
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.resizeComponents()
		return m, nil

	case spinner.TickMsg:
		if !m.ready {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)

	case DSNConfigLoadedMsg:
		data.MapConfigToFullData(msg.Config, &m.appData.FullData)
		m.appData.DSNConfig = msg.Config
		m.resizeComponents()
		return m, tea.Batch(
			loadDSNData(m.httpClient, m.cfg),
			tickRefresh(m.cfg.RefreshInterval),
		)

	case DSNConfigErrorMsg:
		// Fatal error on config load
		return m, tea.Quit

	case tea.MouseMsg:
		return m.handleMouse(msg)

	case DSNDataLoadedMsg:
		data.MapDataToFullData(msg.Data, &m.appData.FullData)
		m.appData.LastError = ""
		m.appData.ConsecutiveErrors = 0
		m.appData.LastUpdated = time.Now()

		// Record signal history for sparklines
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
				dk := dish.Name + ":down"
				dh := append(m.appData.SignalHistory[dk], downVal)
				if len(dh) > 20 {
					dh = dh[len(dh)-20:]
				}
				m.appData.SignalHistory[dk] = dh

				var upVal float64
				for _, sig := range dish.UpSignals {
					if sig.IsActive {
						if v, err := strconv.ParseFloat(sig.Power, 64); err == nil && v >= 0 {
							upVal = v
							break
						}
					}
				}
				uk := dish.Name + ":up"
				uh := append(m.appData.SignalHistory[uk], upVal)
				if len(uh) > 20 {
					uh = uh[len(uh)-20:]
				}
				m.appData.SignalHistory[uk] = uh
			}
		}

		if m.appData.FullData.Stations == nil {
			return m, nil
		}

		if !m.ready {
			m.ready = true
		}

		if m.appData.SelectedStationIdx < 0 {
			m.appData.SelectedStationIdx = 0
			if m.settings != nil && m.settings.LastStation != "" {
				for i, station := range m.appData.FullData.Stations {
					if strings.EqualFold(station.FriendlyName, m.settings.LastStation) {
						m.appData.SelectedStationIdx = i
						break
					}
				}
			}
		}

		m.appData.DetectSignalChanges()
		m.refreshDishList()
		m.refreshCompactTable()
		return m, nil

	case DSNDataErrorMsg:
		m.appData.LastError = msg.Err.Error()
		m.appData.ConsecutiveErrors++
		return m, nil

	case TickRefreshMsg:
		return m, tea.Batch(
			loadDSNData(m.httpClient, m.cfg),
			tickRefresh(m.cfg.RefreshInterval),
		)

	case TickClockMsg:
		return m, nil

	case CopyResultMsg:
		if msg.Err != nil {
			m.statusMessage = "Clipboard unavailable"
		} else {
			m.statusMessage = "Copied to clipboard"
		}
		cmds = append(cmds, clearStatusMessage(2*time.Second))
		return m, tea.Batch(cmds...)

	case StatusMessageExpiredMsg:
		m.statusMessage = ""
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Modal active: only close/quit/scroll
	if m.activeModal != components.ModalNone {
		switch {
		case msg.String() == "q":
			return m, tea.Quit
		case msg.String() == "esc":
			m.activeModal = components.ModalNone
			return m, nil
		default:
			cmd := m.modal.Update(msg)
			return m, cmd
		}
	}

	// Quit always works
	if msg.String() == "q" {
		return m, tea.Quit
	}

	// Not ready: no other keys
	if !m.ready {
		return m, nil
	}

	switch msg.String() {
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
			return m, cmd
		}
		m.navigateDish(-1)
	case "down", "j":
		if m.useCompactView() {
			cmd := m.compactTable.Update(msg)
			return m, cmd
		}
		m.navigateDish(1)
	case "b":
		m.toggleBookmark()
	case "c":
		m.appData.CompactView = !m.appData.CompactView
		if m.appData.CompactView {
			m.refreshCompactTable()
		}
	case "S":
		if m.useCompactView() {
			m.appData.CycleCompactSortMode()
			m.refreshCompactTable()
		}
	case "T":
		themeName := style.CycleTheme()
		m.settings.Theme = themeName
		data.SaveSettings(m.settings)
		m.compactTable.RefreshStyles()
		m.statusMessage = "Theme: " + themeName
		return m, clearStatusMessage(2 * time.Second)
	case "U":
		units := []string{"km", "au", "lmin", "lhour"}
		current := m.settings.DistanceUnit
		next := "au"
		for i, u := range units {
			if u == current {
				next = units[(i+1)%len(units)]
				break
			}
		}
		m.settings.DistanceUnit = next
		data.SaveSettings(m.settings)
	case "y":
		text := m.getVisibleContent()
		return m, copyToClipboard(text)
	case "J":
		dish, ok := m.appData.GetSelectedDish()
		if ok {
			m.activeModal = components.ModalJSONPreview
			content := components.BuildJSONContent(dish)
			m.modal.SetContent("JSON Preview", content)
			m.modal.SetSize(m.width, m.height-3)
		}
	case "i":
		dish, ok := m.appData.GetSelectedDish()
		if ok && dish.Specs != (model.DishSpecification{}) {
			m.activeModal = components.ModalDishSpecs
			content := components.BuildDishSpecsContent(dish.Specs)
			m.modal.SetContent("Antenna Specification", content)
			m.modal.SetSize(m.width, m.height-3)
		}
	case "+", "=":
		newInterval := m.cfg.RefreshInterval + 5*time.Second
		m.cfg.RefreshInterval = newInterval
		m.settings.RefreshIntervalSeconds = int(newInterval.Seconds())
		data.SaveSettings(m.settings)
	case "-":
		newInterval := m.cfg.RefreshInterval - 5*time.Second
		if newInterval < 10*time.Second {
			newInterval = 10 * time.Second
		}
		m.cfg.RefreshInterval = newInterval
		m.settings.RefreshIntervalSeconds = int(newInterval.Seconds())
		data.SaveSettings(m.settings)
	case "?":
		m.activeModal = components.ModalHelp
		content := components.BuildHelpContent(m.cfg.AppVersion, m.cfg.AppGithubURL)
		m.modal.SetContent("Help", content)
		m.modal.SetSize(m.width, m.height-3)
	}

	return m, nil
}

// useCompactView returns true when compact mode should be rendered — either
// because the user toggled it, or because the terminal is too small for the
// detailed layout.
func (m Model) useCompactView() bool {
	return m.appData.CompactView || m.width < 80 || m.height < 20
}

func (m Model) View() string {
	if !m.ready {
		return m.viewLoading()
	}

	statusBar := m.statusBar.View(m.buildStatusParams())

	if m.activeModal != components.ModalNone {
		modalView := m.modal.View()
		return lipgloss.JoinVertical(lipgloss.Left, modalView, statusBar)
	}

	var mainArea string
	if m.useCompactView() {
		mainArea = m.compactTable.View()
	} else {
		mainArea = m.viewDetailed()
	}

	return lipgloss.JoinVertical(lipgloss.Left, mainArea, statusBar)
}

func (m Model) viewLoading() string {
	title := style.TitleStyle.Render("GO DSN NOW")
	subtitle := style.LabelStyle.Render("Go Deep Space Network Monitor")
	version := style.MutedStyle.Render("v" + m.cfg.AppVersion)
	loading := style.DimStyle.Render(m.spinner.View() + " Loading configuration...")
	url := style.MutedStyle.Render("")

	box := lipgloss.JoinVertical(lipgloss.Center,
		"",
		title,
		subtitle,
		version,
		"",
		loading,
		"",
		url,
		"",
	)

	splashStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.DoubleBorder()).
		BorderForeground(style.ColorBorderModal).
		Padding(1, 4)

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		splashStyle.Render(box),
		lipgloss.WithWhitespaceBackground(style.ColorBgDeep),
	)
}

func (m Model) viewDetailed() string {
	// Calculate panel widths
	leftWidth := m.width / 4
	if leftWidth < 20 {
		leftWidth = 20
	}
	rightWidth := m.width - leftWidth
	statusHeight := 3
	mainHeight := m.height - statusHeight

	// Left panel: dish list + station bar
	stationBarLines := 2 + len(m.appData.FullData.Stations)
	dishHeight := mainHeight - stationBarLines
	if dishHeight < 5 {
		dishHeight = 5
	}

	station := m.appData.FullData.Stations[m.appData.SelectedStationIdx]

	dishListView := m.dishList.View(leftWidth, dishHeight, m.appData.SelectedDishIdx, len(station.Dishes))
	stationBarView := m.stationBar.View(m.appData.FullData.Stations, m.appData.SelectedStationIdx)

	leftPanel := lipgloss.JoinVertical(lipgloss.Left, dishListView, stationBarView)
	leftPanel = lipgloss.NewStyle().Width(leftWidth).Height(mainHeight).Render(leftPanel)

	// Right panel: antenna info + target + signals
	dish, _ := m.appData.GetSelectedDish()
	targets, _ := m.appData.GetTargets()
	upSignals, _ := m.appData.GetUpSignals()
	downSignals, _ := m.appData.GetDownSignals()

	upSparkline := style.Sparkline(m.appData.SignalHistory[dish.Name+":up"])
	downSparkline := style.Sparkline(m.appData.SignalHistory[dish.Name+":down"])

	antennaInfo := components.RenderAntennaInfo(dish, rightWidth)
	targetView := components.RenderTarget(targets, m.appData.SelectedTargetIdx, rightWidth, m.settings.DistanceUnit)

	var signalsRow string
	if m.width >= 180 {
		// Wide layout: stack signals vertically at full right-panel width
		upView := components.RenderUpSignal(upSignals, m.appData.SelectedUpSignalIdx, rightWidth, upSparkline)
		downView := components.RenderDownSignal(downSignals, m.appData.SelectedDownSignalIdx, rightWidth, downSparkline)
		signalsRow = lipgloss.JoinVertical(lipgloss.Left, upView, downView)
	} else {
		signalHalfWidth := rightWidth / 2
		upView := components.RenderUpSignal(upSignals, m.appData.SelectedUpSignalIdx, signalHalfWidth, upSparkline)
		downView := components.RenderDownSignal(downSignals, m.appData.SelectedDownSignalIdx, signalHalfWidth, downSparkline)
		signalsRow = lipgloss.JoinHorizontal(lipgloss.Top, upView, downView)
	}

	rightPanel := lipgloss.JoinVertical(lipgloss.Left, antennaInfo, targetView, signalsRow)
	rightPanel = lipgloss.NewStyle().Width(rightWidth).Height(mainHeight).Render(rightPanel)

	return lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)
}

func (m Model) handleMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	if !m.ready {
		return m, nil
	}

	isWheel := msg.Button == tea.MouseButtonWheelUp || msg.Button == tea.MouseButtonWheelDown

	// Forward scroll to modal when open
	if m.activeModal != components.ModalNone {
		if isWheel {
			cmd := m.modal.Update(msg)
			return m, cmd
		}
		return m, nil
	}

	if m.appData.CompactView {
		return m, nil
	}

	leftWidth := m.width / 4
	if leftWidth < 20 {
		leftWidth = 20
	}

	switch {
	case isWheel:
		// Let the list handle scrolling natively, then sync index
		cmd := m.dishList.Update(msg)
		newIdx := m.dishList.Index()
		if newIdx != m.appData.SelectedDishIdx {
			m.appData.SelectedDishIdx = newIdx
			m.appData.SelectedTargetIdx = 0
			m.appData.SelectedUpSignalIdx = 0
			m.appData.SelectedDownSignalIdx = 0
		}
		return m, cmd

	case msg.Button == tea.MouseButtonLeft && msg.Action == tea.MouseActionPress:
		if msg.X < leftWidth && m.appData.SelectedStationIdx >= 0 {
			// Account for the list's scroll offset and the top panel border row
			offset := m.dishList.VisibleOffset()
			dishIdx := offset + msg.Y - 1
			if dishIdx >= 0 {
				station := m.appData.FullData.Stations[m.appData.SelectedStationIdx]
				if dishIdx < len(station.Dishes) {
					m.appData.SelectedDishIdx = dishIdx
					m.appData.SelectedTargetIdx = 0
					m.appData.SelectedUpSignalIdx = 0
					m.appData.SelectedDownSignalIdx = 0
					m.dishList.Select(dishIdx)
				}
			}
		}
	}

	return m, nil
}

func (m *Model) resizeComponents() {
	leftWidth := m.width / 4
	if leftWidth < 20 {
		leftWidth = 20
	}
	statusHeight := 3
	mainHeight := m.height - statusHeight

	// Station bar panel: top border + station lines + bottom border
	stationCount := len(m.appData.FullData.Stations)
	if stationCount == 0 {
		stationCount = 3 // estimate before data loads
	}
	stationBarHeight := stationCount + 2
	// Dish list content height: subtract station bar and dish panel's own 2 border lines
	dishHeight := mainHeight - stationBarHeight - 2
	if dishHeight < 3 {
		dishHeight = 3
	}

	m.dishList.SetSize(leftWidth, dishHeight)
	m.stationBar.SetWidth(leftWidth)
	m.statusBar.SetWidth(m.width)
	m.compactTable.SetSize(m.width, mainHeight)
	m.modal.SetSize(m.width, m.height-3)
}

func (m *Model) refreshDishList() {
	if m.appData.SelectedStationIdx < 0 || m.appData.SelectedStationIdx >= len(m.appData.FullData.Stations) {
		return
	}
	station := m.appData.FullData.Stations[m.appData.SelectedStationIdx]
	m.dishList.SetItems(station.Dishes, m.appData.Bookmarks)

	if m.appData.SelectedDishIdx >= len(station.Dishes) {
		m.appData.SelectedDishIdx = 0
	}
	m.dishList.Select(m.appData.SelectedDishIdx)
}

func (m *Model) refreshCompactTable() {
	var rows []components.CompactRow
	for _, station := range m.appData.FullData.Stations {
		for _, dish := range station.Dishes {
			target := "-"
			if len(dish.Targets) > 0 && dish.Targets[0].Spacecraft.FriendlyName != "" {
				target = dish.Targets[0].Spacecraft.FriendlyName
			}
			rows = append(rows, components.CompactRow{
				Station:     station.FriendlyName,
				Dish:        dish.FriendlyName,
				Target:      target,
				UpSignals:   dish.CountWorkingUpSignals(),
				DownSignals: dish.CountWorkingDownSignals(),
				Activity:    dish.Activity,
			})
		}
	}

	switch m.appData.CompactSortMode {
	case model.CompactSortByActivity:
		sort.SliceStable(rows, func(i, j int) bool {
			return rows[i].Activity < rows[j].Activity
		})
	case model.CompactSortBySignalCount:
		sort.SliceStable(rows, func(i, j int) bool {
			return (rows[i].UpSignals + rows[i].DownSignals) > (rows[j].UpSignals + rows[j].DownSignals)
		})
	case model.CompactSortByTarget:
		sort.SliceStable(rows, func(i, j int) bool {
			return rows[i].Target < rows[j].Target
		})
	}

	m.compactTable.SetSortLabel(m.appData.CompactSortModeLabel())
	m.compactTable.SetRows(rows)
}

func (m *Model) cycleStation() {
	stations := m.appData.FullData.Stations
	if len(stations) == 0 {
		return
	}
	m.appData.SelectedStationIdx = (m.appData.SelectedStationIdx + 1) % len(stations)
	m.appData.SelectedDishIdx = 0
	m.appData.SelectedTargetIdx = 0
	m.appData.SelectedUpSignalIdx = 0
	m.appData.SelectedDownSignalIdx = 0
	m.settings.LastStation = stations[m.appData.SelectedStationIdx].FriendlyName
	data.SaveSettings(m.settings)
	m.refreshDishList()
}

func (m *Model) cycleTarget() {
	targets, ok := m.appData.GetTargets()
	if !ok || len(targets) == 0 {
		return
	}
	m.appData.SelectedTargetIdx = (m.appData.SelectedTargetIdx + 1) % len(targets)
}

func (m *Model) cycleUpSignal() {
	signals, ok := m.appData.GetUpSignals()
	if !ok || len(signals) == 0 {
		return
	}
	m.appData.SelectedUpSignalIdx = (m.appData.SelectedUpSignalIdx + 1) % len(signals)
}

func (m *Model) cycleDownSignal() {
	signals, ok := m.appData.GetDownSignals()
	if !ok || len(signals) == 0 {
		return
	}
	m.appData.SelectedDownSignalIdx = (m.appData.SelectedDownSignalIdx + 1) % len(signals)
}

func (m *Model) navigateDish(direction int) {
	if m.appData.SelectedStationIdx < 0 {
		return
	}
	station := m.appData.FullData.Stations[m.appData.SelectedStationIdx]
	count := len(station.Dishes)
	if count == 0 {
		return
	}

	newIdx := m.appData.SelectedDishIdx + direction
	if newIdx < 0 {
		newIdx = count - 1
	} else if newIdx >= count {
		newIdx = 0
	}

	m.appData.SelectedDishIdx = newIdx
	m.appData.SelectedTargetIdx = 0
	m.appData.SelectedUpSignalIdx = 0
	m.appData.SelectedDownSignalIdx = 0
	m.dishList.Select(newIdx)
}

func (m *Model) toggleBookmark() {
	dish, ok := m.appData.GetSelectedDish()
	if !ok {
		return
	}
	if m.appData.Bookmarks[dish.Name] {
		delete(m.appData.Bookmarks, dish.Name)
	} else {
		m.appData.Bookmarks[dish.Name] = true
	}
	data.SaveBookmarks(m.appData.Bookmarks)
	m.refreshDishList()
}

func (m Model) buildStatusParams() components.StatusBarParams {
	if !m.ready {
		return components.StatusBarParams{}
	}

	connStatus := "connected"
	if m.appData.ConsecutiveErrors >= 3 {
		connStatus = "disconnected"
	} else if m.appData.ConsecutiveErrors >= 1 {
		connStatus = "degraded"
	}

	p := components.StatusBarParams{
		DefaultStatus:   m.activeModal == components.ModalNone,
		LastError:       m.appData.LastError,
		ConnStatus:      connStatus,
		SignalChanges:   m.appData.SignalChanges,
		RefreshInterval: fmt.Sprintf("%ds", int(m.cfg.RefreshInterval.Seconds())),
		StatusMessage:   m.statusMessage,
		DistanceUnit:    m.settings.DistanceUnit,
	}

	if !m.appData.LastUpdated.IsZero() {
		p.LastUpdated = m.appData.LastUpdated.Format("15:04:05")
	}

	if m.activeModal == components.ModalNone {
		if targets, ok := m.appData.GetTargets(); ok {
			p.HasTargets = len(targets) > 1
		}
		if upSignals, ok := m.appData.GetUpSignals(); ok {
			p.HasUpSignals = len(upSignals) > 1
		}
		if downSignals, ok := m.appData.GetDownSignals(); ok {
			p.HasDownSignals = len(downSignals) > 1
		}
		p.HasAntennaSpec = m.appData.HasAntennaSpecs()
	}

	return p
}

func (m Model) getVisibleContent() string {
	if m.appData.CompactView {
		return m.compactTable.GetVisibleContent()
	}

	dish, ok := m.appData.GetSelectedDish()
	if !ok {
		return ""
	}

	var b strings.Builder
	b.WriteString("=== Antenna Information ===\n")
	b.WriteString(fmt.Sprintf("Name: %s\n", style.DashIfEmpty(dish.FriendlyName)))
	b.WriteString(fmt.Sprintf("Type: %s\n", style.DashIfEmpty(dish.Type)))
	b.WriteString(fmt.Sprintf("Activity: %s\n", style.DashIfEmpty(dish.Activity)))
	b.WriteString(fmt.Sprintf("Azimuth: %s\n", style.FormatAngle(dish.AzimuthAngle)))
	b.WriteString(fmt.Sprintf("Elevation: %s\n", style.FormatAngle(dish.ElevationAngle)))
	b.WriteString(fmt.Sprintf("Wind: %s\n", style.FormatWind(dish.WindSpeed)))

	if targets, ok := m.appData.GetTargets(); ok && len(targets) > 0 {
		t := targets[m.appData.SelectedTargetIdx]
		b.WriteString("\n=== Target ===\n")
		name := "-"
		if t.Spacecraft != (model.Spacecraft{}) {
			name = t.Spacecraft.FriendlyName
		}
		b.WriteString(fmt.Sprintf("Spacecraft: %s\n", name))
		b.WriteString(fmt.Sprintf("Range: %s\n", style.FormatRangeInUnit(t.UplegRange, m.settings.DistanceUnit)))
		b.WriteString(fmt.Sprintf("Round-trip light time: %s\n", style.FormatRTLT(t.Rtlt)))
	}

	if upSignals, ok := m.appData.GetUpSignals(); ok && len(upSignals) > 0 {
		sig := upSignals[m.appData.SelectedUpSignalIdx]
		b.WriteString("\n=== Up Signal ===\n")
		source := "-"
		if sig.Spacecraft != (model.Spacecraft{}) {
			source = sig.Spacecraft.FriendlyName
		}
		b.WriteString(fmt.Sprintf("Source: %s\n", source))
		b.WriteString(fmt.Sprintf("Signal type: %s\n", style.DefaultIfEmpty(sig.SignalType, "-")))
		b.WriteString(fmt.Sprintf("Frequency band: %s\n", style.DefaultIfEmpty(sig.Band, "-")))
		b.WriteString(fmt.Sprintf("Power transmitted: %s\n", style.FormatPowerTx(sig.Power)))
	}

	if downSignals, ok := m.appData.GetDownSignals(); ok && len(downSignals) > 0 {
		sig := downSignals[m.appData.SelectedDownSignalIdx]
		b.WriteString("\n=== Down Signal ===\n")
		source := "-"
		if sig.Spacecraft != (model.Spacecraft{}) {
			source = sig.Spacecraft.FriendlyName
		}
		b.WriteString(fmt.Sprintf("Source: %s\n", source))
		b.WriteString(fmt.Sprintf("Signal type: %s\n", style.DefaultIfEmpty(sig.SignalType, "-")))
		b.WriteString(fmt.Sprintf("Frequency band: %s\n", style.DefaultIfEmpty(sig.Band, "-")))
		b.WriteString(fmt.Sprintf("Data rate: %s\n", style.FormatDataRate(sig.DataRate)))
		b.WriteString(fmt.Sprintf("Power received: %s\n", style.FormatPowerRx(sig.Power)))
	}

	return b.String()
}
