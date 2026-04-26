package tui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/RustyDaemon/go-dsn-now/internal/model"
	"github.com/RustyDaemon/go-dsn-now/internal/tui/components"
	"github.com/RustyDaemon/go-dsn-now/internal/tui/style"
)

func (m *Model) refreshCompactTable() {
	rows := m.buildCompactRows()

	switch m.ui.compactSort {
	case compactSortByActivity:
		sort.SliceStable(rows, func(i, j int) bool {
			return rows[i].Activity < rows[j].Activity
		})
	case compactSortBySignalCount:
		sort.SliceStable(rows, func(i, j int) bool {
			return (rows[i].UpSignals + rows[i].DownSignals) > (rows[j].UpSignals + rows[j].DownSignals)
		})
	case compactSortByTarget:
		sort.SliceStable(rows, func(i, j int) bool {
			return rows[i].Target < rows[j].Target
		})
	}

	m.compactTable.SetSortLabel(m.ui.compactSort.label())
	m.compactTable.SetRows(rows)
}

func (m Model) buildCompactRows() []components.CompactRow {
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

	return rows
}

func (m Model) buildStatusParams() components.StatusBarParams {
	if !m.ui.ready {
		return components.StatusBarParams{}
	}

	connStatus := "connected"
	if m.appData.ConsecutiveErrors >= 3 {
		connStatus = "disconnected"
	} else if m.appData.ConsecutiveErrors >= 1 {
		connStatus = "degraded"
	}

	params := components.StatusBarParams{
		DefaultStatus:   m.ui.activeModal == components.ModalNone,
		LastError:       m.appData.LastError,
		ConnStatus:      connStatus,
		SignalChanges:   m.appData.SignalChanges,
		RefreshInterval: fmt.Sprintf("%ds", int(m.cfg.RefreshInterval.Seconds())),
		StatusMessage:   m.ui.statusMessage,
		DistanceUnit:    m.prefs.distanceUnit(),
	}

	if !m.appData.LastUpdated.IsZero() {
		params.LastUpdated = m.appData.LastUpdated.Format("15:04:05")
	}

	if m.ui.activeModal == components.ModalNone {
		if targets, ok := m.selectedTargets(); ok {
			params.HasTargets = len(targets) > 1
		}
		if upSignals, ok := m.selectedUpSignals(); ok {
			params.HasUpSignals = len(upSignals) > 1
		}
		if downSignals, ok := m.selectedDownSignals(); ok {
			params.HasDownSignals = len(downSignals) > 1
		}
	}

	return params
}

func (m Model) getVisibleContent() string {
	if m.useCompactView() {
		return m.compactTable.GetVisibleContent()
	}

	dish, ok := m.selectedDish()
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

	if targets, ok := m.selectedTargets(); ok && len(targets) > 0 {
		target := targets[m.ui.selection.target]
		b.WriteString("\n=== Target ===\n")
		name := "-"
		if target.Spacecraft != (model.Spacecraft{}) {
			name = target.Spacecraft.FriendlyName
		}
		b.WriteString(fmt.Sprintf("Spacecraft: %s\n", name))
		b.WriteString(fmt.Sprintf("Range: %s\n", style.FormatRangeInUnit(target.UplegRange, m.prefs.distanceUnit())))
		b.WriteString(fmt.Sprintf("Round-trip light time: %s\n", style.FormatRTLT(target.Rtlt)))
	}

	if upSignals, ok := m.selectedUpSignals(); ok && len(upSignals) > 0 {
		signal := upSignals[m.ui.selection.upSignal]
		b.WriteString("\n=== Up Signal ===\n")
		source := "-"
		if signal.Spacecraft != (model.Spacecraft{}) {
			source = signal.Spacecraft.FriendlyName
		}
		b.WriteString(fmt.Sprintf("Source: %s\n", source))
		b.WriteString(fmt.Sprintf("Signal type: %s\n", style.DefaultIfEmpty(signal.SignalType, "-")))
		b.WriteString(fmt.Sprintf("Frequency band: %s\n", style.DefaultIfEmpty(signal.Band, "-")))
		b.WriteString(fmt.Sprintf("Power transmitted: %s\n", style.FormatPowerTx(signal.Power)))
	}

	if downSignals, ok := m.selectedDownSignals(); ok && len(downSignals) > 0 {
		signal := downSignals[m.ui.selection.downSignal]
		b.WriteString("\n=== Down Signal ===\n")
		source := "-"
		if signal.Spacecraft != (model.Spacecraft{}) {
			source = signal.Spacecraft.FriendlyName
		}
		b.WriteString(fmt.Sprintf("Source: %s\n", source))
		b.WriteString(fmt.Sprintf("Signal type: %s\n", style.DefaultIfEmpty(signal.SignalType, "-")))
		b.WriteString(fmt.Sprintf("Frequency band: %s\n", style.DefaultIfEmpty(signal.Band, "-")))
		b.WriteString(fmt.Sprintf("Data rate: %s\n", style.FormatDataRate(signal.DataRate)))
		b.WriteString(fmt.Sprintf("Power received: %s\n", style.FormatPowerRx(signal.Power)))
	}

	return b.String()
}
