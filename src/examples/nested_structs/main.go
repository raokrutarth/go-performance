package main

import (
	"log"
)

func main() {
	pointerTest()

	nonPointerTest()

}

func nonPointerTest() {
	tree := makeNestedNonPointerNodes()

	b := marshalNonPointerTree(tree)
	log.Printf("Marshaled: %s\n", string(b))

	treeCopy := unMarshalNonPointerTree(b)
	log.Printf("Unmarshaled: %+v\n", treeCopy)
}

func pointerTest() {
	tree := makeNestedNodesWithPointers()

	b := marshalTreeWithPointers(tree)
	log.Printf("Marshaled: %s\n", string(b))

	treeCopy := unMarshalTreeWithPointers(b)
	log.Printf("Unmarshaled: %+v\n", treeCopy)
}
