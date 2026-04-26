package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/RustyDaemon/go-dsn-now/internal/tui/components"
	"github.com/RustyDaemon/go-dsn-now/internal/tui/style"
)

func (m Model) useCompactView() bool {
	return m.ui.compactView || m.ui.width < 80 || m.ui.height < 20
}

func (m Model) View() string {
	if !m.ui.ready {
		return m.viewLoading()
	}

	statusBar := m.statusBar.View(m.buildStatusParams())
	if m.ui.activeModal != components.ModalNone {
		return lipgloss.JoinVertical(lipgloss.Left, m.modal.View(), statusBar)
	}

	mainArea := m.viewDetailed()
	if m.useCompactView() {
		mainArea = m.compactTable.View()
	}

	return lipgloss.JoinVertical(lipgloss.Left, mainArea, statusBar)
}

func (m Model) viewLoading() string {
	title := style.TitleStyle.Render("GO DSN NOW")
	subtitle := style.LabelStyle.Render("Go Deep Space Network Monitor")
	version := style.MutedStyle.Render("v" + m.cfg.AppVersion)
	loading := style.DimStyle.Render(m.spinner.View() + " Loading configuration...")

	box := lipgloss.JoinVertical(lipgloss.Center,
		"",
		title,
		subtitle,
		version,
		"",
		loading,
		"",
	)

	splashStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.DoubleBorder()).
		BorderForeground(style.ColorBorderModal).
		Padding(1, 4)

	return lipgloss.Place(m.ui.width, m.ui.height,
		lipgloss.Center, lipgloss.Center,
		splashStyle.Render(box),
		lipgloss.WithWhitespaceBackground(style.ColorBgDeep),
	)
}

func (m Model) viewDetailed() string {
	leftWidth := m.leftPanelWidth()
	rightWidth := m.ui.width - leftWidth
	mainHeight := m.mainHeight()

	station, _ := m.selectedStation()
	dishHeight := mainHeight - (2 + len(m.appData.FullData.Stations))
	if dishHeight < 5 {
		dishHeight = 5
	}

	dishListView := m.dishList.View(leftWidth, dishHeight, m.ui.selection.dish, len(station.Dishes))
	stationBarView := m.stationBar.View(m.appData.FullData.Stations, m.ui.selection.station)
	leftPanel := lipgloss.NewStyle().Width(leftWidth).Height(mainHeight).Render(
		lipgloss.JoinVertical(lipgloss.Left, dishListView, stationBarView),
	)

	dish, _ := m.selectedDish()
	targets, _ := m.selectedTargets()
	upSignals, _ := m.selectedUpSignals()
	downSignals, _ := m.selectedDownSignals()

	upSparkline := style.Sparkline(m.appData.SignalHistory[dish.Name+":up"])
	downSparkline := style.Sparkline(m.appData.SignalHistory[dish.Name+":down"])

	antennaInfo := components.RenderAntennaInfo(dish, rightWidth)
	targetView := components.RenderTarget(targets, m.ui.selection.target, rightWidth, m.prefs.distanceUnit())

	var signalsRow string
	if m.ui.width >= 180 {
		upView := components.RenderUpSignal(upSignals, m.ui.selection.upSignal, rightWidth, upSparkline)
		downView := components.RenderDownSignal(downSignals, m.ui.selection.downSignal, rightWidth, downSparkline)
		signalsRow = lipgloss.JoinVertical(lipgloss.Left, upView, downView)
	} else {
		signalHalfWidth := rightWidth / 2
		upView := components.RenderUpSignal(upSignals, m.ui.selection.upSignal, signalHalfWidth, upSparkline)
		downView := components.RenderDownSignal(downSignals, m.ui.selection.downSignal, signalHalfWidth, downSparkline)
		signalsRow = lipgloss.JoinHorizontal(lipgloss.Top, upView, downView)
	}

	rightPanel := lipgloss.NewStyle().Width(rightWidth).Height(mainHeight).Render(
		lipgloss.JoinVertical(lipgloss.Left, antennaInfo, targetView, signalsRow),
	)

	return lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)
}

var _ tea.Model = Model{}
