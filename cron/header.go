package cron

import (
	"api-gw/logging"
)

type Cron struct {
	Name      string
	Func      func() error
	Running   bool
	Finished  bool
	IntervalF float64
	IntervalS string
	AutoStart bool
}

var Crons map[string]*Cron

var logObj logging.Log

type Init struct {
	ErrorIDLength *int
}
