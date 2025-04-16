package route

import (
	"api-gw/config"
	"api-gw/handler"
	"api-gw/jwt"
	"api-gw/metrics"
	"net/http"
)

func (mux MuxObj) Register(stats *metrics.StatsObj, info *metrics.InfoObj) {

	handler := handler.APIxOBJ_Handler{
		Stats:     stats,
		Info:      info,
		RateLimit: handler.APIxOBJ_Config_Ratelimit(config.GetRateLimit(mux.Config.API.RateLimit)),
	}

	mux.RegisterMetrics(handler)
}

func (mux MuxObj) RegisterMetrics(h handler.APIxOBJ_Handler) {

	go mux.HandleFuncCustom(mux.Config.API.Path, jwt.Auth(func(w http.ResponseWriter, r *http.Request) {
		h.QueryMetrics(w, r)
	}, mux.Config.API.Auth))

}
