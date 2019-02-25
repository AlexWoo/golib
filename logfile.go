// Copyright (C) AlexWoo(Wu Jie) wj19840501@gmail.com
//
// golib log file
//  use in accesslog in http server, trace log

package golib

import (
	"fmt"
	"os"
	"sync"
)

// use logfile path as index, if Open same file, will reuse log file instance
var (
	lfm    = make(map[string]*LogFile)
	lfLock sync.Mutex
)

// LogFile
type LogFile struct {
	path string
	file *os.File
	fm   sync.Mutex
	link uint32
}

// New golib LogFile instance, if same path has been open,
// return opened instance
func NewLogFile(path string) (*LogFile, error) {
	lfLock.Lock()
	defer lfLock.Unlock()

	lf := lfm[path]
	if lf != nil {
		lf.link++
		return lf, nil
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	lf = &LogFile{
		path: path,
		file: f,
		link: 1,
	}

	lfm[path] = lf

	return lf, nil
}

// reopen all filelog, for module reopen
func reopenfileLogs() error {
	lfLock.Lock()
	defer lfLock.Unlock()

	for _, lf := range lfm {
		lf.fm.Lock()
		defer lf.fm.Unlock()

		f, err := os.OpenFile(lf.path,
			os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("reopen logfile failed %s", err.Error())
		}

		lf.file.Close()

		lf.file = f
	}

	return nil
}

// Close file only if all instance related to file close
// Otherwise only decrease link num
func (f *LogFile) Close() error {
	f.fm.Lock()
	defer f.fm.Unlock()

	f.link--

	if f.link == 0 {
		return f.file.Close()
	}

	return nil
}

// Fd returns the integer Unix file descriptor referencing the open file.
// The file descriptor is valid only until
// 	f.Close is called or f is garbage collected.
// On Unix systems this will cause the SetDeadline methods to stop working.
func (f *LogFile) Fd() uintptr {
	return f.file.Fd()
}

// Name returns the name of the file as presented to Open.
func (f *LogFile) Name() string {
	return f.file.Name()
}

// Write writes len(b) bytes to the File.
// It returns the number of bytes written and an error, if any.
// Write returns a non-nil error when n != len(b).
func (f *LogFile) Write(b []byte) (n int, err error) {
	f.fm.Lock()
	defer f.fm.Unlock()

	if f.link == 0 {
		return 0, fmt.Errorf("file %s has been closed", f.path)
	}

	return f.file.Write(b)
}

// WriteAt writes len(b) bytes to the File starting at byte offset off.
// It returns the number of bytes written and an error, if any.
// WriteAt returns a non-nil error when n != len(b).
func (f *LogFile) WriteAt(b []byte, off int64) (n int, err error) {
	f.fm.Lock()
	defer f.fm.Unlock()

	if f.link == 0 {
		return 0, fmt.Errorf("file %s has been closed", f.path)
	}

	return f.file.WriteAt(b, off)
}

// WriteString is like Write, but writes the contents of string s rather than
// a slice of bytes.
func (f *LogFile) WriteString(s string) (n int, err error) {
	f.fm.Lock()
	defer f.fm.Unlock()

	if f.link == 0 {
		return 0, fmt.Errorf("file %s has been closed", f.path)
	}

	return f.file.WriteString(s)
}
