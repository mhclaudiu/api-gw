package config

import (
	"api-gw/functions"
	"log"
	"strconv"
	"strings"
)

func GetRateLimit(data string) CFGxAPI_Ratelimit {

	buff := strings.Split(data, ",")

	var burst, seconds int
	var err error

	if burst, err = strconv.Atoi(buff[0]); err != nil {
		log.Panic(err)
	}

	if err, seconds = functions.GetSeconds(buff[1]); err != nil {
		log.Panic(err)
	}

	return CFGxAPI_Ratelimit{
		BurstRate: burst,
		Seconds:   seconds,
	}

}
