package logger

import (
	"fmt"
	"io"
	"time"
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

func New(level string, out io.Writer) *Logger {
	return &Logger{level: level, out: out}
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
