package main

const (
	workloadStr  = "abcdefg"
	workloadNum  = 100
	workloadByte = 0x42
)

type MyObj struct {
	b []byte
	n int
	s string
}

// fillObject mocks the object getting new
// items added to it. similar to when a payload is unmarshalled
// into an object
func fillObject(obj *MyObj) {

	obj.b = append(obj.b, workloadByte)
	obj.n += workloadNum
	obj.s += workloadStr

}

// resetObj ...
func resetObj(obj *MyObj) {
	obj.b = []byte{workloadByte}
	obj.n = 1
	obj.s = "abc"
}
