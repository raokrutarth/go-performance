package main

/**
	Set API that uses string type items and does not provide
	concurrency protection
**/

type UnsafeSet struct {
	data map[string]struct{}
}

func NewUnsafeSet() *UnsafeSet {
	return &UnsafeSet{
		data: make(map[string]struct{}),
	}
}

// Add ...
func (us *UnsafeSet) Add(item string) {
	us.data[item] = struct{}{}
}

// Remove ...
func (us *UnsafeSet) Remove(item string) {
	delete(us.data, item)
}

// IsIn
func (us *UnsafeSet) IsIn(item string) bool {
	_, present := us.data[item]
	return present
}
