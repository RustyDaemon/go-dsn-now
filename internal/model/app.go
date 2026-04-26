package model

import "time"

type AppData struct {
	FullData          FullData
	LastError         string
	LastUpdated       time.Time
	ConsecutiveErrors int
	Bookmarks         map[string]bool
	PrevSignalCounts  map[string]signalCount
	SignalChanges     []string
	DishActiveSince   map[string]time.Time
	SignalHistory     map[string][]float64
}

type signalCount struct {
	Up   int
	Down int
}

func NewAppData() *AppData {
	return &AppData{
		FullData:         FullData{},
		PrevSignalCounts: make(map[string]signalCount),
		Bookmarks:        make(map[string]bool),
		DishActiveSince:  make(map[string]time.Time),
		SignalHistory:    make(map[string][]float64),
	}
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

			if up+down > 0 {
				if _, tracked := data.DishActiveSince[dish.Name]; !tracked {
					data.DishActiveSince[dish.Name] = time.Now()
				}
			} else {
				delete(data.DishActiveSince, dish.Name)
			}
		}
	}
}
