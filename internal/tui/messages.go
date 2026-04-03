package tui

import (
	"time"

	"github.com/RustyDaemon/go-dsn-now/internal/model/response"
)

type DSNConfigLoadedMsg struct{ Config response.DSNConfig }
type DSNConfigErrorMsg struct{ Err error }
type DSNDataLoadedMsg struct{ Data response.DSN }
type DSNDataErrorMsg struct{ Err error }

type TickRefreshMsg time.Time
type TickClockMsg time.Time

type CopyResultMsg struct{ Err error }
type StatusMessageExpiredMsg struct{}
