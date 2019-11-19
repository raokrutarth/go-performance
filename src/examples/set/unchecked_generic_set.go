package main

import (
	"sync"
)

/**
	Set API that does not check the type of item bing added
	to the set. Allowing multiple types of items in the same
	object.
**/

// Set generic set type
type UncheckedGenericSet struct {
	data map[interface{}]struct{}
	rw   *sync.RWMutex
}

// NewSet ...
func NewUncheckedSet() *UncheckedGenericSet {
	ucs := &UncheckedGenericSet{
		data: make(map[interface{}]struct{}),
		rw:   &sync.RWMutex{},
	}
	return ucs
}

// Add ...
func (ucs *UncheckedGenericSet) Add(val interface{}) {

	ucs.rw.Lock()
	defer ucs.rw.Unlock()
	ucs.data[val] = struct{}{}
}

// Remove ...
func (ucs *UncheckedGenericSet) Remove(val interface{}) {

	ucs.rw.Lock()
	defer ucs.rw.Unlock()
	delete(ucs.data, val)
}

// IsIn ...
func (ucs *UncheckedGenericSet) IsIn(val interface{}) bool {

	ucs.rw.RLock()
	defer ucs.rw.RUnlock()

	_, present := ucs.data[val]
	return present
}
