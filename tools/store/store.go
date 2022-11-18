package store

import "sync"

// Store defines a concurrent safe in memory key-value data store.
type Store[T any] struct {
	mux  sync.RWMutex
	data map[string]T
}

// New creates a new Store[T] instance.
func New[T any](data map[string]T) *Store[T] {
	return &Store[T]{data: data}
}

// RemoveAll removes all the existing store entries.
func (s *Store[T]) RemoveAll() {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.data = make(map[string]T)
}

// Remove removes a single entry from the store.
//
// Remove does nothing if key doesn't exist in the store.
func (s *Store[T]) Remove(key string) {
	s.mux.Lock()
	defer s.mux.Unlock()

	delete(s.data, key)
}

// Has checks if element with the specified key exist or not.
func (s *Store[T]) Has(key string) bool {
	s.mux.RLock()
	defer s.mux.RUnlock()

	_, ok := s.data[key]

	return ok
}

// Get returns a single element value from the store.
//
// If key is not set, the zero T value is returned.
func (s *Store[T]) Get(key string) T {
	s.mux.RLock()
	defer s.mux.RUnlock()

	return s.data[key]
}

// Set sets (or overwrite if already exist) a new value for key.
func (s *Store[T]) Set(key string, value T) {
	s.mux.Lock()
	defer s.mux.Unlock()

	if s.data == nil {
		s.data = make(map[string]T)
	}

	s.data[key] = value
}

// SetIfLessThanLimit sets (or overwrite if already exist) a new value for key.
//
// This method is similar to Set() but **it will skip adding new elements**
// to the store if the store length has reached the specified limit.
// `false` is returned if maxAllowedElements limit is reached.
func (s *Store[T]) SetIfLessThanLimit(key string, value T, maxAllowedElements int) bool {
	s.mux.Lock()
	defer s.mux.Unlock()

	// init map if not already
	if s.data == nil {
		s.data = make(map[string]T)
	}

	// check for existing item
	_, ok := s.data[key]

	if !ok && len(s.data) >= maxAllowedElements {
		// cannot add more items
		return false
	}

	// add/overwrite item
	s.data[key] = value

	return true
}
