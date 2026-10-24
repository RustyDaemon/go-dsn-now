package tui

import (
	"reflect"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/RustyDaemon/go-dsn-now/internal/config"
	"github.com/RustyDaemon/go-dsn-now/internal/data"
	"github.com/RustyDaemon/go-dsn-now/internal/model"
	"github.com/RustyDaemon/go-dsn-now/internal/model/response"
)

func testModel() Model {
	app := model.NewAppData()
	voyager := model.Spacecraft{Name: "vgr1", FriendlyName: "Voyager 1"}
	webb := model.Spacecraft{Name: "jwst", FriendlyName: "James Webb Space Telescope"}
	app.FullData.Stations = []model.Station{
		{Name: "gdscc", FriendlyName: "Goldstone", Dishes: []model.Dish{
			{Name: "dss26", FriendlyName: "DSS 26"},
			{Name: "dss14", FriendlyName: "DSS 14", DownSignals: []model.DownSignal{{IsActive: true, Spacecraft: voyager}}, Targets: []model.Target{{Name: "vgr1", Spacecraft: voyager}}},
			{Name: "dss25", FriendlyName: "DSS 25", UpSignals: []model.UpSignal{{IsActive: true, Spacecraft: webb}}},
		}},
		{Name: "mdscc", FriendlyName: "Madrid", Dishes: []model.Dish{
			{Name: "dss63", FriendlyName: "DSS 63", UpSignals: []model.UpSignal{{IsActive: true}, {IsActive: true}}},
		}},
		{Name: "cdscc", FriendlyName: "Canberra", Dishes: []model.Dish{
			{Name: "dss43", FriendlyName: "DSS 43", Targets: []model.Target{{Name: "unmapped-mission"}}},
		}},
	}
	m := NewModel(config.Load(), &data.Settings{}, nil, app)
	m.ui.width, m.ui.height = 120, 36
	m.ui.ready, m.ui.configLoaded, m.ui.loading = true, true, false
	m.ui.selection = selectionState{station: 0}
	m.resizeComponents()
	m.refreshCompactTable()
	return m
}

func press(m Model, key string) (Model, tea.Cmd) {
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)}
	switch key {
	case "enter":
		msg = tea.KeyMsg{Type: tea.KeyEnter}
	case "esc":
		msg = tea.KeyMsg{Type: tea.KeyEsc}
	case "down":
		msg = tea.KeyMsg{Type: tea.KeyDown}
	case "up":
		msg = tea.KeyMsg{Type: tea.KeyUp}
	case "ctrl+c":
		msg = tea.KeyMsg{Type: tea.KeyCtrlC}
	}
	updated, cmd := m.Update(msg)
	return updated.(Model), cmd
}

func requireDish(t *testing.T, m Model, name string) {
	t.Helper()
	dish, ok := m.selectedDish()
	if !ok || dish.Name != name {
		t.Fatalf("selected dish = %q, %v; want %q", dish.Name, ok, name)
	}
	row, ok := m.compactTable.SelectedRow()
	if !ok || row.DishID != name {
		t.Fatalf("table dish = %q, %v; want %q", row.DishID, ok, name)
	}
}

func TestCompactSelectionSurvivesSortRefreshAndViewChanges(t *testing.T) {
	m := testModel()
	m, _ = press(m, "c")
	m, _ = press(m, "down")
	requireDish(t, m, "dss14")
	m, _ = press(m, "S") // activity
	m, _ = press(m, "S") // signal count, with Madrid sorted above Goldstone
	requireDish(t, m, "dss14")
	m, _, _ = m.handleDataMessage(DSNDataLoadedMsg{Data: response.DSN{D: []response.Dish{{Name: "dss14"}}}})
	requireDish(t, m, "dss14")
	m, _ = press(m, "enter")
	if m.useCompactView() || !strings.Contains(m.getVisibleContent(), "DSS 14") {
		t.Fatal("Enter did not open the highlighted dish")
	}
	m, _ = press(m, "down")
	requireDish(t, m, "dss25")
	m, _ = press(m, "c")
	requireDish(t, m, "dss25")
}

func TestBookmarkUsesHighlightedDishAcrossStations(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	m := testModel()
	m, _ = press(m, "c")
	for range 3 {
		m, _ = press(m, "down")
	}
	requireDish(t, m, "dss63")
	m, _ = press(m, "b")
	if !m.appData.Bookmarks["dss63"] || len(m.appData.Bookmarks) != 1 {
		t.Fatalf("wrong dish bookmarked: %v", m.appData.Bookmarks)
	}
	row, _ := m.compactTable.SelectedRow()
	if !row.Bookmarked || !data.LoadBookmarks()["dss63"] {
		t.Fatal("bookmark not displayed or persisted")
	}
	m, _ = press(m, "f")
	requireDish(t, m, "dss63")
	m, _ = press(m, "b")
	if m.compactTable.RowCount() != 0 || !strings.Contains(m.View(), "No matching dishes") {
		t.Fatal("unbookmarking the last filtered row did not show empty results")
	}
	m, _ = press(m, "b") // no hidden dish should be bookmarked
	if len(m.appData.Bookmarks) != 0 {
		t.Fatal("empty results still have an actionable dish")
	}
	m, _ = press(m, "esc")
	if m.compactTable.RowCount() != 5 {
		t.Fatal("clearing filters did not restore rows")
	}
	requireDish(t, m, "dss26")
}

