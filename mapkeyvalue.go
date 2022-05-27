package r9e

import (
	"reflect"
	"sort"
	"sync"
)

type mapKeyValueOptions struct {
	size int
}

// MapKeyValueOptions are the options for MapKeyValue.
type MapKeyValueOptions func(*mapKeyValueOptions)

// WithCapacity sets the initial capacity allocation of the MapKeyValue.
func WithCapacity(size int) MapKeyValueOptions {
	return func(kv *mapKeyValueOptions) {
		kv.size = size
	}
}

// MapKeyValue is a generic key-value store container that is thread-safe.
// This use a golang native map data structure as underlying data structure and a mutex to
// protect the data.
type MapKeyValue[K comparable, T any] struct {
	mu   sync.RWMutex
	data map[K]T
}

// kv is a helper struct to sort the values of the map.
type kv[K comparable, T any] struct {
	key   K
	value T
}

// NewMapKeyValue returns a new MapKeyValue container.
func NewMapKeyValue[K comparable, T any](options ...MapKeyValueOptions) *MapKeyValue[K, T] {
	// parse options
	kvo := mapKeyValueOptions{}
	for _, opt := range options {
		opt(&kvo)
	}

	return &MapKeyValue[K, T]{
		data: make(map[K]T, kvo.size),
	}
}

// Set sets the value associated with the key.
func (r *MapKeyValue[K, T]) Set(key K, value T) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[key] = value
}

// GetCheck returns the value associated with the key if this exist also a if this exist.
func (r *MapKeyValue[K, T]) GetCheck(key K) (T, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	value, ok := r.data[key]
	return value, ok
}

// Get returns the value associated with the key.
func (r *MapKeyValue[K, T]) Get(key K) T {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.data[key]
}

// Delete deletes the value associated with the key.
func (r *MapKeyValue[K, T]) Delete(key K) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.data, key)
}

// Clear deletes all key-value pairs in the container.
func (r *MapKeyValue[K, T]) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data = make(map[K]T)
}

// Size returns the number of key-value pairs stored in the container.
func (r *MapKeyValue[K, T]) Size() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.data)
}

// IsEmpty returns true if the container is empty.
func (r *MapKeyValue[K, T]) IsEmpty() bool {
	return r.Size() == 0
}

// IsFull returns true if the container has elements.
func (r *MapKeyValue[K, T]) IsFull() bool {
	return r.Size() != 0
}

// ContainsKey returns true if the key is in the container.
func (r *MapKeyValue[K, T]) ContainsKey(key K) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, ok := r.data[key]
	return ok
}

// ContainsValue returns true if the value is in the container.
func (r *MapKeyValue[K, T]) ContainsValue(value T) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, v := range r.data {
		if reflect.DeepEqual(v, value) {
			return true
		}
	}

	return false
}

// Get returns the value associated with the key.
func (r *MapKeyValue[K, T]) Key(key K) K {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if _, ok := r.data[key]; ok {
		return key
	}
	var empty K
	return empty
}

// Keys returns all keys stored in the container.
func (r *MapKeyValue[K, T]) Keys() []K {
	r.mu.RLock()
	defer r.mu.RUnlock()

	keys := make([]K, 0, len(r.data))
	for key := range r.data {
		keys = append(keys, key)
	}
	return keys
}

// Values returns all values stored in the container.
func (r *MapKeyValue[K, T]) Values() []T {
	r.mu.RLock()
	defer r.mu.RUnlock()

	values := make([]T, 0, len(r.data))
	for _, value := range r.data {
		values = append(values, value)
	}
	return values
}

// Each calls the given function for each key-value pair in the container.
func (r *MapKeyValue[K, T]) Each(fn func(key K, value T)) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for key, value := range r.data {
		fn(key, value)
	}
}

// EachKey calls the given function for each key in the container.
func (r *MapKeyValue[K, T]) EachKey(fn func(key K)) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for key := range r.data {
		fn(key)
	}
}

// EachValue calls the given function for each value in the container.
func (r *MapKeyValue[K, T]) EachValue(fn func(value T)) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, value := range r.data {
		fn(value)
	}
}

// Clone returns a new MapKeyValue with a copy of the underlying data.
func (r *MapKeyValue[K, T]) Clone() *MapKeyValue[K, T] {
	r.mu.RLock()
	defer r.mu.RUnlock()

	clone := NewMapKeyValue[K, T](WithCapacity(r.Size()))
	for key, value := range r.data {
		clone.Set(key, value)
	}
	return clone
}

// CloneAndClear returns a new MapKeyValue with a copy of the underlying data and clears the container.
func (r *MapKeyValue[K, T]) CloneAndClear() *MapKeyValue[K, T] {
	r.mu.RLock()
	defer r.mu.RUnlock()

	clone := NewMapKeyValue[K, T](WithCapacity(r.Size()))
	for key, value := range r.data {
		clone.Set(key, value)
	}
	r.data = make(map[K]T)
	return clone
}

// DeepEqual returns true if the kv is deep equal to the MapKeyValue container
func (r *MapKeyValue[K, T]) DeepEqual(kv *MapKeyValue[K, T]) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.Size() != kv.Size() {
		return false
	}

	for key, value := range r.data {
		if !reflect.DeepEqual(value, kv.Get(key)) {
			return false
		}
	}

	return true
}

