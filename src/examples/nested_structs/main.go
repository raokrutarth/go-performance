package main

import (
	"log"
	"time"

	"golang.performance.com/telemetry"
)

const (
	shouldLog = false
)

/**
	Application to test the performance of trees constructed with
	nodes of different types (e.g. pointer based, generic value based, etc.)
**/

func main() {

	telemetry.Initialize()

	// start multiple goroutines that each test the desired tree type
	for i := 0; i < 5; i++ {
		go func() {
			for {
				// pointerTest()
				nonPointerTest()
			}
		}()
	}
	time.Sleep(30 * time.Minute)

}

// nonPointerTest
func nonPointerTest() {
	tree := makeNestedNonPointerNodes(0, 0)

	b := marshalNonPointerTree(tree)
	if shouldLog {
		log.Printf("Marshaled: %s\n", string(b))
	}

	treeCopy := unMarshalNonPointerTree(b)
	if shouldLog {
		log.Printf("Unmarshaled: %+v\n", treeCopy)
	}
}

func pointerTest() {
	tree := makeNestedNodesWithPointers(0, 0)

	b := marshalTreeWithPointers(tree)
	if shouldLog {
		log.Printf("Marshaled: %s\n", string(b))
	}

	treeCopy := unMarshalTreeWithPointers(b)
	if shouldLog {
		log.Printf("Unmarshaled: %+v\n", treeCopy)
	}
}
