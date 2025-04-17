package handler

import (
	"api-gw/functions"
	"api-gw/json"
	"api-gw/logging"
	"fmt"
	"net/http"
	"time"
)

func (h *APIxOBJ_Handler) QueryMetrics(w http.ResponseWriter, r *http.Request) {

	tStart := time.Now()

	h.ClientAddr = functions.RemoteHost(r)

	totalRequests++

	user := r.Context().Value("user")
	if user == nil {
		user = "N/A"
	}

	h.User = user.(string)

	if h.LimitRate(w, tStart) {

		return
	}

	rsp := APIxOBJ_Json_Rsp{}

	defer func() {

		rsp.ExecTime = functions.ExecTime(tStart)
		rsp.Updated = h.Stats.Time

		json.Write(w, rsp, http.StatusCreated)

		h.Log.Add(logging.Entry{
			Event: fmt.Sprintf("Client: %s - User: %s | Response Time: %s", h.ClientAddr, user, rsp.ExecTime),
			Code:  logging.CONST_CODE_INFO,
		})

	}()

	rsp = APIxOBJ_Json_Rsp{
		Metrics: APIxOBJ_Json_Metrics{
			Stats: *h.Stats,
			Info:  *h.Info,
		},
		Timestamp:      time.Now(),
		ClientRequests: client[h.ClientAddr].GetRequests(),
		TotalRequests:  totalRequests,
		TotalClients:   len(client),
	}
}

func (h *APIxOBJ_Handler) LimitRate(w http.ResponseWriter, tStart time.Time) bool {

	val, res := client[h.ClientAddr]

	if res {

		if (h.RateLimit.BurstRate > 0 && h.RateLimit.Seconds > 0) && (time.Since(val.Stamp) <= time.Duration(h.RateLimit.Seconds)*time.Second) {

			client[h.ClientAddr].Update(ClientxOBJ{
				Stamp:          tStart,
				Requests:       val.Requests + 1,
				ClientRequests: val.ClientRequests + 1,
			})

			if val.Requests > h.RateLimit.BurstRate {

				json.Write(w, APIxOBJ_Json_Status_Rsp{
					Status: APIxOBJ_Json_Status{
						Err: true,
						Msg: "Sorry .. Rate limit exceeded .. Let's cooldown a bit ..",
					},
				}, http.StatusTooManyRequests)

				h.Log.Add(logging.Entry{
					Event: fmt.Sprintf("Client: %s - User: %s | Rate limit exceeded", h.ClientAddr, h.User),
					Code:  logging.CONST_CODE_WARNING,
				})

				client[h.ClientAddr].Update(ClientxOBJ{
					Stamp:          tStart,
					Requests:       val.Requests,
					ClientRequests: val.ClientRequests,
				})

				return true
			}
		} else {

			client[h.ClientAddr].Update(ClientxOBJ{
				Stamp:          tStart,
				Requests:       1,
				ClientRequests: val.ClientRequests + 1,
			})
		}
	} else {

		client.Add(ClientxOBJ{
			ClientAddr:     h.ClientAddr,
			Stamp:          tStart,
			Requests:       1,
			ClientRequests: 1,
		})
	}

	return false
}
