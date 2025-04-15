package route

import (
	"api-gw/cron"
	"api-gw/handler"
	"api-gw/stats"
	"net/http"
)

func (mux MuxObj) Register(stats *stats.StatsObj, info *stats.InfoObj, cron *cron.Init) {

	handler := handler.APIxOBJ_Handler{
		Stats: stats,
		Info:  info,
		Cron:  cron,
	}

	mux.RegisterMetrics(handler)
}

func (mux MuxObj) RegisterMetrics(h handler.APIxOBJ_Handler) {

	go mux.HandleFuncCustom(mux.Config.API.MainPath+mux.Config.API.MetricsPath, func(w http.ResponseWriter, r *http.Request) {
		h.QueryMetrics(w, r)
	})

}
