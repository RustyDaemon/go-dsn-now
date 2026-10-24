package components

import (
	"fmt"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestCompactSelectionScrollsIntoView(t *testing.T) {
	c := NewCompactTable()
	c.SetSize(100, 12)
	rows := make([]CompactRow, 30)
	for i := range rows {
		name := fmt.Sprintf("dish%02d", i)
		rows[i] = CompactRow{StationID: "station", DishID: name, Dish: name}
	}
	c.SetRows(rows)
	for _, name := range []string{"dish29", "dish00", "dish15"} {
		if !c.SelectDish("station", name) || !strings.Contains(c.View(), name) {
			t.Fatalf("selected dish %s is outside the visible viewport", name)
		}
	}
}

func TestStatusBarConstrainsLongErrors(t *testing.T) {
	s := NewStatusBar()
	s.SetWidth(70)
	view := s.View(StatusBarParams{DefaultStatus: true, Compact: true, LastError: strings.Repeat("feed unavailable ", 30), ConnStatus: "degraded", LastUpdated: "12:00:00", RefreshInterval: "10s"})
	if lipgloss.Width(view) > 70 || lipgloss.Height(view) > 3 {
		t.Fatalf("status bar overflowed: %dx%d", lipgloss.Width(view), lipgloss.Height(view))
	}
}
