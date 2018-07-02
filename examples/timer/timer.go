package main

import (
	"golib"
	"log"
	"time"
)

func handler(p interface{}) {
	log.Println("timer handler")
	wait := p.(chan bool)
	wait <- true
}

func main() {
	log.Println("before timer set")

	wait := make(chan bool)
	t := golib.NewTimer(10*time.Second, handler, wait)
	//t.Stop()
	//t.Reset(5 * time.Second)

	<-wait
	log.Println("after timer set", t)
}
