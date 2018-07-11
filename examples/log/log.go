package main

import (
	"fmt"
	"golib"
)

type mainLogHandle struct {
}

func (h *mainLogHandle) Prefix() string {
	return "[main]"
}

func (h *mainLogHandle) Suffix() string {
	return "[END]"
}

func main() {
	h := &mainLogHandle{}
	logger := golib.NewLog(h, "test.log", golib.LOGINFO)
	if logger == nil {
		fmt.Println("NewLog failed")
	}

	logger.LogDebug("test debug")
	logger.LogInfo("test info")
	logger.LogError("test error")
	logger.LogFatal("test fatal")

	logger.LogError("Normal End")
}
