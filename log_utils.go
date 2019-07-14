package main

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path"
	"runtime"
	"strconv"
)

var loggingLevelMap = map[string]logrus.Level{
	"panic": logrus.PanicLevel,
	"error": logrus.ErrorLevel,
	"warn":  logrus.WarnLevel,
	"info":  logrus.InfoLevel,
	"debug": logrus.DebugLevel,
	"trace": logrus.TraceLevel,
}

const (
	stdLogFile   = "/tmp/vlog_server.log"
	ginLogFile   = "/tmp/vlog_gin.log"
	crashLogFile = "/tmp/vlog_crash.log"
)

var ljStdLogger = &lumberjack.Logger{
	Filename:   stdLogFile,
	MaxSize:    80,
	MaxBackups: 10,
}

var ljGinLogger = &lumberjack.Logger{
	Filename:   ginLogFile,
	MaxSize:    80,
	MaxBackups: 10,
}

var ljCrashLogger = &lumberjack.Logger{
	Filename:   crashLogFile,
	MaxSize:    80,
	MaxBackups: 10,
}

// Logging global customized logging module
var logging *logrus.Logger

func loggingCallerBeautifier(f *runtime.Frame) (funciton string, file string) {
	var funcName = f.Function
	var fileName = path.Base(f.File) + ":" + strconv.Itoa(f.Line)
	return funcName, fileName
}

func loggingInitSetup(sc *ServerConfig) *logrus.Logger {
	var logging = logrus.New()

	// output location
	var mwriter = io.MultiWriter(ljStdLogger, os.Stdout)
	logging.SetOutput(mwriter)

	// logging level
	var loggingLevel = logrus.DebugLevel
	if _, ok := loggingLevelMap[serverConfig.LoggingLevel]; ok {
		loggingLevel = loggingLevelMap[serverConfig.LoggingLevel]
	}
	logging.SetLevel(loggingLevel)

	// logging format based on release mode
	if serverConfig.LoggingReleaseMode == true {
		logging.SetReportCaller(true)
		logging.SetFormatter(&logrus.JSONFormatter{
			CallerPrettyfier: loggingCallerBeautifier,
		})
	} else {
		logging.SetReportCaller(true)
		logging.SetFormatter(&ColorTextFormatter{
			FullTimestamp:   true,
			ForceFormatting: true,
		})
	}

	return logging
}
