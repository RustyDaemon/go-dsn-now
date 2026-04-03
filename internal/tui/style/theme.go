package style

import "github.com/charmbracelet/lipgloss"

// Theme defines a complete color scheme.
type Theme struct {
	Name string

	BgDeep      lipgloss.TerminalColor
	BgPanel     lipgloss.TerminalColor
	BgHighlight lipgloss.TerminalColor

	Primary    lipgloss.TerminalColor
	PrimaryDim lipgloss.TerminalColor

	Secondary    lipgloss.TerminalColor
	SecondaryDim lipgloss.TerminalColor

	SignalUp   lipgloss.TerminalColor
	SignalDown lipgloss.TerminalColor

	StatusOK    lipgloss.TerminalColor
	StatusWarn  lipgloss.TerminalColor
	StatusError lipgloss.TerminalColor

	TextBright lipgloss.TerminalColor
	TextNormal lipgloss.TerminalColor
	TextDim    lipgloss.TerminalColor
	TextMuted  lipgloss.TerminalColor

	Border       lipgloss.TerminalColor
	BorderActive lipgloss.TerminalColor
	BorderModal  lipgloss.TerminalColor
}

var themes = []Theme{
	{
		Name:         "cosmic",
		BgDeep:       lipgloss.Color("#0D0221"),
		BgPanel:      lipgloss.Color("#1A0A3E"),
		BgHighlight:  lipgloss.Color("#2D1B69"),
		Primary:      lipgloss.Color("#00D4FF"),
		PrimaryDim:   lipgloss.Color("#4FC3F7"),
		Secondary:    lipgloss.Color("#BB86FC"),
		SecondaryDim: lipgloss.Color("#7C4DFF"),
		SignalUp:     lipgloss.Color("#FF9800"),
		SignalDown:   lipgloss.Color("#00E676"),
		StatusOK:     lipgloss.Color("#00E676"),
		StatusWarn:   lipgloss.Color("#FFD600"),
		StatusError:  lipgloss.Color("#FF5252"),
		TextBright:   lipgloss.Color("#FFFFFF"),
		TextNormal:   lipgloss.Color("#E0E0E0"),
		TextDim:      lipgloss.Color("#9E9E9E"),
		TextMuted:    lipgloss.Color("#616161"),
		Border:       lipgloss.Color("#4A148C"),
		BorderActive: lipgloss.Color("#7C4DFF"),
		BorderModal:  lipgloss.Color("#BB86FC"),
	},
	{
		Name:         "phosphor",
		BgDeep:       lipgloss.Color("#050505"),
		BgPanel:      lipgloss.Color("#0a0f0a"),
		BgHighlight:  lipgloss.Color("#142814"),
		Primary:      lipgloss.Color("#33ff33"),
		PrimaryDim:   lipgloss.Color("#22aa22"),
		Secondary:    lipgloss.Color("#88ff88"),
		SecondaryDim: lipgloss.Color("#44bb44"),
		SignalUp:     lipgloss.Color("#ffaa00"),
		SignalDown:   lipgloss.Color("#33ff33"),
		StatusOK:     lipgloss.Color("#33ff33"),
		StatusWarn:   lipgloss.Color("#ffcc00"),
		StatusError:  lipgloss.Color("#ff3333"),
		TextBright:   lipgloss.Color("#ccffcc"),
		TextNormal:   lipgloss.Color("#88cc88"),
		TextDim:      lipgloss.Color("#558855"),
		TextMuted:    lipgloss.Color("#334433"),
		Border:       lipgloss.Color("#1a441a"),
		BorderActive: lipgloss.Color("#33ff33"),
		BorderModal:  lipgloss.Color("#88ff88"),
	},
	{
		Name:         "solar",
		BgDeep:       lipgloss.Color("#0f0800"),
		BgPanel:      lipgloss.Color("#1a1000"),
		BgHighlight:  lipgloss.Color("#302000"),
		Primary:      lipgloss.Color("#ffaa00"),
		PrimaryDim:   lipgloss.Color("#cc8800"),
		Secondary:    lipgloss.Color("#ffd700"),
		SecondaryDim: lipgloss.Color("#b8960f"),
		SignalUp:     lipgloss.Color("#ff6600"),
		SignalDown:   lipgloss.Color("#00cc88"),
		StatusOK:     lipgloss.Color("#00cc88"),
		StatusWarn:   lipgloss.Color("#ffdd00"),
		StatusError:  lipgloss.Color("#ff4444"),
		TextBright:   lipgloss.Color("#fff0d0"),
		TextNormal:   lipgloss.Color("#ccbb99"),
		TextDim:      lipgloss.Color("#887755"),
		TextMuted:    lipgloss.Color("#554422"),
		Border:       lipgloss.Color("#553300"),
		BorderActive: lipgloss.Color("#ffaa00"),
		BorderModal:  lipgloss.Color("#ffd700"),
	},
	{
		Name:         "nord",
		BgDeep:       lipgloss.Color("#2E3440"),
		BgPanel:      lipgloss.Color("#3B4252"),
		BgHighlight:  lipgloss.Color("#434C5E"),
		Primary:      lipgloss.Color("#88C0D0"),
		PrimaryDim:   lipgloss.Color("#5E81AC"),
		Secondary:    lipgloss.Color("#81A1C1"),
		SecondaryDim: lipgloss.Color("#4C566A"),
		SignalUp:     lipgloss.Color("#D08770"),
		SignalDown:   lipgloss.Color("#A3BE8C"),
		StatusOK:     lipgloss.Color("#A3BE8C"),
		StatusWarn:   lipgloss.Color("#EBCB8B"),
		StatusError:  lipgloss.Color("#BF616A"),
		TextBright:   lipgloss.Color("#ECEFF4"),
		TextNormal:   lipgloss.Color("#D8DEE9"),
		TextDim:      lipgloss.Color("#81A1C1"),
		TextMuted:    lipgloss.Color("#4C566A"),
		Border:       lipgloss.Color("#4C566A"),
		BorderActive: lipgloss.Color("#88C0D0"),
		BorderModal:  lipgloss.Color("#81A1C1"),
	},
	{
		Name:         "dracula",
		BgDeep:       lipgloss.Color("#1a1826"),
		BgPanel:      lipgloss.Color("#282a36"),
		BgHighlight:  lipgloss.Color("#44475a"),
		Primary:      lipgloss.Color("#bd93f9"),
		PrimaryDim:   lipgloss.Color("#6272a4"),
		Secondary:    lipgloss.Color("#ff79c6"),
		SecondaryDim: lipgloss.Color("#44475a"),
		SignalUp:     lipgloss.Color("#ffb86c"),
		SignalDown:   lipgloss.Color("#50fa7b"),
		StatusOK:     lipgloss.Color("#50fa7b"),
		StatusWarn:   lipgloss.Color("#f1fa8c"),
		StatusError:  lipgloss.Color("#ff5555"),
		TextBright:   lipgloss.Color("#f8f8f2"),
		TextNormal:   lipgloss.Color("#e4e4e4"),
		TextDim:      lipgloss.Color("#bd93f9"),
		TextMuted:    lipgloss.Color("#6272a4"),
		Border:       lipgloss.Color("#44475a"),
		BorderActive: lipgloss.Color("#bd93f9"),
		BorderModal:  lipgloss.Color("#ff79c6"),
	},
	{
		Name:         "amber",
		BgDeep:       lipgloss.Color("#0a0800"),
		BgPanel:      lipgloss.Color("#130f00"),
		BgHighlight:  lipgloss.Color("#1f1800"),
		Primary:      lipgloss.Color("#ffb300"),
		PrimaryDim:   lipgloss.Color("#996600"),
		Secondary:    lipgloss.Color("#e67e00"),
		SecondaryDim: lipgloss.Color("#663300"),
		SignalUp:     lipgloss.Color("#ff6b00"),
		SignalDown:   lipgloss.Color("#ffd600"),
		StatusOK:     lipgloss.Color("#b8860b"),
		StatusWarn:   lipgloss.Color("#ff8c00"),
		StatusError:  lipgloss.Color("#cc0000"),
		TextBright:   lipgloss.Color("#ffcc00"),
		TextNormal:   lipgloss.Color("#cc9900"),
		TextDim:      lipgloss.Color("#996600"),
		TextMuted:    lipgloss.Color("#664400"),
		Border:       lipgloss.Color("#3d2b00"),
		BorderActive: lipgloss.Color("#ffb300"),
		BorderModal:  lipgloss.Color("#e67e00"),
	},
}

