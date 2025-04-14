package model

import "strings"

// todo after the MVP is done, refactor the models: types, fields, etc.

type FullData struct {
	Stations    []Station
	Spacecrafts []Spacecraft
	Timestamp   string
}

func (fd *FullData) GetSpacecraftByName(name string) Spacecraft {
	for i := range fd.Spacecrafts {
		if strings.EqualFold(fd.Spacecrafts[i].Name, name) {
			return fd.Spacecrafts[i]
		}
	}

	return Spacecraft{}
}

func (fd *FullData) GetStationByName(name string) *Station {
	for i := range fd.Stations {
		if strings.EqualFold(fd.Stations[i].Name, name) {
			return &fd.Stations[i]
		}
	}

	return &Station{}
}

func (fd *FullData) GetDishByName(name string) *Dish {
	for i := range fd.Stations {
		for j := range fd.Stations[i].Dishes {
			if strings.EqualFold(fd.Stations[i].Dishes[j].Name, name) {
				return &fd.Stations[i].Dishes[j]
			}
		}
	}

	return &Dish{}
}

type Station struct {
	Name           string
	FriendlyName   string
	Longitude      string
	Latitude       string
	TimeUTC        string
	TimeZoneOffset string
	Flag           string
	Dishes         []Dish
}

type Dish struct {
	Name           string
	FriendlyName   string
	Type           string
	AzimuthAngle   string
	ElevationAngle string
	WindSpeed      string
	IsMSPA         bool
	IsArray        bool
	IsDDOR         bool
	Activity       string
	UpSignals      []UpSignal
	DownSignals    []DownSignal
	Targets        []Target
	Specs          DishSpecification
}

func (d *Dish) CountWorkingUpSignals() int {
	count := 0
	if len(d.UpSignals) == 0 {
		return count
	}

	for _, s := range d.UpSignals {
		if s.IsActive {
			count++
		}
	}

	return count
}

func (d *Dish) CountWorkingDownSignals() int {
	count := 0
	if len(d.DownSignals) == 0 {
		return count
	}

	for _, s := range d.DownSignals {
		if s.IsActive {
			count++
		}
	}

	return count
}

func (d *Dish) CountWorkingTargets() int {
	count := 0
	if len(d.Targets) == 0 {
		return count
	}

	for _, t := range d.Targets {
		if t.Spacecraft.Name != "" && t.UplegRange != "-1" && t.Rtlt != "-1" {
			count++
		}
	}

	return count
}

type signal struct {
	IsActive     bool
	SignalType   string
	DataRate     string
	Frequency    string
	Band         string
	Power        string
	SpacecraftID string
	Spacecraft   Spacecraft
}

type UpSignal signal
type DownSignal signal

type Target struct {
	Name         string
	ID           string
	UplegRange   string
	DownlegRange string
	Rtlt         string
	Spacecraft   Spacecraft
}

type Spacecraft struct {
	Name         string
	FriendlyName string
	ExplorerName string
}
