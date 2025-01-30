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

// GetStationFlag temporary disabled - flag emojis as they shift the text 🇪🇸 🇺🇸 🇦🇺
func (s *Station) GetStationFlag() string {
	name := s.FriendlyName

	if strings.EqualFold(name, "madrid") {
		return "[ESP[]"
	}

	if strings.EqualFold(name, "goldstone") {
		return "[USA[]"
	}

	if strings.EqualFold(name, "canberra") {
		return "[AUS[]"
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

// todo seems like we can remap it from config dsn section
// func (dsn *DSN) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
// 	var currentStation *Station

// 	for {
// 		token, err := d.Token()
// 		if err != nil {
// 			break
// 		}

// 		switch elem := token.(type) {
// 		case xml.StartElement:
// 			if elem.Name.Local == "station" {
// 				var station Station

// 				if err := d.DecodeElement(&station, &elem); err != nil {
// 					return err
// 				}

// 				dsn.Stations = append(dsn.Stations, station)
// 				currentStation = &dsn.Stations[len(dsn.Stations)-1]
// 			} else if elem.Name.Local == "dish" {
// 				var dish Dish

// 				if err := d.DecodeElement(&dish, &elem); err != nil {
// 					return err
// 				}

// 				if currentStation != nil {
// 					currentStation.Dishes = append(currentStation.Dishes, dish)
// 				}
// 			}
// 		}
// 	}

// 	return nil
// }
