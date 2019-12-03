package main

import (
	"fmt"
	"time"

	"golang.performance.com/telemetry"
)

const (
	// telemetry constants
	cStrLenTag = "closure_str_len"
)

func main() {
	telemetry.Initialize()

	// global getter test
	run(GetGlobalObj)

	// closure getter test
	// run(GetClosuredObj)
}

func run(getter func() *MyObj) {
	for i := 0; ; i++ {
		obj := getter()
		fillObject(obj)

		if i == 100000 {

			fmt.Printf("%+v\n", obj)
			resetObj(obj)
			i = 0

			time.Sleep(3 * time.Second)
		}
	}
}
