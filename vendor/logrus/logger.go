package logrus

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type Logger struct {
	// The logs are `io.Copy`'d to this in a mutex. It's common to set this to a
	// file, or leave it default which is `os.Stderr`. You can also set this to
	// something more adventurous, such as logging to Kafka.
	Out io.Writer
	// Hooks for the logger instance. These allow firing events based on logging
	// levels and log entries. For example, to send errors to an error tracking
	// service, log to StatsD or dump the core on fatal errors.
	Hooks LevelHooks
	// All log entries pass through the formatter before logged to Out. The
	// included formatters are `TextFormatter` and `JSONFormatter` for which
	// TextFormatter is the default. In development (when a TTY is attached) it
	// logs with colors, but to a file it wouldn't. You can easily implement your
	// own that implements the `Formatter` interface, see the `README` or included
	// formatters for examples.
	Formatter Formatter

	// Flag for whether to log caller info (off by default)
	ReportCaller bool

	// The logging level the logger should log at. This is typically (and defaults
	// to) `logrus.Info`, which allows Info(), Warn(), Error() and Fatal() to be
	// logged.
	Level Level
	// Used to sync writing to the log. Locking is enabled by Default
	mu MutexWrap
	// Reusable empty entry
	entryPool sync.Pool
	// Function to exit the application, defaults to `os.Exit()`
	ExitFunc exitFunc
	// Module map to control logging enabling
	ModuleEnabledTable map[string]bool
	ModuleRWLock       *sync.RWMutex
}

type exitFunc func(int)

type MutexWrap struct {
	lock     sync.Mutex
	disabled bool
}

func (mw *MutexWrap) Lock() {
	if !mw.disabled {
		mw.lock.Lock()
	}
}

func (mw *MutexWrap) Unlock() {
	if !mw.disabled {
		mw.lock.Unlock()
	}
}

func (mw *MutexWrap) Disable() {
	mw.disabled = true
}

// Creates a new logger. Configuration should be set by changing `Formatter`,
// `Out` and `Hooks` directly on the default logger instance. You can also just
// instantiate your own:
//
//    var log = &Logger{
//      Out: os.Stderr,
//      Formatter: new(JSONFormatter),
//      Hooks: make(LevelHooks),
//      Level: logrus.DebugLevel,
//    }
//
// It's recommended to make this a global instance called `log`.
func New() *Logger {
	return &Logger{
		Out:                os.Stderr,
		Formatter:          new(TextFormatter),
		Hooks:              make(LevelHooks),
		Level:              InfoLevel,
		ExitFunc:           os.Exit,
		ReportCaller:       false,
		ModuleEnabledTable: map[string]bool{}, // value always true
		ModuleRWLock:       &sync.RWMutex{},
	}
}

func (logger *Logger) RegisterModule(module string, enabled bool) error {
	logger.ModuleRWLock.Lock()
	defer logger.ModuleRWLock.Unlock()
	if _, exist := logger.ModuleEnabledTable[module]; exist {
		return fmt.Errorf("cannot register logger module <%s>: already exists", module)
	}
	logger.ModuleEnabledTable[module] = enabled
	return nil
}

func (logger *Logger) UnregisterModule(module string) error {
	logger.ModuleRWLock.Lock()
	defer logger.ModuleRWLock.Unlock()
	if _, exist := logger.ModuleEnabledTable[module]; !exist {
		return fmt.Errorf("cannot unregister logger module <%s>: not exist", module)
	}
	delete(logger.ModuleEnabledTable, module)
	return nil
}

func (logger *Logger) GetModuleEnabled(module string) (bool, error) {
	logger.ModuleRWLock.RLock()
	defer logger.ModuleRWLock.RUnlock()
	var enabled bool = false
	var exist bool = false
	enabled, exist = logger.ModuleEnabledTable[module]
	if !exist {
		return false, fmt.Errorf("logger module <%s> not exist", module)
	}
	return enabled, nil
}

func (logger *Logger) SetModuleEnabled(module string, enabled bool) error {
	logger.ModuleRWLock.Lock()
	defer logger.ModuleRWLock.Unlock()
	if _, exist := logger.ModuleEnabledTable[module]; !exist {
		return fmt.Errorf("cannot set logger module <%s>: not exist", module)
	}
	logger.ModuleEnabledTable[module] = enabled
	return nil
}

func (logger *Logger) EnableAllModules() {
	logger.ModuleRWLock.Lock()
	defer logger.ModuleRWLock.Unlock()
	for module := range logger.ModuleEnabledTable {
		logger.ModuleEnabledTable[module] = true
	}
}

func (logger *Logger) DisableAllModules() {
	logger.ModuleRWLock.Lock()
	defer logger.ModuleRWLock.Unlock()
	for module := range logger.ModuleEnabledTable {
		logger.ModuleEnabledTable[module] = false
	}
}

func (logger *Logger) newEntry() *Entry {
	entry, ok := logger.entryPool.Get().(*Entry)
	if ok {
		return entry
	}
	return NewEntry(logger)
}

func (logger *Logger) releaseEntry(entry *Entry) {
	entry.Data = map[string]interface{}{}
	logger.entryPool.Put(entry)
}

