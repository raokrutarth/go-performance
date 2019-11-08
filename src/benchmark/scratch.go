package main

import "fmt"
import "time"
import "testing"

func BenchmarkTimeFormat(b *testing.B) {
	t := time.Unix(1265346057, 0)
	for i := 0; i < b.N; i++ {
		t.Format("Mon Jan  2 15:04:05 2006")
	}
}

func printTime(t time.Time) {
	fmt.Printf("%s\n", t)
}

func main() {
	// benchmark := &testing.B{N: 50000}
	// BenchmarkTimeFormat(benchmark)
	// fmt.Printf("%+v", benchmark)

	t := time.Now()
	printTime(t)

}
