package gui

import (
	"fmt"
	"github.com/RustyDaemon/go-dsn-now/model"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strconv"
	"strings"
)

//TODO: refactor this file

func NewUI() *UI {
	return &UI{
		uiDetails: &uiDetails{
			antennaView:    &uiAntennaView{},
			targetView:     &uiTargetView{},
			upSignalView:   &uiUpSignalView{},
			downSignalView: &uiDownSignalView{},
		},
		dishSpecsView: &uiDishSpec{},
	}
}

func (u *UI) BuildAppUI(onListItemChanged func(index int)) *tview.Pages {
	u.statusBar = u.buildStatusBar()
	u.dishesList = buildDishesList(onListItemChanged)

	targetView := u.buildTargetsView()
	upSignalView := u.buildUpSignalsView()
	downSignalView := u.buildDownSignalsView()
	details := u.buildDetailsView(targetView, upSignalView, downSignalView)
	dishesMenu := u.buildDishesMenu()
	layout := u.buildLayout(dishesMenu, details)

	u.pages = tview.NewPages().
		AddPage("main", layout, true, true).
		AddPage("initializingModal", buildInitializingModal(), true, true).
		AddPage("jsonPreviewModal", u.buildJSONPreviewModal(), true, false).
		AddPage("dishSpecificationModal", u.buildDishSpecsModal(), true, false)

	return u.pages
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

func (u *UI) OpenDishSpecificationModal(spec model.DishSpecification) {
	u.dishSpecsView.name.SetText(fmt.Sprintf("[yellow]%s[-]", DashIfEmpty(spec.Name)))
	u.dishSpecsView.tp.SetText(fmt.Sprintf("[yellow]%s[-]", DashIfEmpty(spec.Type)))
	u.dishSpecsView.diameter.SetText(fmt.Sprintf("[yellow]%s[-]", DashIfEmpty(spec.Diameter)))
	u.dishSpecsView.txFreq.SetText(fmt.Sprintf("[yellow]%s[-]", DashIfEmpty(spec.TransmittersFrequency)))
	u.dishSpecsView.rxFreq.SetText(fmt.Sprintf("[yellow]%s[-]", DashIfEmpty(spec.ReceiversFrequency)))
	u.dishSpecsView.txPower.SetText(fmt.Sprintf("[yellow]%s[-]", DashIfEmpty(spec.TransmittersPower)))
	u.dishSpecsView.precision.SetText(fmt.Sprintf("[yellow]%s[-]", DashIfEmpty(spec.Precision)))
	u.dishSpecsView.antennaSpeed.SetText(fmt.Sprintf("[yellow]%s[-]", DashIfEmpty(spec.AntennaSpeed)))
	u.dishSpecsView.totalWeight.SetText(fmt.Sprintf("[yellow]%s[-]", DashIfEmpty(spec.TotalWeight)))
	u.dishSpecsView.dishWeight.SetText(fmt.Sprintf("[yellow]%s[-]", DashIfEmpty(spec.DishWeight)))
	u.dishSpecsView.totalPanels.SetText(fmt.Sprintf("[yellow]%s[-]", DashIfEmpty(spec.TotalPanels)))
	u.dishSpecsView.surfaceArea.SetText(fmt.Sprintf("[yellow]%s[-]", DashIfEmpty(spec.SurfaceArea)))
	u.dishSpecsView.opWindResist.SetText(fmt.Sprintf("[yellow]%s[-]", DashIfEmpty(spec.OperationalWindResistance)))
	u.dishSpecsView.maxWindResist.SetText(fmt.Sprintf("[yellow]%s[-]", DashIfEmpty(spec.WindResistance)))
	u.dishSpecsView.builtIn.SetText(fmt.Sprintf("[yellow]%s[-]", DashIfEmpty(spec.BuiltIn)))
	u.dishSpecsView.url.SetText(fmt.Sprintf("[yellow]%s[-]", DashIfEmpty(spec.WebUrl)))

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
	u.upSignalsView.SetTitle(fmt.Sprintf(" [red::b]↑[-:-:-:-] Up Signal: %s ", titles))
}

func (u *UI) UpdateDownSignalsTitleData(titles string) {
	u.downSignalsView.SetTitle(fmt.Sprintf(" [green::b]↓[-:-:-:-] Down Signal: %s ", titles))
}

func (u *UI) UpdateStatusBar(params StatusBarParams) {
	if !params.DefaultStatus {
		u.statusBar.SetText("[green](Esc)[-] close preview")
	} else {
		defaultText := "[green]s[-] station, %t%u%p%s[green]p[-] JSON, [green]q[-] exit"

		if params.HasTargets {
			defaultText = strings.Replace(defaultText, "%t", "[green]t[-] target, ", 1)
		} else {
			defaultText = strings.Replace(defaultText, "%t", "", 1)
		}

		if params.HasUpSignals {
			defaultText = strings.Replace(defaultText, "%u", "[green]u[-] up signal, ", 1)
		} else {
			defaultText = strings.Replace(defaultText, "%u", "", 1)
		}

		if params.HasDownSignals {
			defaultText = strings.Replace(defaultText, "%p", "[green]d[-] down signal, ", 1)
		} else {
			defaultText = strings.Replace(defaultText, "%p", "", 1)
		}

		if params.HasAntennaSpec {
			defaultText = strings.Replace(defaultText, "%s", "[green]i[-] specs, ", 1)
		} else {
			defaultText = strings.Replace(defaultText, "%s", "", 1)
		}

		u.statusBar.SetText(defaultText)
	}
}

func (u *UI) UpdateDetailsText(d model.Dish) {
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

	u.uiDetails.antennaView.name.SetText(fmt.Sprintf("Name: [yellow]%s[-]", DashIfEmpty(d.FriendlyName)))
	u.uiDetails.antennaView.typeT.SetText(fmt.Sprintf("Type: [yellow]%s[-]", DashIfEmpty(d.Type)))
	u.uiDetails.antennaView.activity.SetText(fmt.Sprintf("Activity: [yellow]%s[-]", DashIfEmpty(d.Activity)))
	u.uiDetails.antennaView.azimuth.SetText(fmt.Sprintf("Azimuth: [yellow]%s[-]", azimuth))
	u.uiDetails.antennaView.elevation.SetText(fmt.Sprintf("Elevation: [yellow]%s[-]", elevation))
	u.uiDetails.antennaView.wind.SetText(fmt.Sprintf("Wind: [yellow]%s[-]", wind))
}

func (u *UI) UpdateUpSignalData(upSignal model.UpSignal) {
	source := "-"
	powerTrans := "-"
	isActive := "[red]inactive[-]"
	freqBand := DefaultIfEmpty(upSignal.Band, "-")
	signalType := DefaultIfEmpty(upSignal.SignalType, "-")

	if upSignal.Spacecraft != (model.Spacecraft{}) && (len(upSignal.Spacecraft.FriendlyName) > 0) {
		source = upSignal.Spacecraft.FriendlyName
	}

	if len(upSignal.Power) > 0 {
		powerTrans = fmt.Sprintf("%s kW", upSignal.Power)
	}

	if upSignal.IsActive {
		isActive = "[green]active[-]"
	}

	if upSignal == (model.UpSignal{}) {
		u.upSignalsView.SetBorderColor(tcell.ColorRed)
	} else {
		u.upSignalsView.SetBorderColor(tcell.ColorDefault)
	}

	u.uiDetails.upSignalView.source.SetText(fmt.Sprintf("Source: [yellow]%s[-]", source))
	u.uiDetails.upSignalView.isActive.SetText(fmt.Sprintf("Is active: %s", isActive))
	u.uiDetails.upSignalView.signalType.SetText(fmt.Sprintf("Signal type: [yellow]%s[-]", signalType))
	u.uiDetails.upSignalView.freqBand.SetText(fmt.Sprintf("Frequency band: [yellow]%s[-]", freqBand))
	u.uiDetails.upSignalView.powerTransmitted.SetText(fmt.Sprintf("Power transmitted: [yellow]%s[-]", powerTrans))
}

func (u *UI) UpdateDownSignalData(downSignal model.DownSignal) {
	source := "-"
	powerReceived := "-"
	isActive := "[red]inactive[-]"
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
		isActive = "[green]active[-]"
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
		u.downSignalsView.SetBorderColor(tcell.ColorRed)
	} else {
		u.downSignalsView.SetBorderColor(tcell.ColorDefault)
	}

	u.uiDetails.downSignalView.source.SetText(fmt.Sprintf("Source: [yellow]%s[-]", source))
	u.uiDetails.downSignalView.isActive.SetText(fmt.Sprintf("Is active: %s", isActive))
	u.uiDetails.downSignalView.signalType.SetText(fmt.Sprintf("Signal type: [yellow]%s[-]", signalType))
	u.uiDetails.downSignalView.freqBand.SetText(fmt.Sprintf("Frequency band: [yellow]%s[-]", freqBand))
	u.uiDetails.downSignalView.dataRate.SetText(fmt.Sprintf("Data rate: [yellow]%s[-]", dataRate))
	u.uiDetails.downSignalView.powerReceived.SetText(fmt.Sprintf("Power received: [yellow]%s[-]", powerReceived))
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
		u.targetsView.SetBorderColor(tcell.ColorRed)
	} else {
		u.targetsView.SetBorderColor(tcell.ColorDefault)
	}

	u.uiDetails.targetView.spacecraftName.SetText(fmt.Sprintf("Spacecraft: [yellow]%s[-]", name))
	u.uiDetails.targetView.upleg.SetText(fmt.Sprintf("Range: [yellow]%s[-]", upLeg))
	u.uiDetails.targetView.rtlt.SetText(fmt.Sprintf("Round-trip light time: [yellow]%s[-]", rtlt))
}
