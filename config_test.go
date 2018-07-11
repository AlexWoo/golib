package golib_test

import (
	"fmt"
	"golib"
	"testing"
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
	CBool     bool
	CString   string
	CUint64   uint64
	CInt64    int64
	CSize     golib.Size
	CDuration time.Duration
	CEnum     string
}

func TestConfig(t *testing.T) {
	f, err := ini.Load("test/test.ini")
	if err != nil {
		t.Error("Load test.ini failed", err)
		return
	}

	config := &Config{}
	if !golib.Config(f, "Config", config) {
		t.Error("Parse config failed")
		return
	}

	if !config.CBool {
		t.Error("cbool config failed, expect:", true, "get:", config.CBool)
	}
	fmt.Println("config.CBool", config.CBool)

	if config.CString != "Hello World" {
		t.Error("cstring config failed, expect:", "Hello World",
			"get:", config.CString)
	}
	fmt.Println("config.CString", config.CString)

	if config.CUint64 != 100 {
		t.Error("cuint64 config failed, expect:", 100, "get:", config.CUint64)
	}
	fmt.Println("config.CUint64", config.CUint64)

	if config.CInt64 != -100 {
		t.Error("cint64 config failed, expect:", -100, "get:", config.CInt64)
	}
	fmt.Println("config.CInt64", config.CInt64)

	if config.CSize != 10*golib.MByte {
		t.Error("csize config failed, expect:", "10M", "get:", config.CSize)
	}
	fmt.Println("config.CSize", config.CSize)

	if config.CDuration != 20*time.Second {
		t.Error("cduration config failed, expect:", "20s",
			"get:", config.CDuration)
	}
	fmt.Println("config.CDuration", config.CDuration)

	e := enumT.ConfEnum(config.CEnum, T2)
	if e != T1 {
		t.Error("cenum config failed, expect:", "T1", "get:", e)
	}
	fmt.Println("config.CEnum", e)

	fmt.Println("---------------------------------------------------")
}

type UnConfig struct {
	CBool     bool          `default:"true"`
	CString   string        `default:"Hello World"`
	CUint64   uint64        `default:"100"`
	CInt64    int64         `default:"-100"`
	CSize     golib.Size    `default:"10M"`
	CDuration time.Duration `default:"20s"`
	CEnum     string        `default:"T2"`
}

func TestUnconfig(t *testing.T) {
	f, err := ini.Load("test/test.ini")
	if err != nil {
		t.Error("Load test.ini failed", err)
		return
	}

	config := &UnConfig{}
	if !golib.Config(f, "UnConfig", config) {
		t.Error("Parse config failed")
		return
	}

	if !config.CBool {
		t.Error("cbool config failed, expect:", true, "get:", config.CBool)
	}
	fmt.Println("config.CBool", config.CBool)

	if config.CString != "Hello World" {
		t.Error("cstring config failed, expect:", "Hello World",
			"get:", config.CString)
	}
	fmt.Println("config.CString", config.CString)

	if config.CUint64 != 100 {
		t.Error("cuint64 config failed, expect:", 100, "get:", config.CUint64)
	}
	fmt.Println("config.CUint64", config.CUint64)

	if config.CInt64 != -100 {
		t.Error("cint64 config failed, expect:", -100, "get:", config.CInt64)
	}
	fmt.Println("config.CInt64", config.CInt64)

	if config.CSize != 10*golib.MByte {
		t.Error("csize config failed, expect:", "10M", "get:", config.CSize)
	}
	fmt.Println("config.CSize", config.CSize)

	if config.CDuration != 20*time.Second {
		t.Error("cduration config failed, expect:", "20s",
			"get:", config.CDuration)
	}
	fmt.Println("config.CDuration", config.CDuration)

	e := enumT.ConfEnum(config.CEnum, T3)
	if e != T2 {
		t.Error("cenum config failed, expect:", "T2", "get:", e)
	}
	fmt.Println("config.CEnum", e)

	fmt.Println("---------------------------------------------------")
}

type Unsuppoted struct {
	Unsupport string
}

func TestUnsupportedConf(t *testing.T) {
	f, err := ini.Load("test/test.ini")
	if err != nil {
		t.Error("load test.ini failed", err)
		return
	}

	config := &Unsuppoted{}
	if !golib.Config(f, "Unsuppoted", config) {
		t.Error("parse config failed")
		return
	}

	fmt.Println(config.Unsupport)
}

type UnsuppotedType struct {
	Unsupport map[int]string
}

func TestUnsupportedType(t *testing.T) {
	f, err := ini.Load("test/test.ini")
	if err != nil {
		t.Error("load test.ini failed", err)
		return
	}

	config := &UnsuppotedType{}
	if !golib.Config(f, "UnsuppotedType", config) {
		fmt.Println("parse config failed")
		return
	}

	t.Error("parse successd")
}
