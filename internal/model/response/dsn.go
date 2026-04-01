package response

import "strings"

type DSN struct {
	Timestamp string    `xml:"timestamp"`
	Stations  []Station `xml:"station"`
	D         []Dish    `xml:"dish"`
}

type Station struct {
	Name           string `xml:"name,attr"`
	FriendlyName   string `xml:"friendlyName,attr"`
	TimeUTC        string `xml:"timeUTC,attr"`
	TimeZoneOffset string `xml:"timeZoneOffset,attr"`
}

func (s *Station) GetStationFlag() string {
	name := s.FriendlyName

	if strings.EqualFold(name, "madrid") {
		return "\U0001F1EA\U0001F1F8"
	}

	if strings.EqualFold(name, "goldstone") {
		return "\U0001F1FA\U0001F1F8"
	}

	if strings.EqualFold(name, "canberra") {
		return "\U0001F1E6\U0001F1FA"
	}

	return ""
}

type Dish struct {
	Name           string           `xml:"name,attr"`
	AzimuthAngle   string           `xml:"azimuthAngle,attr"`
	ElevationAngle string           `xml:"elevationAngle,attr"`
	WindSpeed      string           `xml:"windSpeed,attr"`
	IsMSPA         string           `xml:"isMSPA,attr"`
	IsArray        string           `xml:"isArray,attr"`
	IsDDOR         string           `xml:"isDDOR,attr"`
	Activity       string           `xml:"activity,attr"`
	UpSignals      []DishUpSignal   `xml:"upSignal"`
	DownSignals    []DishDownSignal `xml:"downSignal"`
	Target         []DishTarget     `xml:"target"`
}

type DishUpSignal struct {
	IsActive     string `xml:"active,attr"`
	SignalType   string `xml:"signalType,attr"`
	DataRate     string `xml:"dataRate,attr"`
	Frequency    string `xml:"frequency,attr"`
	Band         string `xml:"band,attr"`
	Power        string `xml:"power,attr"`
	Spacecraft   string `xml:"spacecraft,attr"`
	SpacecraftID string `xml:"spacecraftID,attr"`
}

type DishDownSignal struct {
	IsActive     string `xml:"active,attr"`
	SignalType   string `xml:"signalType,attr"`
	DataRate     string `xml:"dataRate,attr"`
	Frequency    string `xml:"frequency,attr"`
	Band         string `xml:"band,attr"`
	Power        string `xml:"power,attr"`
	Spacecraft   string `xml:"spacecraft,attr"`
	SpacecraftID string `xml:"spacecraftID,attr"`
}

type DishTarget struct {
	Name         string `xml:"name,attr"`
	ID           string `xml:"id,attr"`
	UplegRange   string `xml:"uplegRange,attr"`
	DownlegRange string `xml:"downlegRange,attr"`
	Rtlt         string `xml:"rtlt,attr"`
}
