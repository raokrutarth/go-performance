package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"runtime"
	"time"

	"golang.performance.com/telemetry"
	// jsoniter "github.com/json-iterator/go"
)

/**
	test to check the memory performance of a tree (of same size)
	constructed with nested maps and custom structs

**/

// var json = jsoniter.ConfigCompatibleWithStandardLibrary

const (
	maxTreeDepth       = 3
	numChildrenPerNode = 100
)

const (
	mapTreeLeaves    = "map_tree_leaf_nodes"
	structTreeLeaves = "struct_tree_leaf_nodes"
	leafNodeKey      = "leaf"
	letterBytes      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// counter value with this tag increases to correlate events in the test with the memory chart
	eventTag = "eventCounter"

	testCounterTag = "num_tests"

	marshalResultLenTag = "marshal_result_bytes"

	treeBuildCompleteEvent = 50
	postTreeBuildGCFinish  = 100
	preMarshalEvent        = 200
	postMarshalEvent       = 300
)

/**
	main function that infinitely runs the chosen test
**/
func main() {
	telemetry.Initialize()

	telemetry.SetRawValue(testCounterTag, 0)

	for {
		mapTreeTest()

		telemetry.IncreaseRawValue(testCounterTag, 1)
	}

}

/*
	map based spanned trees
*/
func mapTreeTest() {
	root := make(map[string]interface{})
	makeSpannedMapTree(root, 0)

	waitAndMarshal(root, "map")
}

// makeSpannedMapTree constructs the tree of intended dimensions using nested maps
func makeSpannedMapTree(parent map[string]interface{}, depth int) {
	if depth == maxTreeDepth {

		parent[leafNodeKey] = make([]byte, 100)
		telemetry.IncreaseRawValue(mapTreeLeaves, 1)

	} else if depth < maxTreeDepth {

		for i := 0; i < numChildrenPerNode; i++ {

			newNode := make(map[string]interface{})
			parent[getRandomKey()] = newNode
			makeSpannedMapTree(newNode, depth+1)

		}
	}
}

/*
	struct based spanned tree with static types in nodes
*/

// treeNode can wither have children (and no value) or a value
type treeNode struct {
	name     string      `json:"TableName,string"`
	children []*treeNode `json:"Children,omitempty"`
	value    []byte      `json:"Value,omitempty"`
}

func structTreeTest() {
	root := &treeNode{name: "root", children: []*treeNode{}}
	makeSpannedStructTree(root, 0)

	waitAndMarshal(root, "non_interface struct")
}

func makeSpannedStructTree(parent *treeNode, depth int) {
	if depth == maxTreeDepth {
		leafNode := &treeNode{
			name:  leafNodeKey,
			value: make([]byte, 100),
		}
		parent.children = append(parent.children, leafNode)
		telemetry.IncreaseRawValue(structTreeLeaves, 1)

	} else if depth < maxTreeDepth {

		for i := 0; i < numChildrenPerNode; i++ {

			newNode := &treeNode{
				name:     getRandomKey(),
				children: []*treeNode{},
			}
			parent.children = append(parent.children, newNode)
			makeSpannedStructTree(newNode, depth+1)

		}

	}
}

// marshalSpannedStructTree is the custom JSON marshaler to traverse the tree
// can take 30m+ and 10GB+ with a 3.2mil leaf node tree
func marshalSpannedStructTree(root *treeNode) ([]byte, error) {
	res := []byte{}

	dequeue := []*treeNode{root}

	for len(dequeue) > 0 {
		// pop
		currNode, dequeue := dequeue[0], dequeue[1:]

		if currNode.name == leafNodeKey {
			// seeing a child of type byte[]
			res = append(res, currNode.value...)
			continue
		}
		b, err := json.Marshal(*currNode)
		if err != nil {
			return res, err
		}
		res = append(res, b...)

		for _, child := range currNode.children {
			// enqueue
			dequeue = append(dequeue, child)
		}
	}
	return res, nil
}

// waitAndMarshal sets event flags in the graph for correlation of
// events to memory consumption during marshaling
func waitAndMarshal(root interface{}, testName string) {

	telemetry.SetRawValue(eventTag, treeBuildCompleteEvent) // set event identifier in graph
	runtime.GC()                                            // flush the GC so only the tree is occupying memory
	telemetry.SetRawValue(eventTag, postTreeBuildGCFinish)  // set event identifier in graph

	log.Printf("[+] spanned %s tree constructed. Only tree object is occupying memory\n", testName)
	time.Sleep(45 * time.Second)

	log.Printf("[+] pre-marshal wait complete. Marshaling\n")
	telemetry.SetRawValue(eventTag, preMarshalEvent) // set event identifier in graph

	marshaledBytes := doMarshalRuns(root, 5)

	log.Printf("[+] Marshal runs complete. Result length: %d\n", len(marshaledBytes))
	telemetry.SetRawValue(eventTag, postMarshalEvent)                        // set event identifier in graph
	telemetry.SetRawValue(marshalResultLenTag, float64(len(marshaledBytes))) // send result length to graph

	// stay alive with only JSON result in memory so memory stats can be scraped
	log.Printf("[+] Only marshal result is in memory. Waiting...\n")
	time.Sleep(30 * time.Minute)
	log.Printf("len: %d\nValue: \n%+v", len(marshaledBytes), marshaledBytes)
}

// doMarshalRuns marshals the passed tree num times and returns
// the byte array of the final marshal result
func doMarshalRuns(root interface{}, num int) []byte {
	var b []byte
	var err error

	// marshal a few times to get highest memory peak
	for i := 0; i < num; i++ {

		// marshal based on tree type
		if st, ok := root.(*treeNode); ok {
			b, err = marshalSpannedStructTree(st) // BAD
		} else {
			b, err = json.Marshal(root)
		}

		if err != nil {
			log.Fatalf("unable to Marshal. Exiting...")
		}

		telemetry.IncreaseRawValue(eventTag, 10) // add to event identifier in graph

		log.Printf("[+] Finished marshal iteration %d. Flushing GC to let only marshal result persist\n", i)
		runtime.GC() // flush the GC to remove the tree/old marshal result from memory
	}

	return b
}

// getRandomKey returns a random fixed size string
func getRandomKey() string {
	b := make([]byte, 32)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}
