package handler

import (
	"api-gw/cron"
	"api-gw/stats"
	"context"
)

type APIxOBJ_Handler struct {
	Stats      *stats.StatsObj
	Info       *stats.InfoObj
	Cron       *cron.Init
	Ctx        context.Context
	ClientAddr string
}

type APIxOBJ_Json_Rsp struct {
	ExecTime string         `json:"ExecTime"`
	Stats    stats.StatsObj `json:"Stats"`
	Info     stats.InfoObj  `json:"Info"`
}
