package model

import "github.com/RustyDaemon/go-dsn-now/model/response"

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
	}
}

//func (data *AppData) UpdateDSNConfig(config response.DSNConfig) {
//	data.dsnConfig = config
//}
//
//func (data *AppData) GetDSNConfig() response.DSNConfig {
//	return data.dsnConfig
//}

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

func (data *AppData) HasAntennaSpecs() bool {
	dish, ok := data.GetSelectedDish()
	if !ok {
		return false
	}

	return dish.Specs != (DishSpecification{})
}
