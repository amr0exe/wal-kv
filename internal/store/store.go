package store

import (
	"kvstore/internal/types"
	"kvstore/internal/wal"
	"os"
	"sync"
)

// KV defines structure of an in-memory key-value store
// RWmutex is chosen for focusing on read-heavy worloads
type KV struct {
	mu    sync.RWMutex
	store map[string]string
	wal   *wal.WAL

	seq uint32
}

func (kv *KV) nextSeq() uint32 {
	kv.seq++
	return kv.seq
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

	temp := make(map[string]string)
	maxSeq := uint32(0)

	for _, r := range records {
		switch r.Op {
		case types.OpSet:
			temp[r.Key] = r.Value
		case types.OpDel:
			delete(temp, r.Key)
		}

		if r.Seq > maxSeq {
			maxSeq = r.Seq
		}
	}

	kv.store = temp
	kv.seq = maxSeq

	return kv, nil
}

func (kv *KV) SET(k string, v string) error {
	seq := kv.nextSeq()

	if err := kv.wal.Append(seq, types.OpSet, k, v); err != nil {
		return err
	}

	if os.Getenv("CRASH_AFTER_WRITE") == "1" {
		os.Exit(1)
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
	seq := kv.nextSeq()
	if err := kv.wal.Append(seq, types.OpDel, k, ""); err != nil {
		return err
	}

	kv.mu.Lock()
	defer kv.mu.Unlock()

	delete(kv.store, k)
	return nil
}
