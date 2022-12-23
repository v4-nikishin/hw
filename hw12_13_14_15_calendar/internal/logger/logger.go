package logger

import (
	"fmt"
	"io"
	"time"

	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/config"
)

const (
	Error = iota
	Warn
	Info
	Debug

	ErrorStr = "error"
	WarnStr  = "warn"
	InfoStr  = "info"
	DebugStr = "debug"

	ErrorTag = "[ERROR]"
	WarnTag  = "[WARN]"
	InfoTag  = "[INFO]"
	DebugTag = "[DEBUG]"
)

func levelNum(level string) int {
	switch level {
	case ErrorStr:
		return Error
	case WarnStr:
		return Warn
	case InfoStr:
		return Info
	case DebugStr:
		return Debug
	}
	return Error
}

func timeStamp() string {
	return "[" + time.Now().String() + "]"
}

func tmpl(tag string, msg string) string {
	return timeStamp() + " " + tag + " " + msg
}

type Logger struct {
	level string
	out   io.Writer
}

func New(cfg config.LoggerConf, out io.Writer) *Logger {
	return &Logger{level: cfg.Level, out: out}
}

func (l Logger) Error(msg string) {
	fmt.Fprintln(l.out, tmpl(ErrorTag, msg))
}

func (l Logger) Warn(msg string) {
	if levelNum(l.level) >= Warn {
		fmt.Fprintln(l.out, tmpl(WarnTag, msg))
	}
}

func (l Logger) Info(msg string) {
	if levelNum(l.level) >= Info {
		fmt.Fprintln(l.out, tmpl(InfoTag, msg))
	}
}

func (l Logger) Debug(msg string) {
	if levelNum(l.level) >= Debug {
		fmt.Fprintln(l.out, tmpl(DebugTag, msg))
	}
}
