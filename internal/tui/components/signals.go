package components

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/RustyDaemon/go-dsn-now/internal/model"
	"github.com/RustyDaemon/go-dsn-now/internal/tui/style"
)

func RenderUpSignal(signals []model.UpSignal, selectedIdx int, width int) string {
	arrow := style.SignalUpStyle.Render("↑")

	if len(signals) == 0 {
		title := arrow + " " + style.TitleStyle.Render("Up Signal:") + " " + style.DimStyle.Render("No signal")
		return style.RenderTitledPanel(title, style.DimStyle.Render("No signal"), width, style.ColorTextMuted)
	}

	if selectedIdx >= len(signals) {
		selectedIdx = 0
	}

	title := arrow + " " + style.TitleStyle.Render("Up Signal:") + " " + buildIndexIndicators(len(signals), selectedIdx)

	sig := signals[selectedIdx]

	source := "-"
	if sig.Spacecraft != (model.Spacecraft{}) && len(sig.Spacecraft.FriendlyName) > 0 {
		source = sig.Spacecraft.FriendlyName
	}

	activeStr := style.ErrorStyle.Render("inactive")
	if sig.IsActive {
		activeStr = style.PrimaryStyle.Render("active")
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		style.LabelStyle.Render("Source: ")+style.ValueStyle.Render(source),
		style.LabelStyle.Render("Is active: ")+activeStr,
		style.LabelStyle.Render("Signal type: ")+style.ValueStyle.Render(style.DefaultIfEmpty(sig.SignalType, "-")),
		style.LabelStyle.Render("Frequency band: ")+style.ValueStyle.Render(style.DefaultIfEmpty(sig.Band, "-")),
		style.LabelStyle.Render("Power transmitted: ")+style.ValueStyle.Render(style.FormatPowerTx(sig.Power)),
	)

	borderColor := style.ColorBorder
	if sig == (model.UpSignal{}) {
		borderColor = style.ColorTextMuted
	}

	return style.RenderTitledPanel(title, content, width, borderColor)
}

func RenderDownSignal(signals []model.DownSignal, selectedIdx int, width int) string {
	arrow := style.SignalDownStyle.Render("↓")

	if len(signals) == 0 {
		title := arrow + " " + style.TitleStyle.Render("Down Signal:") + " " + style.DimStyle.Render("No signal")
		return style.RenderTitledPanel(title, style.DimStyle.Render("No signal"), width, style.ColorTextMuted)
	}

	if selectedIdx >= len(signals) {
		selectedIdx = 0
	}

	title := arrow + " " + style.TitleStyle.Render("Down Signal:") + " " + buildIndexIndicators(len(signals), selectedIdx)

	sig := signals[selectedIdx]

	source := "-"
	if sig.Spacecraft != (model.Spacecraft{}) && len(sig.Spacecraft.FriendlyName) > 0 {
		source = sig.Spacecraft.FriendlyName
	}

	activeStr := style.ErrorStyle.Render("inactive")
	if sig.IsActive {
		activeStr = style.PrimaryStyle.Render("active")
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		style.LabelStyle.Render("Source: ")+style.ValueStyle.Render(source),
		style.LabelStyle.Render("Is active: ")+activeStr,
		style.LabelStyle.Render("Signal type: ")+style.ValueStyle.Render(style.DefaultIfEmpty(sig.SignalType, "-")),
		style.LabelStyle.Render("Frequency band: ")+style.ValueStyle.Render(style.DefaultIfEmpty(sig.Band, "-")),
		style.LabelStyle.Render("Data rate: ")+style.ValueStyle.Render(style.FormatDataRate(sig.DataRate)),
		style.LabelStyle.Render("Power received: ")+style.ValueStyle.Render(style.FormatPowerRx(sig.Power)),
	)

	borderColor := style.ColorBorder
	if sig == (model.DownSignal{}) {
		borderColor = style.ColorTextMuted
	}

	return style.RenderTitledPanel(title, content, width, borderColor)
}
