package gui

import (
	"fmt"
	"github.com/RustyDaemon/go-dsn-now/data"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (u *UI) buildStatusBar() *tview.TextView {
	statusBar := NewTextView("").SetTextAlign(tview.AlignCenter)
	statusBar.SetBorder(true)

	return statusBar
}

func buildDishesList(onListItemChanged func(index int)) *tview.List {
	dishesList := tview.NewList().
		SetHighlightFullLine(true).
		ShowSecondaryText(false).
		SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
			onListItemChanged(index)
		})
	dishesList.SetBorder(true).SetTitle("Antennas").SetTitleAlign(tview.AlignCenter)

	return dishesList
}

func (u *UI) buildDishesMenu() *tview.Flex {
	dishesMenu := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(u.dishesList, 0, 2, true).
		AddItem(u.buildStationView(), 0, 1, false)

	return dishesMenu
}

func (u *UI) buildLayout(dishesMenu, detailsPane tview.Primitive) *tview.Flex {
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().
			AddItem(dishesMenu, 0, 1, true). //or fixedSize 15
			AddItem(detailsPane, 0, 5, false),
			0, 1, true).
		AddItem(u.statusBar, 3, 1, false)

	return layout
}

func (u *UI) buildTargetsView() *tview.Flex {
	view := tview.NewFlex().SetDirection(tview.FlexRow)

	spacecraftNameView := NewTextView("[yellow]...[-]")
	uplegView := NewTextView("[yellow]...[-]")
	rtltView := NewTextView("[yellow]...[-]")

	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.AddItem(spacecraftNameView, 2, 1, false)
	flex.AddItem(uplegView, 0, 1, false)
	flex.AddItem(rtltView, 0, 1, false)
	flex.SetBorder(true).SetTitle(" Target ").SetBorder(true)

	view.AddItem(flex, 0, 1, false)

	u.targetsView = flex
	u.uiDetails.targetView.spacecraftName = spacecraftNameView
	u.uiDetails.targetView.upleg = uplegView
	u.uiDetails.targetView.rtlt = rtltView

	return view
}

func (u *UI) buildUpSignalsView() *tview.Flex {
	view := tview.NewFlex().SetDirection(tview.FlexRow)

	sourceView := NewTextView("[yellow]...[-]")
	activeView := NewTextView("[yellow]...[-]")
	signalTypeView := NewTextView("[yellow]...[-]")
	freqBandView := NewTextView("[yellow]...[-]")
	powerTransView := NewTextView("[yellow]...[-]")

	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.AddItem(sourceView, 2, 1, false)
	flex.AddItem(activeView, 1, 1, false)
	flex.AddItem(signalTypeView, 1, 1, false)
	flex.AddItem(freqBandView, 1, 1, false)
	flex.AddItem(powerTransView, 1, 1, false)
	flex.SetBorder(true).SetTitle(" [red]↑[-] Up Signal ").SetBorder(true)

	view.AddItem(flex, 9, 1, false)

	u.upSignalsView = flex
	u.uiDetails.upSignalView.isActive = activeView
	u.uiDetails.upSignalView.signalType = signalTypeView
	u.uiDetails.upSignalView.source = sourceView
	u.uiDetails.upSignalView.freqBand = freqBandView
	u.uiDetails.upSignalView.powerTransmitted = powerTransView

	return view
}

func (u *UI) buildDownSignalsView() *tview.Flex {
	view := tview.NewFlex().SetDirection(tview.FlexRow)

	sourceView := NewTextView("[yellow]...[-]")
	activeView := NewTextView("[yellow]...[-]")
	signalTypeView := NewTextView("[yellow]...[-]")
	freqBandView := NewTextView("[yellow]...[-]")
	dataRateView := NewTextView("[yellow]...[-]")
	powerRecView := NewTextView("[yellow]...[-]")

	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.AddItem(sourceView, 2, 1, false)
	flex.AddItem(activeView, 1, 1, false)
	flex.AddItem(signalTypeView, 1, 1, false)
	flex.AddItem(freqBandView, 1, 1, false)
	flex.AddItem(dataRateView, 1, 1, false)
	flex.AddItem(powerRecView, 1, 1, false)
	flex.SetBorder(true).SetTitle(" [green]↓[-] Down Signal ").SetBorder(true)

	view.AddItem(flex, 9, 1, false)

	u.downSignalsView = flex
	u.uiDetails.downSignalView.source = sourceView
	u.uiDetails.downSignalView.isActive = activeView
	u.uiDetails.downSignalView.signalType = signalTypeView
	u.uiDetails.downSignalView.freqBand = freqBandView
	u.uiDetails.downSignalView.dataRate = dataRateView
	u.uiDetails.downSignalView.powerReceived = powerRecView

	return view
}

