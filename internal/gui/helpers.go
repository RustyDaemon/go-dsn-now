package gui

import (
	"fmt"

	"github.com/rivo/tview"
)

func DefaultIfEmpty(value, defaultValue string) string {
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func DashIfEmpty(value string) string {
	if len(value) == 0 {
		return "-"
	}
	return value
}

func NewTextView(text string) *tview.TextView {
	textView := tview.NewTextView()
	textView.SetText(text)
	textView.SetDynamicColors(true)

	return textView
}

func wrapInModal(p tview.Primitive, hPad, vPad, hWeight, vWeight int) tview.Primitive {
	return tview.NewFlex().
		AddItem(nil, hPad, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, vPad, vWeight, false).
			AddItem(p, 0, 5, true).
			AddItem(nil, vPad, vWeight, false),
			0, hWeight, true).
		AddItem(nil, hPad, 1, false)
}

func setSpecField(view *tview.TextView, value string, color string) {
	view.SetText(fmt.Sprintf("[%s]%s[-]", color, DashIfEmpty(value)))
}
