package store

import "sync"

// KV defines structure of an in-memory key-value store
// RWmutex is chosen for focusing on read-heavy worloads
type KV struct {
	mu    sync.RWMutex
	store map[string]string
}

func NewKVStore() *KV {
	return &KV{
		store: make(map[string]string),
	}
}

func (kv *KV) SET(k string, v string) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	kv.store[k] = v
}

func (kv *KV) GET(k string) (string, bool) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	v, ok := kv.store[k]
	return v, ok
}

func (kv *KV) DEL(k string) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	delete(kv.store, k)
}
