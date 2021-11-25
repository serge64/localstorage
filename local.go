package localstorage

import (
	"errors"
	"sync"

	"github.com/serge64/localstorage/cache"
)

var (
	ErrNotFoundKey  = errors.New("not found key")
	ErrNotUniqueKey = errors.New("not unique key")
)

type LocalStorage struct {
	db        map[string]interface{}
	mutex     sync.RWMutex
	cache     cache.Cache
	bufKeys   []string
	bufValues []interface{}
}

// New creates a new instance LocalStorage.
//
// The bufferSize argument sets the initial
// size of the buffers.
func New(bufferSize int) LocalStorage {
	return LocalStorage{
		db:        make(map[string]interface{}),
		mutex:     sync.RWMutex{},
		cache:     cache.New(),
		bufKeys:   make([]string, 0, bufferSize),
		bufValues: make([]interface{}, 0, bufferSize),
	}
}

// the Get returns value by key.
func (s *LocalStorage) Get(key string) (interface{}, bool) {
	s.mutex.RLock()
	value, ok := s.get(key)
	s.mutex.RUnlock()
	return value, ok
}

// the Put adds key and value in storage.
//
// If key is not unique, Put returns ErrKeyNotUnique.
func (s *LocalStorage) Put(key string, value interface{}) error {
	s.mutex.Lock()
	err := s.put(key, value)
	s.mutex.Unlock()
	return err
}

// the Del deletes value by key.
//
// If key is not found, Del returns ErrKeyNotFound.
func (s *LocalStorage) Del(key string) error {
	s.mutex.Lock()
	err := s.del(key)
	s.mutex.Unlock()
	return err
}

// the Keys returns keys array.
func (s *LocalStorage) Keys() []string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if s.cache.Cached() {
		return s.bufKeys
	}
	s.bufKeys = s.bufKeys[:0]
	for k := range s.db {
		s.bufKeys = append(s.bufKeys, k)
	}
	s.cache.Save()
	return s.bufKeys
}

// the Values returns values array.
func (s *LocalStorage) Values() []interface{} {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if s.cache.Cached() {
		return s.bufValues
	}
	s.bufValues = s.bufValues[:0]
	for _, v := range s.db {
		s.bufValues = append(s.bufValues, v)
	}
	s.cache.Save()
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
	s.cache.Reset()
	return nil
}

func (s *LocalStorage) del(key string) error {
	if _, ok := s.get(key); !ok {
		return ErrNotFoundKey
	}
	delete(s.db, key)
	s.cache.Reset()
	return nil
}
