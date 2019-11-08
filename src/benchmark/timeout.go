package main

import (
	"fmt"
	"os"
	"time"
)

/*
	Test whether a select is multi-way blocking in go when
	used without a default case
*/

func main() {
	c := make(chan error)

	go func() {
		t := 10
		time.Sleep(time.Duration(t) * time.Second)
		c <- fmt.Errorf("after %d sec", t)
	}()

	select {
	case e := <-c:
		fmt.Printf("got error %s\n", e)

	case <-time.After(5 * time.Second):
		fmt.Println("Timed out at 5s. Exiting...")
		os.Exit(1)
	}
}
