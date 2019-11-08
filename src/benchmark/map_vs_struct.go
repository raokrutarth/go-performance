package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	// jsoniter "github.com/json-iterator/go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// var json = jsoniter.ConfigCompatibleWithStandardLibrary

const (
	maxTreeDepth       = 5
	numChildrenPerNode = 20
)

func main() {
	go exportPerformance()

	// mapTreeTest()
	structTreeTest()
}

/*
	map based spanned trees
*/
func mapTreeTest() {
	root := make(map[string]interface{})
	makeSpannedMapTree(root, 0)

	waitAndMarshal(root, "map")
}

func makeSpannedMapTree(parent map[string]interface{}, depth int) {
	if depth == maxTreeDepth {

		parent[leafNodeKey] = make([]byte, 100)
		leafCounter.WithLabelValues(mapTreeLeaves).Inc()

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
		leafCounter.WithLabelValues(structTreeLeaves).Inc()

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

/*
	Helper functions and variables to send data to the benchmark tool
*/

const (
	restEndpoint        = "/metrics"
	merticsServePort    = ":35005"
	mapTreeLeaves       = "map_tree_leaf_nodes"
	structTreeLeaves    = "struct_tree_leaf_nodes"
	leafNodeKey         = "leaf"
	letterBytes         = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	eventTag            = "eventCounter" // the counter value with this tag changes to correlate events in the test with the memory chart
	marshalResultLenTag = "marshal_result_bytes"

	treeBuildCompleteEvent = 50
	postTreeBuildGCFinish  = 100
	preMarshalEvent        = 200
	postMarshalEvent       = 300
)

var (
	leafCounter = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "hpe_config_arbitrary_value",
			Help: "Value of arbitrary in-app numbers",
		},
		[]string{"value_name"},
	)
)

// waitAndMarshal sets event flags in the graph for correlation of
// events to memory consumption during marshaling
func waitAndMarshal(root interface{}, testName string) {

	leafCounter.WithLabelValues(eventTag).Set(treeBuildCompleteEvent) // set event identifier in graph
	runtime.GC()                                                      // flush the GC so only the tree is occupying memory
	leafCounter.WithLabelValues(eventTag).Set(postTreeBuildGCFinish)  // set event identifier in graph

	log.Printf("[+] spanned %s tree constructed. Only tree object is occupying memory\n", testName)
	time.Sleep(45 * time.Second)

	log.Printf("[+] pre-marshal wait complete. Marshaling\n")
	leafCounter.WithLabelValues(eventTag).Set(preMarshalEvent) // set event identifier in graph

	b := doMarshalRuns(root, 5)

	log.Printf("[+] Marshal runs complete. Result length: %d\n", len(b))
	leafCounter.WithLabelValues(eventTag).Set(postMarshalEvent)           // set event identifier in graph
	leafCounter.WithLabelValues(marshalResultLenTag).Set(float64(len(b))) // send result length to graph

	// stay alive with only JSON result in memory so memory stats can be scraped
	log.Printf("[+] Only marshal result is in memory. Waiting...\n")
	time.Sleep(30 * time.Minute)
	log.Printf("len: %d\nValue: \n%+v", len(b), b)
}

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

		leafCounter.WithLabelValues(eventTag).Add(10) // add to event identifier in graph

		log.Printf("[+] Finished marshal iteration %d. Flushing GC to let only marshal result persist\n", i)
		runtime.GC() // flush the GC to remove the tree/old marshal result from memory
	}

	return b
}

func getRandomKey() string {
	b := make([]byte, 32)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func exportPerformance() {
	prometheus.MustRegister(leafCounter)

	http.Handle(restEndpoint, promhttp.Handler())
	http.ListenAndServe(merticsServePort, nil)
}
