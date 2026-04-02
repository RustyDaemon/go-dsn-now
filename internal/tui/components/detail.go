package components

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/RustyDaemon/go-dsn-now/internal/model"
	"github.com/RustyDaemon/go-dsn-now/internal/tui/style"
)

func RenderAntennaInfo(dish model.Dish, width int) string {
	nameVal := style.DashIfEmpty(dish.FriendlyName)
	typeVal := style.DashIfEmpty(dish.Type)
	activityVal := style.DashIfEmpty(dish.Activity)
	azimuthVal := style.FormatAngle(dish.AzimuthAngle)
	elevationVal := style.FormatAngle(dish.ElevationAngle)
	windVal := style.FormatWind(dish.WindSpeed)

	innerWidth := width - 4
	if innerWidth < 20 {
		innerWidth = 20
	}
	halfWidth := innerWidth / 2
	thirdWidth := innerWidth / 3

	row1Left := style.LabelStyle.Render("Name: ") + style.ValueStyle.Render(nameVal)
	row1Right := style.LabelStyle.Render("Type: ") + style.ValueStyle.Render(typeVal)

	row1 := lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.NewStyle().Width(halfWidth).Render(row1Left),
		lipgloss.NewStyle().Width(halfWidth).Align(lipgloss.Right).Render(row1Right),
	)

	row2 := lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.NewStyle().Width(thirdWidth).Render(style.LabelStyle.Render("Azimuth: ")+style.ValueStyle.Render(azimuthVal)),
		lipgloss.NewStyle().Width(thirdWidth).Align(lipgloss.Center).Render(style.LabelStyle.Render("Elevation: ")+style.ValueStyle.Render(elevationVal)),
		lipgloss.NewStyle().Width(thirdWidth).Align(lipgloss.Right).Render(style.LabelStyle.Render("Wind: ")+style.ValueStyle.Render(windVal)),
	)

	row3 := style.LabelStyle.Render("Activity: ") + style.ValueStyle.Render(activityVal)

	content := lipgloss.JoinVertical(lipgloss.Left, row1, row2, row3)

	title := style.TitleStyle.Render("Antenna Information")
	return style.RenderTitledPanel(title, content, width, style.ColorBorder)
}
