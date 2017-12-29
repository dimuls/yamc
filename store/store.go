package store

import (
	"sync"
	"time"
)

type Store interface {
	Get(key string) (string, error)
	Set(key string, value string, ttl time.Duration)
	ListGet(key string, index int) (string, error)
	ListSet(key string, list []string, ttl time.Duration)
	DictGet(key string, dkey string) (string, error)
	DictSet(key string, dict map[string]string, ttl time.Duration)
	Remove(key string) error
	Keys() []string
	StartCleaner() error
	StopCleaner() error
}

type store struct {
	mutex   sync.RWMutex
	params  Params
	clock   Clock
	items   map[string]item
	cleaner *ticker
}

func NewStore(p Params, c Clock) (Store, error) {
	if err := p.Validate(); err != nil {
		return nil, ErrInvalidParams.detailed(err.Error())
	}
	if c == nil {
		return nil, ErrNilClock
	}
	s := &store{
		mutex:  sync.RWMutex{},
		params: p,
		clock:  c,
		items:  map[string]item{},
	}
	return s, nil
}

func (s *store) Get(key string) (string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	i, err := s.get(key)
	if err != nil {
		return "", err
	}
	return i.keyValue()
}

func (s *store) Set(key string, value string, ttl time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.items[key] = newKeyItem(value, s.expiry(ttl))
}

func (s *store) ListGet(key string, index int) (string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	i, err := s.get(key)
	if err != nil {
		return "", err
	}
	return i.listValue(index)
}

func (s *store) ListSet(key string, list []string, ttl time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.items[key] = newListItem(list, s.expiry(ttl))
}

func (s *store) DictGet(key string, dkey string) (string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	i, err := s.get(key)
	if err != nil {
		return "", err
	}
	return i.dictValue(dkey)
}

func (s *store) DictSet(key string, dict map[string]string, ttl time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.items[key] = newDictItem(dict, s.expiry(ttl))
}

func (s *store) Remove(key string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if _, exists := s.items[key]; !exists {
		return ErrKeyNotExists
	}
	delete(s.items, key)
	return nil
}

func (s *store) Keys() []string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	var keys []string
	for k, i := range s.items {
		if i.expired(s.clock.Now()) {
			continue
		}
		keys = append(keys, k)
	}
	return keys
}

func (s *store) StartCleaner() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.cleaner == nil {
		c, err := newTicker(s.params.CleaningPeriod, s.clean)
		if err != nil {
			return ErrFailedToCreateCleaner.detailed(err.Error())
		}
		s.cleaner = c
	}
	if err := s.cleaner.start(); err != nil {
		return ErrFailedToStartCleaner.detailed(err.Error())
	}
	return nil
}

func (s *store) StopCleaner() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.cleaner == nil {
		return ErrCleanerNotStartedYet
	}
	if err := s.cleaner.stop(); err != nil {
		return ErrFailedToStopCleaner.detailed(err.Error())
	}
	return nil
}

func (s *store) get(key string) (item, error) {
	i, exists := s.items[key]
	if !exists || i.expired(s.clock.Now()) {
		return nil, ErrKeyNotExists
	}
	return i, nil
}

func (s *store) clean() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for k, i := range s.items {
		if i.expired(s.clock.Now()) {
			delete(s.items, k)
		}
	}
}

func (s *store) expiry(ttl time.Duration) time.Time {
	return s.clock.Now().Add(ttl)
}
