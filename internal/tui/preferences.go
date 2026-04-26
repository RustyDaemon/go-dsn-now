package tui

import (
	"time"

	"github.com/RustyDaemon/go-dsn-now/internal/data"
)

type preferences struct {
	settings *data.Settings
}

func newPreferences(settings *data.Settings) *preferences {
	if settings == nil {
		settings = &data.Settings{}
	}

	return &preferences{settings: settings}
}

func (p *preferences) saveSettings() {
	_ = data.SaveSettings(p.settings)
}

func (p *preferences) setTheme(name string) {
	p.settings.Theme = name
	p.saveSettings()
}

func (p *preferences) distanceUnit() string {
	return p.settings.DistanceUnit
}

func (p *preferences) cycleDistanceUnit() string {
	units := []string{"km", "au", "lmin", "lhour"}
	current := p.settings.DistanceUnit
	next := "au"

	for i, unit := range units {
		if unit == current {
			next = units[(i+1)%len(units)]
			break
		}
	}

	p.settings.DistanceUnit = next
	p.saveSettings()
	return next
}

func (p *preferences) setRefreshInterval(interval time.Duration) {
	p.settings.RefreshIntervalSeconds = int(interval.Seconds())
	p.saveSettings()
}

func (p *preferences) setLastStation(name string) {
	p.settings.LastStation = name
	p.saveSettings()
}

func (p *preferences) lastStation() string {
	return p.settings.LastStation
}

func (p *preferences) toggleBookmark(bookmarks map[string]bool, dishName string) {
	if bookmarks[dishName] {
		delete(bookmarks, dishName)
	} else {
		bookmarks[dishName] = true
	}

	_ = data.SaveBookmarks(bookmarks)
}
