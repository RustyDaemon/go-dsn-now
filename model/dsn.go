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

// todo seems like we can remap it from config dsn section
// func (dsn *DSN) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
// 	var currentStation *DSNStation

// 	for {
// 		token, err := d.Token()
// 		if err != nil {
// 			break
// 		}

// 		switch elem := token.(type) {
// 		case xml.StartElement:
// 			if elem.Name.Local == "station" {
// 				var station DSNStation

// 				if err := d.DecodeElement(&station, &elem); err != nil {
// 					return err
// 				}

// 				dsn.Stations = append(dsn.Stations, station)
// 				currentStation = &dsn.Stations[len(dsn.Stations)-1]
// 			} else if elem.Name.Local == "dish" {
// 				var dish DSNDish

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
