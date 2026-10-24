package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/x/ansi"

	"github.com/RustyDaemon/go-dsn-now/internal/model"
	"github.com/RustyDaemon/go-dsn-now/internal/tui/style"
)

func (m Model) matchesDish(station model.Station, dish model.Dish) bool {
	if m.ui.activeOnly && dish.CountWorkingUpSignals()+dish.CountWorkingDownSignals() == 0 {
		return false
	}
	if m.ui.bookmarkedOnly && !m.appData.Bookmarks[dish.Name] {
		return false
	}
	query := strings.ToLower(strings.TrimSpace(m.search.Value()))
	if query == "" {
		return true
	}
	fields := []string{station.Name, station.FriendlyName, dish.Name, dish.FriendlyName}
	addSpacecraft := func(s model.Spacecraft) {
		fields = append(fields, s.Name, s.FriendlyName, s.ExplorerName)
	}
	for _, target := range dish.Targets {
		fields = append(fields, target.Name, target.ID)
		addSpacecraft(target.Spacecraft)
	}
	for _, signal := range dish.UpSignals {
		addSpacecraft(signal.Spacecraft)
	}
	for _, signal := range dish.DownSignals {
		addSpacecraft(signal.Spacecraft)
	}
	text := strings.ToLower(strings.Join(fields, " "))
	for _, word := range strings.Fields(query) {
		if !strings.Contains(text, word) {
			return false
		}
	}
	return true
}

// Dish list positions refer to filtered rows; selection positions refer to source data.
func (m Model) visibleDishIndices() []int {
	station, ok := m.selectedStation()
	if !ok {
		return nil
	}
	var indices []int
	for i, dish := range station.Dishes {
		if m.matchesDish(station, dish) {
			indices = append(indices, i)
		}
	}
	return indices
}

func (m *Model) clearFilters() {
	m.search.SetValue("")
	m.ui.activeOnly = false
	m.ui.bookmarkedOnly = false
	m.refreshCompactTable()
}

func (m Model) viewFilters() string {
	if m.ui.searchEditing {
		return ansi.Truncate(m.search.View()+"  Enter apply · Esc cancel", max(1, m.ui.width), "…")
	}
	search := "/ search"
	if m.search.Value() != "" {
		search = fmt.Sprintf("/ %q", m.search.Value())
	}
	state := func(enabled bool) string {
		if enabled {
			return "on"
		}
		return "off"
	}
	text := fmt.Sprintf("%s  a active:%s  f starred:%s  %d dishes  Esc clear",
		search, state(m.ui.activeOnly), state(m.ui.bookmarkedOnly), m.compactTable.RowCount())
	return style.LabelStyle.Render(ansi.Truncate(text, max(1, m.ui.width), "…"))
}
