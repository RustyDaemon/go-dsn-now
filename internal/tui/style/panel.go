package style

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func RenderTitledPanel(title, content string, width int, borderColor lipgloss.TerminalColor) string {
	if width < 6 {
		width = 6
	}

	innerWidth := width - 2

	titleRendered := ""
	if title != "" {
		titleRendered = " " + title + " "
	}
	titleVisualWidth := lipgloss.Width(titleRendered)

	topFillLen := innerWidth - 1 - titleVisualWidth
	if topFillLen < 0 {
		topFillLen = 0
	}

	bc := lipgloss.NewStyle().Foreground(borderColor)

	topLine := bc.Render("╭") + bc.Render("─") + titleRendered + bc.Render(strings.Repeat("─", topFillLen)) + bc.Render("╮")

	bottomLine := bc.Render("╰") + bc.Render(strings.Repeat("─", innerWidth)) + bc.Render("╯")

	contentStyle := lipgloss.NewStyle().
		Width(innerWidth).
		PaddingLeft(1).
		PaddingRight(1)

	styledContent := contentStyle.Render(content)

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

	lines := make([]string, 0, len(bodyLines)+2)
	lines = append(lines, topLine)
	lines = append(lines, bodyLines...)
	lines = append(lines, bottomLine)

	return strings.Join(lines, "\n")
}

func RenderPanel(content string, width int, borderColor lipgloss.TerminalColor) string {
	return RenderTitledPanel("", content, width, borderColor)
}
