package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/RustyDaemon/go-dsn-now/gui"
	"github.com/RustyDaemon/go-dsn-now/mapper"
	"github.com/RustyDaemon/go-dsn-now/model"
	"github.com/RustyDaemon/go-dsn-now/model/response"
	"github.com/RustyDaemon/go-dsn-now/network"
	"github.com/gdamore/tcell/v2"

	"github.com/rivo/tview"
)

var (
	app  *tview.Application
	ui   *gui.UI
	data *model.AppData
)

func main() {
	data = model.NewAppData()

	app = tview.NewApplication()
	ui = gui.NewUI()
	appUI := ui.BuildAppUI(onListItemChanged)

	updateStatusBar(true)

	interrupt := make(chan os.Signal, 1)
	chanConfig := make(chan response.DSNConfig)
	chanDSNData := make(chan response.DSN)
	chanError := make(chan error)

	signal.Notify(interrupt, os.Interrupt)

	go network.LoadDSNConfig(chanConfig, chanError)

	select {
	case data.DSNConfig = <-chanConfig:
	case err := <-chanError:
		log.Fatal(err)
	}

	mapper.MapConfigToFullData(data.DSNConfig, &data.FullData)

	go runDSNDataLoader(chanDSNData, chanError, interrupt)

	go func() {
		for {
			select {
			case dsnData := <-chanDSNData:
				onDataReceived(dsnData)
			case err := <-chanError:
				log.Fatal(err)
			}
		}
	}()

	setKeybindings()

	if err := app.SetRoot(appUI, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func onListItemChanged(index int) {
	updateDishDetails(index)
	updateTargetsData()
	updateUpSignalsTitleData()
	updateDownSignalsTitleData()

	if !data.IsPreviewShown && !data.IsSpecsShown {
		updateStatusBar(true)
	}
}

func setKeybindings() {
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if !data.IsPreviewShown && !data.IsSpecsShown {
			switch event.Rune() {
			case 's':
				updateStationSelection()
			case 't':
				updateTargetSelection()
			case 'u':
				updateUpSignalSelection()
			case 'd':
				updateDownSignalSelection()
			case 'p':
				if !data.IsReady {
					break
				}
				data.IsPreviewShown = true
				updateStatusBar(false)
				showPreview()
			case 'i':
				if !data.IsReady {
					break
				}
				data.IsSpecsShown = true
				updateStatusBar(false)
				showDishSpecs()
			}
		}

		if event.Rune() == 'q' {
			app.Stop()
		}

		if event.Key() == tcell.KeyEscape && data.IsPreviewShown {
			ui.CloseJSONPreviewModal()
			data.IsPreviewShown = false
			updateStatusBar(true)
		} else if event.Key() == tcell.KeyEscape && data.IsSpecsShown {
			ui.CloseDishSpecificationModal()
			data.IsSpecsShown = false
			updateStatusBar(true)
		}

		return event
	})
}

func updateStatusBar(defaultStatus bool) {
	if !data.IsReady {
		return
	}

	params := gui.StatusBarParams{DefaultStatus: defaultStatus}

	if !params.DefaultStatus {
		ui.UpdateStatusBar(params)
		return
	}

	if targets, ok := data.GetTargets(); ok {
		params.HasTargets = len(targets) > 1
	}

	if upSignals, ok := data.GetUpSignals(); ok {
		params.HasUpSignals = len(upSignals) > 1
	}

	if downSignals, ok := data.GetDownSignals(); ok {
		params.HasDownSignals = len(downSignals) > 1
	}

	if data.HasAntennaSpecs() {
		params.HasAntennaSpec = true
	}

	ui.UpdateStatusBar(params)
}

func onDataReceived(dsnData response.DSN) {
	mapper.MapDataToFullData(dsnData, &data.FullData)
	app.QueueUpdateDraw(func() {
		if data.FullData.Stations == nil {
			return
		}

		if !data.IsReady {
			data.IsReady = true
			ui.CloseInitializingModal()
		}

		if data.SelectedStationIdx < 0 {
			data.SelectedStationIdx = 0
		}

		populateStationsData()
		updateDishesList()
	})
}

