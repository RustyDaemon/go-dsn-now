package model

type DSN struct {
	Timestamp string       `xml:"timestamp"`
	Stations  []DSNStation `xml:"station"`
	D         []DSNDish    `xml:"dish"`
}

type DSNStation struct {
	Name           string    `xml:"name,attr"`
	FriendlyName   string    `xml:"friendlyName,attr"`
	TimeUTC        string    `xml:"timeUTC,attr"`
	TimeZoneOffset string    `xml:"timeZoneOffset,attr"`
	Dishes         []DSNDish `xml:"-"`
}

type DSNDish struct {
	Name           string              `xml:"name,attr"`
	AzimuthAngle   string              `xml:"azimuthAngle,attr"`
	ElevationAngle string              `xml:"elevationAngle,attr"`
	WindSpeed      string              `xml:"windSpeed,attr"`
	IsMSPA         string              `xml:"isMSPA,attr"`
	IsArray        string              `xml:"isArray,attr"`
	IsDDOR         string              `xml:"isDDOR,attr"`
	Activity       string              `xml:"activity,attr"`
	UpSignal       []DSNDishUpSignal   `xml:"upSignal"`
	DownSignals    []DSNDishDownSignal `xml:"downSignal"`
	Target         []DSNDishTarget     `xml:"target"`
}

type DSNDishUpSignal struct {
	IsActive     string `xml:"active,attr"`
	SignalType   string `xml:"signalType,attr"`
	DataRate     string `xml:"dataRate,attr"`
	Frequency    string `xml:"frequency,attr"`
	Band         string `xml:"band,attr"`
	Power        string `xml:"power,attr"`
	Spacecraft   string `xml:"spacecraft,attr"`
	SpacecraftID string `xml:"spacecraftID,attr"`
}

type DSNDishDownSignal struct {
	IsActive     string `xml:"active,attr"`
	SignalType   string `xml:"signalType,attr"`
	DataRate     string `xml:"dataRate,attr"`
	Frequency    string `xml:"frequency,attr"`
	Band         string `xml:"band,attr"`
	Power        string `xml:"power,attr"`
	Spacecraft   string `xml:"spacecraft,attr"`
	SpacecraftID string `xml:"spacecraftID,attr"`
}

type DSNDishTarget struct {
	Name         string `xml:"name,attr"`
	ID           string `xml:"id,attr"`
	UplegRange   string `xml:"uplegRange,attr"`
	DownlegRange string `xml:"downlegRange,attr"`
	Rtlt         string `xml:"rtlt,attr"`
}
