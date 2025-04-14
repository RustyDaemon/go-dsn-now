package gui

import "github.com/rivo/tview"

type StatusBarParams struct {
	DefaultStatus  bool
	HasTargets     bool
	HasUpSignals   bool
	HasDownSignals bool
	HasAntennaSpec bool
}

type UI struct {
	dishesList      *tview.List
	stationsList    *tview.TextView
	selectedStation *tview.TextView
	pages           *tview.Pages
	uiDetails       *uiDetails
	jsonPreview     *tview.TextArea
	dishSpecsView   *uiDishSpec
	aboutView       *uiAbout
	statusBar       *tview.TextView
	targetsView     *tview.Flex
	upSignalsView   *tview.Flex
	downSignalsView *tview.Flex
}

type uiDishSpec struct {
	name          *tview.TextView
	tp            *tview.TextView
	diameter      *tview.TextView
	txFreq        *tview.TextView
	rxFreq        *tview.TextView
	txPower       *tview.TextView
	precision     *tview.TextView
	antennaSpeed  *tview.TextView
	totalWeight   *tview.TextView
	dishWeight    *tview.TextView
	totalPanels   *tview.TextView
	surfaceArea   *tview.TextView
	opWindResist  *tview.TextView
	maxWindResist *tview.TextView
	builtIn       *tview.TextView
	url           *tview.TextView
}

type uiDetails struct {
	antennaView    *uiAntennaView
	targetView     *uiTargetView
	upSignalView   *uiUpSignalView
	downSignalView *uiDownSignalView
}

type uiAntennaView struct {
	name      *tview.TextView
	typeT     *tview.TextView
	activity  *tview.TextView
	azimuth   *tview.TextView
	elevation *tview.TextView
	wind      *tview.TextView
}

type uiTargetView struct {
	spacecraftName *tview.TextView
	upleg          *tview.TextView
	rtlt           *tview.TextView
}

type uiUpSignalView struct {
	source           *tview.TextView
	isActive         *tview.TextView
	signalType       *tview.TextView
	freqBand         *tview.TextView
	powerTransmitted *tview.TextView
}

type uiDownSignalView struct {
	source        *tview.TextView
	isActive      *tview.TextView
	signalType    *tview.TextView
	freqBand      *tview.TextView
	dataRate      *tview.TextView
	powerReceived *tview.TextView
}

type uiAbout struct {
	version *tview.TextView
	url     *tview.TextView
}
