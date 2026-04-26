package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/RustyDaemon/go-dsn-now/internal/tui/style"
)

type ModalType int

const (
	ModalNone ModalType = iota
	ModalHelp
)

type Modal struct {
	viewport viewport.Model
	width    int
	height   int
	title    string
}

func NewModal() Modal {
	vp := viewport.New(0, 0)
	return Modal{viewport: vp}
}

func (m *Modal) SetSize(w, h int) {
	m.width = w
	m.height = h
	vpW := w - 4
	vpH := h - 4
	if vpW < 40 {
		vpW = 40
	}
	if vpH < 10 {
		vpH = 10
	}
	m.viewport.Width = vpW - 6
	m.viewport.Height = vpH - 6
}

func (m *Modal) SetContent(title, content string) {
	m.title = title
	m.viewport.SetContent(content)
	m.viewport.GotoTop()
}

func (m *Modal) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return cmd
}

func (m Modal) View() string {
	vpW := m.width - 4
	vpH := m.height - 4
	if vpW < 40 {
		vpW = 40
	}
	if vpH < 10 {
		vpH = 10
	}

	content := m.viewport.View()
	title := style.ModalTitleStyle.Render(m.title)

	innerWidth := vpW - 2
	if innerWidth < 4 {
		innerWidth = 4
	}

	bc := lipgloss.NewStyle().Foreground(style.ColorBorderModal)
	titleRendered := " " + title + " "
	titleVisualWidth := lipgloss.Width(titleRendered)
	topFillLen := innerWidth - 1 - titleVisualWidth
	if topFillLen < 0 {
		topFillLen = 0
	}

	topLine := bc.Render("╔") + bc.Render("═") + titleRendered + bc.Render(strings.Repeat("═", topFillLen)) + bc.Render("╗")

	scrollUp := !m.viewport.AtTop()
	scrollDown := !m.viewport.AtBottom()
	scrollStyle := lipgloss.NewStyle().Foreground(style.ColorPrimary)
	var arrowLeft, arrowRight string
	if scrollUp {
		arrowLeft = scrollStyle.Render("▲")
	} else {
		arrowLeft = bc.Render("═")
	}
	if scrollDown {
		arrowRight = scrollStyle.Render("▼")
	} else {
		arrowRight = bc.Render("═")
	}
	bottomFillLen := innerWidth - 2
	if bottomFillLen < 0 {
		bottomFillLen = 0
	}
	bottomLine := bc.Render("╚") + arrowLeft + bc.Render(strings.Repeat("═", bottomFillLen)) + arrowRight + bc.Render("╝")

	contentStyle := lipgloss.NewStyle().
		Width(innerWidth - 4).
		Height(vpH - 4).
		PaddingLeft(2).
		PaddingRight(2)

	styledContent := contentStyle.Render(content)
	contentLines := strings.Split(styledContent, "\n")

	var bodyLines []string
	for _, line := range contentLines {
		lineWidth := lipgloss.Width(line)
		pad := innerWidth - lineWidth
		if pad < 0 {
			pad = 0
		}
		bodyLines = append(bodyLines, bc.Render("║")+line+strings.Repeat(" ", pad)+bc.Render("║"))
	}

	lines := make([]string, 0, len(bodyLines)+2)
	lines = append(lines, topLine)
	lines = append(lines, bodyLines...)
	lines = append(lines, bottomLine)

	box := strings.Join(lines, "\n")

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		box,
		lipgloss.WithWhitespaceBackground(style.ColorBgDeep),
	)
}

func BuildHelpContent(version, githubURL string) string {
	var b strings.Builder
	b.WriteString(style.TitleStyle.Render("Keybindings") + "\n\n")

	keys := []struct{ key, desc string }{
		{"s", "Cycle station"},
		{"t", "Cycle target"},
		{"u", "Cycle up signal"},
		{"d", "Cycle down signal"},
		{"↑ ↓", "Navigate dishes"},
		{"b", "Bookmark dish"},
		{"c", "Toggle compact view"},
		{"S", "Cycle compact sort"},
		{"T", "Cycle theme"},
		{"U", "Cycle distance unit"},
		{"y", "Copy to clipboard"},
		{"+ -", "Adjust refresh interval"},
		{"?", "This help"},
		{"Esc", "Close modal"},
		{"q", "Quit"},
	}

	for _, k := range keys {
		fmt.Fprintf(&b, "  %s  %s\n",
			style.PrimaryStyle.Render(fmt.Sprintf("%-6s", k.key)),
			style.LabelStyle.Render(k.desc),
		)
	}

	b.WriteString("\n" + style.TitleStyle.Render("About") + "\n\n")
	fmt.Fprintf(&b, "  Version  %s\n", style.ValueStyle.Render(version))
	fmt.Fprintf(&b, "  GitHub   %s\n", style.ValueStyle.Render(githubURL))

	return b.String()
}
