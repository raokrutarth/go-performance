package main

import (
	"testing"
)

func BenchmarkClosure(b *testing.B) {

	for i := 0; i < b.N; i++ {
		o := GetClosuredObj()
		fillObject(o)
	}
}

func BenchmarkGlobal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		o := GetGlobalObj()
		fillObject(o)
	}
}
