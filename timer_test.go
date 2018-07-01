package golib_test

import (
	"fmt"
	"testing"
	"time"

	"golib"
)

func TestNilStop(t *testing.T) {
	timer := &golib.Timer{}

	timer.Stop()
}

func TestNilReset(t *testing.T) {
	timer := &golib.Timer{}

	timer.Reset(2 * time.Second)
}

func timer(p interface{}) {
	ret := p.(chan time.Time)

	now := time.Now()
	fmt.Println("timer execute", now)
	ret <- now
}

func TestNormalTimer(t *testing.T) {
	fmt.Println("")
	tick := 2 * time.Second
	ret := make(chan time.Time)

	golib.NewTimer(tick, timer, ret)

	fmt.Println("[TestNormalTimer]before timer execute", time.Now())
	<-ret
	fmt.Println("[TestNormalTimer]after timer execute", time.Now())
}

func TestNormalTimerReset(t *testing.T) {
	fmt.Println("")
	tick := 2 * time.Second
	ret := make(chan time.Time)

	timer := golib.NewTimer(tick, timer, ret)
	fmt.Println("[TestNormalTimerReset]before sleep", time.Now())
	time.Sleep(1 * time.Second)
	fmt.Println("[TestNormalTimerReset]after sleep", time.Now())
	timer.Reset(3 * time.Second)
	<-ret
	fmt.Println("[TestNormalTimerReset]after timer execute", time.Now())
}

func TestNormalTimerStop(t *testing.T) {
	fmt.Println("")
	tick := 2 * time.Second
	ret := make(chan time.Time)

	timer := golib.NewTimer(tick, timer, ret)
	fmt.Println("[TestNormalTimerStop]before stop", time.Now())
	timer.Stop()

	ti := time.After(3 * time.Second)
	select {
	case <-ti:
		fmt.Println("[TestNormalTimerStop]after stop", time.Now())
	case <-ret:
		t.Error("[TestNormalTimerStop]timer not stopped")
	}
}

func TestTimerStopTwice(t *testing.T) {
	fmt.Println("")
	tick := 2 * time.Second
	ret := make(chan time.Time)

	timer := golib.NewTimer(tick, timer, ret)
	fmt.Println("[TestTimerStopTwice]before stop", time.Now())
	timer.Stop()

	ti := time.After(3 * time.Second)
	select {
	case <-ti:
		fmt.Println("[TestTimerStopTwice]after stop", time.Now())
	case <-ret:
		t.Error("[TestTimerStopTwice]timer not stopped")
	}

	timer.Stop()
	fmt.Println("[TestTimerStopTwice]after stop twice", time.Now())
}

func TestTimerStopAfterExpire(t *testing.T) {
	fmt.Println("")
	tick := 2 * time.Second
	ret := make(chan time.Time)

	timer := golib.NewTimer(tick, timer, ret)
	fmt.Println("[TestTimerStopAfterExpire]before execute", time.Now())

	<-ret
	fmt.Println("[TestTimerStopAfterExpire]after execute", time.Now())

	timer.Stop()
	fmt.Println("[TestTimerStopAfterExpire]after stop twice", time.Now())
}

func TestTimerResetAfterStop(t *testing.T) {
	fmt.Println("")
	tick := 2 * time.Second
	ret := make(chan time.Time)

	timer := golib.NewTimer(tick, timer, ret)
	fmt.Println("[TestTimerResetAfterStop]before stop", time.Now())
	timer.Stop()

	ti := time.After(3 * time.Second)
	select {
	case <-ti:
		fmt.Println("[TestTimerResetAfterStop]after stop", time.Now())
	case <-ret:
		t.Error("[TestTimerResetAfterStop]timer not stopped")
	}

	timer.Reset(1 * time.Second)
	fmt.Println("[TestTimerResetAfterStop]before reset", time.Now())
	ti = time.After(3 * time.Second)
	select {
	case <-ti:
		t.Error("[TestTimerResetAfterStop]timer reset failed")
	case <-ret:
		fmt.Println("[TestTimerResetAfterStop]after reset", time.Now())
	}
}

func TestTimeResetAfterExpire(t *testing.T) {
	fmt.Println("")
	tick := 2 * time.Second
	ret := make(chan time.Time)

	timer := golib.NewTimer(tick, timer, ret)
	fmt.Println("[TestTimeResetAfterExpire]before execute", time.Now())

	<-ret
	fmt.Println("[TestTimeResetAfterExpire]after execute", time.Now())

	timer.Reset(1 * time.Second)
	fmt.Println("[TestTimeResetAfterExpire]before reset", time.Now())
	ti := time.After(3 * time.Second)
	select {
	case <-ti:
		t.Error("[TestTimeResetAfterExpire]timer reset failed")
	case <-ret:
		fmt.Println("[TestTimeResetAfterExpire]after reset", time.Now())
	}
}
