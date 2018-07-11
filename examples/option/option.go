package main

import (
	"fmt"
	"golib"
)

func main() {
	opt := golib.NewOptParser()
	for opt.GetOpt("ab:") {
		switch opt.Opt() {
		case 'a':
			fmt.Println("a", opt.OptVal())
		case 'b':
			fmt.Println("b", opt.OptVal())
		case '?':
			fmt.Println("help")
		}
	}
}
