package tui

import "github.com/RustyDaemon/go-dsn-now/internal/tui/components"

type compactSortMode int

const (
	compactSortDefault compactSortMode = iota
	compactSortByActivity
	compactSortBySignalCount
	compactSortByTarget
)

func (m *compactSortMode) cycle() {
	*m = (*m + 1) % 4
}

func (m compactSortMode) label() string {
	switch m {
	case compactSortByActivity:
		return "Activity"
	case compactSortBySignalCount:
		return "Signals"
	case compactSortByTarget:
		return "Target"
	default:
		return "Default"
	}
}

type selectionState struct {
	station    int
	dish       int
	target     int
	upSignal   int
	downSignal int
}

type uiState struct {
	width  int
	height int

	ready         bool
	compactView   bool
	compactSort   compactSortMode
	activeModal   components.ModalType
	statusMessage string
	selection     selectionState
}

func newUIState() uiState {
	return uiState{
		activeModal: components.ModalNone,
		selection: selectionState{
			station: -1,
		},
	}
}
