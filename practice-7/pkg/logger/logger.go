package logger

import (
	"log"
	"os"
)

type Interface interface {
	Info(args ...interface{})
	Error(args ...interface{})
	Debug(args ...interface{})
	Warn(args ...interface{})
}

type Logger struct {
	info  *log.Logger
	erro  *log.Logger
	debug *log.Logger
	warn  *log.Logger
}

func New() *Logger {
	return &Logger{
		info:  log.New(os.Stdout, "[INFO]  ", log.LstdFlags),
		erro:  log.New(os.Stderr, "[ERROR] ", log.LstdFlags),
		debug: log.New(os.Stdout, "[DEBUG] ", log.LstdFlags),
		warn:  log.New(os.Stdout, "[WARN]  ", log.LstdFlags),
	}
}

func (l *Logger) Info(args ...interface{})  { l.info.Println(args...) }
func (l *Logger) Error(args ...interface{}) { l.erro.Println(args...) }
func (l *Logger) Debug(args ...interface{}) { l.debug.Println(args...) }
func (l *Logger) Warn(args ...interface{})  { l.warn.Println(args...) }