func (u *UI) buildDetailsView(target, upSignal, downSignal tview.Primitive) *tview.Flex {
	view := tview.NewFlex().SetDirection(tview.FlexRow)

	nameView := NewTextView("[yellow]...[-]")
	typeView := NewTextView("[yellow]...[-]").SetTextAlign(tview.AlignRight)
	azimuthView := NewTextView("[yellow]...[-]")
	elevationView := NewTextView("[yellow]...[-]").SetTextAlign(tview.AlignCenter)
	windView := NewTextView("[yellow]...[-]").SetTextAlign(tview.AlignRight)
	activityView := NewTextView("[yellow]...[-]")

	//placeholder := tview.NewBox()

	antennaInfoView := tview.NewFlex().SetDirection(tview.FlexRow)
	antennaInfoView.AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(nameView, 0, 1, false).
		AddItem(typeView, 0, 1, false),
		1, 1, false)
	antennaInfoView.AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(azimuthView, 0, 1, false).
		AddItem(elevationView, 0, 1, false).
		AddItem(windView, 0, 1, false),
		2, 1, false)
	antennaInfoView.AddItem(activityView, 1, 1, false)
	antennaInfoView.SetBorder(true).SetTitle(" Antenna information ")

	view.AddItem(antennaInfoView, 6, 1, false)
	//view.AddItem(placeholder, 1, 1, false)
	view.AddItem(target, 6, 1, false)
	//view.AddItem(placeholder, 1, 1, false)

	view.AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(upSignal, 0, 1, false).
		AddItem(downSignal, 0, 1, false),
		0, 1, false)

	u.uiDetails.antennaView.name = nameView
	u.uiDetails.antennaView.typeT = typeView
	u.uiDetails.antennaView.activity = activityView
	u.uiDetails.antennaView.azimuth = azimuthView
	u.uiDetails.antennaView.elevation = elevationView
	u.uiDetails.antennaView.wind = windView

	return view
}

func (u *UI) buildStationView() *tview.Flex {
	view := tview.NewFlex().SetDirection(tview.FlexRow)
	view.SetBorder(true).SetTitle("Stations").SetTitleAlign(tview.AlignCenter)

	selectedStation := NewTextView("").SetTextAlign(tview.AlignCenter)
	stationsList := NewTextView("").SetTextAlign(tview.AlignCenter)

	view.AddItem(selectedStation, 2, 1, false)
	view.AddItem(stationsList, 0, 1, false)

	u.stationsList = stationsList
	u.selectedStation = selectedStation

	return view
}

func buildInitializingModal() tview.Primitive {
	modalInitializingContent := NewTextView("[green]Initializing[-] please wait[::l]...[-:-:-:-]").
		SetTextAlign(tview.AlignCenter)
	modalInitializingContent.SetBorder(true).SetTitleAlign(tview.AlignCenter)

	modal := func(p tview.Primitive, width, height int) tview.Primitive {
		return tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(p, height, 1, true).
				AddItem(nil, 0, 1, false), width, 1, true).
			AddItem(nil, 0, 1, false)
	}

	return modal(modalInitializingContent, 40, 3)
}

