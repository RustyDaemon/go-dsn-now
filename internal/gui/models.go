package gui

import "github.com/rivo/tview"

type StatusBarParams struct {
	DefaultStatus   bool
	HasTargets      bool
	HasUpSignals    bool
	HasDownSignals  bool
	HasAntennaSpec  bool
	LastUpdated     string
	LastError       string
	ConnStatus      string // "connected", "degraded", "disconnected"
	SignalChanges   []string
	RefreshInterval string
}

func GetTheme(name string) *Theme {
	if t, ok := Themes[name]; ok {
		return t
	}
	return DarkTheme
}

type UI struct {
	theme *Theme
	dishesList      *tview.List
	stationsList    *tview.TextView
	selectedStation *tview.TextView
	pages           *tview.Pages
	uiDetails       *uiDetails
	jsonPreview     *tview.TextArea
	dishSpecsView   *uiDishSpec
	aboutView       *uiAbout
	statusBarFlex   *tview.Flex
	statusBarLeft   *tview.TextView
	statusBarCenter *tview.TextView
	statusBarRight  *tview.TextView
	targetsView     *tview.Flex
	upSignalsView   *tview.Flex
	downSignalsView *tview.Flex
	compactTable    *tview.Table
	mainContent     *tview.Flex
	compactContent  *tview.Flex
	layoutContainer *tview.Flex
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
	content *tview.TextView
}

type CompactRow struct {
	Station    string
	Dish       string
	Target     string
	UpSignals  int
	DownSignals int
	Activity   string
}