// Map returns a new MapKeyValue after applying the given function fn to each key-value pair.
func (r *MapKeyValue[K, T]) Map(fn func(key K, value T) (newKey K, newValue T)) *MapKeyValue[K, T] {
	r.mu.RLock()
	defer r.mu.RUnlock()

	m := NewMapKeyValue[K, T](WithCapacity(r.Size()))
	for key, value := range r.data {
		newKey, newValue := fn(key, value)
		m.Set(newKey, newValue)
	}
	return m
}

// MapKey returns a new MapKeyValue after applying the given function fn to each key.
func (r *MapKeyValue[K, T]) MapKey(fn func(key K) K) *MapKeyValue[K, T] {
	r.mu.RLock()
	defer r.mu.RUnlock()

	m := NewMapKeyValue[K, T](WithCapacity(r.Size()))
	for key := range r.data {
		newKey := fn(key)
		m.Set(newKey, r.data[key])
	}
	return m
}

// MapValue returns a new MapKeyValue after applying the given function fn to each value.
func (r *MapKeyValue[K, T]) MapValue(fn func(value T) T) *MapKeyValue[K, T] {
	r.mu.RLock()
	defer r.mu.RUnlock()

	m := NewMapKeyValue[K, T](WithCapacity(r.Size()))
	for key, value := range r.data {
		newValue := fn(value)
		m.Set(key, newValue)
	}
	return m
}

// Filter returns a new MapKeyValue after applying the given function fn to each key-value pair.
func (r *MapKeyValue[K, T]) Filter(fn func(key K, value T) bool) *MapKeyValue[K, T] {
	r.mu.RLock()
	defer r.mu.RUnlock()

	m := NewMapKeyValue[K, T](WithCapacity(r.Size()))
	for key, value := range r.data {
		if fn(key, value) {
			m.Set(key, value)
		}
	}
	return m
}

// FilterKey returns a new MapKeyValue after applying the given function fn to each key.
func (r *MapKeyValue[K, T]) FilterKey(fn func(key K) bool) *MapKeyValue[K, T] {
	r.mu.RLock()
	defer r.mu.RUnlock()

	m := NewMapKeyValue[K, T](WithCapacity(r.Size()))
	for key := range r.data {
		if fn(key) {
			m.Set(key, r.data[key])
		}
	}
	return m
}

// FilterValue returns a new MapKeyValue after applying the given function fn to each value.
func (r *MapKeyValue[K, T]) FilterValue(fn func(value T) bool) *MapKeyValue[K, T] {
	r.mu.RLock()
	defer r.mu.RUnlock()

	m := NewMapKeyValue[K, T](WithCapacity(r.Size()))
	for key, value := range r.data {
		if fn(value) {
			m.Set(key, value)
		}
	}
	return m
}

// Partition returns two new MapKeyValue. One with all the elements that satisfy the predicate and
// another with the rest. The predicate is applied to each element.
func (r *MapKeyValue[K, T]) Partition(fn func(key K, value T) bool) (match, others *MapKeyValue[K, T]) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	match = NewMapKeyValue[K, T](WithCapacity(r.Size()))
	others = NewMapKeyValue[K, T](WithCapacity(r.Size()))
	for key, value := range r.data {
		if fn(key, value) {
			match.Set(key, value)
		} else {
			others.Set(key, value)
		}
	}
	return
}

// PartitionKey returns two new MapKeyValue. One with all the elements that satisfy the predicate and
// another with the rest. The predicate is applied to each key.
func (r *MapKeyValue[K, T]) PartitionKey(fn func(key K) bool) (match, others *MapKeyValue[K, T]) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	match = NewMapKeyValue[K, T](WithCapacity(r.Size()))
	others = NewMapKeyValue[K, T](WithCapacity(r.Size()))
	for key := range r.data {
		if fn(key) {
			match.Set(key, r.data[key])
		} else {
			others.Set(key, r.data[key])
		}
	}
	return
}

// PartitionValue returns two new MapKeyValue. One with all the elements that satisfy the predicate and
// another with the rest. The predicate is applied to each value.
func (r *MapKeyValue[K, T]) PartitionValue(fn func(value T) bool) (match, others *MapKeyValue[K, T]) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	match = NewMapKeyValue[K, T](WithCapacity(r.Size()))
	others = NewMapKeyValue[K, T](WithCapacity(r.Size()))
	for key, value := range r.data {
		if fn(value) {
			match.Set(key, value)
		} else {
			others.Set(key, value)
		}
	}
	return
}

// SortKeys returns a new MapKeyValue after sorting the keys using the given sortFn function.
func (r *MapKeyValue[K, T]) SortKeys(sortFn func(key1, key2 K) bool) *MapKeyValue[K, T] {
	r.mu.RLock()
	defer r.mu.RUnlock()

	keys := r.Keys()

	sort.Slice(keys, func(i, j int) bool {
		return sortFn(keys[i], keys[j])
	})

	m := NewMapKeyValue[K, T](WithCapacity(r.Size()))
	for _, key := range keys {
		m.Set(key, r.data[key])
	}
	return m
}

// SortValues returns a new MapKeyValue after sorting the values using given function sortFn.
func (r *MapKeyValue[K, T]) SortValues(sortFn func(value1, value2 T) bool) *MapKeyValue[K, T] {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var kvs []kv[K, T]
	for key, value := range r.data {
		kvs = append(kvs, kv[K, T]{key, value})
	}

	sort.Slice(kvs, func(i, j int) bool {
		return sortFn(kvs[i].value, kvs[j].value)
	})

	m := NewMapKeyValue[K, T](WithCapacity(r.Size()))
	for _, pair := range kvs {
		m.Set(pair.key, pair.value)
	}

	return m
}
