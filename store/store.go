package store

import (
	"sync"
	"time"
)

// Store is memory store
type Store interface {
	Get(key string) (string, error)
	Set(key string, value string, ttl time.Duration)
	ListGet(key string, index int) (string, error)
	ListSet(key string, list []string, ttl time.Duration)
	DictGet(key string, dkey string) (string, error)
	DictSet(key string, dict map[string]string, ttl time.Duration)
	Remove(key string) error
	Keys() []string
	StartCleaning() error
	StopCleaning() error
	StartDumping() error
	StopDumping() error
}

// store is a store implementation
type store struct {
	mutex    sync.RWMutex
	params   Params
	clock    Clock
	dumper   Dumper
	items    items
	cleaning *ticker
	dumping  *ticker
}

// NewStore constructs new store according params p with clock c. Returns error if params are invalid or clock is nil
func NewStore(p Params, c Clock, d Dumper) (Store, error) {
	err := p.Validate()
	if err != nil {
		return nil, ErrInvalidParams.detailed(err.Error())
	}
	s := &store{
		mutex:  sync.RWMutex{},
		params: p,
		clock:  c,
		dumper: d,
		items:  map[string]item{},
	}
	s.items, err = s.dumper.load()
	if err != nil {
		return nil, err
	}
	return s, nil
}

// Get returns value by key. Errors if key is not exists or item is not simple keyItem
func (s *store) Get(key string) (string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	i, err := s.get(key)
	if err != nil {
		return "", err
	}
	return i.keyValue()
}

// Set sets value by key with time to live ttl. Creates new or overrides old of any type
func (s *store) Set(key string, value string, ttl time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.items[key] = newKeyItem(value, s.expiry(ttl))
}

// Get returns value by key and list index. Errors if key is not exists or key item is not listItem
func (s *store) ListGet(key string, index int) (string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	i, err := s.get(key)
	if err != nil {
		return "", err
	}
	return i.listValue(index)
}

// Set sets list by key with time to live ttl. Creates new or overrides old of any type
func (s *store) ListSet(key string, list []string, ttl time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.items[key] = newListItem(list, s.expiry(ttl))
}

// Get returns value by key and dict key dkey. Errors if key is not exists or key item is not simple dictItem
func (s *store) DictGet(key string, dkey string) (string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	i, err := s.get(key)
	if err != nil {
		return "", err
	}
	return i.dictValue(dkey)
}

// Set sets dict by key with time to live ttl. Creates new or overrides old of any type
func (s *store) DictSet(key string, dict map[string]string, ttl time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.items[key] = newDictItem(dict, s.expiry(ttl))
}

// Remove removes item of any type by key. Errors if key is not exists
func (s *store) Remove(key string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if _, exists := s.items[key]; !exists {
		return ErrKeyNotExists
	}
	delete(s.items, key)
	return nil
}

// Keys returns all keys list, not sorted
func (s *store) Keys() []string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	var keys []string
	for k, i := range s.items {
		if i.expired(s.clock.now()) {
			continue
		}
		keys = append(keys, k)
	}
	return keys
}

// StartCleaning starts periodical expired items cleaning. Can be called multiple times
func (s *store) StartCleaning() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.cleaning == nil {
		c, err := newTicker(s.params.CleaningPeriod, s.clean)
		if err != nil {
			return ErrFailToCreateCleaning.detailed(err.Error())
		}
		s.cleaning = c
	}
	if err := s.cleaning.start(); err != nil {
		return ErrFailToStartCleaning.detailed(err.Error())
	}
	return nil
}

// StopCleaning stops periodical expired items cleaning. Can be called multiple times
func (s *store) StopCleaning() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.cleaning == nil {
		return ErrCleaningNotStartedYet
	}
	if err := s.cleaning.stop(); err != nil {
		return ErrFailToStopCleaning.detailed(err.Error())
	}
	return nil
}

// StartDumping starts periodical file dumping. Can be called multiple times
func (s *store) StartDumping() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.dumping == nil {
		d, err := newTicker(s.params.DumpingPeriod, s.dump)
		if err != nil {
			return ErrFailToCreateDumping.detailed(err.Error())
		}
		s.dumping = d
	}
	if err := s.dumping.start(); err != nil {
		return ErrFailToStartDumping.detailed(err.Error())
	}
	return nil
}

// StopDumping stops periodical file dumping. Can be called multiple times
func (s *store) StopDumping() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.dumping == nil {
		return ErrDumperNotStartedYet
	}
	if err := s.dumping.stop(); err != nil {
		return ErrFailToStopDumping.detailed(err.Error())
	}
	return nil
}

// get is item getter. Returns error if key is not exists
func (s *store) get(key string) (item, error) {
	i, exists := s.items[key]
	if !exists || i.expired(s.clock.now()) {
		return nil, ErrKeyNotExists
	}
	return i, nil
}

// clean removes all expired items
func (s *store) clean() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for k, i := range s.items {
		if i.expired(s.clock.now()) {
			delete(s.items, k)
		}
	}
}

func (s *store) dump() {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if err := s.dumper.dump(s.items); err != nil {
		// TODO: log error
	}
}

// expiry computes expire time according clock's now and given ttl
func (s *store) expiry(ttl time.Duration) time.Time {
	return s.clock.now().Add(ttl)
}
