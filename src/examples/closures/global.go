package main

var globalObj = MyObj{
	b: []byte{42},
	n: 1,
	s: "abc",
}

// GetGlobalObj gives access to the global object
func GetGlobalObj() *MyObj {
	return &globalObj
}
