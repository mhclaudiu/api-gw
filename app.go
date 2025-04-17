package main

import (
	"api-gw/config"
	"api-gw/cron"
	"api-gw/functions"
	"api-gw/jwt"
	"api-gw/logging"
	"api-gw/metrics"
	"api-gw/route"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/fatih/color"
	"github.com/rs/cors"
)

const (
	APP_VERSION = "1.0.0"
)

type App struct {
	Cron   *cron.Init
	Stats  *metrics.StatsObj
	Info   *metrics.InfoObj
	Log    logging.Log
	Config config.CFG
}

func (app *App) Hook() {

	colored := color.New(color.FgHiGreen, color.Italic).SprintFunc()

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	//signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	signal.Notify(sigs, os.Interrupt, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSEGV)

	fmt.Println()

	err, exitSyncTimer := functions.GetSeconds(app.Config.APP.ExitSyncTimer)
	if err != nil {

		fmt.Printf(" --> %s", colored(err))
	}

	go func() {

		sig := <-sigs

		fmt.Printf(" %v\n\n", sig)
		// cleanup code here
		done <- true

		if exitSyncTimer > 0 {

			fmt.Printf(" --> %s", colored("Closing and syncing resources .. Please wait .."))

		} else {

			fmt.Printf(" --> %s", colored("Force exit detected .. Please wait .."))
		}

		fmt.Println()

	}()

	//fmt.Println("awaiting signal")
	<-done

	if exitSyncTimer > 0 {

		functions.TimerWait(exitSyncTimer)
	}

	fmt.Println("\nBye!\n")
}

func (app *App) StartWebAPI() {

	cors := cors.New(cors.Options{
		AllowedOrigins: []string{app.Config.API.CorsFilter},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
	})

	//mux := http.NewServeMux()
	muxObj := route.MuxObj{
		Mux:    http.NewServeMux(),
		Config: app.Config,
	}

	go func() {

		log.Panic(http.ListenAndServe(fmt.Sprintf("%s:%d", app.Config.API.Host, app.Config.API.Port), cors.Handler(muxObj.Mux)))

	}()

	muxObj.Register(app.Stats, app.Info)
}

func (app *App) LoadCRONs() {

	app.Cron = &cron.Init{}
	app.Cron.Init()

	app.Cron.Add("QueryStats", func() error { app.Stats.Query(app.Info); return nil }, "5s", true)
}

func (app *App) MetricsInit() {

	app.Info = &metrics.InfoObj{}
	app.Stats = &metrics.StatsObj{}
}

func (app *App) GenerateTestToken() {

	if app.Config.API.Auth {

		token, err := jwt.CreateToken()

		if err != nil {

			app.Log.Add(logging.Entry{
				Event: fmt.Sprintf("Test Token Error: %s", err.Error()),
			})

		} else {
			app.Log.Add(logging.Entry{
				Event: fmt.Sprintf("Generated Test Token: %s", token),
			})
		}
	}
}
