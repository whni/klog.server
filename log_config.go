package main

import (
	"fmt"
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
	logFilePrefix    = "klog"
	stdLogFile       = "/tmp/" + logFilePrefix + "_server.log"
	ginLogFile       = "/tmp/" + logFilePrefix + "_gin.log"
	errLogFile       = "/tmp/" + logFilePrefix + "_error.log"
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
	logModMain             = "MainModule"
	logModGinContext       = "GinContext"
	logModDBControl        = "DBControl"
	logModInstituteHandler = "InstituteHandler"
	logModClassHandler     = "ClassHandler"
	logModTeacherHandler   = "TeacherHandler"
	logModStudentHandler   = "StudentHandler"
)

var logModEnabledTable = map[string]bool{
	logModMain:             true,
	logModGinContext:       true,
	logModDBControl:        true,
	logModInstituteHandler: true,
	logModClassHandler:     true,
	logModTeacherHandler:   true,
	logModStudentHandler:   true,
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
	if sc.LoggingDestination == "stdout+file" || sc.LoggingDestination == "file+stdout" {
		var mwriter = io.MultiWriter(ljStdLogger, os.Stdout)
		logging.SetOutput(mwriter)
	} else if sc.LoggingDestination == "file" {
		logging.SetOutput(ljStdLogger)
	} else {
		// sc.LoggingDestination == "stdout" or other
		logging.SetOutput(os.Stdout)
	}

	// logging level
	var loggingLevel = logrus.DebugLevel
	if _, ok := loggingLevelMap[sc.LoggingLevel]; ok {
		loggingLevel = loggingLevelMap[sc.LoggingLevel]
	}
	logging.SetLevel(loggingLevel)

	// logging format based on release mode
	if sc.LoggingReleaseMode == true {
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

func loggingRegisterModules(l *logrus.Logger, moduleTable map[string]bool) error {
	for module, enabled := range moduleTable {
		if err := l.RegisterModule(module, enabled); err != nil {
			return err
		}
	}
	return nil
}

func loggingErrRedirect(errFile string) error {
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
		return fmt.Errorf("Failed to open error log file: %v", err.Error())
	}

	// redirect stderr
	err = unix.Dup2(int(f.Fd()), int(os.Stderr.Fd()))
	if err != nil {
		return fmt.Errorf("Failed to redirect stderr to error file: %v", err.Error())
	}

	return nil
}
