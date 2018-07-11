package main

import (
	"fmt"
	"golib"
	"time"

	"github.com/go-ini/ini"
)

const (
	T1 = iota
	T2
	T3
)

var enumT = golib.Enum{
	"T1": T1,
	"T2": T2,
	"T3": T3,
}

type Config struct {
	CBool     bool `default:"true"`
	CString   string
	CUint64   uint64        `default:"100"`
	CInt64    int64         `default:"-100"`
	CSize     golib.Size    `default:"10M"`
	CDuration time.Duration `default:"20s"`
	CEnum     string        `default:"T2"`
}

func main() {
	f, err := ini.Load("test.ini")
	if err != nil {
		fmt.Println("Load test.ini failed", err)
		return
	}

	config := &Config{}
	if !golib.Config(f, "Config", config) {
		fmt.Println("Parse config failed")
		return
	}

	fmt.Println(config)
}
