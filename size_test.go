package golib_test

import (
	"fmt"
	"golib"
	"testing"
)

func TestPrint(t *testing.T) {
	s := 3 * golib.PByte
	fmt.Println(s, int64(s))

	s = 524288
	fmt.Println(s, int64(s))
}

func TestParse(t *testing.T) {
	s, err := golib.ParseSize("1000")
	if err != nil || s != 1000*golib.Byte {
		t.Errorf("Parse Byte error %s", err)
	}

	s, err = golib.ParseSize("1000B")
	if err != nil || s != 1000*golib.Byte {
		t.Errorf("Parse Byte error %s", err)
	}

	s, err = golib.ParseSize("1000k")
	if err != nil || s != 1000*golib.KByte {
		t.Errorf("Parse Byte error %s", err)
	}

	s, err = golib.ParseSize("1000M")
	if err != nil || s != 1000*golib.MByte {
		t.Errorf("Parse Byte error %s", err)
	}

	s, err = golib.ParseSize("-1000g")
	if err != nil || s != -1000*golib.GByte {
		t.Errorf("Parse Byte error %s", err)
	}

	s, err = golib.ParseSize("1000T")
	if err != nil || s != 1000*golib.TByte {
		t.Errorf("Parse Byte error %s", err)
	}

	s, err = golib.ParseSize("1000P")
	if err != nil || s != 1000*golib.PByte {
		t.Errorf("Parse Byte error %s", err)
	}

	s, err = golib.ParseSize("0")
	if err != nil || s != 0*golib.Size(0) {
		t.Errorf("Parse Byte error %s", err)
	}

	s, err = golib.ParseSize("0P")
	if err != nil || s != 0*golib.Size(0) {
		t.Errorf("Parse Byte error %s", err)
	}

	s, err = golib.ParseSize("1000b")
	if err == nil || s != golib.Size(0) {
		t.Errorf("Parse Byte error %s", err)
	}
	fmt.Println(err)
	fmt.Println()

	s, err = golib.ParseSize("1.1b")
	if err == nil || s != golib.Size(0) {
		t.Errorf("Parse Byte error %s", err)
	}
	fmt.Println(err)
	fmt.Println()

	s, err = golib.ParseSize("B")
	if err == nil || s != golib.Size(0) {
		t.Errorf("Parse Byte error %s", err)
	}
	fmt.Println(err)
	fmt.Println()

	s, err = golib.ParseSize("")
	if err == nil || s != golib.Size(0) {
		t.Errorf("Parse Byte error %s", err)
	}
	fmt.Println(err)
	fmt.Println()

	s, err = golib.ParseSize("1000000000000000000P")
	if err == nil || s != golib.Size(0) {
		t.Errorf("Parse Byte error %s", err)
	}
	fmt.Println(err)
	fmt.Println()
}
