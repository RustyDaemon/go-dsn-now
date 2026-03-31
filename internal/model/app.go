package model

import (
	"time"

	"github.com/RustyDaemon/go-dsn-now/internal/model/response"
)

type AppData struct {
	IsReady               bool
	SelectedStationIdx    int
	SelectedDishIdx       int
	SelectedTargetIdx     int
	SelectedUpSignalIdx   int
	SelectedDownSignalIdx int
	FullData              FullData
	DSNConfig             response.DSNConfig
	IsPreviewShown        bool
	IsSpecsShown          bool
	IsAboutShown          bool
	LastError             string
	LastUpdated           time.Time
	ConsecutiveErrors     int
	CompactView           bool
	Bookmarks             map[string]bool // dish name -> bookmarked
	PrevSignalCounts      map[string]signalCount // dish name -> signal counts from last update
	SignalChanges         []string               // recent signal change notifications
}

type signalCount struct {
	Up   int
	Down int
}

func NewAppData() *AppData {
	return &AppData{
		IsReady:               false,
		SelectedStationIdx:    -1,
		SelectedDishIdx:       0,
		SelectedTargetIdx:     0,
		SelectedUpSignalIdx:   0,
		SelectedDownSignalIdx: 0,
		DSNConfig:             response.DSNConfig{},
		FullData:              FullData{},
		IsPreviewShown:        false,
		IsSpecsShown:          false,
		IsAboutShown:          false,
		PrevSignalCounts:      make(map[string]signalCount),
		Bookmarks:             make(map[string]bool),
	}
}

func (data *AppData) IsNoModalShown() bool {
	return !data.IsPreviewShown && !data.IsSpecsShown && !data.IsAboutShown
}

func (data *AppData) GetSelectedDish() (res Dish, ok bool) {
	if data.SelectedStationIdx < 0 || data.SelectedDishIdx < 0 {
		return Dish{}, false
	}

	selectedStation := data.FullData.Stations[data.SelectedStationIdx]
	dish := selectedStation.Dishes[data.SelectedDishIdx]

	return dish, true
}

func (data *AppData) GetDownSignals() (res []DownSignal, ok bool) {
	dish, ok := data.GetSelectedDish()
	if !ok {
		return []DownSignal{}, false
	}

	return dish.DownSignals, true
}

func (data *AppData) GetUpSignals() (res []UpSignal, ok bool) {
	dish, ok := data.GetSelectedDish()
	if !ok {
		return []UpSignal{}, false
	}

	return dish.UpSignals, true
}

func (data *AppData) GetTargets() (res []Target, ok bool) {
	dish, ok := data.GetSelectedDish()
	if !ok {
		return []Target{}, false
	}

	return dish.Targets, true
}

func (data *AppData) DetectSignalChanges() {
	data.SignalChanges = nil
	for i := range data.FullData.Stations {
		for j := range data.FullData.Stations[i].Dishes {
			dish := &data.FullData.Stations[i].Dishes[j]
			up := dish.CountWorkingUpSignals()
			down := dish.CountWorkingDownSignals()

			prev, exists := data.PrevSignalCounts[dish.Name]
			if exists {
				if up > prev.Up {
					data.SignalChanges = append(data.SignalChanges, dish.FriendlyName+" +↑")
				} else if up < prev.Up {
					data.SignalChanges = append(data.SignalChanges, dish.FriendlyName+" -↑")
				}
				if down > prev.Down {
					data.SignalChanges = append(data.SignalChanges, dish.FriendlyName+" +↓")
				} else if down < prev.Down {
					data.SignalChanges = append(data.SignalChanges, dish.FriendlyName+" -↓")
				}
			}
			data.PrevSignalCounts[dish.Name] = signalCount{Up: up, Down: down}
		}
	}
}

func (data *AppData) HasAntennaSpecs() bool {
	dish, ok := data.GetSelectedDish()
	if !ok {
		return false
	}

	return dish.Specs != (DishSpecification{})
}
