package tui

import (
	"strings"

	"github.com/RustyDaemon/go-dsn-now/internal/model"
)

func (m *Model) restoreInitialStationSelection() {
	if m.ui.selection.station >= 0 || len(m.appData.FullData.Stations) == 0 {
		return
	}

	m.ui.selection.station = 0
	lastStation := m.prefs.lastStation()
	if lastStation == "" {
		return
	}

	for i, station := range m.appData.FullData.Stations {
		if strings.EqualFold(station.FriendlyName, lastStation) {
			m.ui.selection.station = i
			return
		}
	}
}

func (m *Model) clampSelection() {
	stations := m.appData.FullData.Stations
	if len(stations) == 0 {
		m.ui.selection = selectionState{station: -1}
		return
	}

	if m.ui.selection.station < 0 {
		m.ui.selection.station = 0
	}
	if m.ui.selection.station >= len(stations) {
		m.ui.selection.station = len(stations) - 1
	}

	dishes := stations[m.ui.selection.station].Dishes
	if len(dishes) == 0 {
		m.ui.selection.dish = 0
		m.resetDishDependents()
		return
	}

	if m.ui.selection.dish < 0 {
		m.ui.selection.dish = 0
	}
	if m.ui.selection.dish >= len(dishes) {
		m.ui.selection.dish = len(dishes) - 1
	}

	targets := dishes[m.ui.selection.dish].Targets
	if len(targets) == 0 {
		m.ui.selection.target = 0
	} else if m.ui.selection.target >= len(targets) {
		m.ui.selection.target = len(targets) - 1
	}

	upSignals := dishes[m.ui.selection.dish].UpSignals
	if len(upSignals) == 0 {
		m.ui.selection.upSignal = 0
	} else if m.ui.selection.upSignal >= len(upSignals) {
		m.ui.selection.upSignal = len(upSignals) - 1
	}

	downSignals := dishes[m.ui.selection.dish].DownSignals
	if len(downSignals) == 0 {
		m.ui.selection.downSignal = 0
	} else if m.ui.selection.downSignal >= len(downSignals) {
		m.ui.selection.downSignal = len(downSignals) - 1
	}
}

func (m *Model) resetDishSelection() {
	m.ui.selection.dish = 0
	m.resetDishDependents()
}

func (m *Model) resetDishDependents() {
	m.ui.selection.target = 0
	m.ui.selection.upSignal = 0
	m.ui.selection.downSignal = 0
}

func (m *Model) syncDishListSelection() {
	m.dishList.Select(m.ui.selection.dish)
}

func (m *Model) selectDish(index int) {
	station, ok := m.selectedStation()
	if !ok || index < 0 || index >= len(station.Dishes) {
		return
	}

	if m.ui.selection.dish != index {
		m.ui.selection.dish = index
		m.resetDishDependents()
	}

	m.syncDishListSelection()
}

func (m *Model) cycleStation() {
	stations := m.appData.FullData.Stations
	if len(stations) == 0 {
		return
	}

	m.ui.selection.station = (m.ui.selection.station + 1) % len(stations)
	m.resetDishSelection()
	m.prefs.setLastStation(stations[m.ui.selection.station].FriendlyName)
	m.refreshDishList()
}

func (m *Model) cycleTarget() {
	targets, ok := m.selectedTargets()
	if !ok || len(targets) == 0 {
		return
	}

	m.ui.selection.target = (m.ui.selection.target + 1) % len(targets)
}

func (m *Model) cycleUpSignal() {
	signals, ok := m.selectedUpSignals()
	if !ok || len(signals) == 0 {
		return
	}

	m.ui.selection.upSignal = (m.ui.selection.upSignal + 1) % len(signals)
}

func (m *Model) cycleDownSignal() {
	signals, ok := m.selectedDownSignals()
	if !ok || len(signals) == 0 {
		return
	}

	m.ui.selection.downSignal = (m.ui.selection.downSignal + 1) % len(signals)
}

func (m *Model) navigateDish(direction int) {
	station, ok := m.selectedStation()
	if !ok || len(station.Dishes) == 0 {
		return
	}

	newIdx := m.ui.selection.dish + direction
	if newIdx < 0 {
		newIdx = len(station.Dishes) - 1
	} else if newIdx >= len(station.Dishes) {
		newIdx = 0
	}

	m.selectDish(newIdx)
}

func (m *Model) refreshDishList() {
	station, ok := m.selectedStation()
	if !ok {
		return
	}

	m.dishList.SetItems(station.Dishes, m.appData.Bookmarks)
	if m.ui.selection.dish >= len(station.Dishes) {
		m.ui.selection.dish = 0
	}
	m.syncDishListSelection()
}

func (m Model) selectedStation() (model.Station, bool) {
	if m.ui.selection.station < 0 || m.ui.selection.station >= len(m.appData.FullData.Stations) {
		return model.Station{}, false
	}

	return m.appData.FullData.Stations[m.ui.selection.station], true
}

func (m Model) selectedDish() (model.Dish, bool) {
	station, ok := m.selectedStation()
	if !ok || m.ui.selection.dish < 0 || m.ui.selection.dish >= len(station.Dishes) {
		return model.Dish{}, false
	}

	return station.Dishes[m.ui.selection.dish], true
}

func (m Model) selectedTargets() ([]model.Target, bool) {
	dish, ok := m.selectedDish()
	if !ok {
		return nil, false
	}

	return dish.Targets, true
}

func (m Model) selectedUpSignals() ([]model.UpSignal, bool) {
	dish, ok := m.selectedDish()
	if !ok {
		return nil, false
	}

	return dish.UpSignals, true
}

func (m Model) selectedDownSignals() ([]model.DownSignal, bool) {
	dish, ok := m.selectedDish()
	if !ok {
		return nil, false
	}

	return dish.DownSignals, true
}
