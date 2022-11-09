package main

import (
	"log"
)

type LoggerInterface interface {
	Debug(v ...any)
	Error(v ...any)
}

type CustomLogger struct {
}

func NewCustomLogger() *CustomLogger {
	return &CustomLogger{}
}

func (l *CustomLogger) Debug(v ...any) {
	log.Println(v...)
}

func (l *CustomLogger) Error(v ...any) {
	log.Println(v...)
}
