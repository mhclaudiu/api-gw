package main

import (
	"fmt"

	"api-gw/logging"
	"api-gw/metrics"
)

var app = App{
	Stats: &metrics.StatsObj{},
	Info:  &metrics.InfoObj{},
}

func main() {

	app.Config.Load("api-gw.conf")

	app.Log.Add(logging.Entry{
		Event: fmt.Sprintf("Starting %s Worker '%s:%d' | Env: %s | Version: %s", app.Config.APP.Name, app.Config.API.Host, app.Config.API.Port, app.Config.APP.Env, APP_VERSION),
	})

	app.GenerateTestToken()

	//app.Metrics.Init()
	app.Info.StartUptime()

	app.Info.Init()

	app.Stats.Query(nil)

	app.StartWebAPI()

	app.LoadCRONs()

	app.Hook()
}
