package store

import (
	"kvstore/internal/types"
	"kvstore/internal/wal"
	"sync"
)

// KV defines structure of an in-memory key-value store
// RWmutex is chosen for focusing on read-heavy worloads
type KV struct {
	mu    sync.RWMutex
	store map[string]string
	wal   *wal.WAL
}

func NewKVStore() (*KV, error) {
	w, err := wal.NewWAL()
	if err != nil {
		return nil, err
	}
	kv := &KV{
		store: make(map[string]string),
		wal:   w,
	}

	records, err := w.Recover()
	if err != nil {
		return nil, err
	}

	for _, r := range records {
		switch r.Op {
		case types.OpSet:
			kv.store[r.Key] = r.Value
		case types.OpDel:
			delete(kv.store, r.Key)
		}
	}

	return kv, nil
}

func (kv *KV) SET(k string, v string) error {
	if err := kv.wal.Append(types.OpSet, k, v); err != nil {
		return err
	}

	kv.mu.Lock()
	defer kv.mu.Unlock()

	kv.store[k] = v
	return nil
}

func (kv *KV) GET(k string) (string, bool) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	v, ok := kv.store[k]
	return v, ok
}

func (kv *KV) DEL(k string) error {
	if err := kv.wal.Append(types.OpDel, k, ""); err != nil {
		return err
	}
	
	kv.mu.Lock()
	defer kv.mu.Unlock()

	delete(kv.store, k)
	return nil
}
