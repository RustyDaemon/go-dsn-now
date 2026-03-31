package gui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/RustyDaemon/go-dsn-now/internal/model"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func NewUI(themeName string) *UI {
	return &UI{
		theme: GetTheme(themeName),
		uiDetails: &uiDetails{
			antennaView:    &uiAntennaView{},
			targetView:     &uiTargetView{},
			upSignalView:   &uiUpSignalView{},
			downSignalView: &uiDownSignalView{},
		},
		dishSpecsView: &uiDishSpec{},
		aboutView:     &uiAbout{},
	}
}

func (u *UI) Theme() *Theme {
	return u.theme
}

func (u *UI) SetTheme(name string) {
	u.theme = GetTheme(name)
}

func (u *UI) CycleTheme() string {
	current := u.theme.Name
	for i, name := range ThemeNames {
		if name == current {
			next := ThemeNames[(i+1)%len(ThemeNames)]
			u.theme = GetTheme(next)
			return next
		}
	}
	u.theme = DarkTheme
	return "dark"
}

func (u *UI) BuildAppUI(onListItemChanged func(index int)) *tview.Pages {
	u.buildStatusBar()
	u.dishesList = buildDishesList(onListItemChanged)

	targetView := u.buildTargetsView()
	upSignalView := u.buildUpSignalsView()
	downSignalView := u.buildDownSignalsView()
	details := u.buildDetailsView(targetView, upSignalView, downSignalView)
	dishesMenu := u.buildDishesMenu()

	// Main detailed view content
	u.mainContent = tview.NewFlex().
		AddItem(dishesMenu, 0, 1, true).
		AddItem(details, 0, 5, false)

	// Compact table view
	u.compactTable = tview.NewTable().
		SetSelectable(true, false).
		SetFixed(1, 0)
	u.compactTable.SetBorder(true).SetTitle(" All Dishes (compact) ").SetTitleAlign(tview.AlignCenter)
	u.compactContent = tview.NewFlex().AddItem(u.compactTable, 0, 1, true)

	// Layout container that swaps between main and compact
	u.layoutContainer = tview.NewFlex().AddItem(u.mainContent, 0, 1, true)

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(u.layoutContainer, 0, 1, true).
		AddItem(u.statusBarFlex, 3, 1, false)

	u.pages = tview.NewPages().
		AddPage("main", layout, true, true).
		AddPage("initializingModal", buildInitializingModal(), true, true).
		AddPage("jsonPreviewModal", u.buildJSONPreviewModal(), true, false).
		AddPage("dishSpecificationModal", u.buildDishSpecsModal(), true, false).
		AddPage("aboutModal", u.buildAboutModal(), true, false)

	return u.pages
}

func (u *UI) ToggleCompactView(compact bool) tview.Primitive {
	if compact {
		u.layoutContainer.RemoveItem(u.mainContent)
		u.layoutContainer.AddItem(u.compactContent, 0, 1, true)
		return u.compactTable
	}
	u.layoutContainer.RemoveItem(u.compactContent)
	u.layoutContainer.AddItem(u.mainContent, 0, 1, true)
	return u.dishesList
}

func (u *UI) UpdateCompactTable(rows []CompactRow) {
	t := u.theme
	u.compactTable.Clear()

	// Header
	headers := []string{"Station", "Dish", "Target", "↑", "↓", "Activity"}
	for i, h := range headers {
		cell := tview.NewTableCell(h).SetSelectable(false).SetAttributes(tcell.AttrBold)
		u.compactTable.SetCell(0, i, cell)
	}

	for i, row := range rows {
		r := i + 1
		u.compactTable.SetCell(r, 0, tview.NewTableCell(row.Station))
		u.compactTable.SetCell(r, 1, tview.NewTableCell(row.Dish))
		u.compactTable.SetCell(r, 2, tview.NewTableCell(row.Target))

		upCell := tview.NewTableCell(fmt.Sprintf("%d", row.UpSignals))
		if row.UpSignals > 0 {
			upCell.SetText(fmt.Sprintf("[%s]%d[-]", t.SignalUp, row.UpSignals))
		}
		u.compactTable.SetCell(r, 3, upCell)

		downCell := tview.NewTableCell(fmt.Sprintf("%d", row.DownSignals))
		if row.DownSignals > 0 {
			downCell.SetText(fmt.Sprintf("[%s]%d[-]", t.SignalDown, row.DownSignals))
		}
		u.compactTable.SetCell(r, 4, downCell)

		actCell := tview.NewTableCell(row.Activity)
		if row.UpSignals == 0 && row.DownSignals == 0 {
			actCell.SetText(fmt.Sprintf("[%s]%s[-]", t.Inactive, row.Activity))
		}
		u.compactTable.SetCell(r, 5, actCell)
	}
}

