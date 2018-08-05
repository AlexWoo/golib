// Copyright (C) AlexWoo(Wu Jie) wj19840501@gmail.com
//
// golib log

package golib

import (
	"log"
	"os"
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

// LogCtx for user defined prefix and suffix add in log
type LogCtx interface {
	Prefix() string
	Suffix() string
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
//	func main() {
//		h := &mainCtx{}
//		logger := golib.NewLog("test.log", golib.LOGINFO)
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
	level  int
	logger *log.Logger
}

func (l *Log) logPrintf(loglv int, c LogCtx, format string, v ...interface{}) {
	format = logLevel[loglv] + c.Prefix() + " " + format + " " + c.Suffix()
	l.logger.Printf(format, v...)
}

// New golib log instance
func NewLog(logPath string, logLevel int) *Log {
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil
	}

	l := &Log{
		path:   logPath,
		level:  logLevel,
		logger: log.New(f, "", log.LstdFlags|log.Lmicroseconds|log.LUTC),
	}

	return l
}

// log a debug level log
func (l *Log) LogDebug(c LogCtx, format string, v ...interface{}) {
	if l.level > LOGDEBUG {
		return
	}

	l.logPrintf(LOGDEBUG, c, format, v...)
}

// log a info level log
func (l *Log) LogInfo(c LogCtx, format string, v ...interface{}) {
	if l.level > LOGINFO {
		return
	}

	l.logPrintf(LOGINFO, c, format, v...)
}

// log a error level log
func (l *Log) LogError(c LogCtx, format string, v ...interface{}) {
	if l.level > LOGERROR {
		return
	}

	l.logPrintf(LOGERROR, c, format, v...)
}

// log a error level log, and exit
func (l *Log) LogFatal(c LogCtx, format string, v ...interface{}) {
	l.logPrintf(LOGFATAL, c, format, v...)
	os.Exit(1)
}
