package config

import (
	"fmt"
	"log"
	"os"

	"ehgm.com.br/url-shortener/domain/ports"
)

type logger struct {
	logInfo  *log.Logger
	logError *log.Logger
}

func NewLogger() ports.Logger {
	return &logger{
		logInfo:  log.New(os.Stdout, "[API] INFO: ", log.Lshortfile),
		logError: log.New(os.Stderr, "[API] ERROR: ", log.Lshortfile),
	}
}

func (l logger) Info(format string, v ...interface{}) {
	l.logInfo.Output(2, fmt.Sprintf(format, v...))
}

func (l logger) Error(format string, v ...interface{}) {
	l.logError.Output(2, fmt.Sprintf(format, v...))
}

func (l logger) Fatal(format string, v ...interface{}) {
	l.logError.Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
}
