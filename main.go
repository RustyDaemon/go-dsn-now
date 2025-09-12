package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/RustyDaemon/go-dsn-now/internal/data"

	"github.com/RustyDaemon/go-dsn-now/internal/gui"
	"github.com/RustyDaemon/go-dsn-now/internal/model"
	"github.com/RustyDaemon/go-dsn-now/internal/model/response"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	app     *tview.Application
	ui      *gui.UI
	appData *model.AppData
)

func main() {
	appData = model.NewAppData()

	app = tview.NewApplication()
	ui = gui.NewUI()
	appUI := ui.BuildAppUI(onListItemChanged)

	updateStatusBar(true)

	interrupt := make(chan os.Signal, 1)
	chanConfig := make(chan response.DSNConfig)
	chanDSNData := make(chan response.DSN)
	chanError := make(chan error)

	signal.Notify(interrupt, os.Interrupt)

	go data.LoadDSNConfig(chanConfig, chanError)

	select {
	case appData.DSNConfig = <-chanConfig:
	case err := <-chanError:
		log.Fatal(err)
	}

	data.MapConfigToFullData(appData.DSNConfig, &appData.FullData)

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

	if appData.IsNoModalShown() {
		updateStatusBar(true)
	}
}

func setKeybindings() {
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if appData.IsNoModalShown() {
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
				if !appData.IsReady {
					break
				}
				appData.IsPreviewShown = true
				updateStatusBar(false)
				showPreview()
			case 'i':
				if !appData.IsReady {
					break
				}
				appData.IsSpecsShown = true
				updateStatusBar(false)
				showDishSpecs()
			case '?':
				if !appData.IsReady {
					break
				}
				appData.IsAboutShown = true
				updateStatusBar(false)
				ui.OpenAboutModal()
			}
		}

		if event.Rune() == 'q' {
			app.Stop()
		}

		if event.Key() == tcell.KeyEscape && appData.IsPreviewShown {
			ui.CloseJSONPreviewModal()
			appData.IsPreviewShown = false
			updateStatusBar(true)
		} else if event.Key() == tcell.KeyEscape && appData.IsSpecsShown {
			ui.CloseDishSpecificationModal()
			appData.IsSpecsShown = false
			updateStatusBar(true)
		} else if event.Key() == tcell.KeyEscape && appData.IsAboutShown {
			ui.CloseAboutModal()
			appData.IsAboutShown = false
			updateStatusBar(true)
		}

		return event
	})
}

func updateStatusBar(defaultStatus bool) {
	if !appData.IsReady {
		return
	}

	params := gui.StatusBarParams{
		DefaultStatus: defaultStatus,
	}

	if !params.DefaultStatus {
		ui.UpdateStatusBar(params)
		return
	}

	if targets, ok := appData.GetTargets(); ok {
		params.HasTargets = len(targets) > 1
	}

	if upSignals, ok := appData.GetUpSignals(); ok {
		params.HasUpSignals = len(upSignals) > 1
	}

	if downSignals, ok := appData.GetDownSignals(); ok {
		params.HasDownSignals = len(downSignals) > 1
	}

	if appData.HasAntennaSpecs() {
		params.HasAntennaSpec = true
	}

	ui.UpdateStatusBar(params)
}

func onDataReceived(dsnData response.DSN) {
	data.MapDataToFullData(dsnData, &appData.FullData)
	app.QueueUpdateDraw(func() {
		if appData.FullData.Stations == nil {
			return
		}

		if !appData.IsReady {
			appData.IsReady = true
			ui.CloseInitializingModal()
		}

		if appData.SelectedStationIdx < 0 {
			appData.SelectedStationIdx = 0
		}

		populateStationsData()
		updateDishesList()
	})
}

