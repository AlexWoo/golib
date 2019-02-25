// Copyright (C) AlexWoo(Wu Jie) wj19840501@gmail.com
//
// golib log

package golib

import (
	"fmt"
	"log"
	"os"
	"sync"
)

// log level const definition
const (
	LOGDEBUG = iota
	LOGINFO
	LOGERROR
	LOGFATAL
)

var logLevel = []string{
	"[debug] ",
	"[info ] ",
	"[error] ",
	"[fatal] ",
}

// log level conf enum
var LoglvEnum = Enum{
	"debug": LOGDEBUG,
	"info":  LOGINFO,
	"error": LOGERROR,
	"fatal": LOGFATAL,
}

// use log file path as index, if Open same file, will reuse Log instance
var (
	logm     = make(map[string]*Log)
	logmLock sync.Mutex
)

// LogCtx for user defined prefix and suffix add in log
type LogCtx interface {
	// Log Prefix
	Prefix() string

	// Log Suffix
	Suffix() string

	// Log Level
	LogLevel() int
}

// golib log struct
// Example:
//
// New a logger and start log:
//	type mainCtx struct {
//	}
//
//	func (h *mainCtx) Prefix() string {
//		return "[main]"
//	}
//
//	func (h *mainCtx) Suffix() string {
//		return "[END]"
//	}
//
//	func (h *mainCtx) LogLevel() int {
//		return golib.LOGINFO
//	}
//
//	func main() {
//		h := &mainCtx{}
//		logger := golib.NewLog("test.log")
//		if logger == nil {
//			fmt.Println("NewLog failed")
//		}
//
//		logger.LogDebug(h, "test debug")
//		logger.LogInfo(h, "test info")
//		logger.LogError(h, "test error")
//		logger.LogFatal(h, "test fatal")
//
//		logger.LogError(h, "Normal End")
//	}
//
// Result in test.log:
//
//	2018/07/11 08:30:14.671825 [info ] [main] test info [END]
//	2018/07/11 08:30:14.671875 [error] [main] test error [END]
//	2018/07/11 08:30:14.671884 [fatal] [main] test fatal [END]
type Log struct {
	path   string
	logger *log.Logger
}

func (l *Log) logPrintf(loglv int, c LogCtx, format string, v ...interface{}) {
	prefix := c.Prefix()
	suffix := c.Suffix()

	fmt := logLevel[loglv]
	if prefix != "" {
		fmt += prefix + " "
	}
	fmt += format
	if suffix != "" {
		fmt += " " + suffix
	}

	l.logger.Printf(fmt, v...)
}

// New golib log instance
func NewLog(logPath string) *Log {
	logmLock.Lock()
	defer logmLock.Unlock()

	reuselog := logm[logPath]
	if reuselog != nil {
		return reuselog
	}

	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil
	}

	l := &Log{
		path:   logPath,
		logger: log.New(f, "", log.LstdFlags|log.Lmicroseconds),
	}

	logm[logPath] = l

	return l
}

// Reopen all logfiles use golib.NewLog to create
func reopenLogs() error {
	logmLock.Lock()
	defer logmLock.Unlock()

	for _, l := range logm {
		f, err := os.OpenFile(l.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return fmt.Errorf("reopen log failed %s", err.Error())
		}

		l.logger = log.New(f, "", log.LstdFlags|log.Lmicroseconds)
	}

	return nil
}

// Return Log Path
func (l *Log) LogPath() string {
	return l.path
}

// Return *log.Logger
func (l *Log) Log() *log.Logger {
	return l.logger
}

// log a debug level log
func (l *Log) LogDebug(c LogCtx, format string, v ...interface{}) {
	if c.LogLevel() > LOGDEBUG {
		return
	}

	l.logPrintf(LOGDEBUG, c, format, v...)
}

// log a info level log
func (l *Log) LogInfo(c LogCtx, format string, v ...interface{}) {
	if c.LogLevel() > LOGINFO {
		return
	}

	l.logPrintf(LOGINFO, c, format, v...)
}

// log a error level log
func (l *Log) LogError(c LogCtx, format string, v ...interface{}) {
	if c.LogLevel() > LOGERROR {
		return
	}

	l.logPrintf(LOGERROR, c, format, v...)
}

// log a error level log, and exit
func (l *Log) LogFatal(c LogCtx, format string, v ...interface{}) {
	l.logPrintf(LOGFATAL, c, format, v...)
	os.Exit(1)
}
