package route

import (
	"api-gw/config"
	"net/http"
)

type MuxObj struct {
	Mux    *http.ServeMux
	Config config.CFG
}

type HandlerFuncCustom func(http.ResponseWriter, *http.Request)

const (
	CONST_HTTP_CODE_APP_STARTING   = http.StatusServiceUnavailable
	CONST_HTTP_CODE_APP_RESTARTING = http.StatusGatewayTimeout
)
