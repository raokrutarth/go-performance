package main

import (
	"sync"
)

type TypedSet struct {
	data map[string]struct{}
	sync.RWMutex
}

func NewTypedSet() *TypedSet {
	return &TypedSet{
		data: make(map[string]struct{}),
	}
}

// Add ...
func (ts *TypedSet) Add(item string) {
	ts.Lock()
	defer ts.Unlock()

	ts.data[item] = struct{}{}
}

// Remove ...
func (ts *TypedSet) Remove(item string) {
	ts.Lock()
	defer ts.Unlock()

	delete(ts.data, item)
}

// IsIn
func (ts *TypedSet) IsIn(item string) bool {
	ts.RLock()
	defer ts.RUnlock()

	if _, present := ts.data[item]; present {
		return true
	}
	return false
}