func (u *UI) OpenJSONPreviewModal(json string) {
	u.jsonPreview.SetText(json, false)
	u.pages.ShowPage("jsonPreviewModal")
}

func (u *UI) CloseJSONPreviewModal() {
	u.pages.HidePage("jsonPreviewModal")
}

func (u *UI) CloseInitializingModal() {
	u.pages.HidePage("initializingModal")
}

func (u *UI) OpenAboutModal(version, githubURL string) {
	t := u.theme
	var b strings.Builder
	fmt.Fprintf(&b, " [::b]Keybindings[-:-:-:-]\n\n")
	fmt.Fprintf(&b, "  [%s]s[-]         Cycle station\n", t.Primary)
	fmt.Fprintf(&b, "  [%s]t[-]         Cycle target\n", t.Primary)
	fmt.Fprintf(&b, "  [%s]u[-]         Cycle up signal\n", t.Primary)
	fmt.Fprintf(&b, "  [%s]d[-]         Cycle down signal\n", t.Primary)
	fmt.Fprintf(&b, "  [%s]↑[-] [%s]↓[-]       Navigate dishes\n", t.Primary, t.Primary)
	fmt.Fprintf(&b, "  [%s]b[-]         Bookmark dish\n", t.Primary)
	fmt.Fprintf(&b, "  [%s]c[-]         Toggle compact view\n", t.Primary)
	fmt.Fprintf(&b, "  [%s]j[-]         JSON preview\n", t.Primary)
	fmt.Fprintf(&b, "  [%s]i[-]         Antenna specs\n", t.Primary)
	fmt.Fprintf(&b, "  [%s]T[-]         Cycle theme\n", t.Primary)
	fmt.Fprintf(&b, "  [%s]?[-]         This help\n", t.Primary)
	fmt.Fprintf(&b, "  [%s]Esc[-]       Close modal\n", t.Primary)
	fmt.Fprintf(&b, "  [%s]q[-]         Quit\n", t.Primary)
	fmt.Fprintf(&b, "\n [::b]About[-:-:-:-]\n\n")
	fmt.Fprintf(&b, "  Version  [%s]%s[-]\n", t.Secondary, version)
	fmt.Fprintf(&b, "  GitHub   [%s:::%s]%s[-:-:-:-]\n", t.Secondary, githubURL, githubURL)
	u.aboutView.content.SetText(b.String())
	u.pages.ShowPage("aboutModal")
}

func (u *UI) CloseAboutModal() {
	u.pages.HidePage("aboutModal")
}

func (u *UI) OpenDishSpecificationModal(spec model.DishSpecification) {
	c := u.theme.Secondary
	setSpecField(u.dishSpecsView.name, spec.Name, c)
	setSpecField(u.dishSpecsView.tp, spec.Type, c)
	setSpecField(u.dishSpecsView.diameter, spec.Diameter, c)
	setSpecField(u.dishSpecsView.txFreq, spec.TransmittersFrequency, c)
	setSpecField(u.dishSpecsView.rxFreq, spec.ReceiversFrequency, c)
	setSpecField(u.dishSpecsView.txPower, spec.TransmittersPower, c)
	setSpecField(u.dishSpecsView.precision, spec.Precision, c)
	setSpecField(u.dishSpecsView.antennaSpeed, spec.AntennaSpeed, c)
	setSpecField(u.dishSpecsView.totalWeight, spec.TotalWeight, c)
	setSpecField(u.dishSpecsView.dishWeight, spec.DishWeight, c)
	setSpecField(u.dishSpecsView.totalPanels, spec.TotalPanels, c)
	setSpecField(u.dishSpecsView.surfaceArea, spec.SurfaceArea, c)
	setSpecField(u.dishSpecsView.opWindResist, spec.OperationalWindResistance, c)
	setSpecField(u.dishSpecsView.maxWindResist, spec.WindResistance, c)
	setSpecField(u.dishSpecsView.builtIn, spec.BuiltIn, c)
	setSpecField(u.dishSpecsView.url, spec.WebUrl, c)

	u.pages.ShowPage("dishSpecificationModal")
}