func runDSNDataLoader(result chan response.DSN, ce chan error, interrupt chan os.Signal) {
	chanDSNData := make(chan response.DSN)
	chanError := make(chan error)

	go data.LoadDSNData(chanDSNData, chanError)

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

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			go data.LoadDSNData(chanDSNData, chanError)

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
	stations := appData.FullData.Stations

	if len(stations) == 0 {
		ui.UpdateStationsList("")
		return
	}

	if appData.SelectedStationIdx < 0 {
		return
	}

	text := ""

	for i, station := range stations {
		if i == appData.SelectedStationIdx {
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
	downSignals, ok := appData.GetDownSignals()
	if !ok {
		return
	}

	if downSignals == nil || len(downSignals) == 0 {
		ui.UpdateDownSignalsTitleData("No signal")
		ui.UpdateDownSignalData(model.DownSignal{})
		return
	}

	titles := ""

	if appData.SelectedDownSignalIdx >= len(downSignals) {
		appData.SelectedDownSignalIdx = 0
	}

	for i := range downSignals {
		if i == appData.SelectedDownSignalIdx {
			titles += fmt.Sprintf("[green::b][%d][-:-:-:-]", i+1)
		} else {
			titles += fmt.Sprintf("[%d]", i+1)
		}
	}

	ui.UpdateDownSignalsTitleData(titles)
	ui.UpdateDownSignalData(downSignals[appData.SelectedDownSignalIdx])
}

func updateUpSignalsTitleData() {
	upSignals, ok := appData.GetUpSignals()
	if !ok {
		return
	}

	if upSignals == nil || len(upSignals) == 0 {
		ui.UpdateUpSignalsTitleData("No signal")
		ui.UpdateUpSignalData(model.UpSignal{})
		return
	}

	titles := ""

	if appData.SelectedUpSignalIdx < 0 && len(upSignals) > 0 {
		appData.SelectedUpSignalIdx = 0
	}

	for i := range upSignals {
		if i == appData.SelectedUpSignalIdx {
			titles += fmt.Sprintf("[green::b][%d][-:-:-:-]", i+1)
		} else {
			titles += fmt.Sprintf("[%d]", i+1)
		}
	}

	ui.UpdateUpSignalsTitleData(titles)
	ui.UpdateUpSignalData(upSignals[appData.SelectedUpSignalIdx])
}

func updateTargetsData() {
	targets, ok := appData.GetTargets()
	if !ok {
		return
	}

	if targets == nil || len(targets) == 0 {
		ui.UpdateTargetsTitleData("No target")
		ui.UpdateTargetData(model.Target{})
		return
	}

	titles := ""

	if appData.SelectedTargetIdx >= len(targets) {
		appData.SelectedTargetIdx = 0
	}

	for i, target := range targets {
		if i == appData.SelectedTargetIdx {
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
	ui.UpdateTargetData(targets[appData.SelectedTargetIdx])
}

func updateDownSignalSelection() {
	if !appData.IsReady {
		return
	}

	downSignals, ok := appData.GetDownSignals()
	if !ok || downSignals == nil || len(downSignals) == 0 {
		return
	}

	appData.SelectedDownSignalIdx += 1

	if appData.SelectedDownSignalIdx >= len(downSignals) {
		appData.SelectedDownSignalIdx = 0
	}

	updateDownSignalsTitleData()
}

func updateUpSignalSelection() {
	if !appData.IsReady {
		return
	}

	upSignals, ok := appData.GetUpSignals()
	if !ok || upSignals == nil || len(upSignals) == 0 {
		return
	}

	appData.SelectedUpSignalIdx += 1

	if appData.SelectedUpSignalIdx >= len(upSignals) {
		appData.SelectedUpSignalIdx = 0
	}

	updateUpSignalsTitleData()
}

func updateTargetSelection() {
	if !appData.IsReady {
		return
	}

	targets, ok := appData.GetTargets()
	if !ok || targets == nil || len(targets) == 0 {
		return
	}

	appData.SelectedTargetIdx += 1

	if appData.SelectedTargetIdx >= len(targets) {
		appData.SelectedTargetIdx = 0
	}

	updateTargetsData()
}

func updateDishDetails(index int) {
	if appData.SelectedStationIdx < 0 {
		return
	}

	appData.SelectedDishIdx = index

	selectedStation := appData.FullData.Stations[appData.SelectedStationIdx]
	dish := selectedStation.Dishes[index]

	ui.UpdateDetailsText(dish)
}

func showPreview() {
	dish, ok := appData.GetSelectedDish()
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
	dish, ok := appData.GetSelectedDish()
	if !ok {
		return
	}

	if dish.Specs == (model.DishSpecification{}) {
		return
	}

	ui.OpenDishSpecificationModal(dish.Specs)
}

func updateStationSelection() {
	if !appData.IsReady {
		return
	}

	appData.SelectedStationIdx += 1

	if appData.SelectedStationIdx >= len(appData.FullData.Stations) {
		appData.SelectedStationIdx = 0
	}

	appData.SelectedDishIdx = 0
	populateStationsData()
	updateDishesList()
}

func updateDishesList() {
	currDishSelected := appData.SelectedDishIdx
	currUpSignalSelected := appData.SelectedUpSignalIdx
	currDownSignalSelected := appData.SelectedDownSignalIdx
	currTargetSelected := appData.SelectedTargetIdx

	ui.ClearDishesList()

	selectedStation := appData.FullData.Stations[appData.SelectedStationIdx]

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

	appData.SelectedDishIdx = currDishSelected
	appData.SelectedUpSignalIdx = currUpSignalSelected
	appData.SelectedDownSignalIdx = currDownSignalSelected
	appData.SelectedTargetIdx = currTargetSelected
	ui.SetSelectedDish(currDishSelected)
}
