package neows

import (
	"bytes"
	"fmt"
	"sync"
	"testing"
	"time"
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

	t.Run("Should append to existing key", func(t *testing.T) {
		c := newInMemoryCache()
		k := "key"
		v1 := []byte("value1")
		v2 := []byte("value2")
		c.set(k, v1)
		err := c.append(k, v2)
		if err != nil {
			t.Errorf("expected to append to key %s, but got error %s", k, err)
		}
		v, ok := c.get(k)
		if !ok {
			t.Errorf("expected to find key %s in cache, but didn't", k)
		}
		if string(v) != string(append(v1, v2...)) {
			t.Errorf("expected value of key %s to be %s, but got %s", k, string(append(v1, v2...)), string(v))
		}
	})

	t.Run("Should fail to append to non-existing key", func(t *testing.T) {
		c := newInMemoryCache()
		v := []byte("value")
		err := c.append("non-existing-key", v)
		if err == nil {
			t.Errorf("expected error while appending to non-existing key, but got nil")
		}
	})

	t.Run("Should set and get small value with expiration", func(t *testing.T) {
		c := newInMemoryCache()
		k := "key"
		v := []byte("value")
		c.setWithTTL(k, v, 500*time.Millisecond)
		_, ok := c.get(k)
		if !ok {
			t.Errorf("expected to find key %s in cache, but didn't", k)
		}
		time.Sleep(500 * time.Millisecond)
		_, ok = c.get(k)
		if ok {
			t.Errorf("key %s should have expired, but didn't", k)
		}
	})

	t.Run("Should not find deleted key", func(t *testing.T) {
		c := newInMemoryCache()
		k := "key"
		v := []byte("value")
		c.set(k, v)
		c.del(k)
		_, ok := c.get(k)
		if ok {
			t.Errorf("expected to not find key %s in cache, but did", k)
		}
	})

	t.Run("Should not find any keys after flush", func(t *testing.T) {
		c := newInMemoryCache()
		c.set("key1", []byte("value1"))
		c.set("key2", []byte("value2"))
		c.flush()
		_, ok := c.get("key1")
		if ok {
			t.Errorf("expected to not find key1 in cache, but did")
		}
		_, ok = c.get("key2")
		if ok {
			t.Errorf("expected to not find key2 in cache, but did")
		}
	})

	t.Run("Should set default TTL", func(t *testing.T) {
		c := newInMemoryCache()
		c.setDefaultTTL(500 * time.Millisecond)
		c.set("key", []byte("value"))
		_, ok := c.get("key")
		if !ok {
			t.Errorf("expected to find key in cache, but didn't")
		}
		time.Sleep(500 * time.Millisecond)
		_, ok = c.get("key")
		if ok {
			t.Errorf("key should have expired, but didn't")
		}
	})
}