// Adds a field to the log entry, note that it doesn't log until you call
// Debug, Print, Info, Warn, Error, Fatal or Panic. It only creates a log entry.
// If you want multiple fields, use `WithFields`.
func (logger *Logger) WithField(key string, value interface{}) *Entry {
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.WithField(key, value)
}

// Adds a struct of fields to the log entry. All it does is call `WithField` for
// each `Field`.
func (logger *Logger) WithFields(fields Fields) *Entry {
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.WithFields(fields)
}

// Add an error as single field to the log entry.  All it does is call
// `WithError` for the given `error`.
func (logger *Logger) WithError(err error) *Entry {
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.WithError(err)
}

// Add a context to the log entry.
func (logger *Logger) WithContext(ctx context.Context) *Entry {
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.WithContext(ctx)
}

// Overrides the time of the log entry.
func (logger *Logger) WithTime(t time.Time) *Entry {
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.WithTime(t)
}

// Logmf with module field
func (logger *Logger) Logmf(level Level, module string, format string, args ...interface{}) {
	if enabled, _ := logger.GetModuleEnabled(module); !enabled {
		return
	}
	if logger.IsLevelEnabled(level) {
		entry := logger.WithField("module", module)
		entry.Logf(level, format, args...)
		logger.releaseEntry(entry)
	}
}

func (logger *Logger) Printmf(module string, format string, args ...interface{}) {
	if enabled, _ := logger.GetModuleEnabled(module); !enabled {
		return
	}
	entry := logger.WithField("module", module)
	entry.Printf(format, args...)
	logger.releaseEntry(entry)
}

func (logger *Logger) Tracemf(module string, format string, args ...interface{}) {
	logger.Logmf(TraceLevel, module, format, args...)
}

func (logger *Logger) Debugmf(module string, format string, args ...interface{}) {
	logger.Logmf(DebugLevel, module, format, args...)
}

func (logger *Logger) Infomf(module string, format string, args ...interface{}) {
	logger.Logmf(InfoLevel, module, format, args...)
}

func (logger *Logger) Warnmf(module string, format string, args ...interface{}) {
	logger.Logmf(WarnLevel, module, format, args...)
}

func (logger *Logger) Warningmf(module string, format string, args ...interface{}) {
	logger.Warnmf(module, format, args...)
}

func (logger *Logger) Errormf(module string, format string, args ...interface{}) {
	logger.Logmf(ErrorLevel, module, format, args...)
}

func (logger *Logger) Fatalmf(module string, format string, args ...interface{}) {
	logger.Logmf(FatalLevel, module, format, args...)
	logger.Exit(1)
}

func (logger *Logger) Panicmf(module string, format string, args ...interface{}) {
	logger.Logmf(PanicLevel, module, format, args...)
}

// Logmln with module logging
func (logger *Logger) Logmln(level Level, module string, args ...interface{}) {
	if enabled, _ := logger.GetModuleEnabled(module); !enabled {
		return
	} else if logger.IsLevelEnabled(level) {
		entry := logger.WithField("module", module)
		entry.Logln(level, args...)
		logger.releaseEntry(entry)
	}
}

func (logger *Logger) Printmln(module string, args ...interface{}) {
	if enabled, _ := logger.GetModuleEnabled(module); !enabled {
		return
	}
	entry := logger.WithField("module", module)
	entry.Println(args...)
	logger.releaseEntry(entry)
}

func (logger *Logger) Tracemln(module string, args ...interface{}) {
	logger.Logmln(TraceLevel, module, args...)
}

func (logger *Logger) Debugmln(module string, args ...interface{}) {
	logger.Logmln(DebugLevel, module, args...)
}

func (logger *Logger) Infomln(module string, args ...interface{}) {
	logger.Logmln(InfoLevel, module, args...)
}

func (logger *Logger) Warnmln(module string, args ...interface{}) {
	logger.Logmln(WarnLevel, module, args...)
}

func (logger *Logger) Warningmln(module string, args ...interface{}) {
	logger.Warnmln(module, args...)
}

func (logger *Logger) Errormln(module string, args ...interface{}) {
	logger.Logmln(ErrorLevel, module, args...)
}

func (logger *Logger) Fatalmln(module string, args ...interface{}) {
	logger.Logmln(FatalLevel, module, args...)
	logger.Exit(1)
}

func (logger *Logger) Panicmln(module string, args ...interface{}) {
	logger.Logmln(PanicLevel, module, args...)
}

func (logger *Logger) Logf(level Level, format string, args ...interface{}) {
	if logger.IsLevelEnabled(level) {
		entry := logger.newEntry()
		entry.Logf(level, format, args...)
		logger.releaseEntry(entry)
	}
}

func (logger *Logger) Tracef(format string, args ...interface{}) {
	logger.Logf(TraceLevel, format, args...)
}

func (logger *Logger) Debugf(format string, args ...interface{}) {
	logger.Logf(DebugLevel, format, args...)
}

func (logger *Logger) Infof(format string, args ...interface{}) {
	logger.Logf(InfoLevel, format, args...)
}

