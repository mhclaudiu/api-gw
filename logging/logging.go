package logging

import (
	"api-gw/functions"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/fatih/color"
)

func (l *Log) Add(data Entry) error {

	if data.Err != nil {

		data.Event = data.Err.Error()

		data.Code = functions.PointerTo(CONST_CODE_ERROR)
	}

	if (len(data.Event)) < 2 {

		return nil
	}

	_, file, lineNo, ok := runtime.Caller(1)

	if !ok {

		log.Println("[LOG::ERR] Runtime.Caller() failed ..")
		return nil
	}

	fileName := path.Base(file)

	dir := filepath.Base(filepath.Dir(file))
	buff := strings.TrimSuffix(fileName, filepath.Ext(fileName))

	if dir != buff {

		buff = fmt.Sprintf("%s-%s", dir, buff)
	}

	data.Trigger = strings.ToUpper(fmt.Sprintf("%s:%d", buff, lineNo))

	buffIP := l.IP
	if len(l.IP) > 2 {
		buffIP = fmt.Sprintf(" | IP: '%s'", l.IP)
	}

	noCode := false

	if data.Code == nil {

		data.Code = functions.PointerTo(0)

		noCode = true
	}

	if eventID := l.ProcessEventID(&data.Event, *data.Code); eventID != nil {

		data.EventID = eventID
	}

	code := GetSeverity(*data.Code)

	/*if data.Code != nil {

		code = GetSeverity(*data.Code)
	}*/

	re, err := regexp.Compile(`error|Err|Error|SQL\sErr|SQL\sError`)
	if err != nil {

		log.Println(err)
		return nil
	}

	if code != CONST_CODE_ERROR && re.MatchString(data.Event) {

		code = CONST_CODE_ERROR
	}

	var colorType color.Attribute

	switch code {

	case CONST_CODE_ERROR:

		colorType = color.FgHiRed

	case CONST_CODE_WARNING:

		colorType = color.FgHiYellow

	case CONST_CODE_INFO:

		colorType = color.FgHiCyan
	}

	if noCode {

		colorType = color.FgHiBlack

	} else {

		if *data.Code == 0 {

			colorType = color.FgHiMagenta
		}
	}

	colored := color.New(colorType).SprintfFunc()

	log.Printf("%s", colored("[%s] %+v%s", data.Trigger, data.Event, buffIP))

	if data.Exit {

		os.Exit(0)
	}

	return data.EventID
}
