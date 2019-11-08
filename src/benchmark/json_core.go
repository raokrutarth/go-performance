package main

import (
	"fmt"
)

func foo(s []byte) (t []byte, ok bool) {
	if len(s) < 2 || s[0] != '"' || s[len(s)-1] != '"' {
		return
	}

	s = s[1 : len(s)-1]

	return s, true
}

func main() {

	r := make([]byte, 3, 10)

	fmt.Printf("%d %d", len(r), cap(r))

	r, b := foo(r)

	fmt.Printf("%+v %t\n", r, b)

	if r == nil {
		fmt.Println("r = nil")
	}

}
