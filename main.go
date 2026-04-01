package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/atotto/clipboard"

	"github.com/RustyDaemon/go-dsn-now/internal/config"
	"github.com/RustyDaemon/go-dsn-now/internal/data"

	"github.com/RustyDaemon/go-dsn-now/internal/gui"
	"github.com/RustyDaemon/go-dsn-now/internal/model"
	"github.com/RustyDaemon/go-dsn-now/internal/model/response"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	app               *tview.Application
	ui                *gui.UI
	appData           *model.AppData
	cfg               *config.Config
	httpClient        *http.Client
	cancel            context.CancelFunc
	refreshIntervalCh = make(chan time.Duration, 1)
	appSettings       *data.Settings
)

func main() {
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorDefault

	ctx, cancelFunc := signal.NotifyContext(context.Background(), os.Interrupt)
	cancel = cancelFunc
	defer cancel()

	cfg = config.Load()

	appSettings = data.LoadSettings()
	if appSettings != nil {
		if appSettings.RefreshIntervalSeconds > 0 {
			interval := time.Duration(appSettings.RefreshIntervalSeconds) * time.Second
			if interval < 10*time.Second {
				interval = 10 * time.Second
			}
			cfg.RefreshInterval = interval
		}
		if appSettings.Theme != "" {
			cfg.Theme = appSettings.Theme
		}
	} else {
		appSettings = &data.Settings{
			RefreshIntervalSeconds: int(cfg.RefreshInterval.Seconds()),
			Theme:                  cfg.Theme,
		}
	}

	httpClient = data.NewHTTPClient(cfg)

	appData = model.NewAppData()
	appData.Bookmarks = data.LoadBookmarks()

	app = tview.NewApplication()
	ui = gui.NewUI(cfg.Theme)
	appUI := ui.BuildAppUI(onListItemChanged)

	ui.SetStationClickedFunc(func(index int) {
		if !appData.IsReady || index < 0 || index >= len(appData.FullData.Stations) {
			return
		}
		appData.SelectedStationIdx = index
		appData.SelectedDishIdx = 0
		populateStationsData()
		updateDishesList()
	})

	updateStatusBar(true)

	chanConfig := make(chan response.DSNConfig)
	chanDSNData := make(chan response.DSN)
	chanError := make(chan error)

	go data.LoadDSNConfig(ctx, httpClient, cfg, chanConfig, chanError)

	select {
	case appData.DSNConfig = <-chanConfig:
	case err := <-chanError:
		log.Fatal(err)
	case <-ctx.Done():
		return
	}

	data.MapConfigToFullData(appData.DSNConfig, &appData.FullData)

	go runDSNDataLoader(ctx, chanDSNData, chanError)

	go func() {
		for {
			select {
			case dsnData := <-chanDSNData:
				onDataReceived(dsnData)
			case err := <-chanError:
				app.QueueUpdateDraw(func() {
					appData.LastError = err.Error()
					appData.ConsecutiveErrors++
					updateStatusBar(true)
				})
			case <-ctx.Done():
				return
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
			case 'j':
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
			case 'b':
				if !appData.IsReady {
					break
				}
				if dish, ok := appData.GetSelectedDish(); ok {
					if appData.Bookmarks[dish.Name] {
						delete(appData.Bookmarks, dish.Name)
					} else {
						appData.Bookmarks[dish.Name] = true
					}
					data.SaveBookmarks(appData.Bookmarks)
					updateDishesList()
				}
			case 'c':
				if !appData.IsReady {
					break
				}
				appData.CompactView = !appData.CompactView
				focus := ui.ToggleCompactView(appData.CompactView)
				app.SetFocus(focus)
				if appData.CompactView {
					updateCompactView()
				}
			case 'T':
				newTheme := ui.CycleTheme()
				appSettings.Theme = newTheme
				data.SaveSettings(appSettings)
				if appData.IsReady {
					populateStationsData()
					updateDishesList()
					if appData.CompactView {
						updateCompactView()
					}
				}
			case '+', '=':
				newInterval := cfg.RefreshInterval + 5*time.Second
				cfg.RefreshInterval = newInterval
				appSettings.RefreshIntervalSeconds = int(newInterval.Seconds())
				data.SaveSettings(appSettings)
				refreshIntervalCh <- newInterval
				updateStatusBar(true)
			case '-':
				newInterval := cfg.RefreshInterval - 5*time.Second
				if newInterval < 10*time.Second {
					newInterval = 10 * time.Second
				}
				cfg.RefreshInterval = newInterval
				appSettings.RefreshIntervalSeconds = int(newInterval.Seconds())
				data.SaveSettings(appSettings)
				refreshIntervalCh <- newInterval
				updateStatusBar(true)
			case 'y':
				if !appData.IsReady {
					break
				}
				text := ui.GetVisibleContent(appData.CompactView)
				if err := clipboard.WriteAll(text); err != nil {
					ui.SetStatusBarMessage("Clipboard unavailable")
				} else {
					ui.SetStatusBarMessage("Copied to clipboard")
				}
				go func() {
					time.Sleep(2 * time.Second)
					app.QueueUpdateDraw(func() {
						updateStatusBar(true)
					})
				}()
			case '?':
				if !appData.IsReady {
					break
				}
				appData.IsAboutShown = true
				updateStatusBar(false)
				ui.OpenAboutModal(cfg.AppVersion, cfg.AppGithubURL)
			}
		}

		if event.Rune() == 'q' {
			cancel()
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

	connStatus := "connected"
	if appData.ConsecutiveErrors >= 3 {
		connStatus = "disconnected"
	} else if appData.ConsecutiveErrors >= 1 {
		connStatus = "degraded"
	}

	params := gui.StatusBarParams{
		DefaultStatus:   defaultStatus,
		LastError:       appData.LastError,
		ConnStatus:      connStatus,
		SignalChanges:   appData.SignalChanges,
		RefreshInterval: fmt.Sprintf("%ds", int(cfg.RefreshInterval.Seconds())),
	}

	if !appData.LastUpdated.IsZero() {
		params.LastUpdated = appData.LastUpdated.Format("15:04:05")
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
	app.QueueUpdateDraw(func() {
		data.MapDataToFullData(dsnData, &appData.FullData)
		appData.LastError = ""
		appData.ConsecutiveErrors = 0
		appData.LastUpdated = time.Now()
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

		appData.DetectSignalChanges()
		populateStationsData()
		updateDishesList()
		if appData.CompactView {
			updateCompactView()
		}
	})
}

func runDSNDataLoader(ctx context.Context, result chan response.DSN, ce chan error) {
	chanDSNData := make(chan response.DSN)
	chanError := make(chan error)

	go data.LoadDSNData(ctx, httpClient, cfg, chanDSNData, chanError)

	select {
	case dsnData := <-chanDSNData:
		result <- dsnData
	case err := <-chanError:
		ce <- err
		return
	case <-ctx.Done():
		return
	}

	ticker := time.NewTicker(cfg.RefreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			go data.LoadDSNData(ctx, httpClient, cfg, chanDSNData, chanError)

			select {
			case dsnData := <-chanDSNData:
				result <- dsnData
			case err := <-chanError:
				ce <- err
			case <-ctx.Done():
				return
			}
		case newInterval := <-refreshIntervalCh:
			ticker.Reset(newInterval)
		case <-ctx.Done():
			return
		}
	}
}

func cycleSelection(currentIdx *int, count int) {
	*currentIdx++
	if *currentIdx >= count {
		*currentIdx = 0
	}
}

func buildIndexTitles(count, selectedIdx int) string {
	t := ui.Theme()
	var b strings.Builder
	for i := 0; i < count; i++ {
		if i == selectedIdx {
			fmt.Fprintf(&b, "[%s::b][%d][-:-:-:-]", t.Primary, i+1)
		} else {
			fmt.Fprintf(&b, "[%d]", i+1)
		}
	}
	return b.String()
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

	t := ui.Theme()
	var b strings.Builder

	for i, station := range stations {
		fmt.Fprintf(&b, `["%d"]`, i)
		if i == appData.SelectedStationIdx {
			fmt.Fprintf(&b, "[%s::b]%s[-:-:-:-] %s", t.Primary, station.Name, strings.ToLower(station.Flag))
			ui.UpdateSelectedStation(fmt.Sprintf("[%s::b]%s[-:-:-:-]", t.Primary, station.FriendlyName))
		} else {
			fmt.Fprintf(&b, "%s %s", station.Name, strings.ToLower(station.Flag))
		}
		fmt.Fprintf(&b, `[""]`)

		if i < len(stations)-1 {
			b.WriteString("\n")
		}
	}

	ui.UpdateStationsList(b.String())
}

func updateDownSignalsTitleData() {
	downSignals, ok := appData.GetDownSignals()
	if !ok {
		return
	}

	if len(downSignals) == 0 {
		ui.UpdateDownSignalsTitleData("No signal")
		ui.UpdateDownSignalData(model.DownSignal{})
		return
	}

	if appData.SelectedDownSignalIdx >= len(downSignals) {
		appData.SelectedDownSignalIdx = 0
	}

	ui.UpdateDownSignalsTitleData(buildIndexTitles(len(downSignals), appData.SelectedDownSignalIdx))
	ui.UpdateDownSignalData(downSignals[appData.SelectedDownSignalIdx])
}

func updateUpSignalsTitleData() {
	upSignals, ok := appData.GetUpSignals()
	if !ok {
		return
	}

	if len(upSignals) == 0 {
		ui.UpdateUpSignalsTitleData("No signal")
		ui.UpdateUpSignalData(model.UpSignal{})
		return
	}

	if appData.SelectedUpSignalIdx < 0 {
		appData.SelectedUpSignalIdx = 0
	}

	ui.UpdateUpSignalsTitleData(buildIndexTitles(len(upSignals), appData.SelectedUpSignalIdx))
	ui.UpdateUpSignalData(upSignals[appData.SelectedUpSignalIdx])
}

func updateTargetsData() {
	targets, ok := appData.GetTargets()
	if !ok {
		return
	}

	if len(targets) == 0 {
		ui.UpdateTargetsTitleData("No target")
		ui.UpdateTargetData(model.Target{})
		return
	}

	if appData.SelectedTargetIdx >= len(targets) {
		appData.SelectedTargetIdx = 0
	}

	t := ui.Theme()
	var b strings.Builder
	for i, target := range targets {
		if i == appData.SelectedTargetIdx {
			fmt.Fprintf(&b, "[%s::b]%s[-:-:-:-]", t.Primary, target.Name)
		} else {
			b.WriteString(target.Name)
		}
		if i < len(targets)-1 {
			b.WriteString(" - ")
		} else {
			b.WriteString(" ")
		}
	}

	ui.UpdateTargetsTitleData(b.String())
	ui.UpdateTargetData(targets[appData.SelectedTargetIdx])
}

func updateDownSignalSelection() {
	if !appData.IsReady {
		return
	}

	downSignals, ok := appData.GetDownSignals()
	if !ok || len(downSignals) == 0 {
		return
	}

	cycleSelection(&appData.SelectedDownSignalIdx, len(downSignals))
	updateDownSignalsTitleData()
}

func updateUpSignalSelection() {
	if !appData.IsReady {
		return
	}

	upSignals, ok := appData.GetUpSignals()
	if !ok || len(upSignals) == 0 {
		return
	}

	cycleSelection(&appData.SelectedUpSignalIdx, len(upSignals))
	updateUpSignalsTitleData()
}

func updateTargetSelection() {
	if !appData.IsReady {
		return
	}

	targets, ok := appData.GetTargets()
	if !ok || len(targets) == 0 {
		return
	}

	cycleSelection(&appData.SelectedTargetIdx, len(targets))
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

func updateCompactView() {
	var rows []gui.CompactRow
	for _, station := range appData.FullData.Stations {
		for _, dish := range station.Dishes {
			target := "-"
			if len(dish.Targets) > 0 && dish.Targets[0].Spacecraft.FriendlyName != "" {
				target = dish.Targets[0].Spacecraft.FriendlyName
			}
			rows = append(rows, gui.CompactRow{
				Station:     station.FriendlyName,
				Dish:        dish.FriendlyName,
				Target:      target,
				UpSignals:   dish.CountWorkingUpSignals(),
				DownSignals: dish.CountWorkingDownSignals(),
				Activity:    dish.Activity,
			})
		}
	}
	ui.UpdateCompactTable(rows)
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

	dt := ui.Theme()
	for _, dish := range selectedStation.Dishes {
		var b strings.Builder

		if appData.Bookmarks[dish.Name] {
			fmt.Fprintf(&b, "[%s]★[-] ", dt.Accent)
		}

		workingDownSignals := dish.CountWorkingDownSignals()
		workingUpSignals := dish.CountWorkingUpSignals()
		workingTargets := dish.CountWorkingTargets()

		if workingDownSignals == 0 && workingUpSignals == 0 && workingTargets == 0 {
			fmt.Fprintf(&b, "[%s]%s[-] ", dt.Inactive, dish.FriendlyName)
			fmt.Fprintf(&b, "[%s]✘[-]", dt.Inactive)
		} else if workingDownSignals > 0 && workingUpSignals > 0 && workingTargets > 0 {
			fmt.Fprintf(&b, "[%s]%s[-] ", dt.Primary, dish.FriendlyName)
		} else {
			fmt.Fprintf(&b, "%s ", dish.FriendlyName)
		}

		for i := 0; i < workingUpSignals; i++ {
			fmt.Fprintf(&b, "[%s::b]↑[-:-:-:-]", dt.SignalUp)
		}
		for i := 0; i < workingDownSignals; i++ {
			fmt.Fprintf(&b, "[%s::b]↓[-:-:-:-]", dt.SignalDown)
		}

		ui.AddNewDish(b.String())
	}

	appData.SelectedDishIdx = currDishSelected
	appData.SelectedUpSignalIdx = currUpSignalSelected
	appData.SelectedDownSignalIdx = currDownSignalSelected
	appData.SelectedTargetIdx = currTargetSelected
	ui.SetSelectedDish(currDishSelected)
}
