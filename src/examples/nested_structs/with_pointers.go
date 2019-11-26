package main

import "encoding/json"

// treeNode can wither have children (and no value) or a value
type PointerNode struct {
	Name     string         `json:"TableName,string"`
	Children []*PointerNode `json:"Children,omitempty"`
	Value    []byte         `json:"Value,omitempty"` // only for the leaf nodes
}

func makeNestedNodesWithPointers(treeHeight, numChildren int) *PointerNode {
	root := &PointerNode{Name: "root", Children: []*PointerNode{}}
	left := &PointerNode{Name: "left_leaf", Value: []byte{1}}
	right := &PointerNode{Name: "right_leaf", Value: []byte{2}}

	root.Children = append(root.Children, left, right)

	return root
}

func marshalTreeWithPointers(tree *PointerNode) []byte {

	b, err := json.Marshal(tree)
	if err != nil {
		panic(err)
	}
	return b

}

// unMarshalTreeWithPointers needs to go through each child pointer
// and create a new child object instead of pointing to the old one
func unMarshalTreeWithPointers(b []byte) *PointerNode {
	oldTree := &PointerNode{}
	err := json.Unmarshal(b, &oldTree)
	if err != nil {
		panic(err)
	}

	newTree := PointerNode{
		Name: oldTree.Name,
	}

	newTree.Children = append(newTree.Children, oldTree.Children...)

	return &newTree
}
