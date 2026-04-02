package tui

import (
	"time"

	"github.com/RustyDaemon/go-dsn-now/internal/model/response"
)

// Data lifecycle messages
type DSNConfigLoadedMsg struct{ Config response.DSNConfig }
type DSNConfigErrorMsg struct{ Err error }
type DSNDataLoadedMsg struct{ Data response.DSN }
type DSNDataErrorMsg struct{ Err error }

// Timer messages
type TickRefreshMsg time.Time
type TickClockMsg time.Time

// Action result messages
type CopyResultMsg struct{ Err error }
type StatusMessageExpiredMsg struct{}
