package main

import (
	"fmt"

	"golib"
)

func main() {
	lf1, _ := golib.NewLogFile("test1.log")
	lf2, _ := golib.NewLogFile("test1.log")
	lf3, _ := golib.NewLogFile("test2.log")

	fmt.Println(lf1.Fd())
	fmt.Println(lf2.Fd())
	fmt.Println(lf3.Fd())

	lf3.WriteString("test1....")
	lf3.Close()
	fmt.Println(lf3.Fd())
	lf3.WriteString("test2....")
}