var currentThemeIdx int

// Color palette (updated by ApplyTheme)
var (
	ColorBgDeep      lipgloss.TerminalColor
	ColorBgPanel     lipgloss.TerminalColor
	ColorBgHighlight lipgloss.TerminalColor

	ColorPrimary    lipgloss.TerminalColor
	ColorPrimaryDim lipgloss.TerminalColor

	ColorSecondary    lipgloss.TerminalColor
	ColorSecondaryDim lipgloss.TerminalColor

	ColorSignalUp   lipgloss.TerminalColor
	ColorSignalDown lipgloss.TerminalColor

	ColorStatusOK    lipgloss.TerminalColor
	ColorStatusWarn  lipgloss.TerminalColor
	ColorStatusError lipgloss.TerminalColor

	ColorTextBright lipgloss.TerminalColor
	ColorTextNormal lipgloss.TerminalColor
	ColorTextDim    lipgloss.TerminalColor
	ColorTextMuted  lipgloss.TerminalColor

	ColorBorder       lipgloss.TerminalColor
	ColorBorderActive lipgloss.TerminalColor
	ColorBorderModal  lipgloss.TerminalColor
)

// Panel styles (updated by ApplyTheme)
var (
	PanelStyle       lipgloss.Style
	ActivePanelStyle lipgloss.Style
	ModalStyle       lipgloss.Style
)

