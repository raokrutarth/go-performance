package main

import (
	"encoding/json"
	"fmt"
)

// treeNode can wither have children (and no value) or a value
type treeNode struct {
	Name     string      `json:"TableName,string"`
	Children []*treeNode `json:"Children,omitempty"`
	Value    []byte      `json:"Value,omitempty"` // only for the leaf nodes
}

func makeNestedNodesWithPointers() *treeNode {
	root := &treeNode{Name: "root", Children: []*treeNode{}}
	left := &treeNode{Name: "left_leaf", Value: []byte{1}}
	right := &treeNode{Name: "right_leaf", Value: []byte{2}}

	root.Children = append(root.Children, left, right)

	return root
}

func main() {
	tree := makeNestedNodesWithPointers()

	fmt.Printf("Tree: %+v\n\n", tree)

	b, err := json.Marshal(&tree)
	if err != nil || len(b) == 0 {
		panic(err)
	}

	fmt.Printf("JSON string: %s\n", string(b))

}
