package main

// https://www.calhoun.io/5-useful-ways-to-use-closures-in-go/
// write closures vs. global variables benchmark

func getFactory() func() *MyObj {
	o := MyObj{
		b: []byte{42},
		n: 1,
		s: "abc",
	}

	return func() *MyObj {
		return &o
	}
}

// GetClosuredObj gives access to the globally scoped object protected
// using a closure
var GetClosuredObj = getFactory()
