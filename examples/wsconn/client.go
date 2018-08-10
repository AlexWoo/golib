package main

import (
	"fmt"
	"golib"
	"strconv"
	"time"
)

func send(c golib.Conn) {
	for i := 0; i < 100000; i++ {
		c.Send([]byte(strconv.Itoa(i)))
		time.Sleep(10 * time.Millisecond)
	}
}

func handler(c golib.Conn, data []byte) {
	fmt.Println(string(data))
}

func main() {
	logger := golib.NewLog("client.log", golib.LOGINFO)

	conn := golib.NewWSClient("a", "ws://127.0.0.1:8080", 3*time.Second, 3,
		1024, handler, logger)
	if conn == nil {
		return
	}

	wait := make(chan bool)

	go send(conn)

	<-wait
}
