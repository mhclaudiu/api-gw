package main

import (
	"fmt"

	"api-gw/logging"
)

var app App

func main() {

	app.Config.Load("api-gw.conf")

	app.Log.Add(logging.Entry{
		Event: fmt.Sprintf("Starting %s Worker '%s:%d' | Env: %s | Version: %s", app.Config.APP.Name, app.Config.API.Host, app.Config.API.Port, app.Config.APP.Env, APP_VERSION),
	})

	app.InitMaps()

	app.Info.Init()

	app.Stats.Query(app.Info)

	app.StartWebAPI()

	app.LoadCRONs()

	app.Hook()
}
