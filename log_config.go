package main

import (
	"fmt"
	"io"
	"logrus"
	"os"
	"path"
	"runtime"
	"strconv"

	"golang.org/x/sys/unix"
	"gopkg.in/natefinch/lumberjack.v2"
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
	logModMain              = "MAIN_MODDULE"
	logModGinContext        = "GIN_CONTEXT"
	logModDBControl         = "DB_CONTROL"
	logModInstituteMgmt     = "INSTITUTE_MGMT"
	logModTeacherMgmt       = "TEACHER_MGMT"
	logModCourseMgmt        = "COURSE_MGMT"
	logModCourseRecordMgmt  = "COURSE_RECORD_MGMT"
	logModCourseCommentMgmt = "COURSE_COMMENT_MGMT"
	logModRelativeMgmt      = "RELATIVE_MGMT"
	logModStudentMgmt       = "STUDENT_MGMT"
	logModCloudMediaMgmt    = "CLOUDMEDIA_MGMT"
	logModReferenceMgmt     = "REFERENCE_MGMT"
	logModTemplateMgmt      = "TEMPLATE_MGMT"
	logModStoryMgmt         = "STORY_MGMT"
)

var logModEnabledTable = map[string]bool{
	logModMain:              true,
	logModGinContext:        true,
	logModDBControl:         true,
	logModInstituteMgmt:     true,
	logModTeacherMgmt:       true,
	logModCourseMgmt:        true,
	logModCourseRecordMgmt:  true,
	logModCourseCommentMgmt: true,
	logModRelativeMgmt:      true,
	logModStudentMgmt:       true,
	logModCloudMediaMgmt:    true,
	logModReferenceMgmt:     true,
	logModTemplateMgmt:      true,
	logModStoryMgmt:         true,
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
