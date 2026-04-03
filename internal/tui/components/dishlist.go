package components

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/RustyDaemon/go-dsn-now/internal/model"
	"github.com/RustyDaemon/go-dsn-now/internal/tui/style"
)

// DishItem implements list.Item for a dish entry.
type DishItem struct {
	Dish       model.Dish
	Bookmarked bool
}

func (d DishItem) FilterValue() string { return d.Dish.FriendlyName }
func (d DishItem) Title() string       { return d.Dish.FriendlyName }
func (d DishItem) Description() string { return "" }

// DishDelegate renders each dish list item.
type DishDelegate struct{}

func (d DishDelegate) Height() int                             { return 1 }
func (d DishDelegate) Spacing() int                            { return 0 }
func (d DishDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d DishDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(DishItem)
	if !ok {
		return
	}

	var b strings.Builder

	isSelected := index == m.Index()
	if isSelected {
		b.WriteString(style.PrimaryStyle.Render("▸ "))
	} else {
		b.WriteString("  ")
	}

	if item.Bookmarked {
		b.WriteString(style.AccentStyle.Render("★") + " ")
	}

	workingDown := item.Dish.CountWorkingDownSignals()
	workingUp := item.Dish.CountWorkingUpSignals()
	workingTargets := item.Dish.CountWorkingTargets()

	name := item.Dish.FriendlyName
	if workingDown == 0 && workingUp == 0 && workingTargets == 0 {
		b.WriteString(style.MutedStyle.Render(name))
		b.WriteString(" ")
		b.WriteString(style.MutedStyle.Render("✘"))
	} else if workingDown > 0 && workingUp > 0 && workingTargets > 0 {
		if isSelected {
			b.WriteString(style.PrimaryBoldStyle.Render(name))
		} else {
			b.WriteString(style.PrimaryStyle.Render(name))
		}
	} else {
		if isSelected {
			b.WriteString(lipgloss.NewStyle().Foreground(style.ColorTextBright).Bold(true).Render(name))
		} else {
			b.WriteString(lipgloss.NewStyle().Foreground(style.ColorTextNormal).Render(name))
		}
	}

	b.WriteString(" ")
	for i := 0; i < workingUp; i++ {
		b.WriteString(style.SignalUpStyle.Render("↑"))
	}
	for i := 0; i < workingDown; i++ {
		b.WriteString(style.SignalDownStyle.Render("↓"))
	}

	line := b.String()

	if isSelected {
		// Pad to full width and apply highlight background
		lineWidth := lipgloss.Width(line)
		availWidth := m.Width()
		if lineWidth < availWidth {
			line = line + strings.Repeat(" ", availWidth-lineWidth)
		}
		line = lipgloss.NewStyle().Background(style.ColorBgHighlight).Render(line)
	}

	fmt.Fprint(w, line)
}

// DishList wraps a bubbles/list.Model for dish navigation.
type DishList struct {
	list list.Model
}

func NewDishList() DishList {
	delegate := DishDelegate{}
	l := list.New(nil, delegate, 0, 0)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetShowFilter(false)
	l.SetShowHelp(false)
	l.SetShowPagination(false)
	l.DisableQuitKeybindings()
	l.Styles.NoItems = style.DimStyle

	return DishList{list: l}
}

func (d *DishList) SetSize(w, h int) {
	d.list.SetWidth(w - 4) // account for panel border + padding
	d.list.SetHeight(h)
}

func (d *DishList) SetItems(dishes []model.Dish, bookmarks map[string]bool) {
	items := make([]list.Item, len(dishes))
	for i, dish := range dishes {
		items[i] = DishItem{
			Dish:       dish,
			Bookmarked: bookmarks[dish.Name],
		}
	}
	d.list.SetItems(items)
}

func (d *DishList) Index() int {
	return d.list.Index()
}

func (d *DishList) Select(index int) {
	d.list.Select(index)
}

// VisibleOffset returns the index of the first visible item in the list.
func (d *DishList) VisibleOffset() int {
	return d.list.Paginator.Page * d.list.Paginator.PerPage
}

func (d *DishList) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	d.list, cmd = d.list.Update(msg)
	return cmd
}

func (d DishList) View(width, height int, selectedIdx, total int) string {
	title := style.TitleStyle.Render("Antennas")
	if total > 0 {
		title = style.TitleStyle.Render(fmt.Sprintf("Antennas (%d/%d)", selectedIdx+1, total))
	}

	content := d.list.View()

	return style.RenderTitledPanel(title, content, width, style.ColorBorderActive)
}
