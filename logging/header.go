package logging

import (
	"context"
	"sync"

	"github.com/natefinch/lumberjack"
)

const (
	CONST_CODE_INFO    = 100
	CONST_CODE_WARNING = 200
	CONST_CODE_ERROR   = 300

	CONST_LOG_ADD_CONV_ERR = "SYS Conv Err"

	CONST_LOG_EVENTID = "Oops .. Something went wrong .. Refference Error ID '%s'"

	CONST_ERROR_ID_HASH_LEN = 20
)

var Mu sync.Mutex

var WG sync.WaitGroup

type Log struct {
	UserID        *int64
	IP            string
	Hash          string
	Ctx           context.Context
	ErrorIDLength *int
}
type Entry struct {
	Event   string
	Code    int
	Err     error
	Trigger string
	EventID error
	Exit    bool
}

type FILExOBJ struct {
	MaxSize int
	MaxDays int
	Path    string
	Enabled bool
}

var fileLogger *lumberjack.Logger
