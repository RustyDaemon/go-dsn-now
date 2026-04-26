package tui

func (m Model) leftPanelWidth() int {
	width := m.ui.width / 4
	if width < 20 {
		width = 20
	}
	return width
}

func (m Model) mainHeight() int {
	return m.ui.height - 3
}

func (m *Model) resizeComponents() {
	leftWidth := m.leftPanelWidth()
	mainHeight := m.mainHeight()

	stationCount := len(m.appData.FullData.Stations)
	if stationCount == 0 {
		stationCount = 3
	}

	dishHeight := mainHeight - (stationCount + 2) - 2
	if dishHeight < 3 {
		dishHeight = 3
	}

	m.dishList.SetSize(leftWidth, dishHeight)
	m.stationBar.SetWidth(leftWidth)
	m.statusBar.SetWidth(m.ui.width)
	m.compactTable.SetSize(m.ui.width, mainHeight)
	m.modal.SetSize(m.ui.width, m.ui.height-3)
}
