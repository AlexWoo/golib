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

	return nil
}

func (m *HTTPServerModule) Init() error {
	s, err := golib.NewHTTPServer(m.addr, "", "", m.location, m.clientTimerout,
		m.keepalived, m.log, m.handle, m.accesslog)
	if err != nil {
		return err
	}

	m.server = s

	return nil
}

func (m *HTTPServerModule) PreMainloop() error {
	return nil
}

func (m *HTTPServerModule) Mainloop() error {
	err := m.server.Start()
	fmt.Println(err)

	return err
}

func (m *HTTPServerModule) Reload() error {
	return nil
}

func (m *HTTPServerModule) Reopen() error {
	return nil
}

func (m *HTTPServerModule) Exit() {
	m.server.Close()
}

func main() {
	ms := golib.NewModules()

	ms.AddModule("httpserver", &HTTPServerModule{})

	ms.Start()
}
