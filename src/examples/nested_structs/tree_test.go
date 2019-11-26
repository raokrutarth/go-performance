package main

import (
	"testing"
)

/**
	Benchmarks to compare the performance of trees constructed with
	nodes of different typee (e.g. pointer based, generic value based, etc.)
**/
func BenchmarkPointerTreeMarshal(b *testing.B) {

	for i := 0; i < b.N; i++ {
		tree := makeNestedNodesWithPointers(0, 0)
		marshalTreeWithPointers(tree)
	}

}

func BenchmarkNonPointerTreeMarshal(b *testing.B) {

	for i := 0; i < b.N; i++ {
		tree := makeNestedNonPointerNodes(0, 0)
		marshalNonPointerTree(tree)
	}

}

func TestTreesAreEqual(t *testing.T) {
	tree1 := makeNestedNodesWithPointers(0, 0)
	tree2 := makeNestedNonPointerNodes(0, 0)

	b1 := marshalTreeWithPointers(tree1)
	b2 := marshalNonPointerTree(tree2)

	if string(b1) != string(b2) {
		t.FailNow()
	}

}