func (u *UI) CloseDishSpecificationModal() {
	u.pages.HidePage("dishSpecificationModal")
}

func (u *UI) UpdateStationsList(text string) {
	u.stationsList.SetText(text)
}

func (u *UI) UpdateSelectedStation(text string) {
	u.selectedStation.SetText(text)
}

func (u *UI) AddNewDish(name string) {
	u.dishesList.AddItem(name, "", 0, nil)
}

func (u *UI) ClearDishesList() {
	u.dishesList.Clear()
}

func (u *UI) SetSelectedDish(index int) {
	u.dishesList.SetCurrentItem(index)
}

func (u *UI) UpdateTargetsTitleData(titles string) {
	u.targetsView.SetTitle(fmt.Sprintf(" Target: %s", titles))
}

func (u *UI) UpdateUpSignalsTitleData(titles string) {
	u.upSignalsView.SetTitle(fmt.Sprintf(" [%s::b]↑[-:-:-:-] Up Signal: %s ", u.theme.SignalUp, titles))
}

func (u *UI) UpdateDownSignalsTitleData(titles string) {
	u.downSignalsView.SetTitle(fmt.Sprintf(" [%s::b]↓[-:-:-:-] Down Signal: %s ", u.theme.SignalDown, titles))
}

func (u *UI) UpdateStatusBar(params StatusBarParams) {
	t := u.theme
	// Left section: keybindings
	if !params.DefaultStatus {
		u.statusBarLeft.SetText(fmt.Sprintf("[%s](Esc)[-] close", t.Primary))
	} else {
		var b strings.Builder
		fmt.Fprintf(&b, "[%s]s[-]tation, ", t.Primary)
		if params.HasTargets {
			fmt.Fprintf(&b, "[%s]t[-]arget, ", t.Primary)
		}
		if params.HasUpSignals {
			fmt.Fprintf(&b, "[%s]u[-]p signal, ", t.Primary)
		}
		if params.HasDownSignals {
			fmt.Fprintf(&b, "[%s]d[-]own signal, ", t.Primary)
		}
		if params.HasAntennaSpec {
			fmt.Fprintf(&b, "[%s]i[-]nfo, ", t.Primary)
		}
		fmt.Fprintf(&b, "[%s]j[-]SON, [%s]?[-], [%s]q[-]uit", t.Primary, t.Primary, t.Primary)
		u.statusBarLeft.SetText(b.String())
	}

	// Center section: error display or signal change notifications
	if params.LastError != "" {
		u.statusBarCenter.SetText(fmt.Sprintf("[%s]Error: %s[-]", t.Error, params.LastError))
	} else if len(params.SignalChanges) > 0 {
		u.statusBarCenter.SetText(fmt.Sprintf("[%s]%s[-]", t.Accent, strings.Join(params.SignalChanges, ", ")))
	} else {
		u.statusBarCenter.SetText("")
	}

	// Right section: connection indicator + last updated timestamp
	var right strings.Builder
	switch params.ConnStatus {
	case "connected":
		fmt.Fprintf(&right, "[%s]●[-] ", t.StatusConnected)
	case "degraded":
		fmt.Fprintf(&right, "[%s]●[-] ", t.StatusDegraded)
	case "disconnected":
		fmt.Fprintf(&right, "[%s]●[-] ", t.StatusError)
	}
	if params.LastUpdated != "" {
		fmt.Fprintf(&right, "Updated: [%s]%s[-]", t.Secondary, params.LastUpdated)
	}
	u.statusBarRight.SetText(right.String())
}

