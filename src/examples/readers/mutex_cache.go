package main

import "sync"

type MutexCache struct {
	m map[string]int
	sync.RWMutex
}

func newMutexCache() *MutexCache {
	return &MutexCache{
		m: make(map[string]int),
	}
}

func (pc *MutexCache) get(key string) int {
	pc.RLock()
	defer pc.RUnlock()

	if value, ok := pc.m[key]; ok {
		return value
	}

	return -1
}

func (pc *MutexCache) set(key string, value int) {
	pc.Lock()
	defer pc.Unlock()

	pc.m[key] = value
}
