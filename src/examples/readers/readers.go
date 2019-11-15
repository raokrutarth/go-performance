package main

import (
	"math/rand"
	"sync"
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

type ProtectedCache struct {
	m map[string]int
	sync.RWMutex
}

func (pc *ProtectedCache) get(key string) int {
	pc.RLock()
	defer pc.RUnlock()

	if value, ok := pc.m[key]; ok {
		return value
	}

	return -1
}

func (pc *ProtectedCache) set(key string, value int) {
	pc.Lock()
	defer pc.Unlock()

	pc.m[key] = value
}

func reader(id int, pc *ProtectedCache, key string) {
	var start time.Time
	for {
		start = time.Now()

		pc.get(key)

		telemetry.ExportVariableValue(readWaitTimeTag, float64(time.Since(start).Nanoseconds()))
		telemetry.IncreaseRawValue(readCounterTag, 1)

	}
}

func increasor(id int, pc *ProtectedCache, key string) {
	var start time.Time
	for {
		start = time.Now()

		val := pc.get(key)
		val++
		pc.set(key, val)

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

func setter(id int, pc *ProtectedCache, key string, value int) {
	var start time.Time
	for {
		start = time.Now()

		pc.set(key, value)

		telemetry.ExportVariableValue(setWaitTimeTag, float64(time.Since(start).Nanoseconds()))
		telemetry.IncreaseRawValue(setCounterTag, 1)
	}
}

// getRandomKey returns a random fixed size string
func getRandomKey() string {
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	bucket := make([]byte, 32000)
	capacity := int64(len(letterBytes))

	for i := range bucket {
		bucket[i] = letterBytes[rand.Int63()%capacity]
	}

	return string(bucket)
}

func main() {
	pc := &ProtectedCache{
		m: make(map[string]int),
	}

	// setup reader, setter and increasors on same keys
	for i := 0; i < 50; i++ {
		key := getRandomKey()

		go reader(i, pc, key)
		go increasor(i, pc, key)
		go setter(i, pc, key, rand.Intn(100))
	}

	for {
		time.Sleep(60 * time.Minute)
	}
}
