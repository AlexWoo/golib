// Copyright (C) AlexWoo(Wu Jie) wj19840501@gmail.com
//
// golib module

package golib

import (
	"log"
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
	Mainloop() error

	// Reload Config. Notice: not all config can be reload
	Reload() error

	// ReOpen logs for logrotate
	Reopen() error

	// Exit mainloop
	Exit()
}

type modulectx struct {
	m Module
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

// Start module system
func (ms *Modules) Start() {
	ms.preInit()

	ms.init()

	ms.preMainloop()

	ms.mainloop()
}

func (ms *Modules) preInit() {
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
			log.Fatal("module", name, "pre init error", err)
		}

		log.Println("module", name, "pre init successd")
	}
}

func (ms *Modules) init() {
	for _, name := range ms.callseq {
		err := ms.modules[name].m.Init()
		if err != nil {
			log.Fatal("module", name, "init error", err)
		}

		log.Println("module", name, "init successd")
	}
}

func (ms *Modules) preMainloop() {
	for _, name := range ms.callseq {
		err := ms.modules[name].m.PreMainloop()
		if err != nil {
			log.Fatal("module", name, "pre mainloop error", err)
		}

		log.Println("module", name, "pre mainloop successd")
	}
}

func (ms *Modules) wrap(name string) {
	err := ms.modules[name].m.Mainloop()
	if err != nil {
		log.Println("module", name, "mainloop error", err)
	}

	log.Println("module", name, "exit")
	ms.nModule--
	ms.closeModule <- name
}

func (ms *Modules) mainloop() {
	for _, name := range ms.callseq {
		go ms.wrap(name)
	}

	var t time.Timer
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
			log.Println("get signal:", s)

			switch s {
			case syscall.SIGINT, syscall.SIGQUIT:
				t := time.NewTimer(5 * time.Second)
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
		case <-t.C:
			exit = true
		}
	}

	log.Println("system exit")
}

func (ms *Modules) reload() {
	for name, mctx := range ms.modules {
		err := mctx.m.Reload()
		if err != nil {
			log.Println("module", name, "reload error", err)
		}
	}
}

func (ms *Modules) reopen() {
	for name, mctx := range ms.modules {
		err := mctx.m.Reopen()
		if err != nil {
			log.Println("module", name, "reopen error", err)
		}
	}
}

func (ms *Modules) exit() {
	for _, mctx := range ms.modules {
		mctx.m.Exit()
	}
}
