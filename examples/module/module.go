package main

import (
	"fmt"
	"golib"
	"net/http"
	"time"
)

type HTTPServerModule struct {
	addr           string
	location       string
	clientTimerout time.Duration
	keepalived     time.Duration
	accesslog      string
	log            *golib.Log

	server *golib.HTTPServer
}

var ms *golib.Modules

func (m *HTTPServerModule) handle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("Hello World"))
}

func (m *HTTPServerModule) PreInit() error {
	m.addr = ":8080"
	m.location = "/"
	m.clientTimerout = 10 * time.Second
	m.keepalived = 60 * time.Second
	m.accesslog = "access.log"
	m.log = golib.NewLog("error.log")

	ms.SetLog("error.log", golib.LOGINFO)

	return nil
}

func (m *HTTPServerModule) Init() error {
	s, err := golib.NewHTTPServer(m.addr, "", "", m.location, m.clientTimerout,
		m.keepalived, m.log, m.handle, m.accesslog)
	if err != nil {
		return err
	}

	m.server = s

	golib.AddReloader("HTTPServer", m)

	return nil
}

func (m *HTTPServerModule) PreMainloop() error {
	return nil
}

func (m *HTTPServerModule) Mainloop() {
	err := m.server.Start()
	fmt.Printf("HTTPServer Stop %V", err)
	//log.Println("before sleep", err)
	//time.Sleep(10 * time.Second)
	//log.Println("after sleep")
}

func (m *HTTPServerModule) Reload() error {
	fmt.Println("HTTP Server Reload")
	return fmt.Errorf("test error")
}

func (m *HTTPServerModule) Exit() {
	m.server.Close()
}

func main() {
	ms = golib.NewModules()

	ms.AddModule("httpserver", &HTTPServerModule{})

	ms.Start()
}
