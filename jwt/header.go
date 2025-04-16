package jwt

import "api-gw/logging"

type APIxOBJ_Json_Status struct {
	Err bool   `json:"Error"`
	Msg string `json:"Message,omitempty"`
}

type APIxOBJ_Json_Status_Rsp struct {
	Status APIxOBJ_Json_Status `json:"Status"`
}

var appLog logging.Log

var secret = []byte("test123321test")
