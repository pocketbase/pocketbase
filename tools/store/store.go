package store

import (
	"encoding/json"
	"sync"
)

// @todo remove after https://github.com/golang/go/issues/20135
const ShrinkThreshold = 200 // the number is arbitrary chosen

// Store defines a concurrent safe in memory key-value data store.
type Store[K comparable, T any] struct {
	data    map[K]T
	mu      sync.RWMutex
	deleted int64
}

// New creates a new Store[T] instance with a shallow copy of the provided data (if any).
func New[K comparable, T any](data map[K]T) *Store[K, T] {
	s := &Store[K, T]{}

	s.Reset(data)

	return s
}

// Reset clears the store and replaces the store data with a
// shallow copy of the provided newData.
func (s *Store[K, T]) Reset(newData map[K]T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(newData) > 0 {
		s.data = make(map[K]T, len(newData))
		for k, v := range newData {
			s.data[k] = v
		}
	} else {
		s.data = make(map[K]T)
	}

	s.deleted = 0
}

// Length returns the current number of elements in the store.
func (s *Store[K, T]) Length() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.data)
}

// RemoveAll removes all the existing store entries.
func (s *Store[K, T]) RemoveAll() {
	s.Reset(nil)
}

// Remove removes a single entry from the store.
//
// Remove does nothing if key doesn't exist in the store.
func (s *Store[K, T]) Remove(key K) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, key)
	s.deleted++

	// reassign to a new map so that the old one can be gc-ed because it doesn't shrink
	//
	// @todo remove after https://github.com/golang/go/issues/20135
	if s.deleted >= ShrinkThreshold {
		newData := make(map[K]T, len(s.data))
		for k, v := range s.data {
			newData[k] = v
		}
		s.data = newData
		s.deleted = 0
	}
}

// Has checks if element with the specified key exist or not.
func (s *Store[K, T]) Has(key K) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.data[key]

	return ok
}

// Get returns a single element value from the store.
//
// If key is not set, the zero T value is returned.
func (s *Store[K, T]) Get(key K) T {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.data[key]
}

// GetOk is similar to Get but returns also a boolean indicating whether the key exists or not.
func (s *Store[K, T]) GetOk(key K) (T, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	v, ok := s.data[key]

	return v, ok
}

// GetAll returns a shallow copy of the current store data.
func (s *Store[K, T]) GetAll() map[K]T {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var clone = make(map[K]T, len(s.data))

	for k, v := range s.data {
		clone[k] = v
	}

	return clone
}

// Values returns a slice with all of the current store values.
func (s *Store[K, T]) Values() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var values = make([]T, 0, len(s.data))

	for _, v := range s.data {
		values = append(values, v)
	}

	return values
}

// Set sets (or overwrite if already exists) a new value for key.
func (s *Store[K, T]) Set(key K, value T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.data == nil {
		s.data = make(map[K]T)
	}

	s.data[key] = value
}

// SetFunc sets (or overwrite if already exists) a new value resolved
// from the function callback for the provided key.
//
// The function callback receives as argument the old store element value (if exists).
// If there is no old store element, the argument will be the T zero value.
//
// Example:
//
//	s := store.New[string, int](nil)
//	s.SetFunc("count", func(old int) int {
//	    return old + 1
//	})
func (s *Store[K, T]) SetFunc(key K, fn func(old T) T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.data == nil {
		s.data = make(map[K]T)
	}

	s.data[key] = fn(s.data[key])
}

// GetOrSet retrieves a single existing value for the provided key
// or stores a new one if it doesn't exist.
func (s *Store[K, T]) GetOrSet(key K, setFunc func() T) T {
	// lock only reads to minimize locks contention
	s.mu.RLock()
	v, ok := s.data[key]
	s.mu.RUnlock()

	if !ok {
		s.mu.Lock()
		v = setFunc()
		if s.data == nil {
			s.data = make(map[K]T)
		}
		s.data[key] = v
		s.mu.Unlock()
	}

	return v
}

// SetIfLessThanLimit sets (or overwrite if already exist) a new value for key.
//
// This method is similar to Set() but **it will skip adding new elements**
// to the store if the store length has reached the specified limit.
// false is returned if maxAllowedElements limit is reached.
func (s *Store[K, T]) SetIfLessThanLimit(key K, value T, maxAllowedElements int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.data == nil {
		s.data = make(map[K]T)
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

// UnmarshalJSON implements [json.Unmarshaler] and imports the
// provided JSON data into the store.
//
// The store entries that match with the ones from the data will be overwritten with the new value.
func (s *Store[K, T]) UnmarshalJSON(data []byte) error {
	raw := map[K]T{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.data == nil {
		s.data = make(map[K]T)
	}

	for k, v := range raw {
		s.data[k] = v
	}

	return nil
}

// MarshalJSON implements [json.Marshaler] and export the current
// store data into valid JSON.
func (s *Store[K, T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.GetAll())
}