func (logger *Logger) Printf(format string, args ...interface{}) {
	entry := logger.newEntry()
	entry.Printf(format, args...)
	logger.releaseEntry(entry)
}

func (logger *Logger) Warnf(format string, args ...interface{}) {
	logger.Logf(WarnLevel, format, args...)
}

func (logger *Logger) Warningf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

func (logger *Logger) Errorf(format string, args ...interface{}) {
	logger.Logf(ErrorLevel, format, args...)
}

func (logger *Logger) Fatalf(format string, args ...interface{}) {
	logger.Logf(FatalLevel, format, args...)
	logger.Exit(1)
}

func (logger *Logger) Panicf(format string, args ...interface{}) {
	logger.Logf(PanicLevel, format, args...)
}

func (logger *Logger) Log(level Level, args ...interface{}) {
	if logger.IsLevelEnabled(level) {
		entry := logger.newEntry()
		entry.Log(level, args...)
		logger.releaseEntry(entry)
	}
}

func (logger *Logger) Trace(args ...interface{}) {
	logger.Log(TraceLevel, args...)
}

func (logger *Logger) Debug(args ...interface{}) {
	logger.Log(DebugLevel, args...)
}

func (logger *Logger) Info(args ...interface{}) {
	logger.Log(InfoLevel, args...)
}

func (logger *Logger) Print(args ...interface{}) {
	entry := logger.newEntry()
	entry.Print(args...)
	logger.releaseEntry(entry)
}

func (logger *Logger) Warn(args ...interface{}) {
	logger.Log(WarnLevel, args...)
}

func (logger *Logger) Warning(args ...interface{}) {
	logger.Warn(args...)
}

func (logger *Logger) Error(args ...interface{}) {
	logger.Log(ErrorLevel, args...)
}

func (logger *Logger) Fatal(args ...interface{}) {
	logger.Log(FatalLevel, args...)
	logger.Exit(1)
}

func (logger *Logger) Panic(args ...interface{}) {
	logger.Log(PanicLevel, args...)
}

func (logger *Logger) Logln(level Level, args ...interface{}) {
	if logger.IsLevelEnabled(level) {
		entry := logger.newEntry()
		entry.Logln(level, args...)
		logger.releaseEntry(entry)
	}
}

func (logger *Logger) Traceln(args ...interface{}) {
	logger.Logln(TraceLevel, args...)
}

func (logger *Logger) Debugln(args ...interface{}) {
	logger.Logln(DebugLevel, args...)
}

func (logger *Logger) Infoln(args ...interface{}) {
	logger.Logln(InfoLevel, args...)
}

func (logger *Logger) Println(args ...interface{}) {
	entry := logger.newEntry()
	entry.Println(args...)
	logger.releaseEntry(entry)
}

func (logger *Logger) Warnln(args ...interface{}) {
	logger.Logln(WarnLevel, args...)
}

func (logger *Logger) Warningln(args ...interface{}) {
	logger.Warnln(args...)
}

func (logger *Logger) Errorln(args ...interface{}) {
	logger.Logln(ErrorLevel, args...)
}

func (logger *Logger) Fatalln(args ...interface{}) {
	logger.Logln(FatalLevel, args...)
	logger.Exit(1)
}

func (logger *Logger) Panicln(args ...interface{}) {
	logger.Logln(PanicLevel, args...)
}

func (logger *Logger) Exit(code int) {
	runHandlers()
	if logger.ExitFunc == nil {
		logger.ExitFunc = os.Exit
	}
	logger.ExitFunc(code)
}

//When file is opened with appending mode, it's safe to
//write concurrently to a file (within 4k message on Linux).
//In these cases user can choose to disable the lock.
func (logger *Logger) SetNoLock() {
	logger.mu.Disable()
}

func (logger *Logger) level() Level {
	return Level(atomic.LoadUint32((*uint32)(&logger.Level)))
}

// SetLevel sets the logger level.
func (logger *Logger) SetLevel(level Level) {
	atomic.StoreUint32((*uint32)(&logger.Level), uint32(level))
}

// GetLevel returns the logger level.
func (logger *Logger) GetLevel() Level {
	return logger.level()
}

// AddHook adds a hook to the logger hooks.
func (logger *Logger) AddHook(hook Hook) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.Hooks.Add(hook)
}

// IsLevelEnabled checks if the log level of the logger is greater than the level param
func (logger *Logger) IsLevelEnabled(level Level) bool {
	return logger.level() >= level
}

// SetFormatter sets the logger formatter.
func (logger *Logger) SetFormatter(formatter Formatter) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.Formatter = formatter
}

// SetOutput sets the logger output.
func (logger *Logger) SetOutput(output io.Writer) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.Out = output
}

func (logger *Logger) SetReportCaller(reportCaller bool) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.ReportCaller = reportCaller
}

// ReplaceHooks replaces the logger hooks and returns the old ones
func (logger *Logger) ReplaceHooks(hooks LevelHooks) LevelHooks {
	logger.mu.Lock()
	oldHooks := logger.Hooks
	logger.Hooks = hooks
	logger.mu.Unlock()
	return oldHooks
}
