package main_test

import (
	"fmt"
	"reflect"
	"testing"
)

/**
	Test file that implements Benchmarks that do sequential and parallel access to
	different types of sets
**/

const numItems = 500

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

// func BenchmarkCheckedGenericSetParallel(b *testing.B) {
// 	set := NewCheckedSet(reflect.String)
// 	benchmarkGenericSetParallel(set, b)
// }

// func BenchmarkUnCheckedGenericSetParallel(b *testing.B) {
// 	set := NewUncheckedSet(reflect.String)
// 	benchmarkGenericSetParallel(set, b)
// }

func benchmarkSetParallel(set Set, b *testing.B) {
	items := []string{}

	for i := 0; i < numItems; i++ {
		items = append(items, GenerateItem(numItems))
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

	for i := 0; i < numItems; i++ {
		items = append(items, GenerateItem(numItems))
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

	for i := 0; i < numItems; i++ {
		items = append(items, GenerateItem(numItems))
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

	for i := 0; i < numItems; i++ {
		items = append(items, GenerateItem(numItems))
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
