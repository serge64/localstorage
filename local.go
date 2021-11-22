package localstorage

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	ErrNotFoundKey  = errors.New("not found key")
	ErrNotUniqueKey = errors.New("not unique key")
)

type wrapper struct {
	value interface{}
	ttl   time.Time
}

type LocalStorage struct {
	db        map[string]wrapper
	mutex     *sync.RWMutex
	bufKeys   []string
	bufValues []interface{}
}

func NewLocalStorage(ctx context.Context, bufferSize int) LocalStorage {
	s := LocalStorage{
		db:        make(map[string]wrapper),
		mutex:     &sync.RWMutex{},
		bufKeys:   make([]string, 0, bufferSize),
		bufValues: make([]interface{}, 0, bufferSize),
	}
	go s.scheduler(ctx)
	return s
}

func (s LocalStorage) scheduler(ctx context.Context) {
	tick := time.NewTicker(time.Duration(50) * time.Microsecond)
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			s.mutex.Lock()
			for k, v := range s.db {
				if time.Now().After(v.ttl) {
					s.del(k)
				}
			}
			s.mutex.Unlock()
		}
	}
}

func (s LocalStorage) Get(key string) (interface{}, bool) {
	s.mutex.RLock()
	wrap, ok := s.get(key)
	s.mutex.RUnlock()
	return wrap.value, ok
}

func (s LocalStorage) Put(key string, value interface{}, ttl time.Duration) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if ttl > 0 {
		return s.put(key, wrapper{value: value, ttl: time.Now().Add(ttl)})
	}
	return s.put(key, wrapper{value: value, ttl: time.Time{}})

}

func (s LocalStorage) Del(key string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if _, ok := s.get(key); !ok {
		return ErrNotFoundKey
	}
	s.del(key)
	return nil
}

func (s LocalStorage) Keys() []string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	s.bufKeys = s.bufKeys[:0]
	for k := range s.db {
		s.bufKeys = append(s.bufKeys, k)
	}
	return s.bufKeys
}

func (s LocalStorage) Values() []interface{} {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	s.bufValues = s.bufValues[:0]
	for _, v := range s.db {
		s.bufValues = append(s.bufValues, v.value)
	}
	return s.bufValues
}

func (s LocalStorage) get(key string) (wrapper, bool) {
	wrap, ok := s.db[key]
	return wrap, ok
}

func (s LocalStorage) put(key string, value wrapper) error {
	if _, ok := s.get(key); ok {
		return ErrNotUniqueKey
	}
	s.db[key] = value
	return nil
}

func (s LocalStorage) del(key string) {
	delete(s.db, key)
}