// Text styles (updated by ApplyTheme)
var (
	TitleStyle       lipgloss.Style
	ValueStyle       lipgloss.Style
	LabelStyle       lipgloss.Style
	DimStyle         lipgloss.Style
	MutedStyle       lipgloss.Style
	SignalUpStyle    lipgloss.Style
	SignalDownStyle  lipgloss.Style
	ErrorStyle       lipgloss.Style
	AccentStyle      lipgloss.Style
	PrimaryStyle     lipgloss.Style
	PrimaryBoldStyle lipgloss.Style
)

// Status indicators (updated by ApplyTheme)
var (
	ConnectedDot    string
	DegradedDot     string
	DisconnectedDot string
)

// Modal title style (updated by ApplyTheme)
var ModalTitleStyle lipgloss.Style

func init() {
	ApplyTheme(themes[0])
}

// ApplyTheme updates all global color and style variables.
func ApplyTheme(t Theme) {
	ColorBgDeep = t.BgDeep
	ColorBgPanel = t.BgPanel
	ColorBgHighlight = t.BgHighlight
	ColorPrimary = t.Primary
	ColorPrimaryDim = t.PrimaryDim
	ColorSecondary = t.Secondary
	ColorSecondaryDim = t.SecondaryDim
	ColorSignalUp = t.SignalUp
	ColorSignalDown = t.SignalDown
	ColorStatusOK = t.StatusOK
	ColorStatusWarn = t.StatusWarn
	ColorStatusError = t.StatusError
	ColorTextBright = t.TextBright
	ColorTextNormal = t.TextNormal
	ColorTextDim = t.TextDim
	ColorTextMuted = t.TextMuted
	ColorBorder = t.Border
	ColorBorderActive = t.BorderActive
	ColorBorderModal = t.BorderModal

	PanelStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorder).
		Padding(0, 1)
	ActivePanelStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorderActive).
		Padding(0, 1)
	ModalStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.DoubleBorder()).
		BorderForeground(ColorBorderModal).
		Padding(1, 2)

	TitleStyle = lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true)
	ValueStyle = lipgloss.NewStyle().Foreground(ColorSecondary)
	LabelStyle = lipgloss.NewStyle().Foreground(ColorTextNormal)
	DimStyle = lipgloss.NewStyle().Foreground(ColorTextDim)
	MutedStyle = lipgloss.NewStyle().Foreground(ColorTextMuted)
	SignalUpStyle = lipgloss.NewStyle().Foreground(ColorSignalUp).Bold(true)
	SignalDownStyle = lipgloss.NewStyle().Foreground(ColorSignalDown).Bold(true)
	ErrorStyle = lipgloss.NewStyle().Foreground(ColorStatusError)
	AccentStyle = lipgloss.NewStyle().Foreground(ColorSecondary)
	PrimaryStyle = lipgloss.NewStyle().Foreground(ColorPrimary)
	PrimaryBoldStyle = lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true)

	ConnectedDot = lipgloss.NewStyle().Foreground(ColorStatusOK).Render("●")
	DegradedDot = lipgloss.NewStyle().Foreground(ColorStatusWarn).Render("●")
	DisconnectedDot = lipgloss.NewStyle().Foreground(ColorStatusError).Render("●")

	ModalTitleStyle = lipgloss.NewStyle().Foreground(ColorSecondary).Bold(true)
}

// CycleTheme advances to the next theme and returns its name.
func CycleTheme() string {
	currentThemeIdx = (currentThemeIdx + 1) % len(themes)
	ApplyTheme(themes[currentThemeIdx])
	return themes[currentThemeIdx].Name
}

// SetThemeByName applies a theme by name. Does nothing if name is not found.
func SetThemeByName(name string) {
	for i, t := range themes {
		if t.Name == name {
			currentThemeIdx = i
			ApplyTheme(t)
			return
		}
	}
}

// CurrentThemeName returns the name of the active theme.
func CurrentThemeName() string {
	return themes[currentThemeIdx].Name
}
