package main

import (
	"reflect"
	"sync"
)

// Set generic set type
type UncheckedGenericSet struct {
	data map[interface{}]struct{}
	rw   *sync.RWMutex
	t    reflect.Kind
}

// NewSet ...
func NewUncheckedSet(t reflect.Kind) *UncheckedGenericSet {
	ucs := &UncheckedGenericSet{
		data: make(map[interface{}]struct{}),
		rw:   &sync.RWMutex{},
		t:    t,
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

	if _, present := ucs.data[val]; present {
		return true
	}

	return false
}
