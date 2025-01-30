package converter

import (
	"github.com/RustyDaemon/go-dsn-now/model"
	"github.com/RustyDaemon/go-dsn-now/model/response"
	"strconv"
)

func ConvertDownSignals(signals []response.DishDownSignal, dest *model.FullData) []model.DownSignal {
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

func ConvertUpSignals(signals []response.DishUpSignal, dest *model.FullData) []model.UpSignal {
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

func ConvertTargets(targets []response.DishTarget, dest *model.FullData) []model.Target {
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

func ConvertDishes(dishes []response.DSNConfigSiteDish) []model.Dish {
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
