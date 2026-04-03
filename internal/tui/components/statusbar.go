package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/RustyDaemon/go-dsn-now/internal/tui/style"
)

type StatusBarParams struct {
	DefaultStatus   bool
	HasTargets      bool
	HasUpSignals    bool
	HasDownSignals  bool
	HasAntennaSpec  bool
	LastUpdated     string
	LastError       string
	ConnStatus      string
	SignalChanges   []string
	RefreshInterval string
	StatusMessage   string
	DistanceUnit    string
}

type StatusBar struct {
	width int
}

func NewStatusBar() StatusBar {
	return StatusBar{}
}

func (s *StatusBar) SetWidth(w int) {
	s.width = w
}

func (s StatusBar) View(p StatusBarParams) string {
	left := s.renderLeft(p)
	center := s.renderCenter(p)
	right := s.renderRight(p)

	innerWidth := s.width - 4
	if innerWidth < 10 {
		innerWidth = 10
	}

	rightWidth := lipgloss.Width(right)
	centerWidth := lipgloss.Width(center)

	leftWidth := innerWidth - rightWidth - centerWidth
	if leftWidth < 0 {
		leftWidth = 0
	}

	bar := lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.NewStyle().Width(leftWidth).Render(left),
		lipgloss.NewStyle().Width(centerWidth).Render(center),
		lipgloss.NewStyle().Width(rightWidth).Render(right),
	)

	return style.RenderTitledPanel("", bar, s.width, style.ColorBorder)
}

func (s StatusBar) renderLeft(p StatusBarParams) string {
	if !p.DefaultStatus {
		return style.PrimaryStyle.Render("(Esc)") + " close"
	}

	var parts []string
	parts = append(parts, style.PrimaryStyle.Render("s")+"tation")
	if p.HasTargets {
		parts = append(parts, style.PrimaryStyle.Render("t")+"arget")
	}
	if p.HasUpSignals {
		parts = append(parts, style.PrimaryStyle.Render("u")+"p")
	}
	if p.HasDownSignals {
		parts = append(parts, style.PrimaryStyle.Render("d")+"own")
	}
	if p.HasAntennaSpec {
		parts = append(parts, style.PrimaryStyle.Render("i")+"nfo")
	}
	parts = append(parts, style.PrimaryStyle.Render("J")+"SON")
	parts = append(parts, style.PrimaryStyle.Render("?"))
	parts = append(parts, style.PrimaryStyle.Render("q")+"uit")

	return strings.Join(parts, " ")
}

func (s StatusBar) renderCenter(p StatusBarParams) string {
	if p.StatusMessage != "" {
		return " " + style.PrimaryStyle.Render(p.StatusMessage) + " "
	}
	if p.LastError != "" {
		return " " + style.ErrorStyle.Render("Error: "+p.LastError) + " "
	}
	if len(p.SignalChanges) > 0 {
		return " " + style.AccentStyle.Render(strings.Join(p.SignalChanges, ", ")) + " "
	}
	return ""
}

func (s StatusBar) renderRight(p StatusBarParams) string {
	var parts []string

	switch p.ConnStatus {
	case "connected":
		parts = append(parts, style.ConnectedDot)
	case "degraded":
		parts = append(parts, style.DegradedDot)
	case "disconnected":
		parts = append(parts, style.DisconnectedDot)
	}

	if p.LastUpdated != "" {
		parts = append(parts, "upd "+style.ValueStyle.Render(p.LastUpdated))
	}
	if p.RefreshInterval != "" {
		parts = append(parts, style.ValueStyle.Render(p.RefreshInterval))
	}
	unit := p.DistanceUnit
	if unit == "" {
		unit = "km"
	}
	parts = append(parts, style.MutedStyle.Render(unit))

	return " " + strings.Join(parts, " ") + " "
}

