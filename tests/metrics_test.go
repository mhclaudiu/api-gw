package main

import (
	"api-gw/config"
	"api-gw/handler"
	"api-gw/metrics"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

var cfg config.CFG

func TestMetrics(t *testing.T) {

	cfg.Load("../api-gw.conf")

	objHandler := handler.APIxOBJ_Handler{
		Stats:     &metrics.StatsObj{},
		Info:      &metrics.InfoObj{},
		RateLimit: handler.APIxOBJ_Config_Ratelimit(config.GetRateLimit(cfg.API.RateLimit)),
	}

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		objHandler.QueryMetrics(w, r)
	})

	t.Run("Metrics Query", func(t *testing.T) {

		req := httptest.NewRequest("GET", "/", nil)
		rr := httptest.NewRecorder()

		h.ServeHTTP(rr, req)

		var rsp map[string]interface{}
		if err := json.NewDecoder(rr.Body).Decode(&rsp); err != nil {
			t.Errorf("Error Decoding Response: %v", err)
		}

		status, ok := rsp["Status"].(map[string]interface{})

		if !ok {
			t.Errorf("Could not create Status Map")
		}

		/*if reflect.ValueOf(status["Error"]).Bool() {
			t.Errorf("Status Error")
		}*/
		if err, _ := status["Error"].(bool); err {
			t.Errorf("Expected: %v | Received: %v", !err, err)
		}

		if rr.Code != http.StatusCreated {
			t.Errorf("Expected: %d | Received: %v", http.StatusCreated, rr.Code)
		}
	})

	t.Run("Metrics Rate Limit", func(t *testing.T) {

		var req *http.Request
		var rr *httptest.ResponseRecorder

		for i := 0; i <= objHandler.RateLimit.BurstRate+2; i++ {

			req = httptest.NewRequest("GET", "/", nil)
			rr = httptest.NewRecorder()

			h.ServeHTTP(rr, req)
		}

		var rsp map[string]interface{}
		if err := json.NewDecoder(rr.Body).Decode(&rsp); err != nil {
			t.Errorf("Error Decoding Response: %v", err)
		}

		status, ok := rsp["Status"].(map[string]interface{})

		if !ok {
			t.Errorf("Could not create Status Map")
		}

		if err, _ := status["Error"].(bool); !err {
			t.Errorf("Expected: %v | Received: %v", !err, err)
		}

		if rr.Code != http.StatusTooManyRequests {
			t.Errorf("Expected: %d | Received: %v", http.StatusTooManyRequests, rr.Code)
		}
	})
}
