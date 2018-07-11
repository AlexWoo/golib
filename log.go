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

// LogHandle for user defined prefix and suffix add in log
type LogHandle interface {
	Prefix() string
	Suffix() string
}

// golib log struct
// Example:
//
// New a logger and start log:
//	type mainLogHandle struct {
//	}
//
//	func (h *mainLogHandle) Prefix() string {
//		return "[main]"
//	}
//
//	func (h *mainLogHandle) Suffix() string {
//		return "[END]"
//	}
//
//	func main() {
//		h := &mainLogHandle{}
//		logger := golib.NewLog(h, "test.log", golib.LOGINFO)
//		if logger == nil {
//			fmt.Println("NewLog failed")
//		}
//
//		logger.LogDebug("test debug")
//		logger.LogInfo("test info")
//		logger.LogError("test error")
//		logger.LogFatal("test fatal")
//
//		logger.LogError("Normal End")
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
	handle LogHandle
}

func (l *Log) logPrintf(loglv int, format string, v ...interface{}) {
	format = logLevel[loglv] + l.handle.Prefix() + " " + format +
		" " + l.handle.Suffix()
	l.logger.Printf(format, v...)
}

// New golib log instance
func NewLog(handle LogHandle, logPath string, logLevel int) *Log {
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil
	}

	l := &Log{
		path:   logPath,
		level:  logLevel,
		handle: handle,
		logger: log.New(f, "", log.LstdFlags|log.Lmicroseconds|log.LUTC),
	}

	return l
}

// log a debug level log
func (l *Log) LogDebug(format string, v ...interface{}) {
	if l.level > LOGDEBUG {
		return
	}

	l.logPrintf(LOGDEBUG, format, v...)
}

// log a info level log
func (l *Log) LogInfo(format string, v ...interface{}) {
	if l.level > LOGINFO {
		return
	}

	l.logPrintf(LOGINFO, format, v...)
}

// log a error level log
func (l *Log) LogError(format string, v ...interface{}) {
	if l.level > LOGERROR {
		return
	}

	l.logPrintf(LOGERROR, format, v...)
}

// log a error level log, and exit
func (l *Log) LogFatal(format string, v ...interface{}) {
	l.logPrintf(LOGFATAL, format, v...)
	os.Exit(1)
}
