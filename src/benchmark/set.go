package main

import (
	"fmt"
	"reflect"
	"sync"
)

// Set generic set type
type Set struct {
	data map[interface{}]struct{}
	rw   *sync.RWMutex
	t    reflect.Kind
}

// NewSet ...
func NewSet(t reflect.Kind) *Set {
	s := &Set{
		data: make(map[interface{}]struct{}),
		rw:   &sync.RWMutex{},
		t:    t,
	}
	return s
}

// Add ...
func (s *Set) Add(val interface{}) error {
	if reflect.ValueOf(val).Kind() != s.t {
		return fmt.Errorf("Invalid type %T passed to Set addition. Expected type %s", val, s.t.String())
	}

	s.rw.Lock()
	defer s.rw.Unlock()
	s.data[val] = struct{}{}
	return nil
}

// Remove ...
func (s *Set) Remove(val interface{}) error {
	if reflect.ValueOf(val).Kind() != s.t {
		return fmt.Errorf("Invalid type %T passed to Set removal. Expected type %s", val, s.t.String())
	}

	s.rw.Lock()
	defer s.rw.Unlock()
	delete(s.data, val)
	return nil
}

// IsIn ...
func (s *Set) IsIn(val interface{}) (bool, error) {
	if reflect.ValueOf(val).Kind() != s.t {
		return false, fmt.Errorf("Invalid type %T passed to Set removal. Expected type %s", val, s.t.String())
	}

	s.rw.RLock()
	defer s.rw.RUnlock()
	if _, present := s.data[val]; present {
		return true, nil
	}
	return false, nil
}

func main() {
	fmt.Println("Hello world")
	set := NewSet(reflect.String)
	k1, k2, k3 := "aaa", "bbb", "a"
	var a, b, c bool
	var err error

	err = set.Add(k1)
	err = set.Add(k2)
	if err != nil {
		fmt.Printf("%s\n\n", err)
		return
	}

	a, err = set.IsIn(k1)
	b, err = set.IsIn(k2)
	c, err = set.IsIn(k3)
	if err != nil {
		fmt.Printf("%s\n\n", err)
	}
	fmt.Printf("%s in set: %v, %s in set: %v, %s in set: %v\n", k1, a, k2, b, k3, c)

	err = set.Remove(k1)
	if err != nil {
		fmt.Printf("%s\n\n", err)
		return
	}

	a, err = set.IsIn(k1)
	b, err = set.IsIn(k2)
	c, err = set.IsIn(k3)
	if err != nil {
		fmt.Printf("%s\n\n", err)
	}

	fmt.Printf("%s in set: %v, %s in set: %v, %s in set: %v\n", k1, a, k2, b, k3, c)

	err = set.Add(42)
	if err != nil {
		fmt.Printf("%s\n\n", err)
	}
}
