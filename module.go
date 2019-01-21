// Copyright (C) AlexWoo(Wu Jie) wj19840501@gmail.com
//
// golib module

package golib

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Module interface
type Module interface {
	// Do something before Module init
	PreInit() error

	// Module init, such as load config, init control paras
	Init() error

	// Do something before enter in mainloop
	PreMainloop() error

	// Run mainloop
	Mainloop()

	// Exit mainloop
	Exit()
}

type modulectx struct {
	m     Module
	timer *Timer
}

// Module manage
type Modules struct {
	modules map[string]*modulectx
	callseq []string

	// Module alived number
	closeModule chan string
	nModule     uint

	// Signals
	signals chan os.Signal

	// log
	log      *Log
	loglevel int
}

var modules *Modules

// New a module manager instance, it is singleton
func NewModules() *Modules {
	if modules != nil {
		return modules
	}

	modules = &Modules{
		modules:     make(map[string]*modulectx),
		closeModule: make(chan string),
		signals:     make(chan os.Signal),
	}

	return modules
}

// Add User Module in module manager
func (ms *Modules) AddModule(name string, m Module) {
	mctx := &modulectx{
		m: m,
	}

	ms.modules[name] = mctx
	ms.callseq = append(ms.callseq, name)
	ms.nModule++
}

func (ms *Modules) Prefix() string {
	return ""
}

func (ms *Modules) Suffix() string {
	return ""
}

func (ms *Modules) LogLevel() int {
	return ms.loglevel
}

// Start module system
func (ms *Modules) Start() {
	ms.preInit()
	ms.log.LogError(ms, "start system ...")

	ms.log.LogInfo(ms, "init ...")
	ms.init()

	ms.log.LogInfo(ms, "pre mainloop ...")
	ms.preMainloop()

	ms.log.LogInfo(ms, "mainloop ...")
	ms.mainloop()
}

// Set log
func (ms *Modules) SetLog(log string, loglevel int) {
	ms.log = NewLog(log)
	ms.loglevel = loglevel
}

func (ms *Modules) preInit() error {
	// init signals
	// quit gracefully
	signal.Notify(ms.signals, syscall.SIGQUIT, syscall.SIGINT)

	// quit directly
	signal.Notify(ms.signals, syscall.SIGTERM)

	// reload config
	signal.Notify(ms.signals, syscall.SIGHUP)

	// reopen logs
	signal.Notify(ms.signals, syscall.SIGUSR1)

	// ignore signals
	signal.Ignore(syscall.SIGALRM)

	// pre init user modules
	for _, name := range ms.callseq {
		err := ms.modules[name].m.PreInit()
		if err != nil {
			return err
		}
	}

	if ms.log == nil {
		ms.SetLog("error.log", LOGINFO)
	}

	return nil
}

func (ms *Modules) init() {
	for _, name := range ms.callseq {
		err := ms.modules[name].m.Init()
		if err != nil {
			ms.log.LogFatal(ms, "module %s init error %s", name, err.Error())
		}

		ms.log.LogInfo(ms, "module %s init successd", name)
	}
}

func (ms *Modules) preMainloop() {
	for _, name := range ms.callseq {
		err := ms.modules[name].m.PreMainloop()
		if err != nil {
			ms.log.LogFatal(ms, "module %s pre mainloop error %s",
				name, err.Error())
		}

		ms.log.LogInfo(ms, "module %s pre mainloop successd", name)
	}
}

func (ms *Modules) close(name string) {
	ms.nModule--
	ms.closeModule <- name
}

func (ms *Modules) wrap(name string) {
	ms.modules[name].m.Mainloop()

	t := ms.modules[name].timer
	if t != nil {
		t.Stop()
	}

	ms.close(name)
}

func (ms *Modules) closeTimeout(n interface{}) {
	ms.close(n.(string))
}

func (ms *Modules) mainloop() {
	for _, name := range ms.callseq {
		go ms.wrap(name)
	}

	exit := false

	for {
		if ms.nModule == 0 {
			exit = true
		}

		if exit {
			break
		}

		select {
		case s := <-ms.signals:
			ms.log.LogInfo(ms, "get signal: %s", s.String())

			switch s {
			case syscall.SIGINT, syscall.SIGQUIT:
				ms.exit()
			case syscall.SIGTERM:
				exit = true
			case syscall.SIGHUP:
				ms.reload()
			case syscall.SIGUSR1:
				ms.reopen()
			}
		case <-ms.closeModule:
			break
		}
	}

	ms.log.LogError(ms, "system exit")
}

func (ms *Modules) reload() {
	ms.log.LogInfo(ms, "reload ...")

	if err := Reload(""); err != nil {
		ms.log.LogError(ms, "reload failed: %s", err.Error())
	}
}

func (ms *Modules) reopen() {
	ms.log.LogInfo(ms, "reopen ...")

	if err := reopenLogs(); err != nil {
		ms.log.LogError(ms, "%s", err.Error())
	}

	if err := reopenHTTPServer(); err != nil {
		ms.log.LogError(ms, "%s", err.Error())
	}
}

func (ms *Modules) exit() {
	ms.log.LogError(ms, "exiting ...")

	for name, mctx := range ms.modules {
		mctx.timer = NewTimer(5*time.Second, ms.closeTimeout, name)
		mctx.m.Exit()
	}
}
