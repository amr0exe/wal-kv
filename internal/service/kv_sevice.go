package service

import (
	"errors"
	"kvstore/internal/store"
	ty "kvstore/internal/types"
)

var ErrNoKV = errors.New("key/value found to be empty ...")

type KVService struct {
	db         *store.KV
	onMutation func(store.Mutation)
}

func NewKVService(db *store.KV) *KVService {
	return &KVService{db: db}
}

func (s *KVService) SetOnMutation(cb func(store.Mutation)) {
	s.onMutation = cb
}

func (s *KVService) Set(key, value string) error {
	if key == "" || value == "" {
		return ErrNoKV
	}

	mut, err := s.db.SET(key, value)
	if err != nil {
		return err
	}
	if s.onMutation != nil {
		s.onMutation(mut)
	}
	return nil
}

func (s *KVService) Get(key string) (error, string) {
	if key == "" {
		return ErrNoKV, ""
	}

	value, ok := s.db.GET(key)
	if !ok {
		return ErrNoKV, ""
	}

	return nil, value
}

func (s *KVService) Del(key string) error {
	if key == "" {
		return ErrNoKV
	}

	mut, err := s.db.DEL(key)
	if err != nil {
		return err
	}
	if s.onMutation != nil {
		s.onMutation(mut)
	}
	return nil
}

func (s *KVService) ApplyMutation(mut store.Mutation) error {
	return s.db.Apply(mut)
}

func (s *KVService) GetSnapshotState() ([]ty.SnapshotRecord, error) {
	storedMutations := s.db.GetAllEntries()
	records := make([]ty.SnapshotRecord, len(storedMutations))
	for i, m := range storedMutations {
		records[i] = ty.SnapshotRecord{
			Seq:   m.Seq,
			Op:    m.Op,
			Key:   m.Key,
			Value: m.Value,
		}
	}
	return records, nil
}
