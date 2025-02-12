package logger

import (
	"log"
	"os"
)

type Logger struct {
	*log.Logger
}

func NewLogger() *Logger {
	return &Logger{
		Logger: log.New(os.Stderr, "", log.LstdFlags),
	}
}

func (l *Logger) Info(v ...interface{}) {
	l.Println(append([]interface{}{"[INFO]"}, v...)...)
}

func (l *Logger) Error(v ...interface{}) {
	l.Println(append([]interface{}{"[ERROR]"}, v...)...)
}

func (l *Logger) Fatal(v ...interface{}) {
	l.Println(append([]interface{}{"[FATAL]"}, v...)...)
	os.Exit(1)
}
