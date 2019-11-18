package main

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

const numItems = 500

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func BenchmarkTypedSet(b *testing.B) {
	set := NewTypedSet()
	benchmarkSet(set, b)
}

func BenchmarkCheckedGenericSet(b *testing.B) {
	set := NewCheckedSet(reflect.String)
	benchmarkGenericSet(set, b)
}

func BenchmarkUnCheckedGenericSet(b *testing.B) {
	set := NewUncheckedSet(reflect.String)
	benchmarkGenericSet(set, b)
}

func BenchmarkTypedSetParallel(b *testing.B) {
	set := NewTypedSet()
	benchmarkSetParallel(set, b)
}

func BenchmarkCheckedGenericSetParallel(b *testing.B) {
	set := NewCheckedSet(reflect.String)
	benchmarkGenericSetParallel(set, b)
}

func BenchmarkUnCheckedGenericSetParallel(b *testing.B) {
	set := NewUncheckedSet(reflect.String)
	benchmarkGenericSetParallel(set, b)
}

func benchmarkSetParallel(set Set, b *testing.B) {
	items := []string{}

	for i := 0; i < 500; i++ {
		items = append(items, generateItem(numItems))
	}

	for i := 0; i < b.N; i++ {
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
	}
}

func benchmarkSet(set Set, b *testing.B) {

	items := []string{}

	for i := 0; i < 500; i++ {
		items = append(items, generateItem(numItems))
	}

	for i := 0; i < b.N; i++ {
		for _, item := range items {
			set.Add(item)
		}

		for _, item := range items {
			if !set.IsIn(item) {
				panic(fmt.Errorf("expected key not present"))
			}
		}

		for _, item := range items {
			set.Remove(item)
		}

		for _, item := range items {
			if set.IsIn(item) {
				panic(fmt.Errorf("unexpected item present"))
			}
		}
	}
}

func benchmarkGenericSet(set GenericSet, b *testing.B) {

	items := []string{}

	for i := 0; i < 500; i++ {
		items = append(items, generateItem(numItems))
	}

	for i := 0; i < b.N; i++ {
		for _, item := range items {
			set.Add(item)
		}

		for _, item := range items {
			if !set.IsIn(item) {
				panic(fmt.Errorf("expected key not present"))
			}
		}

		for _, item := range items {
			set.Remove(item)
		}

		for _, item := range items {
			if set.IsIn(item) {
				panic(fmt.Errorf("unexpected item present"))
			}
		}
	}

}

func benchmarkGenericSetParallel(set GenericSet, b *testing.B) {
	items := []string{}

	for i := 0; i < 500; i++ {
		items = append(items, generateItem(numItems))
	}

	for i := 0; i < b.N; i++ {
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
	}
}

func generateItem(size int) string {
	b := make([]byte, size)
	numChars := len(charset)

	for i := range b {
		b[i] = charset[seededRand.Intn(numChars)]
	}

	return string(b)
}
