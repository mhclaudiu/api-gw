package handler

import (
	"api-gw/functions"
	"api-gw/json"
	"net/http"
	"time"
)

func (h APIxOBJ_Handler) QueryMetrics(w http.ResponseWriter, r *http.Request) {

	tStart := time.Now()

	h.ClientAddr = functions.RemoteHost(r)
	h.Ctx = r.Context()

	rsp := APIxOBJ_Json_Rsp{}

	defer func() {

		rsp.ExecTime = functions.ExecTime(tStart)

		json.Write(w, rsp)

	}()

	rsp = APIxOBJ_Json_Rsp{
		Stats: *h.Stats,
		Info:  *h.Info,
	}

}
