package main

import (
	"fmt"
	"sync"
	"time"
)

/*
	verifying anonymous cache field is not by value
*/

type Cache struct {
	m map[string]uint64
	sync.RWMutex
}

func main() {

	cache := &Cache{m: make(map[string]uint64)}

	for i := 0; i < 20; i++ {
		go func() {
			for {
				cache.Lock()
				cache.m["a"]++
				cache.Unlock()
			}
		}()
	}

	for {
		time.Sleep(3 * time.Second)
		cache.RLock()
		fmt.Println(cache.m["a"])
		cache.RUnlock()
	}
}