func TestSearchAndCombinedFilters(t *testing.T) {
	for _, tc := range []struct {
		name, query     string
		active, starred bool
		want            []string
	}{
		{name: "station name", query: "mAdRiD", want: []string{"dss63"}},
		{name: "station ID", query: "MDSCC", want: []string{"dss63"}},
		{name: "dish ID", query: "dss14", want: []string{"dss14"}},
		{name: "spacecraft name", query: "voyager", want: []string{"dss14"}},
		{name: "signal-only spacecraft", query: "webb", want: []string{"dss25"}},
		{name: "unknown spacecraft", query: "unmapped", want: []string{"dss43"}},
		{name: "multiple words", query: " goldstone  voyager ", want: []string{"dss14"}},
		{name: "active signals", active: true, want: []string{"dss14", "dss25", "dss63"}},
		{name: "bookmarks", starred: true, want: []string{"dss14", "dss43"}},
		{name: "combined", query: "goldstone", active: true, starred: true, want: []string{"dss14"}},
		{name: "no matches", query: "saturn"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			m := testModel()
			m.appData.Bookmarks = map[string]bool{"dss14": true, "dss43": true}
			m.search.SetValue(tc.query)
			m.ui.activeOnly, m.ui.bookmarkedOnly = tc.active, tc.starred
			m.refreshCompactTable()
			var got []string
			for _, row := range m.buildCompactRows() {
				got = append(got, row.DishID)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("matches = %v; want %v", got, tc.want)
			}
			if len(tc.want) == 0 {
				if _, ok := m.selectedDish(); ok || !strings.Contains(m.View(), "No matching dishes") {
					t.Fatal("empty results expose a hidden dish")
				}
			}
		})
	}
}

func TestDetailedNavigationUsesFilteredIndices(t *testing.T) {
	m := testModel()
	m.search.SetValue("goldstone")
	m, _ = press(m, "a")
	requireDish(t, m, "dss14")
	if m.dishList.Index() != 0 {
		t.Fatal("filtered list cursor is not aligned")
	}
	m, _ = press(m, "down")
	requireDish(t, m, "dss25")
	m, _ = press(m, "down")
	requireDish(t, m, "dss14")
	updated, _ := m.Update(tea.MouseMsg{X: 2, Y: 3, Button: tea.MouseButtonLeft, Action: tea.MouseActionPress})
	m = updated.(Model)
	requireDish(t, m, "dss25")
	m, _, _ = m.handleDataMessage(DSNDataLoadedMsg{Data: response.DSN{D: []response.Dish{{Name: "dss25"}}}})
	requireDish(t, m, "dss14") // the selected dish no longer passes active-only
}

func TestSearchEditingDoesNotTriggerShortcuts(t *testing.T) {
	m := testModel()
	m, _ = press(m, "/")
	for _, letter := range "qaf" {
		var cmd tea.Cmd
		m, cmd = press(m, string(letter))
		if cmd != nil {
			t.Fatal("typing a search triggered a command")
		}
	}
	if m.search.Value() != "qaf" || m.ui.activeOnly || m.ui.bookmarkedOnly {
		t.Fatal("search text triggered a shortcut")
	}
	m, _ = press(m, "esc")
	if m.search.Value() != "" || m.ui.searchEditing || m.compactTable.RowCount() != 5 {
		t.Fatal("cancel did not restore the prior search")
	}
	m, _ = press(m, "/")
	m, _ = press(m, "voyager")
	m, _ = press(m, "enter")
	requireDish(t, m, "dss14")
	if m.ui.searchEditing || !m.useCompactView() {
		t.Fatal("applying a search should keep the results overview")
	}
	m, _ = press(m, "/")
	m, _ = press(m, "x")
	m, _ = press(m, "esc")
	if m.search.Value() != "voyager" {
		t.Fatal("cancel did not preserve the applied search")
	}
	m, _ = press(m, "esc")
	if m.search.Value() != "" || m.compactTable.RowCount() != 5 {
		t.Fatal("Esc did not clear the applied search")
	}
}

func TestCompactMouseNavigationAndSmallTerminal(t *testing.T) {
	m := testModel()
	m.ui.compactView = true
	updated, _ := m.Update(tea.MouseMsg{Button: tea.MouseButtonWheelDown})
	m = updated.(Model)
	requireDish(t, m, "dss14")
	updated, _ = m.Update(tea.WindowSizeMsg{Width: 70, Height: 18})
	m = updated.(Model)
	m, cmd := press(m, "enter")
	if !m.useCompactView() || cmd == nil || !strings.Contains(m.ui.statusMessage, "80×20") {
		t.Fatal("small terminal should explain why details cannot open")
	}
	_ = m.View()
}

func TestCompactViewFitsTerminal(t *testing.T) {
	for _, size := range [][2]int{{120, 36}, {80, 24}, {70, 18}} {
		m := testModel()
		m.ui.compactView = true
		m.ui.width, m.ui.height = size[0], size[1]
		m.resizeComponents()
		view := m.View()
		if lipgloss.Width(view) > size[0] || lipgloss.Height(view) > size[1] {
			t.Errorf("%dx%d view is %dx%d", size[0], size[1], lipgloss.Width(view), lipgloss.Height(view))
		}
	}
}

func BenchmarkCompactView(b *testing.B) {
	for _, discardedDetails := range []bool{false, true} {
		name := "Current"
		if discardedDetails {
			name = "WithDiscardedDetails"
		}
		b.Run(name, func(b *testing.B) {
			m := testModel()
			m.ui.compactView = true
			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				if discardedDetails {
					_ = m.viewDetailed()
				}
				_ = m.View()
			}
		})
	}
}
