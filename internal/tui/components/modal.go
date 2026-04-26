package components

import (
	"fmt"
	"strings"

	"encoding/json"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/RustyDaemon/go-dsn-now/internal/model"
	"github.com/RustyDaemon/go-dsn-now/internal/tui/style"
)

type ModalType int

const (
	ModalNone ModalType = iota
	ModalHelp
	ModalJSONPreview
	ModalDishSpecs
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

func BuildJSONContent(dish model.Dish) string {
	j, err := json.MarshalIndent(dish, "", "  ")
	if err != nil {
		return "Error formatting JSON"
	}
	return colorizeJSON(string(j))
}

func colorizeJSON(s string) string {
	keyStyle := lipgloss.NewStyle().Foreground(style.ColorPrimary).Bold(true)
	stringStyle := lipgloss.NewStyle().Foreground(style.ColorSignalDown)
	numberStyle := lipgloss.NewStyle().Foreground(style.ColorSignalUp)
	boolStyle := lipgloss.NewStyle().Foreground(style.ColorSecondary).Bold(true)
	nullStyle := lipgloss.NewStyle().Foreground(style.ColorTextMuted).Italic(true)
	punctStyle := lipgloss.NewStyle().Foreground(style.ColorTextDim)

	var b strings.Builder
	for i := 0; i < len(s); {
		c := s[i]
		switch {
		case c == ' ' || c == '\n' || c == '\t' || c == '\r':
			b.WriteByte(c)
			i++
		case c == '"':
			j := i + 1
			for j < len(s) {
				if s[j] == '\\' && j+1 < len(s) {
					j += 2
					continue
				}
				if s[j] == '"' {
					j++
					break
				}
				j++
			}
			tok := s[i:j]
			k := j
			for k < len(s) && (s[k] == ' ' || s[k] == '\t') {
				k++
			}
			if k < len(s) && s[k] == ':' {
				b.WriteString(keyStyle.Render(tok))
			} else {
				b.WriteString(stringStyle.Render(tok))
			}
			i = j
		case c == '{' || c == '}' || c == '[' || c == ']' || c == ',' || c == ':':
			b.WriteString(punctStyle.Render(string(c)))
			i++
		case c == '-' || (c >= '0' && c <= '9'):
			j := i + 1
			for j < len(s) {
				ch := s[j]
				if ch == '.' || ch == 'e' || ch == 'E' || ch == '+' || ch == '-' || (ch >= '0' && ch <= '9') {
					j++
					continue
				}
				break
			}
			b.WriteString(numberStyle.Render(s[i:j]))
			i = j
		case strings.HasPrefix(s[i:], "true"):
			b.WriteString(boolStyle.Render("true"))
			i += 4
		case strings.HasPrefix(s[i:], "false"):
			b.WriteString(boolStyle.Render("false"))
			i += 5
		case strings.HasPrefix(s[i:], "null"):
			b.WriteString(nullStyle.Render("null"))
			i += 4
		default:
			b.WriteByte(c)
			i++
		}
	}
	return b.String()
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
		{"J", "JSON preview"},
		{"i", "Antenna specs"},
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

func BuildDishSpecsContent(spec model.DishSpecification) string {
	var b strings.Builder

	fields := []struct{ label, value string }{
		{"Name", spec.Name},
		{"Type", spec.Type},
		{"Diameter", spec.Diameter},
		{"", ""},
		{"Transmitters frequency", spec.TransmittersFrequency},
		{"Receivers frequency", spec.ReceiversFrequency},
		{"Transmitters power", spec.TransmittersPower},
		{"Precision", spec.Precision},
		{"Antenna speed", spec.AntennaSpeed},
		{"", ""},
		{"Total weight", spec.TotalWeight},
		{"Dish weight", spec.DishWeight},
		{"Total panels", spec.TotalPanels},
		{"Surface area", spec.SurfaceArea},
		{"", ""},
		{"Operational wind resistance", spec.OperationalWindResistance},
		{"Max wind resistance", spec.WindResistance},
		{"Built in", spec.BuiltIn},
		{"Web URL", spec.WebUrl},
	}

	for _, f := range fields {
		if f.label == "" {
			b.WriteString("\n")
			continue
		}
		fmt.Fprintf(&b, "  %-30s %s\n",
			style.LabelStyle.Render(f.label+":"),
			style.ValueStyle.Render(style.DashIfEmpty(f.value)),
		)
	}

	return b.String()
}
