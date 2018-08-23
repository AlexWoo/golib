// Copyright (C) AlexWoo(Wu Jie) wj19840501@gmail.com
//
// golib http server

package golib

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// access log record content
type record struct {
	RemoteAddr    string `json:"remote_addr"`
	RequestTime   string `json:"request_time"`
	RequestMethod string `json:"request_method"`
	Host          string `json:"host"`
	RequestURI    string `json:"request_uri"`
	Protocol      string `json:"protocol"`
	Referrer      string `json:"referrer"`
	UserAgent     string `json:"user_agent"`
	Status        int    `json:"status"`
	ElapsedTime   string `json:"elapsed_time"`
	SendBytes     int64  `json:"send_bytes"`
}

// An http ResponseWriter implementation
type httpWriter struct {
	http.ResponseWriter
	requestTime time.Time
	rec         record
}

// A ResponseWriter interface is used by an HTTP handler
// to construct an HTTP response.
// Write writes the data to the connection as part of an HTTP reply.
func (hw *httpWriter) Write(p []byte) (int, error) {
	if hw.rec.Status == 0 {
		hw.rec.Status = http.StatusOK
	}

	b, err := hw.ResponseWriter.Write(p)
	if err == nil {
		hw.rec.SendBytes += int64(b)
	}

	return b, err
}

// A ResponseWriter interface is used by an HTTP handler
// to construct an HTTP response.
// WriteHeader sends an HTTP response header with the provided status code.
func (hw *httpWriter) WriteHeader(status int) {
	hw.rec.Status = status
	hw.ResponseWriter.WriteHeader(status)
}

// The CloseNotifier interface is implemented by ResponseWriters
// which allow detecting when the underlying connection has gone away.
func (hw *httpWriter) CloseNotify() <-chan bool {
	if cn, ok := hw.ResponseWriter.(http.CloseNotifier); ok {
		return cn.CloseNotify()
	}

	return make(chan bool)
}

// The Hijacker interface is implemented by ResponseWriters
// that allow an HTTP handler to take over the connection.
func (hw *httpWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if h, ok := hw.ResponseWriter.(http.Hijacker); ok {
		return h.Hijack()
	}

	return nil, nil, errors.New("Doesn't support Hijacker interface")
}

// The Flusher interface is implemented by ResponseWriters
// that allow an HTTP handler to flush buffered data to the client.
func (hw *httpWriter) Flush() {
	if f, ok := hw.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// Pusher is the interface implemented by ResponseWriters
// that support HTTP/2 server push.
func (hw *httpWriter) Push(target string, opts *http.PushOptions) error {
	if p, ok := hw.ResponseWriter.(http.Pusher); ok {
		return p.Push(target, opts)
	}

	return errors.New("Doesn't support Pusher interface")
}

type HTTPServer struct {
	serveMux *http.ServeMux
	server   *http.Server
	handle   func(w http.ResponseWriter, req *http.Request)

	tls      bool
	certfile string
	keyfile  string

	accesslog *os.File
}

// New a HTTP Server
// addr: address HTTP Server listen, could be 127.0.0.1:8080 or :8080
// cert: TLS Certification, if use no TLS, set to empty string
// key: TLS Key, if use no TLS, set to empty string
// location: HTTPServer bind location entry
// handle: user interface for deal with http request
// accesslog: path for HTTP Server record accesslog
func NewHTTPServer(addr string, cert string, key string, location string,
	clientHeaderTimeout time.Duration, keepalived time.Duration, log *Log,
	handle func(w http.ResponseWriter, req *http.Request),
	accesslog string) (*HTTPServer, error) {

	s := &HTTPServer{
		serveMux: &http.ServeMux{},
		server: &http.Server{
			Addr:              addr,
			ReadHeaderTimeout: clientHeaderTimeout,
			IdleTimeout:       keepalived,
			ErrorLog:          log.logger,
		},
		handle:   handle,
		tls:      false,
		certfile: cert,
		keyfile:  key,
	}

	if cert != "" || key != "" {
		s.tls = true

		_, err := os.Stat(cert)
		if err != nil {
			return nil, fmt.Errorf("TLS certfication(%s) error: %s", cert, err)
		}

		_, err = os.Stat(key)
		if err != nil {
			return nil, fmt.Errorf("TLS key(%s) error: %s", key, err)
		}
	}

	s.serveMux.HandleFunc(location, s.handler)
	s.server.Handler = s.serveMux

	al, err := os.OpenFile(accesslog, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	s.accesslog = al

	return s, nil
}

// record access log
func (s *HTTPServer) log(hw *httpWriter) {
	elapsedTime := time.Since(hw.requestTime).Seconds()
	hw.rec.RequestTime = hw.requestTime.Format("02/Jan/2006 03:04:05.000")

	hw.rec.ElapsedTime = strconv.FormatFloat(elapsedTime, 'f', 6, 64)

	j, _ := json.Marshal(hw.rec)
	msg := strings.TrimSpace(string(j))
	msg += "\n"

	s.accesslog.WriteString(msg)
}

// handle a new http request
func (s *HTTPServer) handler(w http.ResponseWriter, req *http.Request) {
	writer := &httpWriter{
		ResponseWriter: w,
		requestTime:    time.Now(),
		rec: record{
			RemoteAddr:    req.RemoteAddr,
			RequestMethod: req.Method,
			Host:          req.Host,
			RequestURI:    req.RequestURI,
			Protocol:      req.Proto,
			Referrer:      req.Referer(),
			UserAgent:     req.UserAgent(),
		},
	}

	s.handle(writer, req)

	s.log(writer)
}

// Start HTTP Server, if close normal, will return nil, otherwise return error
func (s *HTTPServer) Start() error {
	var err error

	if s.certfile == "" || s.keyfile == "" { // http
		err = s.server.ListenAndServe()
	} else { // https
		err = s.server.ListenAndServeTLS(s.certfile, s.keyfile)
	}

	if err == http.ErrServerClosed {
		return nil
	} else {
		return err
	}
}

// Close HTTP Server
func (s *HTTPServer) Close() {
	s.server.Shutdown(context.Background())
}
