package config

import (
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
)

func TextLogFormatter() log.Formatter {
	return &log.TextFormatter{
		ForceColors:               true,
		DisableColors:             false,
		ForceQuote:                true,
		DisableQuote:              false,
		EnvironmentOverrideColors: true,
		DisableTimestamp:          false,
		FullTimestamp:             true,
		TimestampFormat:           "2006/01/02 15:04:05.000",
		DisableSorting:            true,
		SortingFunc:               nil,
		DisableLevelTruncation:    false,
		PadLevelText:              true,
		QuoteEmptyFields:          false,
		FieldMap: log.FieldMap{
			"FieldKeyTime":  "@timestamp",
			"FieldKeyLevel": "@level",
			"FieldKeyFile":  "@file",
			"FieldKeyMsg":   "@message"},
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			if strings.Contains(frame.File, "src/") {
				srcIndex := strings.Index(frame.File, "src/")
				packagePath := frame.File[srcIndex:]
				packageName := strings.Replace(packagePath, "/", ".", -1)
				return packageName, frame.Function + "[" + strconv.FormatInt(int64(frame.Line), 10) + "]"
			}
			return frame.Function, frame.File
		},
	}
}

func JsonLogFormatter() log.Formatter {
	return &log.JSONFormatter{
		TimestampFormat:   "2006/01/02 15:04:05.000",
		DisableTimestamp:  false,
		DisableHTMLEscape: false,
		DataKey:           "",
		FieldMap:          nil,
		CallerPrettyfier:  nil,
		PrettyPrint:       true,
	}
}

func MakeNewLog() *log.Logger {
	logger := log.New()

	writerWithStdout := os.Stdout
	writerWithLogFile, err := os.OpenFile("stand.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		log.Fatalf("create file log.txt failed: %v", err)
	}
	logger.SetOutput(io.MultiWriter(writerWithStdout, writerWithLogFile))
	logger.SetReportCaller(true)
	logger.SetLevel(log.DebugLevel)
	return logger
}

func InitLog() {

	writerWithStdout := os.Stdout
	writerWithLogFile, err := os.OpenFile("stand.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		log.Fatalf("create file log.txt failed: %v", err)
	}
	log.SetOutput(io.MultiWriter(writerWithStdout, writerWithLogFile))
	log.SetReportCaller(true)
	log.SetLevel(log.DebugLevel)
}
