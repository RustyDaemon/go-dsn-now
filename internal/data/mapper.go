package data

import (
	"github.com/RustyDaemon/go-dsn-now/internal/model"
	"github.com/RustyDaemon/go-dsn-now/internal/model/response"
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
			Dishes:       convertDishes(s.Dishes),
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
			dn.DownSignals = convertDownSignals(d.DownSignals, dest)
			dn.UpSignals = convertUpSignals(d.UpSignals, dest)
			dn.Targets = convertTargets(d.Target, dest)
			dn.Specs = dn.GetDishSpecification()
		}
	}
}

func convertDownSignals(signals []response.DishDownSignal, dest *model.FullData) []model.DownSignal {
	var result []model.DownSignal
	for _, s := range signals {
		spacecraft := dest.GetSpacecraftByName(s.Spacecraft)
		isActive, err := strconv.ParseBool(s.IsActive)
		if err != nil {
			isActive = false
		}

		result = append(result, model.DownSignal{
			IsActive:     isActive,
			SignalType:   s.SignalType,
			DataRate:     s.DataRate,
			Frequency:    s.Frequency,
			Band:         s.Band,
			Power:        s.Power,
			SpacecraftID: s.SpacecraftID,
			Spacecraft:   spacecraft,
		})
	}
	return result
}

func convertUpSignals(signals []response.DishUpSignal, dest *model.FullData) []model.UpSignal {
	var result []model.UpSignal
	for _, s := range signals {
		spacecraft := dest.GetSpacecraftByName(s.Spacecraft)
		isActive, err := strconv.ParseBool(s.IsActive)
		if err != nil {
			isActive = false
		}

		result = append(result, model.UpSignal{
			IsActive:     isActive,
			SignalType:   s.SignalType,
			DataRate:     s.DataRate,
			Frequency:    s.Frequency,
			Band:         s.Band,
			Power:        s.Power,
			SpacecraftID: s.SpacecraftID,
			Spacecraft:   spacecraft,
		})
	}
	return result
}

func convertTargets(targets []response.DishTarget, dest *model.FullData) []model.Target {
	var result []model.Target
	for _, t := range targets {
		spacecraft := dest.GetSpacecraftByName(t.Name)

		result = append(result, model.Target{
			Name:         t.Name,
			ID:           t.ID,
			UplegRange:   t.UplegRange,
			DownlegRange: t.DownlegRange,
			Rtlt:         t.Rtlt,
			Spacecraft:   spacecraft,
		})
	}
	return result
}

func convertDishes(dishes []response.DSNConfigSiteDish) []model.Dish {
	var result []model.Dish
	for _, dish := range dishes {
		result = append(result, model.Dish{
			Name:         dish.Name,
			FriendlyName: dish.FriendlyName,
			Type:         dish.Type,
		})
	}
	return result
}
