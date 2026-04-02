package components

import (
	"strings"

	"github.com/RustyDaemon/go-dsn-now/internal/model"
	"github.com/RustyDaemon/go-dsn-now/internal/tui/style"
)

type StationBar struct {
	width int
}

func NewStationBar() StationBar {
	return StationBar{}
}

func (s *StationBar) SetWidth(w int) {
	s.width = w
}

func (s StationBar) View(stations []model.Station, selectedIdx int) string {
	if len(stations) == 0 {
		return ""
	}

	var lines []string
	for i, station := range stations {
		if i == selectedIdx {
			lines = append(lines, style.PrimaryBoldStyle.Render("▸ "+station.FriendlyName)+" "+station.Flag)
		} else {
			lines = append(lines, style.DimStyle.Render("  "+station.FriendlyName)+" "+station.Flag)
		}
	}

	content := strings.Join(lines, "\n")
	title := style.TitleStyle.Render("Stations")
	return style.RenderTitledPanel(title, content, s.width, style.ColorBorder)
}
