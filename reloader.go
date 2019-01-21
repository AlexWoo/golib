// Copyright (C) AlexWoo(Wu Jie) wj19840501@gmail.com
//
// golib reloader

package golib

import (
	"fmt"
	"sync"
)

// Reloader interface
type Reloader interface {
	Reload() error
}

var (
	reloaderm     = make(map[string]Reloader)
	reloadermLock sync.Mutex
)

// Register a reloader, if reloader has been loaded, do nothing
func AddReloader(name string, r Reloader) {
	reloadermLock.Lock()
	defer reloadermLock.Unlock()

	if _, ok := reloaderm[name]; ok { // reloader has been loaded
		return
	}

	reloaderm[name] = r
}

// Call Reload whose name is name, if name is "", reload all reloaders
func Reload(name string) error {
	reloadermLock.Lock()
	defer reloadermLock.Unlock()

	if name == "" {
		for n, r := range reloaderm {
			err := r.Reload()
			if err != nil {
				return fmt.Errorf("Reload %s failed, %s", n, err.Error())
			}
		}

		return nil
	}

	r, ok := reloaderm[name]
	if ok {
		return r.Reload()
	}

	return fmt.Errorf("%s not register as a Reloader", name)
}
