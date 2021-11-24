package localstorage

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrNotFoundKey  = errors.New("not found key")
	ErrNotUniqueKey = errors.New("not unique key")
)

type LocalStorage struct {
	db        map[string]interface{}
	mutex     sync.RWMutex
	cache     atomic.Value
	bufKeys   []string
	bufValues []interface{}
}

// New creates a new instance LocalStorage.
//
// The bufferSize argument sets the initial
// size of the buffers.
func New(ctx context.Context, bufferSize int) LocalStorage {
	return LocalStorage{
		db:        make(map[string]interface{}),
		mutex:     sync.RWMutex{},
		cache:     atomic.Value{},
		bufKeys:   make([]string, 0, bufferSize),
		bufValues: make([]interface{}, 0, bufferSize),
	}
}

// the Get returns value by key.
func (s *LocalStorage) Get(key string) (interface{}, bool) {
	s.mutex.RLock()
	wrap, ok := s.get(key)
	s.mutex.RUnlock()
	return wrap, ok
}

// the Put adds key and value in storage.
func (s *LocalStorage) Put(key string, value interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.put(key, value)

}

// the Del deletes value by key.
func (s *LocalStorage) Del(key string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if _, ok := s.get(key); !ok {
		return ErrNotFoundKey
	}
	s.del(key)
	return nil
}

// the Keys returns keys array.
func (s *LocalStorage) Keys() []string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if s.cached() {
		return s.bufKeys
	}
	s.bufKeys = s.bufKeys[:0]
	for k := range s.db {
		s.bufKeys = append(s.bufKeys, k)
	}
	s.newCache()
	return s.bufKeys
}

// the Values returns values array.
func (s *LocalStorage) Values() []interface{} {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if s.cached() {
		return s.bufValues
	}
	s.bufValues = s.bufValues[:0]
	for _, v := range s.db {
		s.bufValues = append(s.bufValues, v)
	}
	s.newCache()
	return s.bufValues
}

func (s *LocalStorage) get(key string) (interface{}, bool) {
	value, ok := s.db[key]
	return value, ok
}

func (s *LocalStorage) put(key string, value interface{}) error {
	if _, ok := s.get(key); ok {
		return ErrNotUniqueKey
	}
	s.db[key] = value
	s.resetCache()
	return nil
}

func (s *LocalStorage) del(key string) {
	delete(s.db, key)
	s.resetCache()
}

func (s *LocalStorage) newCache() {
	s.cache.Store(true)
}

func (s *LocalStorage) resetCache() {
	s.cache.Store(false)
}

func (s *LocalStorage) cached() bool {
	cache := s.cache.Load()
	if cache == nil {
		return false
	}
	return cache.(bool)
}