func (u *UI) UpdateDetailsText(d model.Dish) {
	t := u.theme
	azimuth := "-"
	elevation := "-"
	wind := "-"

	if len(d.AzimuthAngle) > 0 {
		azimuth = fmt.Sprintf("%s˚", d.AzimuthAngle)
	}

	if len(d.ElevationAngle) > 0 {
		elevation = fmt.Sprintf("%s˚", d.ElevationAngle)
	}

	if len(d.WindSpeed) > 0 {
		wind = fmt.Sprintf("%s km/h", d.WindSpeed)
	}

	u.uiDetails.antennaView.name.SetText(fmt.Sprintf("Name: [%s]%s[-]", t.Secondary, DashIfEmpty(d.FriendlyName)))
	u.uiDetails.antennaView.typeT.SetText(fmt.Sprintf("Type: [%s]%s[-]", t.Secondary, DashIfEmpty(d.Type)))
	u.uiDetails.antennaView.activity.SetText(fmt.Sprintf("Activity: [%s]%s[-]", t.Secondary, DashIfEmpty(d.Activity)))
	u.uiDetails.antennaView.azimuth.SetText(fmt.Sprintf("Azimuth: [%s]%s[-]", t.Secondary, azimuth))
	u.uiDetails.antennaView.elevation.SetText(fmt.Sprintf("Elevation: [%s]%s[-]", t.Secondary, elevation))
	u.uiDetails.antennaView.wind.SetText(fmt.Sprintf("Wind: [%s]%s[-]", t.Secondary, wind))
}

func (u *UI) UpdateUpSignalData(upSignal model.UpSignal) {
	t := u.theme
	source := "-"
	powerTrans := "-"
	isActive := fmt.Sprintf("[%s]inactive[-]", t.Error)
	freqBand := DefaultIfEmpty(upSignal.Band, "-")
	signalType := DefaultIfEmpty(upSignal.SignalType, "-")

	if upSignal.Spacecraft != (model.Spacecraft{}) && (len(upSignal.Spacecraft.FriendlyName) > 0) {
		source = upSignal.Spacecraft.FriendlyName
	}

	if len(upSignal.Power) > 0 {
		powerTrans = fmt.Sprintf("%s kW", upSignal.Power)
	}

	if upSignal.IsActive {
		isActive = fmt.Sprintf("[%s]active[-]", t.Primary)
	}

	if upSignal == (model.UpSignal{}) {
		u.upSignalsView.SetBorderColor(t.ErrorBorder)
	} else {
		u.upSignalsView.SetBorderColor(tcell.ColorDefault)
	}

	u.uiDetails.upSignalView.source.SetText(fmt.Sprintf("Source: [%s]%s[-]", t.Secondary, source))
	u.uiDetails.upSignalView.isActive.SetText(fmt.Sprintf("Is active: %s", isActive))
	u.uiDetails.upSignalView.signalType.SetText(fmt.Sprintf("Signal type: [%s]%s[-]", t.Secondary, signalType))
	u.uiDetails.upSignalView.freqBand.SetText(fmt.Sprintf("Frequency band: [%s]%s[-]", t.Secondary, freqBand))
	u.uiDetails.upSignalView.powerTransmitted.SetText(fmt.Sprintf("Power transmitted: [%s]%s[-]", t.Secondary, powerTrans))
}

