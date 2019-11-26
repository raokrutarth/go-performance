package main

import "encoding/json"

// treeNode can wither have children (and no value) or a value
type NonPointerNode struct {
	Name     string           `json:"TableName,string"`
	Children []NonPointerNode `json:"Children,omitempty"`
	Value    []byte           `json:"Value,omitempty"` // only for the leaf nodes
}

func makeNestedNonPointerNodes(treeHeight, numChildren int) NonPointerNode {
	root := NonPointerNode{Name: "root", Children: []NonPointerNode{}}
	left := NonPointerNode{Name: "left_leaf", Value: []byte{1}}
	right := NonPointerNode{Name: "right_leaf", Value: []byte{2}}

	root.Children = append(root.Children, left, right)

	return root
}

func marshalNonPointerTree(root NonPointerNode) []byte {

	b, err := json.Marshal(root)
	if err != nil {
		panic(err)
	}
	return b
}

func unMarshalNonPointerTree(b []byte) NonPointerNode {

	tree := NonPointerNode{}
	err := json.Unmarshal(b, &tree)
	if err != nil {
		panic(err)
	}

	return tree
}
