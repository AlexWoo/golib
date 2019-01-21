package main

import (
	"fmt"
	"golib"
	"os"
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
	//logger.LogFatal(h, "test fatal")

	logger.LogError(h, "Normal End")

	fmt.Printf("!!!!!logger1: %p\n", logger)

	logger2 := golib.NewLog("test.log")
	fmt.Printf("!!!!!logger2: %p\n", logger2)

	os.Rename("test.log", "test.bak.log")

	golib.ReopenLogs()

	logger.LogInfo(h, "test fatal 20")
	logger2.LogInfo(h, "test fatal 21")
	logger2.LogInfo(h, "test fatal 22")
}
