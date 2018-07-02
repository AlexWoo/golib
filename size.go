// Copyright (C) AlexWoo(Wu Jie) wj19840501@gmail.com
//
// golib size

package golib

import (
	"errors"
	"strconv"
)

// A Size represents size type such as Byte, KByte, MByte, GByte, TByte, PByte
type Size int64

const step = 1024

// To conver integer number of units to a Size, multiply:
//	size := 10
//	fmt.Println(golib.Size(size)*golib.MByte) // prints 10Mbyte
const (
	Byte  Size = 1
	KByte      = step * Byte
	MByte      = step * KByte
	GByte      = step * MByte
	TByte      = step * GByte
	PByte      = step * TByte
)

var unitMap = map[string]int64{
	"B": int64(Byte),
	"k": int64(KByte),
	"K": int64(KByte),
	"m": int64(MByte),
	"M": int64(MByte),
	"g": int64(GByte),
	"G": int64(GByte),
	"t": int64(TByte),
	"T": int64(TByte),
	"p": int64(PByte),
	"P": int64(PByte),
}

var unitBase = []string{"B", "K", "M", "G", "T", "P"}

func (s Size) String() string {
	ret := int64(s)
	base := 0

	for {
		if len(unitBase) == base {
			break
		}

		if ret%step == 0 {
			ret /= step
			base++
			continue
		}

		break
	}

	return strconv.FormatInt(int64(ret), 10) + unitBase[base]
}

func ParseSize(s string) (Size, error) {
	if s == "" {
		return 0, errors.New("null string")
	}

	orig := s
	slen := len(orig)
	u := s[slen-1:]

	unit, ok := unitMap[u]
	if ok {
		s = s[:slen-1]
	} else {
		unit = int64(Byte)
	}

	d, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}

	if d == 0 {
		return 0, nil
	}

	d = d * unit
	if d == 0 {
		return 0, errors.New("over flow")
	}

	return Size(d), nil
}
