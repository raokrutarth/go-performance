package main

import (
	"math/rand"
	"time"

	"golang.performance.com/telemetry"
)

const (
	readCounterTag     = "num_read"
	increaseCounterTag = "num_increase"
	setCounterTag      = "num_set"

	readWaitTimeTag     = "read_wait_ns"
	increaseWaitTimeTag = "increase_wait_time_ns"
	setWaitTimeTag      = "set_wait_time_ns"
)

/*
	Test to see if mutexes have better CPU usage and I/O waiting
	than channel protection.
*/
type Cache interface {
	get(key string) int
	set(key string, value int)
}

func main() {
	telemetry.Initialize()
	// cache := newMutexCache()
	cache := newChanCache()

	// setup reader, setter and increasors on same keys
	for i := 0; i < 300; i++ {
		key := getRandomKey()

		go setter(i, cache, key, rand.Intn(100))
		go reader(i, cache, key)
		go increasor(i, cache, key)
	}

	for {
		time.Sleep(60 * time.Minute)
	}
}

// reader reads a value for the given key from the cache
func reader(id int, cache Cache, key string) {
	var start time.Time
	for {
		start = time.Now()

		cache.get(key)

		telemetry.ExportVariableValue(readWaitTimeTag, float64(time.Since(start).Nanoseconds()))
		telemetry.IncreaseRawValue(readCounterTag, 1)

	}
}

// increasor increases a value for the given key in the cache
func increasor(id int, cache Cache, key string) {
	var start time.Time
	for {
		start = time.Now()

		cache.set(key, cache.get(key)+1)

		telemetry.ExportVariableValue(increaseWaitTimeTag, float64(time.Since(start).Nanoseconds()))
		telemetry.IncreaseRawValue(increaseCounterTag, 1)
	}

}

// func remover(id int, sc *SafeCache, key string) {
// 	for {
// 		fmt.Printf("remover %d removing key %s\n", id, key)
// 		sc.remove(key)
// 		fmt.Printf("remover %d removed key: %s\n", id, key)
// 		// to slow down terminal output
// 		time.Sleep(time.Duration(duration*3) * time.Second)
// 	}

// }

func setter(id int, cache Cache, key string, value int) {
	var start time.Time
	for {
		start = time.Now()

		cache.set(key, value)

		telemetry.ExportVariableValue(setWaitTimeTag, float64(time.Since(start).Nanoseconds()))
		telemetry.IncreaseRawValue(setCounterTag, 1)
	}
}

// getRandomKey returns a random fixed size string
func getRandomKey() string {
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	bucket := make([]byte, 1000000)
	capacity := int64(len(letterBytes))

	for i := range bucket {
		bucket[i] = letterBytes[rand.Int63()%capacity]
	}

	return string(bucket)
}
