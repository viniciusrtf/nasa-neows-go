package neows

import (
	"fmt"
	"sync"
	"time"
)

type cache interface {
	get(key string) ([]byte, bool)
	set(key string, value []byte)
	del(key string)
	flush()
}

type entry struct {
	ttl   time.Time
	value []byte
}

type inMemoryCache struct {
	m          map[string]entry
	mu         sync.Mutex
	defaultTTL time.Duration
}

func (imc *inMemoryCache) get(key string) ([]byte, bool) {
	imc.mu.Lock()
	defer imc.mu.Unlock()
	entry, ok := imc.m[key]
	if !ok || time.Now().After(entry.ttl) {
		delete(imc.m, key)
		return nil, false
	}
	return entry.value, ok
}

func (imc *inMemoryCache) set(key string, value []byte) {
	imc.mu.Lock()
	defer imc.mu.Unlock()
	imc.m[key] = entry{
		ttl:   time.Now().Add(imc.defaultTTL),
		value: value,
	}
}

func (imc *inMemoryCache) append(key string, value []byte) error {
	imc.mu.Lock()
	defer imc.mu.Unlock()
	if entry, ok := imc.m[key]; ok {
		entry.value = append(entry.value, value...)
		imc.m[key] = entry
		return nil
	}
	return fmt.Errorf("key %s not found in cache", key)
}

func (imc *inMemoryCache) setWithTTL(key string, value []byte, ttl time.Duration) {
	imc.mu.Lock()
	defer imc.mu.Unlock()
	imc.m[key] = entry{
		ttl:   time.Now().Add(ttl),
		value: value,
	}
}

func (imc *inMemoryCache) del(key string) {
	imc.mu.Lock()
	defer imc.mu.Unlock()
	delete(imc.m, key)
}

func (imc *inMemoryCache) flush() {
	imc.mu.Lock()
	defer imc.mu.Unlock()
	imc.m = map[string]entry{}
}

func (imc *inMemoryCache) setDefaultTTL(ttl time.Duration) {
	imc.mu.Lock()
	defer imc.mu.Unlock()
	imc.defaultTTL = ttl
}

func newInMemoryCache() *inMemoryCache {
	return &inMemoryCache{
		m:          map[string]entry{},
		defaultTTL: 24 * time.Hour,
	}
}
