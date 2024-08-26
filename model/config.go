package model

import "time"

type APIAlertsConfig struct {
	Logging bool
	Timeout time.Duration
	Debug   bool
}
