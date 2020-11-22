package coderun

import (
	"log"
	"os"
	"strings"
)

var Logger *LoggerType

type LoggerType struct {
	Debug *log.Logger
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger
}

func SetupLogger(level string) *LoggerType {
	l := &LoggerType{}

	devNull, err := os.Open(os.DevNull)
	if err != nil {
		panic(err)
	}

	var debugOut, infoOut, warnOut, errorOut *os.File

	// Send log output to /dev/null if we don't set the right flag
	switch strings.ToLower(level) {
	case "Debug":
		errorOut, warnOut, infoOut, debugOut = os.Stderr, os.Stderr, os.Stderr, os.Stderr
	case "Info":
		errorOut, warnOut, infoOut, debugOut = os.Stderr, os.Stderr, os.Stderr, devNull
	case "Warn":
		errorOut, warnOut, infoOut, debugOut = os.Stderr, os.Stderr, devNull, devNull
	case "Error":
		errorOut, warnOut, infoOut, debugOut = os.Stderr, devNull, devNull, devNull
	default:
		panic("Not a valid log setting")
	}

	l.Debug = log.New(debugOut, "[DEBUG] ", log.Ldate|log.Ltime)
	l.Info = log.New(infoOut, "[INFO]  ", log.Ldate|log.Ltime)
	l.Warn = log.New(warnOut, "[WARN]  ", log.Ldate|log.Ltime)
	l.Error = log.New(errorOut, "[ERROR] ", log.Ldate|log.Ltime)

	l.Debug.Print("Set up logger")
	return l
}
