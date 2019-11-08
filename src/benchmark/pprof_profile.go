package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"
)

var m = &sync.RWMutex{}

func main() {
	runtime.SetBlockProfileRate(1)
	runtime.SetMutexProfileFraction(1)

	profile := "profile.p"
	f, err := os.Create(profile)
	if err != nil {
		log.Printf("Error creating file for CPU Profile %v", err)
	}
	defer f.Close()
	// pp := pprof.Lookup("block")
	pp := pprof.Lookup("mutex")

	ch := make(chan bool, 3)
	ch <- true
	ch <- true
	ch <- true
	for i := 0; i < 50; i++ {
		go func() {
			m.Lock()
			for k := 0; k < 50000; k++ {
				s := fmt.Sprintf("%d", k)
				if s != "" {
					<-ch
					ch <- true
					continue
				}
			}
			m.Unlock()
		}()
	}

	fmt.Println("Loop ended")

	if pp != nil {
		time.Sleep(5 * time.Second)
		pp.WriteTo(f, 2)
		fmt.Printf("count: %d\n", pp.Count())
	}

}
