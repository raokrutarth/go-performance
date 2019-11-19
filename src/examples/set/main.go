package main

import (
	"math/rand"
	"reflect"
	"sync"
	"time"

	"golang.performance.com/telemetry"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func main() {
	numItems := 10000
	itemSize := 500000

	telemetry.Initialize()

	set := NewCheckedSet(reflect.String)
	// set := NewUncheckedSet()
	// set := NewTypedSet()

	items := []string{}

	for i := 0; i < numItems; i++ {
		items = append(items, GenerateItem(itemSize))
	}

	for {
		var waitgroup sync.WaitGroup

		go func() {
			waitgroup.Add(1)
			for _, item := range items {
				set.Add(item)
			}
			waitgroup.Done()
		}()

		go func() {
			waitgroup.Add(1)
			for _, item := range items {
				set.IsIn(item)
			}
			waitgroup.Done()
		}()

		go func() {
			waitgroup.Add(1)
			for _, item := range items {
				set.Remove(item)
			}
			waitgroup.Done()
		}()

		// wait until all the goroutines are completed
		waitgroup.Wait()

		telemetry.IncreaseRawValue("num_cycles", 1)

		// wait for the add/remove cycle to finish
		time.Sleep(10 * time.Second)

	}

}

func GenerateItem(size int) string {
	b := make([]byte, size)
	numChars := len(charset)

	for i := range b {
		b[i] = charset[seededRand.Intn(numChars)]
	}

	return string(b)
}
