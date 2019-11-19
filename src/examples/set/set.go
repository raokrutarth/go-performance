package main


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