// Copyright (C) AlexWoo(Wu Jie) wj19840501@gmail.com
//
// golib config

package golib

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-ini/ini"
)

// String to int for config Enum
//	var enumT = golib.Enum{
//		"T1": T1,
//		"T2": T2,
//		"T3": T3,
//	}
//
//	e := enumT.ConfEnum("T2", T3)
type Enum map[string]int

// Convert config string to enum int
func (e Enum) ConfEnum(value string, defaultVal int) int {
	ret, ok := e[value]

	if ok {
		return ret
	}

	return defaultVal
}

func confValue(key string, s *ini.Section) string {
	return s.Key(key).Value()
}

func strToBoolean(value string) (bool, bool) {
	if strings.ToLower(value) == "true" {
		return true, true
	} else if strings.ToLower(value) == "false" {
		return false, true
	}

	return false, false
}

func defaultBoolean(value string) bool {
	ret, ok := strToBoolean(value)

	if ok {
		return ret
	}

	return false
}

func confBoolean(value string, defaultVal bool) bool {
	ret, ok := strToBoolean(value)

	if ok {
		return ret
	}

	return defaultVal
}

func confString(value string, defaultVal string) string {
	if value == "" {
		return defaultVal
	}

	return value
}

func strToUint64(value string) (uint64, bool) {
	ret, err := strconv.ParseUint(value, 10, 64)

	return ret, err == nil
}

func defaultUint64(value string) uint64 {
	ret, ok := strToUint64(value)

	if ok {
		return ret
	}

	return 0
}

func confUint64(value string, defaultVal uint64) uint64 {
	ret, ok := strToUint64(value)

	if ok {
		return ret
	}

	return defaultVal
}

func strToInt64(value string) (int64, bool) {
	ret, err := strconv.ParseInt(value, 10, 64)

	return ret, err == nil
}

func defaultInt64(value string) int64 {
	ret, ok := strToInt64(value)

	if ok {
		return ret
	}

	return 0
}

func confInt64(value string, defaultVal int64) int64 {
	ret, ok := strToInt64(value)

	if ok {
		return ret
	}

	return defaultVal
}

func strToSize(value string) (Size, bool) {
	ret, err := ParseSize(value)
	if err == nil {
		return ret, true
	}

	return ret, false
}

func defaultSize(value string) Size {
	ret, ok := strToSize(value)

	if ok {
		return ret
	}

	return 0
}

func confSize(value string, defaultVal Size) Size {
	ret, ok := strToSize(value)

	if ok {
		return ret
	}

	return defaultVal
}

func strToTimeDuration(value string) (time.Duration, bool) {
	if value == "" {
		return time.Duration(0), false
	}

	ret, err := time.ParseDuration(value)
	if err == nil {
		return ret, true
	}

	return time.Duration(0), false
}

func defaultTimeDuration(value string) time.Duration {
	ret, ok := strToTimeDuration(value)

	if ok {
		return ret
	}

	return time.Duration(0)
}

func confTimeDuration(value string, defaultVal time.Duration) time.Duration {
	ret, ok := strToTimeDuration(value)

	if ok {
		return ret
	}

	return defaultVal
}

// Reflect ini config into struct
//
// Example:
//
// test.ini:
//	[Config]
//	cstring = Hello World
//
// Parse test.ini:
//	const (
//		T1 = iota
//		T2
//		T3
//	)
//
//	var enumT = golib.Enum{
//		"T1": T1,
//		"T2": T2,
//		"T3": T3,
//	}
//
//	type Config struct {
//		CBool     bool `default:"true"`
//		CString   string
//		CUint64   uint64        `default:"100"`
//		CInt64    int64         `default:"-100"`
//		CSize     golib.Size    `default:"10M"`
//		CDuration time.Duration `default:"20s"`
//		CEnum     string        `default:"T2"`
//	}
//
//	func main() {
//		f, err := ini.Load("test.ini")
//		if err != nil {
//			fmt.Println("Load test.ini failed", err)
//			return
//		}
//
//		config := &Config{}
//		if !golib.Config(f, "Config", config) {
//			fmt.Println("Parse config failed")
//			return
//		}
//
//		fmt.Println(config)
//	}
//
// Result:
//	&{true Hello World 100 -100 10M 20s T2}
func Config(f *ini.File, secName string, it interface{}) bool {
	if secName == "DEFAULT" {
		secName = ""
	}

	s := f.Section(secName)

	t := reflect.TypeOf(it).Elem()
	v := reflect.ValueOf(it).Elem()
	n := t.NumField()

	for i := 0; i < n; i++ {
		field := t.Field(i)
		value := v.Field(i)
		fn := field.Name
		ft := field.Type.Name()
		fd := field.Tag.Get("default")

		switch ft {
		case "bool":
			confV := confValue(strings.ToLower(fn), s)
			value.SetBool(confBoolean(confV, defaultBoolean(fd)))
		case "string":
			confV := confValue(strings.ToLower(fn), s)
			value.SetString(confString(confV, fd))
		case "uint64":
			confV := confValue(strings.ToLower(fn), s)
			value.SetUint(confUint64(confV, defaultUint64(fd)))
		case "int64":
			confV := confValue(strings.ToLower(fn), s)
			value.SetInt(confInt64(confV, defaultInt64(fd)))
		case "Size":
			confV := confValue(strings.ToLower(fn), s)
			value.SetInt(int64(confSize(confV, defaultSize(fd))))
		case "Duration":
			confV := confValue(strings.ToLower(fn), s)
			value.SetInt(int64(confTimeDuration(confV,
				defaultTimeDuration(fd))))
		default:
			fmt.Printf("Unsuppoted config, secName: %s, name: %s, type: %s\n",
				secName, fn, ft)
			return false
		}
	}

	return true
}
