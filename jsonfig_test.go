package golib_test

import (
	"fmt"
	"golib"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/tidwall/gjson"
)

type JConfig struct {
	CBool     bool
	CString   string
	CUint64   uint64
	CInt64    int64
	CSize     golib.Size
	CDuration time.Duration
	CEnum     string
}

func TestJConfig(t *testing.T) {
	f, err := os.Open("test/config.json")
	if err != nil {
		t.Error("Load config.json failed", err)
		return
	}

	j, _ := ioutil.ReadAll(f)
	if !gjson.ValidBytes(j) {
		t.Error("Invalid json", string(j))
	}

	config := &JConfig{}
	err = golib.JsonConfig(string(j), config)
	if err != nil {
		t.Error("Parse config failed:", err)
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

type JUnConfig struct {
	CBool     bool          `json:"cbool" default:"true"`
	CString   string        `default:"Hello World"`
	CUint64   uint64        `default:"100"`
	CInt64    int64         `default:"-100"`
	CSize     golib.Size    `default:"10M"`
	CDuration time.Duration `default:"20s"`
	CEnum     string        `default:"T2"`
}

func TestJUnconfig(t *testing.T) {
	config := &JUnConfig{}
	err := golib.JsonConfigFile("test/unconfig.json", config)
	if err != nil {
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

type JUnsuppoted struct {
	Unsupport string
}

func TestJUnsupportedConf(t *testing.T) {
	config := &JUnsuppoted{}
	err := golib.JsonConfigFile("test/unsupported.json", config)
	if err != nil {
		t.Error("parse config failed", err)
		return
	}

	fmt.Println(config.Unsupport)
}

type JUnsuppotedType struct {
	Unsupport map[int]string
}

func TestJUnsupportedType(t *testing.T) {
	config := &JUnsuppotedType{}
	err := golib.JsonConfigFile("test/unsuppotedtype.json", config)
	if err != nil {
		fmt.Println("parse config failed")
		return
	}

	t.Error("parse successd")
}
