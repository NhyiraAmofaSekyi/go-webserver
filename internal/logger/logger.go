package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

type Logger struct {
	debug   *log.Logger
	info    *log.Logger
	error   *log.Logger
	isDebug bool
}

var std *Logger

const (
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorReset  = "\033[0m"
)

func Init(debugMode bool) {
	file, err := os.OpenFile(
		fmt.Sprintf("logs/server_%s.log", time.Now().Format("2006-01-02")),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0666,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Write to both file and stdout
	multiWriter := io.MultiWriter(os.Stdout, file)

	// Common flags for all loggers
	const flags = log.Ldate | log.Ltime | log.Lshortfile

	std = &Logger{
		debug:   log.New(multiWriter, colorYellow+"DEBUG: "+colorReset, flags),
		info:    log.New(multiWriter, colorGreen+"INFO: "+colorReset, flags),
		error:   log.New(multiWriter, colorRed+"ERROR: "+colorReset, flags),
		isDebug: debugMode,
	}
}

func Debug(format string, v ...interface{}) {
	if std.isDebug {
		std.debug.Output(2, fmt.Sprintf(format, v...))
	}
}

func Info(format string, v ...interface{}) {
	std.info.Output(2, fmt.Sprintf(format, v...))
}

func Error(format string, v ...interface{}) {
	std.error.Output(2, fmt.Sprintf(format, v...))
}

func Fatal(format string, v ...interface{}) {
	std.error.Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
}
