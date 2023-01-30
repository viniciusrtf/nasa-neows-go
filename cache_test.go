package neows

import (
	"bytes"
	"fmt"
	"sync"
	"testing"
)

func TestInMemoryCache(t *testing.T) {
	t.Run("Should set and get small value", func(t *testing.T) {
		// Initialize cache
		cache := newInMemoryCache()

		// Set key1
		key := "key"
		value := []byte("value")
		cache.set(key, value)

		// Test get key1
		v, ok := cache.get(key)
		if !ok {
			t.Errorf("expected to find \"%s\" in cache, but didn't", key)
		}
		if string(v) != string(value) {
			t.Errorf("expected %s value to be %s, but got %s", key, value, v)
		}
	})

	t.Run("Should set and get large value", func(t *testing.T) {
		// Initialize cache
		cache := newInMemoryCache()

		// Set 10Kb value
		key := "key"
		value := make([]byte, 10*1024)
		cache.set(key, value)

		v, ok := cache.get(key)
		if !ok {
			t.Errorf("expected to find \"%s\" in cache, but didn't", key)
		}
		if !bytes.Equal(v, value) {
			t.Errorf("expected value of length %d, got %d", len(value), len(v))
		}
	})

	t.Run("Should not find non existing key", func(t *testing.T) {
		// Initialize cache
		cache := newInMemoryCache()

		// Test get non-existing-key
		_, ok := cache.get("non-existing-key")
		if ok {
			t.Error("expected to not find \"non-existing-key\" in cache, but did")
		}
	})

	t.Run("Should be safe to read and write concurrently", func(t *testing.T) {
		c := newInMemoryCache()
		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				c.set(fmt.Sprintf("key%d", i), []byte(fmt.Sprintf("value%d", i)))
			}
		}()

		go func() {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				c.get(fmt.Sprintf("key%d", i))
			}
		}()

		wg.Wait()

		for i := 0; i < 1000; i++ {
			v, ok := c.get(fmt.Sprintf("key%d", i))
			if !ok {
				t.Errorf("expected key%d to exist in the cache", i)
			}
			if string(v) != fmt.Sprintf("value%d", i) {
				t.Errorf("expected value of key%d to be %s, but got %s", i, fmt.Sprintf("value%d", i), string(v))
			}
		}
	})
}
