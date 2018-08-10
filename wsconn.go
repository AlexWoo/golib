// Copyright (C) AlexWoo(Wu Jie) wj19840501@gmail.com
//
// golib websocket connection

package golib

import (
	"net"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Common Connection Type
type Conn interface {
	Accept()
	Send(data []byte)
	Close()

	// for log ctx
	Prefix() string
	Suffix() string
}

// Webscoket Connection, support websocket channel reuse with same name,
// support for data buffer and resent if websocket connection reconnect
type WSConn struct {
	conn        *websocket.Conn
	name        string
	url         string
	connTimeout time.Duration
	maxRetries  int
	sendq       chan []byte
	recvq       chan []byte
	reconnect   chan bool
	quit        chan bool
	buf         []byte
	log         *Log
	handler     func(c Conn, data []byte)
}

var (
	wsconns     = make(map[string]*WSConn)
	wsconnsLock sync.Mutex
)

func (c *WSConn) connect() bool {
	dialer := websocket.DefaultDialer
	dialer.HandshakeTimeout = c.connTimeout
	retries := 0

	for {
		if retries >= c.maxRetries {
			c.Close()
			return false
		}

		conn, _, err := dialer.Dial(c.url, nil)
		if err == nil { // connect successd
			go c.read()
			go c.loop()
			c.conn = conn

			return true
		}

		c.log.LogError(c, "connect err: %s", err)

		if websocket.IsCloseError(err) {
			e := err.(*websocket.CloseError)

			switch e.Code {
			case websocket.CloseTLSHandshake:
				if strings.HasPrefix(c.url, "ws://") {
					strings.Replace(c.url, "ws://", "wss://", 1)
					continue
				}
			}
		} else {
			e, ok := err.(net.Error)
			if ok {
				if e.Timeout() {
					retries++
					continue
				}
			}
		}

		retries++
		time.Sleep(c.connTimeout)
	}
}

func (c *WSConn) read() {
	for {
		conn := c.conn
		_, data, err := conn.ReadMessage()
		if err == nil {
			c.recvq <- data
			continue
		}

		c.log.LogError(c, "read err: %s", err)

		c.recvq <- []byte{}

		return
	}
}

func (c *WSConn) write(data []byte) bool {
	err := c.conn.WriteMessage(websocket.TextMessage, data)
	if err == nil {
		c.buf = []byte{}
		return true
	}

	c.log.LogError(c, "write err: %s", err)

	c.buf = data

	if websocket.IsCloseError(err) {
		e := err.(*websocket.CloseError)

		switch e.Code {
		case websocket.CloseMessageTooBig:
			c.buf = []byte{}
		}
	}

	return false
}

func (c *WSConn) loop() {
	var data []byte

	if len(c.buf) != 0 {
		data = c.buf
		if !c.write(data) {
			return
		}
	}

	for {
		select {
		case data = <-c.recvq:
			if len(data) == 0 {
				c.reconnect <- true
				return
			}
			c.handler(c, data)
		case data = <-c.sendq:
			if !c.write(data) {
				c.reconnect <- true
				return
			}
		}
	}
}

func (c *WSConn) dial() {
	go func() {
		if !c.connect() {
			return
		}
		c.log.LogInfo(c, "connect successd")

		for {
			select {
			case <-c.reconnect:
				if !c.connect() {
					return
				}
				c.log.LogInfo(c, "reconnect successd")
			case <-c.quit:
				c.conn.Close()
				return
			}
		}
	}()
}

// New and return a websocket client connection instance
//
// name is connection name.
// url is url websocket client connect to, must start with ws:// or wss://.
// connTimeout is connect timeout websocket client connect to server.
// maxRetries is max times to try if websocket client connect to server failed.
// qsize is send and receive channel size.
// handler is callback function when websocket client receive data
func NewWSClient(name string, url string, connTimeout time.Duration,
	maxRetries int, qsize uint64, handler func(c Conn, data []byte),
	log *Log) *WSConn {

	if !strings.HasPrefix(url, "ws://") && !strings.HasPrefix(url, "wss://") {
		return nil
	}

	wsconnsLock.Lock()

	conn := wsconns[name]
	if conn != nil {
		wsconnsLock.Unlock()
		return conn
	}

	conn = &WSConn{
		name:        name,
		url:         url,
		connTimeout: connTimeout,
		maxRetries:  maxRetries,
		sendq:       make(chan []byte, qsize),
		recvq:       make(chan []byte, qsize),
		reconnect:   make(chan bool),
		quit:        make(chan bool, 1),
		log:         log,
		handler:     handler,
	}
	wsconns[name] = conn
	wsconnsLock.Unlock()

	conn.dial()

	return conn
}

// New and return a websocket server connection instance
//
// name is connection name.
// c is websocket.Conn instance bind with websocket server.
// qsize is send and receive channel size.
// handler is callback function when websocket server receive data
func NewWSServer(name string, c *websocket.Conn, qsize uint64,
	handler func(c Conn, data []byte), log *Log) *WSConn {

	wsconnsLock.Lock()

	conn := wsconns[name]
	if conn != nil {
		conn.conn = c
	} else {
		conn = &WSConn{
			conn:      c,
			name:      name,
			sendq:     make(chan []byte, qsize),
			recvq:     make(chan []byte, qsize),
			reconnect: make(chan bool),
			quit:      make(chan bool, 1),
			log:       log,
			handler:   handler,
		}
		wsconns[name] = conn
	}
	wsconnsLock.Unlock()

	return conn
}

// Wait for connection
func (c *WSConn) Accept() {
	go c.read()
	go c.loop()

	for {
		select {
		case <-c.reconnect:
			return
		case <-c.quit:
			c.conn.Close()
			return
		}
	}
}

// Send data through websocket connection
// data need to send, now only support text message
func (c *WSConn) Send(data []byte) {
	c.sendq <- data
}

// Close websocket connection
func (c *WSConn) Close() {
	wsconnsLock.Lock()
	delete(wsconns, c.name)
	wsconnsLock.Unlock()

	c.quit <- true
}

func (c *WSConn) Prefix() string {
	return "[websocket]"
}

func (c *WSConn) Suffix() string {
	suf := ", Websocket[" + c.name + "]"
	if c.url != "" { // websocket client
		suf += " Url: " + c.url + " Client: " + c.conn.LocalAddr().String() +
			" Server: " + c.conn.RemoteAddr().String()
	} else {
		suf += " Client: " + c.conn.RemoteAddr().String() +
			" Server: " + c.conn.LocalAddr().String()
	}

	return suf
}
