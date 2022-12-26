package logger

import (
	"io"
	"log"

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

	ErrorTag = "[ERRO] "
	WarnTag  = "[WARN] "
	InfoTag  = "[INFO] "
	DebugTag = "[DEBU] "
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

type Logger struct {
	logg  *log.Logger
	level string
}

func New(cfg config.LoggerConf, out io.Writer) *Logger {
	logg := log.New(out, "", log.Ldate|log.Lmsgprefix|log.Lmicroseconds)
	return &Logger{logg: logg, level: cfg.Level}
}

func (l Logger) Error(msg string) {
	l.logg.SetPrefix(ErrorTag)
	l.logg.Println(msg)
}

func (l Logger) Warn(msg string) {
	if levelNum(l.level) >= Warn {
		l.logg.SetPrefix(WarnTag)
		l.logg.Println(msg)
	}
}

func (l Logger) Info(msg string) {
	if levelNum(l.level) >= Info {
		l.logg.SetPrefix(InfoTag)
		l.logg.Println(msg)
	}
}

func (l Logger) Debug(msg string) {
	if levelNum(l.level) >= Debug {
		l.logg.SetPrefix(DebugTag)
		l.logg.Println(msg)
	}
}
