package gui

import "github.com/rivo/tview"

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
