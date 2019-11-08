package main

import (
	"fmt"
)

type S struct {
	num  int
	word string
}

var check *S

func main() {
	fmt.Println(check == nil)
}
