package main

import (
	"fmt"
	"golib"
)

type traceObj struct {
	traceid   string
	currentid string
	parentid  string
	protocol  string
	abstract  string
	detail    string
}

func (t *traceObj) GetTraceId() string {
	return t.traceid
}

func (t *traceObj) GetCurrentId() string {
	return t.currentid
}

func (t *traceObj) GetParentId() string {
	return t.parentid
}

func (t *traceObj) GetType() string {
	return t.protocol
}

func (t *traceObj) GetAbstract() string {
	return t.abstract
}

func (t *traceObj) GetDetail() string {
	return t.detail
}

func main() {
	log, err := golib.NewTraceLog("trace.log")
	if err != nil {
		fmt.Println(err)
		return
	}

	trace1 := &traceObj{
		traceid:   golib.NewRandId(),
		currentid: golib.NewRandId(),
		parentid:  golib.NewBlankId(),
		protocol:  "Test",
		abstract:  "INVITE",
		detail:    "Hello World",
	}

	trace2 := &traceObj{
		traceid:   trace1.traceid,
		currentid: golib.NewRandId(),
		parentid:  trace1.currentid,
		protocol:  "Test",
		abstract:  "INVITE_200",
		detail:    "Hello World1",
	}

	log.Trace(trace1)
	log.Trace(trace2)
}
