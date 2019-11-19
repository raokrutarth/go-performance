package main

import (
	"fmt"
	"reflect"
	"sync"
)

/**
	Set API that checks the type of the item being added
	so only one type of item can be added to the set
**/

// Set generic set type
type CheckedSet struct {
	data map[interface{}]struct{}
	sync.RWMutex
	t reflect.Kind
}

// NewCheckedSet ...
func NewCheckedSet(t reflect.Kind) *CheckedSet {
	cs := &CheckedSet{
		data: make(map[interface{}]struct{}),
		t:    t,
	}
	return cs
}

// Add ...
func (cs *CheckedSet) Add(val interface{}) {
	if reflect.ValueOf(val).Kind() != cs.t {
		panic(fmt.Errorf("Invalid type %T passed to Set addition. Expected type %s", val, cs.t.String()))
	}

	cs.Lock()
	defer cs.Unlock()
	cs.data[val] = struct{}{}
}

// Remove ...
func (cs *CheckedSet) Remove(val interface{}) {
	if reflect.ValueOf(val).Kind() != cs.t {
		panic(fmt.Errorf("Invalid type %T passed to Set removal. Expected type %s", val, cs.t.String()))
	}

	cs.Lock()
	defer cs.Unlock()

	delete(cs.data, val)
}

// IsIn ...
func (cs *CheckedSet) IsIn(val interface{}) bool {
	if reflect.ValueOf(val).Kind() != cs.t {
		panic(fmt.Errorf("Invalid type %T passed to Set removal. Expected type %s", val, cs.t.String()))
	}

	cs.RLock()
	defer cs.RUnlock()
	_, present := cs.data[val]
	return present
}
