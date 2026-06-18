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
	seq   uint32
}

type Mutation struct {
	Seq   uint32
	Op    types.OpType
	Key   string
	Value string
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

func (kv *KV) Apply(m Mutation) error {
	if kv.wal != nil {
		if err := kv.wal.Append(m.Seq, m.Op, m.Key, m.Value); err != nil {
			return err
		}
	}

	kv.mu.Lock()
	defer kv.mu.Unlock()

	switch m.Op {
	case types.OpSet:
		kv.store[m.Key] = m.Value
	case types.OpDel:
		delete(kv.store, m.Key)
	}

	kv.seq = m.Seq

	return nil
}

func (kv *KV) SET(k, v string) (Mutation, error) {
	mut := Mutation{
		Seq:   kv.nextSeq(),
		Op:    types.OpSet,
		Key:   k,
		Value: v,
	}

	if err := kv.Apply(mut); err != nil {
		return Mutation{}, err
	}

	return mut, nil
}

func (kv *KV) GET(k string) (string, bool) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	v, ok := kv.store[k]
	return v, ok
}

func (kv *KV) DEL(k string) (Mutation, error) {
	mut := Mutation{
		Seq: kv.nextSeq(),
		Op:  types.OpDel,
		Key: k,
	}

	if err := kv.Apply(mut); err != nil {
		return Mutation{}, err
	}
	return mut, nil
}

func (kv *KV) GetAllEntries() []Mutation {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	var index uint32 = 1
	mutations := make([]Mutation, 0, len(kv.store))

	for k, v := range kv.store {
		mutations = append(mutations, Mutation{
			Seq:   index,
			Op:    types.OpSet,
			Key:   k,
			Value: v,
		})

		index++
	}

	return mutations
}

func NewKVStoreInMemory() *KV {
	return &KV{
		store: make(map[string]string),
	}
}

func (kv *KV) LoadSnapshot(mutations []Mutation) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	for _, m := range mutations {
		switch m.Op {
		case types.OpSet:
			kv.store[m.Key] = m.Value
		case types.OpDel:
			delete(kv.store, m.Key)
		}
		if m.Seq > kv.seq {
			kv.seq = m.Seq
		}
	}
}
