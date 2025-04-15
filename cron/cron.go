package cron

import (
	"fmt"
	"strings"
	"time"

	"api-gw/functions"
	"api-gw/logging"
)

func (ci *Init) Init() {

	Crons = make(map[string]*Cron)

	logObj = logging.Log{
		ErrorIDLength: ci.ErrorIDLength,
	}
}

func (ci *Init) Add(name string, f func() error, period string, autoStart bool) {

	name = strings.TrimSpace(name)
	period = strings.TrimSpace(period)

	found, _ := ci.Exists(name)

	if found {

		logObj.Add(logging.Entry{
			Event: fmt.Sprintf("'%s' already exists - every '%s'", name, period),
		})

		ci.Update(name, period)

		return
	}

	err, f64 := GetMiliSeconds(period)

	if err != nil {

		logObj.Add(logging.Entry{
			Event: fmt.Sprintf("'%s' - '%s' | ADD Duration Err - '%v'", name, period, err),
		})

		return
	}

	cron := &Cron{Name: name,
		Func:      f,
		IntervalF: f64,
		IntervalS: period,
		AutoStart: autoStart,
	}

	Crons[name] = cron

	if f64 < 1 {

		logObj.Add(logging.Entry{
			Event: fmt.Sprintf("Adding '%s'", name),
		})

		return
	}

	buff := "and starting "

	if !autoStart {

		buff = ""
	}

	logObj.Add(logging.Entry{
		Event: fmt.Sprintf("Adding %s'%s' at every '%s'", buff, name, period),
	})

	if autoStart {

		cron.Start(ci)
	}
}

func (cron *Cron) Start(ci *Init) {

	if cron.IsRunning() {

		logObj.Add(logging.Entry{
			Event: fmt.Sprintf("'%s' already running - every '%s'", cron.Name, cron.IntervalS),
			Code:  functions.PointerTo(logging.CONST_CODE_WARNING),
		})

		return
	}

	if !cron.AutoStart {

		logObj.Add(logging.Entry{
			Event: fmt.Sprintf("Starting '%s' at every '%s'", cron.Name, cron.IntervalS),
		})

	}

	cron.Running = true
	cron.Finished = true

	go cron.Loop(ci)
}

func (cron *Cron) Loop(ci *Init) {

	for {

		time.Sleep(time.Duration(cron.IntervalF) * time.Millisecond)

		if !cron.Running {

			return
		}

		if !cron.Finished {

			logObj.Add(logging.Entry{
				Event: fmt.Sprintf("'%s' did not yet finished execution ..", cron.Name),
				Code:  functions.PointerTo(logging.CONST_CODE_WARNING),
			})

			return
		}

		//logging.Log.Printf("%v - %v", c.Name, c.Running)
		go cron.Exec(ci)
	}
}

func (cron *Cron) Exec(ci *Init) {

	cron.Finished = false

	err := cron.Func()

	if err != nil {

		logObj.Add(logging.Entry{
			Err:  err,
			Code: functions.PointerTo(logging.CONST_CODE_ERROR),
		})
	}

	cron.Finished = true
}

func (cron *Cron) Stop(ci *Init) {

	if !cron.IsRunning() {

		logObj.Add(logging.Entry{
			Event: fmt.Sprintf("'%s' already stopped ..", cron.Name),
			Code:  functions.PointerTo(logging.CONST_CODE_WARNING),
		})

		return
	}

	logObj.Add(logging.Entry{
		Event: fmt.Sprintf("Stopping '%s'", cron.Name),
	})

	cron.Running = false
	cron.Finished = true
}

func (ci *Init) Update(name, period string) {

	found, cron := ci.Exists(name)

	if !found {

		logObj.Add(logging.Entry{
			Event: fmt.Sprintf("'%s' does not exist..", name),
			Code:  functions.PointerTo(logging.CONST_CODE_WARNING),
		})

		return
	}

	err, f64 := GetMiliSeconds(period)

	if err != nil {

		logObj.Add(logging.Entry{
			Event: fmt.Sprintf("UPDATE Duration Err - '%v'", err),
		})

		return
	}

	if f64 < 1 {

		cron.Stop(ci)

		cron.IntervalF = f64
		cron.IntervalS = period

		return
	}

	if cron.IntervalF != f64 {

		logObj.Add(logging.Entry{
			Event: fmt.Sprintf("Updating '%s' from '%s' to '%s'", cron.Name, cron.IntervalS, period),
			Code:  functions.PointerTo(logging.CONST_CODE_INFO),
		})

		cron.IntervalF = f64
		cron.IntervalS = period
	}

	if !cron.IsRunning() {

		logObj.Add(logging.Entry{
			Event: fmt.Sprintf("Starting '%s' at every '%s'", cron.Name, cron.IntervalS),
		})

		cron.Start(ci)
	}
}

func (cron *Cron) IsRunning() bool {

	return cron.Running
}

func (ci *Init) Exists(name string) (bool, *Cron) {

	if cron, ok := Crons[name]; ok {

		return true, cron
	}

	return false, &Cron{}
}
