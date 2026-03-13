package pokecache

import (
	"sync"
	"time"
	"fmt"
)

type Cache struct {
	Mu *sync.Mutex  // works properly only with pointer
	Entries map[string]cacheEntry
	Interval time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val []byte
}

func NewCache(intervalSeconds time.Duration) *Cache {
	cache := Cache{
        Mu: &sync.Mutex{},
		Entries: map[string]cacheEntry{},
		Interval: intervalSeconds * time.Second,
    }

	// concurrently clears old entries
	go cache.reapLoop()

	return &cache
}

func (c Cache) Add(key string, val []byte) {
	c.Mu.Lock()
	c.Entries[key] = cacheEntry{val: val, createdAt: time.Now()}
	c.Mu.Unlock()
}

func (c Cache) Get(key string) ([]byte, bool) {
	c.Mu.Lock()
	entry, ex := c.Entries[key]
	c.Mu.Unlock()
	return entry.val, ex
}

// deletes old enough entries
func (c Cache) reapLoop() {
	for range time.Tick(c.Interval) {
		// https://stackoverflow.com/questions/26285735/subtracting-time-duration-from-time-in-go
		borderTime := time.Now().Add(time.Duration(-c.Interval))
		for key, entry := range c.Entries {
			if entry.createdAt.Before(borderTime) {
				c.Mu.Lock()
				delete(c.Entries, key)
				c.Mu.Unlock()
			}
		}
	}
}

// print .createdAt for each entry in cache (for debug)
func (c Cache) PrintEntriesTime() {
	for _, entry := range c.Entries{
		fmt.Printf("%v\n", entry.createdAt)
	}
}


// func main() {
// 	c := map[string]int{"a": 1, "b": 2, "c": 3}
// 	for _, some := range c {
// 		fmt.Printf("%v\n", some)
// 	}
// 	now := time.Now()
// 	prev := now.Add(time.Duration(-time.Minute * 3600))
// 	fmt.Printf("%v\n", now)
// 	fmt.Printf("%v\n", prev)
// 	fmt.Printf("%v\n", prev.Before(now))
// }
