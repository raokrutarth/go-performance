package main

import (
	"fmt"
	"runtime/debug"
)

func foo(i int) {
	a := 0
	for {
		a += i
		s := fmt.Sprintf("%d", a)
		if s == "" {
			break
		}
	}
}

func main() {
	debug.SetTraceback("crash")

	go foo(4)
	go foo(8)
	go foo(10)

	x := []int{1, 2, 3}
	fmt.Println(x[5])
}
