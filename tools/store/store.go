package store

import (
	"encoding/json"
	"sync"
)

// @todo remove after https://github.com/golang/go/issues/20135
const ShrinkThreshold = 200 // the number is arbitrary chosen

// Store defines a concurrent safe in memory key-value data store.
type Store[T any] struct {
	data    map[string]T
	mu      sync.RWMutex
	deleted int64
}

// New creates a new Store[T] instance with a shallow copy of the provided data (if any).
func New[T any](data map[string]T) *Store[T] {
	s := &Store[T]{}

	s.Reset(data)

	return s
}

// Reset clears the store and replaces the store data with a
// shallow copy of the provided newData.
func (s *Store[T]) Reset(newData map[string]T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(newData) > 0 {
		s.data = make(map[string]T, len(newData))
		for k, v := range newData {
			s.data[k] = v
		}
	} else {
		s.data = make(map[string]T)
	}

	s.deleted = 0
}

// Length returns the current number of elements in the store.
func (s *Store[T]) Length() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.data)
}

// RemoveAll removes all the existing store entries.
func (s *Store[T]) RemoveAll() {
	s.Reset(nil)
}

// Remove removes a single entry from the store.
//
// Remove does nothing if key doesn't exist in the store.
func (s *Store[T]) Remove(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, key)
	s.deleted++

	// reassign to a new map so that the old one can be gc-ed because it doesn't shrink
	//
	// @todo remove after https://github.com/golang/go/issues/20135
	if s.deleted >= ShrinkThreshold {
		newData := make(map[string]T, len(s.data))
		for k, v := range s.data {
			newData[k] = v
		}
		s.data = newData
		s.deleted = 0
	}
}

// Has checks if element with the specified key exist or not.
func (s *Store[T]) Has(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.data[key]

	return ok
}

// Get returns a single element value from the store.
//
// If key is not set, the zero T value is returned.
func (s *Store[T]) Get(key string) T {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.data[key]
}

// GetOk is similar to Get but returns also a boolean indicating whether the key exists or not.
func (s *Store[T]) GetOk(key string) (T, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	v, ok := s.data[key]

	return v, ok
}

// GetAll returns a shallow copy of the current store data.
func (s *Store[T]) GetAll() map[string]T {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var clone = make(map[string]T, len(s.data))

	for k, v := range s.data {
		clone[k] = v
	}

	return clone
}

// Values returns a slice with all of the current store values.
func (s *Store[T]) Values() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var values = make([]T, 0, len(s.data))

	for _, v := range s.data {
		values = append(values, v)
	}

	return values
}

// Set sets (or overwrite if already exist) a new value for key.
func (s *Store[T]) Set(key string, value T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.data == nil {
		s.data = make(map[string]T)
	}

	s.data[key] = value
}

// GetOrSet retrieves a single existing value for the provided key
// or stores a new one if it doesn't exist.
func (s *Store[T]) GetOrSet(key string, setFunc func() T) T {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.data == nil {
		s.data = make(map[string]T)
	}

	v, ok := s.data[key]
	if !ok {
		v = setFunc()
		s.data[key] = v
	}

	return v
}

// SetIfLessThanLimit sets (or overwrite if already exist) a new value for key.
//
// This method is similar to Set() but **it will skip adding new elements**
// to the store if the store length has reached the specified limit.
// false is returned if maxAllowedElements limit is reached.
func (s *Store[T]) SetIfLessThanLimit(key string, value T, maxAllowedElements int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

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

// UnmarshalJSON implements [json.Unmarshaler] and imports the
// provided JSON data into the store.
//
// The store entries that match with the ones from the data will be overwritten with the new value.
func (s *Store[T]) UnmarshalJSON(data []byte) error {
	raw := map[string]T{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.data == nil {
		s.data = make(map[string]T)
	}

	for k, v := range raw {
		s.data[k] = v
	}

	return nil
}

// MarshalJSON implements [json.Marshaler] and export the current
// store data into valid JSON.
func (s *Store[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.GetAll())
}
