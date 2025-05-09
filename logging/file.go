package logging

import (
	"fmt"
	"time"

	"github.com/natefinch/lumberjack"
)

func (data *FILExOBJ) New() {

	fileLogger = &lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/%s.json", data.Path, time.Now().Format("2006-01-02")),
		MaxSize:    data.MaxSize, // Megabytes
		MaxBackups: 1,
		MaxAge:     data.MaxDays, // Days
		Compress:   false,        // Enable compression
	}

}

func (l *FILExOBJ) Write(data string) error {

	if fileLogger == nil {

		l.New()
	}

	_, err := fileLogger.Write([]byte(fmt.Sprintf("%s\n", data)))

	return err
}
