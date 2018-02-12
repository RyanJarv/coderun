package coderun

import (
	"log"
	"os"
	"strings"
)

var Logger *LoggerType

type LoggerType struct {
	debug *log.Logger
	info  *log.Logger
	warn  *log.Logger
	error *log.Logger
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
	case "debug":
		errorOut, warnOut, infoOut, debugOut = os.Stderr, os.Stderr, os.Stderr, os.Stderr
	case "info":
		errorOut, warnOut, infoOut, debugOut = os.Stderr, os.Stderr, os.Stderr, devNull
	case "warn":
		errorOut, warnOut, infoOut, debugOut = os.Stderr, os.Stderr, devNull, devNull
	case "error":
		errorOut, warnOut, infoOut, debugOut = os.Stderr, devNull, devNull, devNull
	default:
		panic("Not a valid log setting")
	}

	l.debug = log.New(debugOut, "DEBUG", log.Ldate|log.Ltime)
	l.info = log.New(infoOut, "INFO", log.Ldate|log.Ltime)
	l.warn = log.New(warnOut, "WARN", log.Ldate|log.Ltime)
	l.error = log.New(errorOut, "ERROR", log.Ldate|log.Ltime)

	l.debug.Print("Set up logger")
	return l
}
