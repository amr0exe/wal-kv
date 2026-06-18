package service

import (
	"errors"
	"kvstore/internal/store"
	ty "kvstore/internal/types"
)

// common errors
var ErrNoKV = errors.New("key/value found to be empty ...")

// Service layer acts as middleman between handler and store/db
// It handles business logic and validation
type KVService struct {
	db *store.KV
}

func NewKVService(db *store.KV) *KVService {
	return &KVService{db: db}
}

func (s *KVService) Set(key, value string) error {
	if key == "" || value == "" {
		return ErrNoKV
	}

	err := s.db.SET(key, value)
	if err != nil {
		return err
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

	if err := s.db.DEL(key); err != nil {
		return err
	}
	return nil
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
