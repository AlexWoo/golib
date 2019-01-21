package main

import (
	"fmt"
	"golib"
)

type Test struct {
}

func (t *Test) Reload() error {
	fmt.Println("Reload")
	return nil
}

func main() {
	t1 := &Test{}
	t2 := &Test{}

	golib.AddReloader("Test1", t1)
	golib.AddReloader("Test2", t2)

	golib.Reload("")

	fmt.Println("-----------------------")

	golib.Reload("Test1")
}
