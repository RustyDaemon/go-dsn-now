package tui

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	CycleStation      key.Binding
	CycleTarget       key.Binding
	CycleUpSignal     key.Binding
	CycleDownSignal   key.Binding
	NavigateUp        key.Binding
	NavigateDown      key.Binding
	Bookmark          key.Binding
	CompactView       key.Binding
	CompactSort       key.Binding
	Clipboard         key.Binding
	CycleTheme        key.Binding
	CycleDistanceUnit key.Binding
	IncreaseRefresh   key.Binding
	DecreaseRefresh   key.Binding
	Help              key.Binding
	CloseModal        key.Binding
	Quit              key.Binding
}

var Keys = KeyMap{
	CycleStation: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "cycle station"),
	),
	CycleTarget: key.NewBinding(
		key.WithKeys("t"),
		key.WithHelp("t", "cycle target"),
	),
	CycleUpSignal: key.NewBinding(
		key.WithKeys("u"),
		key.WithHelp("u", "cycle up signal"),
	),
	CycleDownSignal: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "cycle down signal"),
	),
	NavigateUp: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "navigate up"),
	),
	NavigateDown: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "navigate down"),
	),
	Bookmark: key.NewBinding(
		key.WithKeys("b"),
		key.WithHelp("b", "bookmark dish"),
	),
	CompactView: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "compact view"),
	),
	CompactSort: key.NewBinding(
		key.WithKeys("S"),
		key.WithHelp("S", "cycle sort"),
	),
	Clipboard: key.NewBinding(
		key.WithKeys("y"),
		key.WithHelp("y", "copy to clipboard"),
	),
	CycleTheme: key.NewBinding(
		key.WithKeys("T"),
		key.WithHelp("T", "cycle theme"),
	),
	CycleDistanceUnit: key.NewBinding(
		key.WithKeys("U"),
		key.WithHelp("U", "cycle distance unit"),
	),
	IncreaseRefresh: key.NewBinding(
		key.WithKeys("+", "="),
		key.WithHelp("+", "increase refresh"),
	),
	DecreaseRefresh: key.NewBinding(
		key.WithKeys("-"),
		key.WithHelp("-", "decrease refresh"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	CloseModal: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "close"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	),
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.CycleStation, k.CycleTarget, k.CycleUpSignal, k.CycleDownSignal,
		k.Help, k.Quit,
	}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.CycleStation, k.CycleTarget, k.CycleUpSignal, k.CycleDownSignal},
		{k.NavigateUp, k.NavigateDown, k.Bookmark, k.CompactView},
		{k.CompactSort, k.Clipboard, k.CycleTheme, k.CycleDistanceUnit},
		{k.IncreaseRefresh, k.DecreaseRefresh, k.Help, k.CloseModal},
		{k.Quit},
	}
}