func runDSNDataLoader(result chan response.DSN, ce chan error, interrupt chan os.Signal) {
	chanDSNData := make(chan response.DSN)
	chanError := make(chan error)

	go network.LoadDSNData(chanDSNData, chanError)

	select {
	case dsnData := <-chanDSNData:
		result <- dsnData
	case err := <-chanError:
		ce <- err
		return
	case <-interrupt:
		ce <- fmt.Errorf("interrupted")
		return
	}

	// Run the loader repeatedly
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			go network.LoadDSNData(chanDSNData, chanError)

			select {
			case dsnData := <-chanDSNData:
				result <- dsnData
			case err := <-chanError:
				ce <- err
			case <-interrupt:
				ce <- fmt.Errorf("interrupted")
			}
		case <-interrupt:
			ce <- fmt.Errorf("interrupted")
			return
		}
	}
}

func populateStationsData() {
	stations := data.FullData.Stations

	if len(stations) == 0 {
		ui.UpdateStationsList("")
		return
	}

	if data.SelectedStationIdx < 0 {
		return
	}

	text := ""

	for i, station := range stations {
		if i == data.SelectedStationIdx {
			stationName := fmt.Sprintf("[green::b]%s[-:-:-:-]", station.Name)
			text = fmt.Sprintf("%s%s %s", text, stationName, strings.ToLower(station.Flag))
			info := fmt.Sprintf("[green::b]%s[-:-:-:-]", station.FriendlyName)

			ui.UpdateSelectedStation(info)
		} else {
			text += fmt.Sprintf("%s %s", station.Name, strings.ToLower(station.Flag))
		}

		if i < len(stations)-1 {
			text += "\n"
		}
	}

	ui.UpdateStationsList(text)
}

func updateDownSignalsTitleData() {
	downSignals, ok := data.GetDownSignals()
	if !ok {
		return
	}

	if downSignals == nil || len(downSignals) == 0 {
		ui.UpdateDownSignalsTitleData("No signal")
		ui.UpdateDownSignalData(model.DownSignal{})
		return
	}

	titles := ""

	if data.SelectedDownSignalIdx >= len(downSignals) {
		data.SelectedDownSignalIdx = 0
	}

	for i := range downSignals {
		if i == data.SelectedDownSignalIdx {
			titles += fmt.Sprintf("[green::b][%d][-:-:-:-]", i+1)
		} else {
			titles += fmt.Sprintf("[%d]", i+1)
		}
	}

	ui.UpdateDownSignalsTitleData(titles)
	ui.UpdateDownSignalData(downSignals[data.SelectedDownSignalIdx])
}

func updateUpSignalsTitleData() {
	upSignals, ok := data.GetUpSignals()
	if !ok {
		return
	}

	if upSignals == nil || len(upSignals) == 0 {
		ui.UpdateUpSignalsTitleData("No signal")
		ui.UpdateUpSignalData(model.UpSignal{})
		return
	}

	titles := ""

	if data.SelectedUpSignalIdx < 0 && len(upSignals) > 0 {
		data.SelectedUpSignalIdx = 0
	}

	for i := range upSignals {
		if i == data.SelectedUpSignalIdx {
			titles += fmt.Sprintf("[green::b][%d][-:-:-:-]", i+1)
		} else {
			titles += fmt.Sprintf("[%d]", i+1)
		}
	}

	ui.UpdateUpSignalsTitleData(titles)
	ui.UpdateUpSignalData(upSignals[data.SelectedUpSignalIdx])
}

func updateTargetsData() {
	targets, ok := data.GetTargets()
	if !ok {
		return
	}

	if targets == nil || len(targets) == 0 {
		ui.UpdateTargetsTitleData("No target")
		ui.UpdateTargetData(model.Target{})
		return
	}

	titles := ""

	if data.SelectedTargetIdx >= len(targets) {
		data.SelectedTargetIdx = 0
	}

	for i, target := range targets {
		if i == data.SelectedTargetIdx {
			titles += fmt.Sprintf("[green::b]%s[-:-:-:-]", target.Name)
		} else {
			titles += target.Name
		}

		if i < len(targets)-1 {
			titles += " - "
		} else {
			titles += " "
		}
	}

	ui.UpdateTargetsTitleData(titles)
	ui.UpdateTargetData(targets[data.SelectedTargetIdx])
}

