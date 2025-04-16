package handler

import (
	"api-gw/logging"
	"api-gw/metrics"
	"sync"
	"time"
)

type APIxOBJ_Config_Ratelimit struct {
	BurstRate int
	Seconds   int
}
type APIxOBJ_Handler struct {
	Stats      *metrics.StatsObj
	Info       *metrics.InfoObj
	ClientAddr string
	Log        logging.Log
	RateLimit  APIxOBJ_Config_Ratelimit
}

type APIxOBJ_Json_Metrics struct {
	Stats metrics.StatsObj `json:"Stats"`
	Info  metrics.InfoObj  `json:"Info"`
}

type APIxOBJ_Json_Status struct {
	Err bool   `json:"Error"`
	Msg string `json:"Message,omitempty"`
}

type APIxOBJ_Json_Rsp struct {
	Status         APIxOBJ_Json_Status  `json:"Status"`
	ExecTime       string               `json:"ResponseTime"`
	Metrics        APIxOBJ_Json_Metrics `json:"Metrics"`
	Updated        string               `json:"Updated"`
	Timestamp      time.Time            `json:"TimeStamp"`
	ClientRequests int                  `json:"CurrentClientRequests"`
	TotalRequests  int                  `json:"TotalRequests"`
	TotalClients   int                  `json:"TotalClients"`
}

type APIxOBJ_Json_Status_Rsp struct {
	Status APIxOBJ_Json_Status `json:"Status"`
}

var totalRequests int

var Mu sync.RWMutex

type ClientxOBJ struct {
	Stamp          time.Time
	Requests       int
	ClientAddr     string
	ClientRequests int
}

type ClientMap map[string]*ClientxOBJ

var client = make(ClientMap)
