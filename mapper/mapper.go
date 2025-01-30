package mapper

import (
	"github.com/RustyDaemon/go-dsn-now/converter"
	"github.com/RustyDaemon/go-dsn-now/model"
	"github.com/RustyDaemon/go-dsn-now/model/response"
	"strconv"
)

func MapConfigToFullData(c response.DSNConfig, dest *model.FullData) {
	for _, s := range c.Spacecrafts {
		dest.Spacecrafts = append(dest.Spacecrafts, model.Spacecraft{
			Name:         s.Name,
			FriendlyName: s.FriendlyName,
			ExplorerName: s.ExplorerName,
		})
	}

	for _, s := range c.Sites {
		dest.Stations = append(dest.Stations, model.Station{
			Name:         s.Name,
			FriendlyName: s.FriendlyName,
			Longitude:    s.Longitude,
			Latitude:     s.Latitude,
			Dishes:       converter.ConvertDishes(s.Dishes),
		})
	}
}

func MapDataToFullData(d response.DSN, dest *model.FullData) {
	dest.Timestamp = d.Timestamp

	for _, s := range d.Stations {
		sn := dest.GetStationByName(s.Name)
		if sn.Name != "" {
			sn.TimeUTC = s.TimeUTC
			sn.TimeZoneOffset = s.TimeZoneOffset
			sn.Flag = s.GetStationFlag()
		}
	}

	for _, d := range d.D {
		isMSPA, err := strconv.ParseBool(d.IsMSPA)
		if err != nil {
			isMSPA = false
		}
		isArray, err := strconv.ParseBool(d.IsArray)
		if err != nil {
			isArray = false
		}
		isDDOR, err := strconv.ParseBool(d.IsDDOR)
		if err != nil {
			isDDOR = false
		}

		dn := dest.GetDishByName(d.Name)
		if dn.Name != "" {
			dn.AzimuthAngle = d.AzimuthAngle
			dn.ElevationAngle = d.ElevationAngle
			dn.WindSpeed = d.WindSpeed
			dn.IsMSPA = isMSPA
			dn.IsArray = isArray
			dn.IsDDOR = isDDOR
			dn.Activity = d.Activity
			dn.DownSignals = converter.ConvertDownSignals(d.DownSignals, dest)
			dn.UpSignals = converter.ConvertUpSignals(d.UpSignals, dest)
			dn.Targets = converter.ConvertTargets(d.Target, dest)
			dn.Specs = dn.GetDishSpecification()
		}
	}
}