func updateDownSignalSelection() {
	if !data.IsReady {
		return
	}

	downSignals, ok := data.GetDownSignals()
	if !ok || downSignals == nil || len(downSignals) == 0 {
		return
	}

	data.SelectedDownSignalIdx += 1

	if data.SelectedDownSignalIdx >= len(downSignals) {
		data.SelectedDownSignalIdx = 0
	}

	updateDownSignalsTitleData()
}

func updateUpSignalSelection() {
	if !data.IsReady {
		return
	}

	upSignals, ok := data.GetUpSignals()
	if !ok || upSignals == nil || len(upSignals) == 0 {
		return
	}

	data.SelectedUpSignalIdx += 1

	if data.SelectedUpSignalIdx >= len(upSignals) {
		data.SelectedUpSignalIdx = 0
	}

	updateUpSignalsTitleData()
}

func updateTargetSelection() {
	if !data.IsReady {
		return
	}

	targets, ok := data.GetTargets()
	if !ok || targets == nil || len(targets) == 0 {
		return
	}

	data.SelectedTargetIdx += 1

	if data.SelectedTargetIdx >= len(targets) {
		data.SelectedTargetIdx = 0
	}

	updateTargetsData()
}

func updateDishDetails(index int) {
	if data.SelectedStationIdx < 0 {
		return
	}

	data.SelectedDishIdx = index

	selectedStation := data.FullData.Stations[data.SelectedStationIdx]
	dish := selectedStation.Dishes[index]

	ui.UpdateDetailsText(dish)
}

func showPreview() {
	dish, ok := data.GetSelectedDish()
	if !ok {
		return
	}

	j, err := json.MarshalIndent(dish, "", "  ")
	if err != nil {
		return
	}

	ui.OpenJSONPreviewModal(string(j))
}

func showDishSpecs() {
	dish, ok := data.GetSelectedDish()
	if !ok {
		return
	}

	if dish.Specs == (model.DishSpecification{}) {
		return
	}

	ui.OpenDishSpecificationModal(dish.Specs)
}

func updateStationSelection() {
	if !data.IsReady {
		return
	}

	data.SelectedStationIdx += 1

	if data.SelectedStationIdx >= len(data.FullData.Stations) {
		data.SelectedStationIdx = 0
	}

	data.SelectedDishIdx = 0
	populateStationsData()
	updateDishesList()
}

func updateDishesList() {
	currDishSelected := data.SelectedDishIdx
	currUpSignalSelected := data.SelectedUpSignalIdx
	currDownSignalSelected := data.SelectedDownSignalIdx
	currTargetSelected := data.SelectedTargetIdx

	ui.ClearDishesList()

	selectedStation := data.FullData.Stations[data.SelectedStationIdx]

	for _, dish := range selectedStation.Dishes {
		upSignalText, downSignalText, nothing := "", "", ""

		workingDownSignals := dish.CountWorkingDownSignals()
		for i := 0; i < workingDownSignals; i++ {
			downSignalText += "[green::b]↓[-:-:-:-]"
		}

		workingUpSignals := dish.CountWorkingUpSignals()
		for i := 0; i < workingUpSignals; i++ {
			upSignalText += "[red::b]↑[-:-:-:-]"
		}

		workingTargets := dish.CountWorkingTargets()
		if workingDownSignals == 0 && workingUpSignals == 0 && workingTargets == 0 {
			nothing = "[red]✘[-]"
			dish.FriendlyName = fmt.Sprintf("[red]%s[-]", dish.FriendlyName)
		}

		if workingDownSignals > 0 && workingUpSignals > 0 && workingTargets > 0 {
			dish.FriendlyName = fmt.Sprintf("[green]%s[-]", dish.FriendlyName)
		}

		text := fmt.Sprintf("%s %s%s%s", dish.FriendlyName, upSignalText, downSignalText, nothing)

		ui.AddNewDish(text)
	}

	data.SelectedDownSignalIdx = currDishSelected
	data.SelectedUpSignalIdx = currUpSignalSelected
	data.SelectedDownSignalIdx = currDownSignalSelected
	data.SelectedTargetIdx = currTargetSelected
	ui.SetSelectedDish(currDishSelected)
}
