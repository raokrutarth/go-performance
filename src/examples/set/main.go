package main

import (
	"golang.performance.com/telemetry"
)

type GenericSet interface {
	Add(item interface{})
	IsIn(item interface{}) bool
	Remove(item interface{})
}

type Set interface {
	Add(string)
	IsIn(string) bool
	Remove(string)
}

func main() {

	telemetry.Initialize()

	for {

	}

}
