// Copyright (C) AlexWoo(Wu Jie) wj19840501@gmail.com
//
// golib trace

package golib

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

type TraceObject interface {
	GetTraceId() string
	GetCurrentId() string
	GetParentId() string
	GetType() string
	GetAbstract() string
	GetDetail() string
}

// TraceFile
type TraceLog struct {
	lf     *LogFile
	ids    map[string]TraceObject
	idLock sync.RWMutex
}

// Create a 40 bytes blank string 40 * '0'
func NewBlankId() string {
	b := make([]byte, 20)

	return fmt.Sprintf("%02x", b)
}

// Create a 40 bytes random string
func NewRandId() string {
	r1 := rand.Int63()
	r2 := rand.Int63()
	r3 := rand.Int31()

	return fmt.Sprintf("%016x%016x%08x", r1, r2, r3)
}

// Create a trance log instance
func NewTraceLog(path string) (*TraceLog, error) {
	lf, err := NewLogFile(path)
	if err != nil {
		return nil, err
	}

	tf := &TraceLog{
		lf: lf,
	}

	return tf, nil
}

// Write Trace Log
func (tl *TraceLog) Trace(to TraceObject) {
	trace := make(map[string]interface{})

	trace["time"] = time.Now().UnixNano() / 1000000 // Use millisecond
	trace["traceid"] = to.GetTraceId()
	trace["id"] = to.GetCurrentId()
	trace["pid"] = to.GetParentId()
	trace["type"] = to.GetType()
	trace["abstract"] = to.GetAbstract()
	trace["detail"] = to.GetDetail()

	b, _ := json.Marshal(trace)

	msg := strings.TrimSpace(string(b))
	msg += "\n"

	tl.lf.WriteString(msg)
}
