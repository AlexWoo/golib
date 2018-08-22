package main

import (
	"fmt"
	"golib"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func handle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("test"))
}

func main() {
	logger := golib.NewLog("error.log")

	s, err := golib.NewHTTPServer(":8080", "", "", "/", 10*time.Second,
		60*time.Second, logger, handle, "access.log")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	close := make(chan bool)
	closed := false

	go func() {
		err := s.Start()
		fmt.Println("!!!!!!!!", err)
		close <- true
	}()

	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT)

	for {
		select {
		case sig := <-signals:
			fmt.Println(sig)
			s.Close()
		case <-close:
			closed = true
		}

		if closed {
			break
		}
	}
}
