package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/RustyDaemon/go-dsn-now/internal/model"
	"github.com/RustyDaemon/go-dsn-now/internal/tui/style"
)

func RenderTarget(targets []model.Target, selectedIdx int, width int, distanceUnit string) string {
	if len(targets) == 0 {
		title := style.TitleStyle.Render("Target:") + " " + style.DimStyle.Render("No target")
		return style.RenderTitledPanel(title, style.DimStyle.Render("No target data"), width, style.ColorTextMuted)
	}

	if selectedIdx >= len(targets) {
		selectedIdx = 0
	}

	title := style.TitleStyle.Render("Target:") + " " + buildIndexIndicators(len(targets), selectedIdx)

	target := targets[selectedIdx]

	name := "-"
	if target.Spacecraft != (model.Spacecraft{}) && len(target.Spacecraft.FriendlyName) > 0 {
		name = target.Spacecraft.FriendlyName
	}

	rangeVal := style.FormatRangeInUnit(target.UplegRange, distanceUnit)
	rtlt := style.FormatRTLT(target.Rtlt)

	content := lipgloss.JoinVertical(lipgloss.Left,
		style.LabelStyle.Render("Spacecraft: ")+style.ValueStyle.Render(name),
		style.LabelStyle.Render("Range: ")+style.ValueStyle.Render(rangeVal),
		style.LabelStyle.Render("Round-trip light time: ")+style.ValueStyle.Render(rtlt),
	)

	isEmpty := target == (model.Target{}) || (target.UplegRange == "-1" && target.Rtlt == "-1" && target.Spacecraft == (model.Spacecraft{}))
	borderColor := style.ColorBorder
	if isEmpty {
		borderColor = style.ColorTextMuted
	}

	return style.RenderTitledPanel(title, content, width, borderColor)
}

func buildIndexIndicators(count, selectedIdx int) string {
	var parts []string
	for i := 0; i < count; i++ {
		label := fmt.Sprintf("[%d]", i+1)
		if i == selectedIdx {
			parts = append(parts, style.PrimaryBoldStyle.Render(label))
		} else {
			parts = append(parts, style.DimStyle.Render(label))
		}
	}
	return strings.Join(parts, "")
}
