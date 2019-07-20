package main

import (
	"golang.org/x/sys/unix"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"logrus"
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
	stdLogFile       = "/tmp/vlog_server.log"
	ginLogFile       = "/tmp/vlog_gin.log"
	errLogFile       = "/tmp/vlog_error.log"
	logFileMaxSize   = 80
	logFileMaxbackup = 10
)

var ljStdLogger = &lumberjack.Logger{
	Filename:   stdLogFile,
	MaxSize:    logFileMaxSize,
	MaxBackups: logFileMaxbackup,
}

var ljGinLogger = &lumberjack.Logger{
	Filename:   ginLogFile,
	MaxSize:    logFileMaxSize,
	MaxBackups: logFileMaxbackup,
}

const (
	logModMain       = "MainModule"
	logModGinContext = "GinContext"
)

var logModEnabledTable = map[string]bool{
	logModMain:       true,
	logModGinContext: true,
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
	if sc.LoggingWithStdout {
		var mwriter = io.MultiWriter(ljStdLogger, os.Stdout)
		logging.SetOutput(mwriter)
	} else {
		logging.SetOutput(ljStdLogger)
	}

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

func loggingRegisterModules(moduleTable map[string]bool) {
	for module, enabled := range moduleTable {
		logging.RegisterModule(module, enabled)
	}
}

func loggingErrRedirect(errFile string) {
	// rotate error log file
	var errFileSize int64 = 0
	if fs, err := os.Stat(errFile); err == nil {
		errFileSize = fs.Size()
	}
	if errFileSize > logFileMaxSize*1024*1024 {
		os.Rename(errFile, errFile+".old")
	}

	f, err := os.OpenFile(errFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		logging.Fatalf("Failed to open error log file: %v", err)
	}

	// redirect stderr
	err = unix.Dup2(int(f.Fd()), int(os.Stderr.Fd()))
	if err != nil {
		logging.Fatalf("Failed to redirect stderr to error file: %v", err)
	}
}
