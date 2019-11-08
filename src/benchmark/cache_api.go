package main

import (
	"fmt"
	"sync"
	"time"
)

/*
	Simple thread safe cache api that allowes locked reads
*/

const (
	// adjust speed of terminal output by changing duration
	duration = 2
)

type myFavObjType struct {
	field1 string
	field2 *int64
}

type SafeCache struct {
	data map[string]myFavObjType
	rw   *sync.RWMutex
}

func (sc *SafeCache) get(key string) (myFavObjType, error) {
	sc.rw.RLock()
	defer sc.rw.RUnlock()
	if res, present := sc.data[key]; present {
		return res, nil
	}
	return myFavObjType{"", new(int64)},
		fmt.Errorf("some error that says a default value had to be generated")
}

func (sc *SafeCache) set(key string, newItem myFavObjType) {
	sc.rw.Lock()
	defer sc.rw.Unlock()
	sc.data[key] = newItem
}

func (sc *SafeCache) remove(key string) error {
	sc.rw.Lock()
	defer sc.rw.Unlock()
	if _, present := sc.data[key]; present {
		delete(sc.data, key)
		return nil
	}
	return fmt.Errorf("Object not present. No delete needed")
}

func (sc *SafeCache) getAndLock(key string) (myFavObjType, chan struct{}, error) {
	/*
		returning a done channel to be consumed by the Unlocker reminds the
		programmer to call unlock and enforces the unlock
	*/
	// log here for deadlock detection
	sc.rw.Lock()
	fmt.Println("got full lock in getAndLock")
	done := make(chan struct{}, 1)
	defer func() { done <- struct{}{} }()
	if res, present := sc.data[key]; present {
		return res, done, nil
	}
	return myFavObjType{"", new(int64)}, done,
		fmt.Errorf("some error that says a default value had to be generated and object is locked")
}

func (sc *SafeCache) setAndUnlock(key string, newItem myFavObjType, lockerDone chan struct{}) {
	/*
		consuming a done channel from a locker makes sure we never unlock the
		cache unless told to do so by this unlocker's locker
	*/

	// log here for deadlock detection
	fmt.Println("waiting for done in setAndUnlock")
	<-lockerDone
	fmt.Println("locker done in setAndUnlock")
	defer sc.rw.Unlock()
	sc.data[key] = newItem
	fmt.Println("releasing full lock in setAndUnlock")
}

/*
	Examples of users of the cache
*/

func reader(id int, sc *SafeCache, key string) {
	for {
		fmt.Printf("reader %d reading key %s\n", id, key)
		val, _ := sc.get(key)
		fmt.Printf("reader %d got result: %+v\n", id, val)
		// to slow down terminal output
		time.Sleep(time.Duration(duration) * time.Second)
	}
}

func increasor(id int, sc *SafeCache, key string) {
	for {
		fmt.Printf("increasor %d increasing key %s\n", id, key)
		val, done, _ := sc.getAndLock(key)
		val.field1 += "0"
		*val.field2 += 1.0
		fmt.Printf("increasor %d setting value: %+v\n", id, val)
		sc.setAndUnlock(key, val, done)
		// to slow down terminal output
		time.Sleep(time.Duration(duration) * time.Second)
	}

}

func remover(id int, sc *SafeCache, key string) {
	for {
		fmt.Printf("remover %d removing key %s\n", id, key)
		sc.remove(key)
		fmt.Printf("remover %d removed key: %s\n", id, key)
		// to slow down terminal output
		time.Sleep(time.Duration(duration*3) * time.Second)
	}

}

func setter(id int, sc *SafeCache, key string) {
	for {
		fmt.Printf("setter %d setting key %s\n", id, key)
		val := myFavObjType{"", new(int64)}
		sc.set(key, val)
		fmt.Printf("setter %d set value %+v\n", id, val)
		// to slow down terminal output
		time.Sleep(time.Duration(duration*3) * time.Second)
	}
}

func main() {
	sc := &SafeCache{
		data: make(map[string]myFavObjType),
		rw:   &sync.RWMutex{},
	}
	for i := 0; i < 5; i++ {
		key := fmt.Sprintf("_%x_", i*255)
		go reader(i, sc, key)
		go increasor(i, sc, key)
		go setter(i, sc, key)
		go remover(i, sc, key)
	}
	for {
		time.Sleep(time.Duration(duration) * time.Second)
	}
}
