package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/RustyDaemon/go-dsn-now/internal/tui/style"
)

type CompactRow struct {
	Station     string
	Dish        string
	Target      string
	UpSignals   int
	DownSignals int
	Activity    string
}

type CompactTable struct {
	table     table.Model
	width     int
	height    int
	sortLabel string
}

func NewCompactTable() CompactTable {
	columns := []table.Column{
		{Title: "Station", Width: 12},
		{Title: "Dish", Width: 12},
		{Title: "Target", Width: 20},
		{Title: "↑", Width: 4},
		{Title: "↓", Width: 4},
		{Title: "Activity", Width: 14},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(style.ColorBorder).
		BorderBottom(true).
		Bold(true).
		Foreground(style.ColorPrimary)
	s.Selected = s.Selected.
		Foreground(style.ColorTextBright).
		Background(style.ColorBgHighlight).
		Bold(true)
	s.Cell = s.Cell.
		Foreground(style.ColorTextNormal)

	t.SetStyles(s)

	return CompactTable{table: t, sortLabel: "Default"}
}

func (c *CompactTable) SetSize(w, h int) {
	c.width = w
	c.height = h
	c.table.SetWidth(w - 4)
	c.table.SetHeight(h - 6)

	available := w - 8
	if available < 40 {
		available = 40
	}
	c.table.SetColumns([]table.Column{
		{Title: "Station", Width: available * 15 / 100},
		{Title: "Dish", Width: available * 15 / 100},
		{Title: "Target", Width: available * 30 / 100},
		{Title: "↑", Width: available * 8 / 100},
		{Title: "↓", Width: available * 8 / 100},
		{Title: "Activity", Width: available * 20 / 100},
	})
}

func (c *CompactTable) SetSortLabel(label string) {
	c.sortLabel = label
}

func (c *CompactTable) SetRows(rows []CompactRow) {
	tableRows := make([]table.Row, len(rows))
	for i, r := range rows {
		upStr := fmt.Sprintf("%d", r.UpSignals)
		downStr := fmt.Sprintf("%d", r.DownSignals)
		tableRows[i] = table.Row{r.Station, r.Dish, r.Target, upStr, downStr, r.Activity}
	}
	c.table.SetRows(tableRows)
}

func (c *CompactTable) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	c.table, cmd = c.table.Update(msg)
	return cmd
}

func (c CompactTable) View() string {
	title := style.TitleStyle.Render(fmt.Sprintf("All Dishes (compact) | Sort: %s", c.sortLabel))
	content := c.table.View()
	return style.RenderTitledPanel(title, content, c.width, style.ColorBorderActive)
}

func (c *CompactTable) RefreshStyles() {
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(style.ColorBorder).
		BorderBottom(true).
		Bold(true).
		Foreground(style.ColorPrimary)
	s.Selected = s.Selected.
		Foreground(style.ColorTextBright).
		Background(style.ColorBgHighlight).
		Bold(true)
	s.Cell = s.Cell.
		Foreground(style.ColorTextNormal)
	c.table.SetStyles(s)
}

func (c CompactTable) GetVisibleContent() string {
	var b strings.Builder
	b.WriteString("Station\tDish\tTarget\t↑\t↓\tActivity\n")
	for _, row := range c.table.Rows() {
		b.WriteString(strings.Join(row, "\t") + "\n")
	}
	return b.String()
}
