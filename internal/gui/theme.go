package gui

import "github.com/gdamore/tcell/v2"

type Theme struct {
	Name            string
	Primary         string // selected items, highlights
	Secondary       string // data values
	Error           string // errors, inactive
	SignalUp        string // up signal color
	SignalDown      string // down signal color
	Accent          string // modal borders, titles
	Inactive        string // inactive items
	StatusConnected string // connection indicator: connected
	StatusDegraded  string // connection indicator: degraded
	StatusError     string // connection indicator: disconnected
	ModalBorder     tcell.Color
	ModalTitle      tcell.Color
	ErrorBorder     tcell.Color
}

var Themes = map[string]*Theme{
	"dark":      DarkTheme,
	"solarized": SolarizedTheme,
	"nord":      NordTheme,
	"light":     LightTheme,
}

var ThemeNames = []string{"dark", "solarized", "nord", "light"}

var DarkTheme = &Theme{
	Name:            "dark",
	Primary:         "green",
	Secondary:       "yellow",
	Error:           "red",
	SignalUp:        "red",
	SignalDown:      "green",
	Accent:          "yellow",
	Inactive:        "red",
	StatusConnected: "green",
	StatusDegraded:  "yellow",
	StatusError:     "red",
	ModalBorder:     tcell.ColorYellow,
	ModalTitle:      tcell.ColorYellow,
	ErrorBorder:     tcell.ColorRed,
}

var SolarizedTheme = &Theme{
	Name:            "solarized",
	Primary:         "#268bd2",
	Secondary:       "#b58900",
	Error:           "#dc322f",
	SignalUp:        "#cb4b16",
	SignalDown:      "#859900",
	Accent:          "#2aa198",
	Inactive:        "#586e75",
	StatusConnected: "#859900",
	StatusDegraded:  "#b58900",
	StatusError:     "#dc322f",
	ModalBorder:     tcell.NewRGBColor(42, 161, 152),
	ModalTitle:      tcell.NewRGBColor(42, 161, 152),
	ErrorBorder:     tcell.NewRGBColor(220, 50, 47),
}

var NordTheme = &Theme{
	Name:            "nord",
	Primary:         "#88c0d0",
	Secondary:       "#ebcb8b",
	Error:           "#bf616a",
	SignalUp:        "#d08770",
	SignalDown:      "#a3be8c",
	Accent:          "#81a1c1",
	Inactive:        "#4c566a",
	StatusConnected: "#a3be8c",
	StatusDegraded:  "#ebcb8b",
	StatusError:     "#bf616a",
	ModalBorder:     tcell.NewRGBColor(129, 161, 193),
	ModalTitle:      tcell.NewRGBColor(129, 161, 193),
	ErrorBorder:     tcell.NewRGBColor(191, 97, 106),
}

var LightTheme = &Theme{
	Name:            "light",
	Primary:         "#0087af",
	Secondary:       "#8700af",
	Error:           "#af0000",
	SignalUp:        "#af5f00",
	SignalDown:      "#005f00",
	Accent:          "#005faf",
	Inactive:        "#808080",
	StatusConnected: "#005f00",
	StatusDegraded:  "#af5f00",
	StatusError:     "#af0000",
	ModalBorder:     tcell.NewRGBColor(0, 95, 175),
	ModalTitle:      tcell.NewRGBColor(0, 95, 175),
	ErrorBorder:     tcell.NewRGBColor(175, 0, 0),
}
