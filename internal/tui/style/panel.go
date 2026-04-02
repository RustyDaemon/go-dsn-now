package style

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderTitledPanel draws a panel with a title embedded in the top border.
// It manually constructs the border to avoid ANSI escape code corruption
// that occurs when injecting styled text into lipgloss border strings.
func RenderTitledPanel(title, content string, width int, borderColor lipgloss.TerminalColor) string {
	if width < 6 {
		width = 6
	}

	innerWidth := width - 2 // left + right border columns

	// Build top border: ╭─ Title ─────────╮
	titleRendered := ""
	if title != "" {
		titleRendered = " " + title + " "
	}
	titleVisualWidth := lipgloss.Width(titleRendered)

	topFillLen := innerWidth - 1 - titleVisualWidth // -1 for the dash before title
	if topFillLen < 0 {
		topFillLen = 0
	}

	bc := lipgloss.NewStyle().Foreground(borderColor)

	topLine := bc.Render("╭") + bc.Render("─") + titleRendered + bc.Render(strings.Repeat("─", topFillLen)) + bc.Render("╮")

	// Build bottom border: ╰─────────────────╯
	bottomLine := bc.Render("╰") + bc.Render(strings.Repeat("─", innerWidth)) + bc.Render("╯")

	// Style and pad content lines
	contentStyle := lipgloss.NewStyle().
		Width(innerWidth).
		PaddingLeft(1).
		PaddingRight(1)

	styledContent := contentStyle.Render(content)

	// Wrap each content line with side borders
	contentLines := strings.Split(styledContent, "\n")
	var bodyLines []string
	for _, line := range contentLines {
		lineWidth := lipgloss.Width(line)
		pad := innerWidth - lineWidth
		if pad < 0 {
			pad = 0
		}
		bodyLines = append(bodyLines, bc.Render("│")+line+strings.Repeat(" ", pad)+bc.Render("│"))
	}

	// Assemble
	lines := make([]string, 0, len(bodyLines)+2)
	lines = append(lines, topLine)
	lines = append(lines, bodyLines...)
	lines = append(lines, bottomLine)

	return strings.Join(lines, "\n")
}

// RenderPanel draws a simple bordered panel without a title.
func RenderPanel(content string, width int, borderColor lipgloss.TerminalColor) string {
	return RenderTitledPanel("", content, width, borderColor)
}
