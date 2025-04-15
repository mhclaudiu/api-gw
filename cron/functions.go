package cron

import (
	"strconv"
	"strings"
	"time"
)

func GetDuration(period string) (error, time.Duration) {

	period = strings.Trim(period, " ")

	if len(period) < 1 || period[0] == '0' {

		return nil, time.Duration(0)
	}

	buff := period[:len(period)-1]

	atoi, err := strconv.Atoi(buff)

	if err != nil {

		return err, 0
	}

	switch period[len(period)-1:] {

	case "m":

		return nil, time.Duration(atoi * 60)

	case "h":

		return nil, time.Duration(atoi * 3600)

	}

	return nil, time.Duration(atoi)
}

func GetSeconds(period string) (error, int) {

	period = strings.Trim(period, " ")

	if len(period) < 1 || period[0] == '0' {

		return nil, 0
	}

	buff := period[:len(period)-1]

	atoi, err := strconv.Atoi(buff)

	if err != nil {

		return err, 0
	}

	switch period[len(period)-1:] {

	case "m":

		return nil, atoi * 60

	case "h":

		return nil, atoi * 3600

	}

	return nil, atoi
}

func GetMiliSeconds(period string) (error, float64) {

	period = strings.Trim(period, " ")

	if len(period) < 1 || period[0] == '0' {

		return nil, 0
	}

	buff := period[:len(period)-1]

	ftoi, err := strconv.ParseFloat(buff, 64)

	if err != nil {

		return err, 0
	}

	ftoi *= 1000

	switch period[len(period)-1:] {

	case "m":

		return nil, ftoi * 60

	case "h":

		return nil, ftoi * 3600

	}

	return nil, ftoi
}
