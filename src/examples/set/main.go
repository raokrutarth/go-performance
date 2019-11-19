package main

import (
	"math/rand"
	"time"

	"golang.performance.com/telemetry"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func main() {
	numItems := 100
	itemSize := 10

	telemetry.Initialize()

	set := NewUncheckedSet()

	items := []string{}

	for i := 0; i < numItems; i++ {
		items = append(items, GenerateItem(itemSize))
	}

	for {
		go func() {
			for _, item := range items {
				set.Add(item)
			}
		}()

		go func() {
			for _, item := range items {
				set.IsIn(item)
			}
		}()

		go func() {
			for _, item := range items {
				set.Remove(item)
			}
		}()

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
