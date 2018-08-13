package main

import (
	"fmt"
	"golib"
)

type mainCtx struct {
}

func (h *mainCtx) Prefix() string {
	return "[main]"
}

func (h *mainCtx) Suffix() string {
	return "[END]"
}

func (h *mainCtx) LogLevel() int {
	return golib.LOGINFO
}

func main() {
	h := &mainCtx{}
	logger := golib.NewLog("test.log")
	if logger == nil {
		fmt.Println(h, "NewLog failed")
	}

	logger.LogDebug(h, "test debug")
	logger.LogInfo(h, "test info")
	logger.LogError(h, "test error")
	logger.LogFatal(h, "test fatal")

	logger.LogError(h, "Normal End")
}
