package main

import (
	"testing"
	"time"
	"internal/pokecache"
)

func TestCaching(t *testing.T) {
	const waitTimeSeconds = 1
	cache := pokecache.NewCache(waitTimeSeconds)
	cache.Add("I am", []byte("based"))
	
	_, ok := cache.Get("I am")
	if !ok {
		t.Errorf("expected to find key")
		return
	}
}

func TestReapLoop(t *testing.T) {
	const baseTimeSeconds = 1
	const waitTime = (baseTimeSeconds + 1) * time.Second
	cache := pokecache.NewCache(baseTimeSeconds)
	cache.Add("https://example.com", []byte("testdata"))

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)

	entry, ok := cache.Get("https://example.com")
	if ok {
		t.Errorf("expected to not find key %v", entry)
		return
	}
}

