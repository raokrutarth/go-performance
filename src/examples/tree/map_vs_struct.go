package main

import (
	"log"
	"math/rand"
	"runtime"
	"runtime/debug"
	"time"

	json "github.com/json-iterator/go"
	// "encoding/json"
	"golang.performance.com/telemetry"
)

/**
	test to check the memory performance of a tree (of same dimension and data)
	constructed with nested maps vs. custom structs

	jsoniter usage:
	https://github.com/sudo-suhas/bulk-marshal/blob/378738a02807145a41d50e82fd8a31caf87236f2/jsonutil/jsoniter_wrapper.go
**/

const (
	treeHeight         = 4
	numChildrenPerNode = 50

	nodeKeySize   = 50  // size of node ID in bytes
	leafValueSize = 100 // size of leaf node's value field in bytes

	testMarshalRuns = true // runs json.Marshal on the trees when true
)

const (
	mapTreeLeaves    = "map_tree_leaf_nodes"
	structTreeLeaves = "struct_tree_leaf_nodes"
	leafNodeKey      = "leaf"
	letterBytes      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// counter value with this tag increases to correlate events in the test with the memory chart
	eventTag            = "eventCounter"
	testCounterTag      = "num_tests"
	marshalResultLenTag = "marshal_result_bytes_"

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

	for i := 0; ; i++ {

		if i%2 == 0 {
			// alternate between map and struct tests
			mapTreeTest()
		} else {
			structTreeTest()
		}

		telemetry.IncreaseRawValue(testCounterTag, 1)

		telemetry.SetRawValue(mapTreeLeaves, 0)
		telemetry.SetRawValue(structTreeLeaves, 0)
		telemetry.SetRawValue(marshalResultLenTag, 0)

		telemetry.SetRawValue(eventTag, 0)

		debug.FreeOSMemory()
		time.Sleep(10 * time.Second)
	}

}

/*
	map based spanned trees
*/
func mapTreeTest() {
	root := make(map[string]interface{})
	makeSpannedMapTree(root, 0)

	if testMarshalRuns {
		MarshalsAndWait(root, "map_tree")
	}

}

// makeSpannedMapTree constructs the tree of intended dimensions using nested maps
func makeSpannedMapTree(parent map[string]interface{}, depth int) {
	if depth == treeHeight {

		parent[leafNodeKey] = getRandomValue(leafValueSize)
		telemetry.IncreaseRawValue(mapTreeLeaves, 1)

	} else if depth < treeHeight {

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
	Name     string      `json:"TableName,string"`
	Children []*treeNode `json:"Children,omitempty"`
	Value    []byte      `json:"Value,omitempty"` // only for the leaf nodes
}

func structTreeTest() {
	root := &treeNode{Name: "root", Children: []*treeNode{}}
	makeSpannedStructTree(root, 0)

	if testMarshalRuns {
		MarshalsAndWait(root, "struct_tree")
	}
}

func makeSpannedStructTree(parent *treeNode, depth int) {
	if depth == treeHeight {
		leafNode := &treeNode{
			Name:  leafNodeKey,
			Value: getRandomValue(leafValueSize),
		}
		parent.Children = append(parent.Children, leafNode)
		telemetry.IncreaseRawValue(structTreeLeaves, 1)

	} else if depth < treeHeight {

		for i := 0; i < numChildrenPerNode; i++ {

			newNode := &treeNode{
				Name:     getRandomKey(),
				Children: []*treeNode{},
			}
			parent.Children = append(parent.Children, newNode)
			makeSpannedStructTree(newNode, depth+1)

		}

	}
}

// MarshalsAndWait sets event flags in the graph for correlation of
// events to memory consumption during marshaling
func MarshalsAndWait(root interface{}, testName string) {

	postTreeConstructionActions(testName)
	log.Printf("[+] pre-marshal wait complete. Marshaling\n")

	marshaledBytes := doMarshalRuns(root, 5)

	log.Printf("[+] Marshal runs complete. Result length: %d\n", len(marshaledBytes))
	telemetry.SetRawValue(eventTag, postMarshalEvent)                                 // set event identifier in graph
	telemetry.SetRawValue(marshalResultLenTag+testName, float64(len(marshaledBytes))) // send result length to graph

	// stay alive with only JSON result in memory so memory stats can be scraped
	log.Printf("[+] Only marshal result is in memory. Waiting...\n")
	time.Sleep(2 * time.Minute)

	if len(marshaledBytes) == 0 {
		// dummy reference to keep marshal result in memory
		log.Println("[+] Got empty marshal result for test " + testName)
	}
}

func postTreeConstructionActions(testName string) {
	telemetry.SetRawValue(eventTag, treeBuildCompleteEvent) // set event identifier in graph
	debug.FreeOSMemory()                                    // flush the GC so only the tree is occupying memory
	telemetry.SetRawValue(eventTag, postTreeBuildGCFinish)  // set event identifier in graph

	log.Printf("[+] spanned %s tree constructed. Only tree object is occupying memory\n", testName)
	time.Sleep(45 * time.Second)

	telemetry.SetRawValue(eventTag, preMarshalEvent) // set event identifier in graph
}

// doMarshalRuns marshals the passed tree num times and returns
// the byte array of the final marshal result
func doMarshalRuns(root interface{}, num int) []byte {
	var b []byte
	var err error

	// marshal a few times to get highest memory peak
	for i := 0; i < num; i++ {
		b, err = json.Marshal(root)

		log.Printf("[+] Finished marshal iteration %d. Flushing GC to let only marshal result persist\n", i)
		runtime.GC() // flush the GC to remove the tree/old marshal result from memory

		if err != nil {
			panic(err)
		}

		telemetry.IncreaseRawValue(eventTag, 5) // add to event identifier in graph
	}

	return b
}

// getRandomKey returns a random fixed size string
func getRandomKey() string {
	b := make([]byte, nodeKeySize)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func getRandomValue(size int) []byte {
	b := make([]byte, size)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return b
}
