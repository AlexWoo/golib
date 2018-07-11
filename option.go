// Copyright (C) AlexWoo(Wu Jie) wj19840501@gmail.com
//
// go option parser

package golib

import (
	"os"
)

// Go getopt implementation as getopt in Standard C Library
//
// Example:
//
// Parse option:
//	func main() {
//		opt := golib.NewOptParser()
//		for opt.GetOpt("ab:") {
//			switch opt.Opt() {
//			case 'a':
//				fmt.Println("a", opt.OptVal())
//			case 'b':
//				fmt.Println("b", opt.OptVal())
//			case '?':
//				fmt.Println("help")
//			}
//		}
//	}
//
// Result:
//	# go run option.go -s
//	help
//
//	# go run option.go -a
//	a
//
//	# go run option.go -a 2
//	a
//	help
//
//	# go run option.go -b
//	help
//
//	# go run option.go -a -b 3
//	a
//	b 3
type OptParser struct {
	opt    byte
	optval string
	optidx int
}

// New Option Parser
func NewOptParser() *OptParser {
	return &OptParser{optidx: 1}
}

// Get Option Value
func (opt *OptParser) OptVal() string {
	return opt.optval
}

// Get Option Key, if GetOpt parse option string error, Opt will get '?'
func (opt *OptParser) Opt() byte {
	return opt.opt
}

// Parse Opt, the option string optstring may contain the following elements:
// individual characters, and characters followed by a colon to indicate an
// option argument is to follow.  For example, an option string "x" recognizes
// an option ``-x'', and an option string "x:" recognizes an option and argument
// ``-x argument''.  It does not matter to getopt() if a following argument has
// leading white space.
func (opt *OptParser) GetOpt(optstring string) bool {
	if opt.optidx >= len(os.Args) {
		return false
	}

	optstr := []byte(optstring)
	arg := []byte(os.Args[opt.optidx])
	if arg[0] != '-' { // not an option
		goto failed
	}

	if !((arg[1] == '?') || (arg[1] >= '0' && arg[1] <= '9') ||
		(arg[1] >= 'a' && arg[1] <= 'z') || (arg[1] >= 'A' && arg[1] <= 'Z')) {

		goto failed
	}

	for i := 0; i < len(optstr); i++ {
		if optstr[i] == arg[1] {
			if i+1 == len(optstr) || optstr[i+1] != ':' { // no argument
				if len(arg) != 2 { // has argument
					goto failed
				}

				opt.opt = arg[1]
				opt.optval = ""
				opt.optidx++
				return true
			} else { // has argument
				if len(arg) != 2 { // optval stick with option
					opt.opt = arg[1]
					opt.optval = string(arg[2:])
					opt.optidx++
					return true
				}

				if opt.optidx+1 == len(os.Args) {
					goto failed
				}

				opt.opt = arg[1]
				opt.optval = os.Args[opt.optidx+1]
				opt.optidx += 2
				return true
			}
		}
	}

failed:
	opt.opt = '?'
	opt.optval = ""
	opt.optidx++

	return true
}