func (u *UI) buildJSONPreviewModal() tview.Primitive {
	modalPreviewContent := tview.NewTextArea().
		SetText("", false)
	modalPreviewContent.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyUp, tcell.KeyDown, tcell.KeyLeft, tcell.KeyRight:
			return event
		default:
			return nil
		}
	})
	modalPreviewContent.SetBorder(true).SetBorderColor(tcell.ColorYellow).
		SetTitleAlign(tview.AlignCenter).SetTitle(" JSON Preview ").
		SetTitleColor(tcell.ColorYellow)

	modal := func(p tview.Primitive) tview.Primitive {
		return tview.NewFlex().
			AddItem(nil, 3, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(p, 0, 5, true).
				AddItem(nil, 0, 1, false),
				0, 4, true).
			AddItem(nil, 3, 1, false)
	}

	u.jsonPreview = modalPreviewContent

	return modal(modalPreviewContent)
}

func (u *UI) buildDishSpecsModal() tview.Primitive {
	view := tview.NewFlex()

	nameView := NewTextView("Name:")
	nameValueView := NewTextView("[yellow]...[-]")

	diameterView := NewTextView("Diameter:")
	diameterValueView := NewTextView("[yellow]...[-]")

	typeView := NewTextView("Type:")
	typeValueView := NewTextView("[yellow]...[-]")

	txFreqView := NewTextView("Transmitters frequency:")
	txFreqValueView := NewTextView("[yellow]...[-]")

	rxFreqView := NewTextView("Receivers frequency:")
	rxFreqValueView := NewTextView("[yellow]...[-]")

	txPowerView := NewTextView("Transmitters power:")
	txPowerValueView := NewTextView("[yellow]...[-]")

	precisionView := NewTextView("Precision:")
	precisionValueView := NewTextView("[yellow]...[-]")

	antennaSpeedView := NewTextView("Antenna speed:")
	antennaSpeedValueView := NewTextView("[yellow]...[-]")

	totalWeightView := NewTextView("Total weight:")
	totalWeightValueView := NewTextView("[yellow]...[-]")

	dishWeightView := NewTextView("Dish weight:")
	dishWeightValueView := NewTextView("[yellow]...[-]")

	totalPanelsView := NewTextView("Total panels:")
	totalPanelsValueView := NewTextView("[yellow]...[-]")

	surfaceAreaView := NewTextView("Surface area:")
	surfaceAreaValueView := NewTextView("[yellow]...[-]")

	opWindResistView := NewTextView("Operational wind resistance:")
	opWindResistValueView := NewTextView("[yellow]...[-]")

	maxWindResistView := NewTextView("Max wind resistance:")
	maxWindResistValueView := NewTextView("[yellow]...[-]")

	builtInView := NewTextView("Built in:")
	builtInValueView := NewTextView("[yellow]...[-]")

	urlView := NewTextView("Web url:")
	urlValueView := NewTextView("[yellow]...[-]")

	separator := tview.NewBox()

	dataView := tview.NewFlex().SetDirection(tview.FlexRow)
	dataView.AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(nameView, 0, 1, false).
		AddItem(nameValueView, 0, 1, false),
		0, 1, false)
	dataView.AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(typeView, 0, 1, false).
		AddItem(typeValueView, 0, 1, false),
		0, 1, false)
	dataView.AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(diameterView, 0, 1, false).
		AddItem(diameterValueView, 0, 1, false),
		0, 1, false)
	dataView.AddItem(separator, 1, 1, false)
	dataView.AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(txFreqView, 0, 1, false).
		AddItem(txFreqValueView, 0, 1, false),
		0, 1, false)
	dataView.AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(rxFreqView, 0, 1, false).
		AddItem(rxFreqValueView, 0, 1, false),
		0, 1, false)
	dataView.AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(txPowerView, 0, 1, false).
		AddItem(txPowerValueView, 0, 1, false),
		0, 1, false)
	//dataView.AddItem(separator, 1, 1, false)
	dataView.AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(precisionView, 0, 1, false).
		AddItem(precisionValueView, 0, 1, false),
		0, 1, false)
	dataView.AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(antennaSpeedView, 0, 1, false).
		AddItem(antennaSpeedValueView, 0, 1, false),
		0, 1, false)
	//dataView.AddItem(separator, 1, 1, false)
	dataView.AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(totalWeightView, 0, 1, false).
		AddItem(totalWeightValueView, 0, 1, false),
		0, 1, false)
	dataView.AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(dishWeightView, 0, 1, false).
		AddItem(dishWeightValueView, 0, 1, false),
		0, 1, false)
	dataView.AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(totalPanelsView, 0, 1, false).
		AddItem(totalPanelsValueView, 0, 1, false),
		0, 1, false)
	dataView.AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(surfaceAreaView, 0, 1, false).
		AddItem(surfaceAreaValueView, 0, 1, false),
		0, 1, false)
	dataView.AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(opWindResistView, 0, 1, false).
		AddItem(opWindResistValueView, 0, 1, false),
		0, 1, false)
	dataView.AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(maxWindResistView, 0, 1, false).
		AddItem(maxWindResistValueView, 0, 1, false),
		0, 1, false)
	dataView.AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(builtInView, 0, 1, false).
		AddItem(builtInValueView, 0, 1, false),
		0, 1, false)
	//dataView.AddItem(separator, 1, 1, false)
	dataView.AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(urlView, 0, 1, false).
		AddItem(urlValueView, 0, 1, false),
		0, 1, false)

	view.AddItem(dataView, 0, 1, false)

	view.SetBorder(true).SetTitle(" Antenna Specification ").
		SetTitleColor(tcell.ColorYellow).SetTitleAlign(tview.AlignCenter).
		SetBorderColor(tcell.ColorYellow)

	modal := func(p tview.Primitive) tview.Primitive {
		return tview.NewFlex().
			AddItem(nil, 3, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 1, 1, false).
				AddItem(p, 0, 6, true).
				AddItem(nil, 3, 1, false),
				0, 4, true).
			AddItem(nil, 3, 1, false)
	}

	u.dishSpecsView.name = nameValueView
	u.dishSpecsView.tp = typeValueView
	u.dishSpecsView.diameter = diameterValueView
	u.dishSpecsView.txFreq = txFreqValueView
	u.dishSpecsView.rxFreq = rxFreqValueView
	u.dishSpecsView.txPower = txPowerValueView
	u.dishSpecsView.precision = precisionValueView
	u.dishSpecsView.antennaSpeed = antennaSpeedValueView
	u.dishSpecsView.totalWeight = totalWeightValueView
	u.dishSpecsView.dishWeight = dishWeightValueView
	u.dishSpecsView.totalPanels = totalPanelsValueView
	u.dishSpecsView.surfaceArea = surfaceAreaValueView
	u.dishSpecsView.opWindResist = opWindResistValueView
	u.dishSpecsView.maxWindResist = maxWindResistValueView
	u.dishSpecsView.builtIn = builtInValueView
	u.dishSpecsView.url = urlValueView

	return modal(view)
}

