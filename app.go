package main

import (
	"api-gw/config"
	"api-gw/cron"
	"api-gw/functions"
	"api-gw/logging"
	"api-gw/route"
	"api-gw/stats"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"github.com/rs/cors"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

const (
	APP_VERSION = "1.0.0"
)

type App struct {
	Cron   *cron.Init
	Stats  *stats.StatsObj
	Info   *stats.InfoObj
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

			fmt.Printf(" --> %s", colored("Closing and syncing resources .. please wait .."))

		} else {

			fmt.Printf(" --> %s", colored("Force exit detected .. please wait .."))
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

	muxObj.Register(app.Stats, app.Info, app.Cron)
}

func (app *App) LoadCRONs() {

	app.Cron = &cron.Init{}
	app.Cron.Init()

	app.Cron.Add("QueryStats", func() error { app.Stats.Query(app.Info); return nil }, "5s", true)
}

func (app *App) InitMaps() {

	app.Stats = &stats.StatsObj{}
	app.Info = &stats.InfoObj{}
}

func (app *App) InitInfo() {

	cpuInfo, _ := cpu.Info()

	memory, _ := mem.VirtualMemory()

	hostInfo, _ := host.Info()

	app.Info = &stats.InfoObj{

		CPU: stats.INFOxCPU{
			Cores: fmt.Sprint(cpuInfo[0].Cores),
			Freq:  fmt.Sprintf("%.0f MHz", cpuInfo[0].Mhz),
			Model: fmt.Sprint(cpuInfo[0].ModelName),
		},
		MEM: stats.INFOxMEM{
			Total: humanize.Bytes(memory.Total),
		},
		HOST: stats.INFOxHOST{
			Name:            hostInfo.Hostname,
			OS:              hostInfo.OS,
			Platform:        hostInfo.Platform,
			PlatformVersion: hostInfo.PlatformVersion,
			KernelVersion:   hostInfo.KernelVersion,
			KernelArch:      hostInfo.KernelArch,
			Uptime:          functions.FormatUptime(int(hostInfo.Uptime)),
			Proccesses:      fmt.Sprint(hostInfo.Procs),
		},
	}

	//app.Info.StartUptime()
}