func (u *UI) UpdateDownSignalData(downSignal model.DownSignal) {
	t := u.theme
	source := "-"
	powerReceived := "-"
	isActive := fmt.Sprintf("[%s]inactive[-]", t.Error)
	dataRate := "-"
	freqBand := DefaultIfEmpty(downSignal.Band, "-")
	signalType := DefaultIfEmpty(downSignal.SignalType, "-")

	if downSignal.Spacecraft != (model.Spacecraft{}) && (len(downSignal.Spacecraft.FriendlyName) > 0) {
		source = downSignal.Spacecraft.FriendlyName
	}

	if len(downSignal.Power) > 0 {
		powerReceived = fmt.Sprintf("%s dBm", downSignal.Power)
	}

	if downSignal.IsActive {
		isActive = fmt.Sprintf("[%s]active[-]", t.Primary)
	}

	if len(downSignal.DataRate) > 0 {
		val, err := strconv.ParseFloat(downSignal.DataRate, 64)
		if err == nil {
			if val < 0 {
				dataRate = downSignal.DataRate
			} else if val > 1000000 {
				val = val / 1000000
				dataRate = fmt.Sprintf("%.2f Mb/s", val)
			} else if val > 1000 {
				val = val / 1000
				dataRate = fmt.Sprintf("%.2f kb/s", val)
			} else {
				dataRate = fmt.Sprintf("%.2f b/s", val)
			}
		} else {
			dataRate = downSignal.DataRate
		}
	}

	if downSignal == (model.DownSignal{}) {
		u.downSignalsView.SetBorderColor(t.ErrorBorder)
	} else {
		u.downSignalsView.SetBorderColor(tcell.ColorDefault)
	}

	u.uiDetails.downSignalView.source.SetText(fmt.Sprintf("Source: [%s]%s[-]", t.Secondary, source))
	u.uiDetails.downSignalView.isActive.SetText(fmt.Sprintf("Is active: %s", isActive))
	u.uiDetails.downSignalView.signalType.SetText(fmt.Sprintf("Signal type: [%s]%s[-]", t.Secondary, signalType))
	u.uiDetails.downSignalView.freqBand.SetText(fmt.Sprintf("Frequency band: [%s]%s[-]", t.Secondary, freqBand))
	u.uiDetails.downSignalView.dataRate.SetText(fmt.Sprintf("Data rate: [%s]%s[-]", t.Secondary, dataRate))
	u.uiDetails.downSignalView.powerReceived.SetText(fmt.Sprintf("Power received: [%s]%s[-]", t.Secondary, powerReceived))
}

func (u *UI) UpdateTargetData(target model.Target) {
	name := "-"
	upLeg := "-"
	rtlt := "-"

	if target.Spacecraft != (model.Spacecraft{}) && (len(target.Spacecraft.FriendlyName) > 0) {
		name = target.Spacecraft.FriendlyName
	}

	if len(target.UplegRange) > 0 {
		val, err := strconv.ParseFloat(target.UplegRange, 64)
		if err == nil {
			if val < 0 {
				upLeg = target.UplegRange
			} else if val > 1000000000 {
				val = val / 1000000000
				upLeg = fmt.Sprintf("%.2f billion km", val)
			} else if val > 1000000 {
				val = val / 1000000
				upLeg = fmt.Sprintf("%.2f million km", val)
			} else if val > 1000 {
				val = val / 1000
				upLeg = fmt.Sprintf("%.2f thousand km", val)
			} else {
				upLeg = fmt.Sprintf("%.2f km", val)
			}
		} else {
			upLeg = target.UplegRange
		}
	}

	if len(target.Rtlt) > 0 {
		val, err := strconv.ParseFloat(target.Rtlt, 64)
		if err == nil {
			if val < 0 {
				rtlt = target.Rtlt
			} else if val > 3600 {
				val = val / 3600
				rtlt = fmt.Sprintf("%.2f hours", val)
			} else if val > 60 {
				val = val / 60
				rtlt = fmt.Sprintf("%.2f minutes", val)
			} else {
				rtlt = fmt.Sprintf("%.2f seconds", val)
			}
		} else {
			rtlt = target.Rtlt
		}
	}

	if target == (model.Target{}) || (target.UplegRange == "-1" && target.Rtlt == "-1" && target.Spacecraft == (model.Spacecraft{})) {
		u.targetsView.SetBorderColor(u.theme.ErrorBorder)
	} else {
		u.targetsView.SetBorderColor(tcell.ColorDefault)
	}

	t := u.theme
	u.uiDetails.targetView.spacecraftName.SetText(fmt.Sprintf("Spacecraft: [%s]%s[-]", t.Secondary, name))
	u.uiDetails.targetView.upleg.SetText(fmt.Sprintf("Range: [%s]%s[-]", t.Secondary, upLeg))
	u.uiDetails.targetView.rtlt.SetText(fmt.Sprintf("Round-trip light time: [%s]%s[-]", t.Secondary, rtlt))
}