func (u *UI) buildAboutModal() tview.Primitive {
	view := tview.NewFlex()

	versionView := NewTextView("Version:")
	versionValueView := NewTextView(fmt.Sprintf("[yellow]%s[-]", data.AppVersion))

	githubView := NewTextView("GitHub:")
	githubViewView := NewTextView(fmt.Sprintf("[yellow:::%s]%s[-:-:-:-]", data.AppGithubUrl, data.AppGithubUrl))

	dataView := tview.NewFlex().SetDirection(tview.FlexRow)
	dataView.AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(versionView, 0, 1, false).
		AddItem(versionValueView, 0, 2, false),
		0, 1, false)
	dataView.AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(githubView, 0, 1, false).
		AddItem(githubViewView, 0, 2, false),
		0, 1, false)
	view.AddItem(dataView, 0, 1, false)

	view.SetBorder(true).SetTitle(" About ").
		SetTitleColor(tcell.ColorYellow).SetTitleAlign(tview.AlignCenter).
		SetBorderColor(tcell.ColorYellow)

	modal := func(p tview.Primitive) tview.Primitive {
		return tview.NewFlex().
			AddItem(nil, 3, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 3, false).
				AddItem(p, 4, 1, true).
				AddItem(nil, 0, 3, false),
				0, 3, true).
			AddItem(nil, 3, 1, false)
	}

	u.aboutView.version = versionValueView
	u.aboutView.url = githubViewView

	return modal(view)
}
