// Copyright (C) AlexWoo(Wu Jie) wj19840501@gmail.com
//
// golib json config

package golib

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

func jsonBoolean(key string, m map[string]interface{}, dv bool) bool {
	v, ok := m[key].(bool)
	if !ok {
		return dv
	}

	return v
}

func jsonString(key string, m map[string]interface{}, dv string) string {
	v, ok := m[key].(string)
	if !ok {
		return dv
	}

	return v
}

func jsonUint64(key string, m map[string]interface{}, dv uint64) uint64 {
	v, ok := m[key].(float64)
	if !ok {
		return dv
	}

	if v < 0 {
		return dv
	}

	return uint64(v)
}

func jsonInt64(key string, m map[string]interface{}, dv int64) int64 {
	v, ok := m[key].(float64)
	if !ok {
		return dv
	}

	return int64(v)
}

func jsonSize(key string, m map[string]interface{}, dv Size) Size {
	v, ok := m[key].(string)
	if !ok {
		return dv
	}

	return confSize(v, dv)
}

func jsonDuration(key string, m map[string]interface{},
	dv time.Duration) time.Duration {

	v, ok := m[key].(string)
	if !ok {
		return dv
	}

	return confTimeDuration(v, dv)
}

// Reflect json config into struct
//
// Example:
//
// test.ini:
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
//		config := &Config{}
//		err := golib.JsonConfigFile("test/config.json", config)
//		if err != nil {
//			fmt.Println("Parse config failed", err)
//			return
//		}
//
//		fmt.Println(config)
//	}
//
// Result:
//	&{true Hello World 100 -100 10M 20s T2}
func JsonConfig(json string, it interface{}) error {
	s, ok := gjson.Parse(json).Value().(map[string]interface{})
	if !ok {
		return fmt.Errorf("not a json map: %s", string(json))
	}

	t := reflect.TypeOf(it).Elem()
	v := reflect.ValueOf(it).Elem()
	n := t.NumField()

	for i := 0; i < n; i++ {
		field := t.Field(i)
		value := v.Field(i)
		fn := strings.ToLower(field.Name)
		ft := field.Type.Name()
		fd := field.Tag.Get("default")
		fj := field.Tag.Get("json")
		if fj != "" {
			fn = fj
		}

		switch ft {
		case "bool":
			value.SetBool(jsonBoolean(fn, s, defaultBoolean(fd)))
		case "string":
			value.SetString(jsonString(fn, s, fd))
		case "uint64":
			value.SetUint(jsonUint64(fn, s, defaultUint64(fd)))
		case "int64":
			value.SetInt(jsonInt64(fn, s, defaultInt64(fd)))
		case "Size":
			value.SetInt(int64(jsonSize(fn, s, defaultSize(fd))))
		case "Duration":
			value.SetInt(int64(jsonDuration(fn, s, defaultTimeDuration(fd))))
		default:
			return fmt.Errorf("Unsuppoted json config, name: %s, type: %s\n",
				fn, ft)
		}
	}

	return nil
}

func JsonConfigFile(path string, it interface{}) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	json, _ := ioutil.ReadAll(f)
	if !gjson.ValidBytes(json) {
		return fmt.Errorf("Invalid json: %s", string(json))
	}

	return JsonConfig(string(json), it)
}
