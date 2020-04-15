package log

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"github.com/fatih/color"
)

var stdout io.Writer = color.Output
var stderr io.Writer = color.Error
var logf *os.File = nil
var mtx_log *sync.Mutex = &sync.Mutex{}
var enableOutput = true
var verboseLevel = INFO

const (
	DEBUG = iota
	INFO
	IMPORTANT
	WARNING
	ERROR
	FATAL
	SUCCESS
)

const (
	MOD_NONE = iota
	MOD_SUCCESS
)

var LogLabels = map[int]string{
	DEBUG:     "dbg",
	INFO:      "inf",
	IMPORTANT: "imp",
	WARNING:   "war",
	ERROR:     "err",
	FATAL:     "!!!",
}

func EnableOutput(enable bool) {
	enableOutput = enable
}

func SetOutput(o io.Writer) {
	stdout = o
}

func SetVerbosityLevel(lvl int) {
	verboseLevel = lvl
}

func SetLogFile(path string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	logf = f
	return nil
}

func NullLogger() *log.Logger {
	return log.New(ioutil.Discard, "", 0)
}

func Debug(format string, args ...interface{}) {
	mtx_log.Lock()
	defer mtx_log.Unlock()

	log_message(DEBUG, MOD_NONE, format, args...)
}

func Info(format string, args ...interface{}) {
	mtx_log.Lock()
	defer mtx_log.Unlock()

	log_message(INFO, MOD_NONE, format, args...)
}

func Important(format string, args ...interface{}) {
	mtx_log.Lock()
	defer mtx_log.Unlock()

	log_message(IMPORTANT, MOD_NONE, format, args...)
}

func Warning(format string, args ...interface{}) {
	mtx_log.Lock()
	defer mtx_log.Unlock()

	log_message(WARNING, MOD_NONE, format, args...)
}

func Error(format string, args ...interface{}) {
	mtx_log.Lock()
	defer mtx_log.Unlock()

	log_message(ERROR, MOD_NONE, format, args...)
}

func Fatal(format string, args ...interface{}) {
	mtx_log.Lock()
	defer mtx_log.Unlock()

	log_message(FATAL, MOD_NONE, format, args...)
}

func Success(format string, args ...interface{}) {
	mtx_log.Lock()
	defer mtx_log.Unlock()

	log_message(INFO, MOD_SUCCESS, format, args...)
}

func format_msg(lvl int, mod int, use_color bool, format string, args ...interface{}) string {
	var sign, msg *color.Color
	label := LogLabels[lvl]
	switch lvl {
	case DEBUG:
		sign = color.New(color.FgBlack, color.BgHiBlack)
		msg = color.New(color.Reset, color.FgHiBlack)
	case INFO:
		sign = color.New(color.FgGreen, color.BgBlack)
		msg = color.New(color.Reset)
	case IMPORTANT:
		sign = color.New(color.FgWhite, color.BgHiBlue)
		msg = color.New(color.Reset)
	case WARNING:
		sign = color.New(color.FgBlack, color.BgYellow)
		msg = color.New(color.Reset)
	case ERROR:
		sign = color.New(color.FgWhite, color.BgRed)
		msg = color.New(color.Reset, color.FgRed)
	case FATAL:
		sign = color.New(color.FgBlack, color.BgRed)
		msg = color.New(color.Reset, color.FgRed, color.Bold)
	}

	if mod > MOD_NONE {
		switch mod {
		case MOD_SUCCESS:
			sign = color.New(color.FgWhite, color.BgGreen)
			msg = color.New(color.Reset, color.FgGreen)
			label = "+++"
		}
	}

	if use_color {
		return "[" + sign.Sprintf("%s", label) + "] " + msg.Sprintf(format, args...)
	} else {
		return "[" + fmt.Sprintf("%s", label) + "] " + fmt.Sprintf(format, args...)
	}
}

func log_message(lvl int, mod int, format string, args ...interface{}) {
	if lvl < verboseLevel {
		return
	}

	t := time.Now()
	time_clr := color.New(color.Reset)

	output := "[" + time_clr.Sprintf("%d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second()) + "] " + format_msg(lvl, mod, true, format, args...)
	log_output := "[" + fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second()) + "] " + format_msg(lvl, mod, false, format, args...)

	if enableOutput {
		fmt.Fprint(stdout, output+"\n")
	}
	if logf != nil {
		logf.WriteString(log_output + "\n")
		logf.Sync()
	}
}
