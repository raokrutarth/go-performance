package main

type ChannelCache struct {
	c chan map[string]int
}

func newChanCache() *ChannelCache {
	cache := &ChannelCache{
		c: make(chan map[string]int),
	}
	go func() {
		cache.c <- make(map[string]int)
	}()

	return cache
}

func (cache *ChannelCache) get(key string) int {
	m := <-cache.c

	if value, ok := m[key]; ok {
		cache.c <- m
		return value
	}
	cache.c <- m

	return -1
}

func (cache *ChannelCache) set(key string, value int) {
	m := <-cache.c
	m[key] = value
	cache.c <- m
}
