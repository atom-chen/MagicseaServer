package log

import (
	"os"
	"log"
	"strings"
	"time"
	"fmt"
	"path"
)

// levels
const (
	debugLevel = 0
	infoLevel  = 1
	errorLevel = 2
	fatalLevel = 3
)

const (
	printDebugLevel = "[DEBUG] "
	printInfoLevel  = "[INFO] "
	printErrorLevel = "[ERROR] "
	printFatalLevel = "[FATAL] "
)

type Logger struct {
	level      int
	baseLogger *log.Logger
	baseFile   *os.File
}

func New(strLevel string, pathhead string, pathname string, flag int) (*Logger, error) {
	// level
	var level int
	switch strings.ToLower(strLevel) {
	case "debug":
		level = debugLevel
	case "info":
		level = infoLevel
	case "error":
		level = errorLevel
	case "fatal":
		level = fatalLevel
	default:
		//todo error
		return nil, nil

	}

	// logger
	var baseLogger *log.Logger
	var baseFile *os.File
	if pathname != "" {
		now := time.Now()
		filename := fmt.Sprintf("%v%d%02d%02d.log", pathhead,
			now.Year(),
			now.Month(),
			now.Day())

		file, err := os.OpenFile(path.Join(pathname, filename), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666) //Create(path.Join(pathname, filename))
		if err != nil {
			return nil, err
		}
		baseLogger = log.New(file, "", flag)
		baseFile = file
	}else {
		baseLogger = log.New(os.Stdout, "", flag)
	}

	// new
	logger := new(Logger)
	logger.level = level
	logger.baseLogger = baseLogger
	logger.baseFile = baseFile

	return logger, nil
}

// It's dangerous to call the method on logging
func (logger *Logger) Close() {
	if logger.baseFile != nil {
		logger.baseFile.Close()
	}

	logger.baseLogger = nil
	logger.baseFile = nil
}

func (logger *Logger) doPrintf(level int, printLevel string, format string, a ...interface{}) {
	if level < logger.level {
		return
	}
	if logger.baseLogger == nil {
		panic("logger closed")
	}

	format = printLevel + format
	logger.baseLogger.Output(4, fmt.Sprintf(format, a...))

	if level == fatalLevel {
		os.Exit(1)
	}
}

func (logger *Logger) Debug(format string, a ...interface{}) {
	logger.doPrintf(debugLevel, printDebugLevel, format, a...)
}

func (logger *Logger) Info(format string, a ...interface{}) {
	logger.doPrintf(infoLevel, printInfoLevel, format, a...)
}

func (logger *Logger) Error(format string, a ...interface{}) {
	logger.doPrintf(errorLevel, printErrorLevel, format, a...)
}

func (logger *Logger) Fatal(format string, a ...interface{}) {
	logger.doPrintf(fatalLevel, printFatalLevel, format, a...)
}

func (logger *Logger) Println(format string, a ...interface{}) {
	logger.Println(format,a...)
}