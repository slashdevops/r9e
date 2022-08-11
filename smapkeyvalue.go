package r9e

import (
	"reflect"
	"sort"
	"sync"
	"sync/atomic"
)

// SMapKeyValue is a generic key-value store container that is thread-safe.
// This use a golang native sync.Map data structure as underlying data structure.
type SMapKeyValue[K comparable, T any] struct {
	count atomic.Uint64
	data  sync.Map
}

// skv is a helper struct to sort the values of the SMapKeyValue container.
type skv[K comparable, T any] struct {
	key   K
	value T
}

// NewSMapKeyValue returns a new SMapKeyValue container.
func NewSMapKeyValue[K comparable, T any]() *SMapKeyValue[K, T] {
	return &SMapKeyValue[K, T]{
		data: sync.Map{},
	}
}

// Set sets the value associated with the key.
func (r *SMapKeyValue[K, T]) Set(key K, value T) {
	r.count.Add(1)
	r.data.Store(key, value)
}

// GetAndCheck returns the value associated with the key if this exist also a
// boolean value if this exist of not.
func (r *SMapKeyValue[K, T]) GetAndCheck(key K) (T, bool) {
	value, ok := r.data.Load(key)

	switch value := value.(type) {
	case T:
		return value, ok
	default:
		var t T
		return t, ok
	}
}

// Get returns the value associated with the key.
// If the key does not exist, return zero value of the type.
func (r *SMapKeyValue[K, T]) Get(key K) T {
	value, _ := r.data.Load(key)

	switch value := value.(type) {
	case T:
		return value
	default:
		var t T
		return t
	}
}

// GetAnDelete returns the value associated with the key and delete it if the key exist
// if the key doesn't exist return the given key value false
func (r *SMapKeyValue[K, T]) GetAnDelete(key K) (T, bool) {
	value, ok := r.data.LoadAndDelete(key)
	if ok {
		r.count.Swap(r.count.Load() - 1)

	}

	switch value := value.(type) {
	case T:
		return value, ok
	default:
		var t T
		return t, ok
	}
}

// Delete deletes the value associated with the key.
func (r *SMapKeyValue[K, T]) Delete(key K) {
	if _, ok := r.data.LoadAndDelete(key); ok {
		r.count.Swap(r.count.Load() - 1)
	}
}

// Clear deletes all key-value pairs stored in the container.
func (r *SMapKeyValue[K, T]) Clear() {
	r.data = sync.Map{}
	r.count.Swap(0)
}

// Size returns the number of key-value pairs stored in the container.
func (r *SMapKeyValue[K, T]) Size() int {
	return int(r.count.Load())
}

// IsEmpty returns true if the container is empty.
func (r *SMapKeyValue[K, T]) IsEmpty() bool {
	return r.Size() == 0
}

// IsFull returns true if the container has elements.
func (r *SMapKeyValue[K, T]) IsFull() bool {
	return r.Size() != 0
}

// ContainsKey returns true if the key is in the container.
func (r *SMapKeyValue[K, T]) ContainsKey(key K) bool {
	_, ok := r.data.Load(key)
	return ok
}

// ContainsValue returns true if the value is in the container.
func (r *SMapKeyValue[K, T]) ContainsValue(value T) bool {
	var ret bool

	r.data.Range(func(key, v any) bool {
		if reflect.DeepEqual(v, value) {
			ret = true
			return false
		}
		return true
	})
	return ret
}

// Get returns the key value associated with the key.
func (r *SMapKeyValue[K, T]) Key(key K) K {
	if _, ok := r.data.Load(key); ok {
		return key
	}
	var empty K
	return empty
}

// Keys returns all keys stored in the container.
func (r *SMapKeyValue[K, T]) Keys() []K {
	keys := make([]K, 0, r.Size())
	r.data.Range(func(key, value any) bool {
		keys = append(keys, key.(K))
		return true
	})
	return keys
}

// Values returns all values stored in the container.
func (r *SMapKeyValue[K, T]) Values() []T {
	values := make([]T, 0, r.Size())
	r.data.Range(func(key, value any) bool {
		values = append(values, value.(T))
		return true
	})
	return values
}

// ForEach calls the given function for each key-value pair in the container.
func (r *SMapKeyValue[K, T]) ForEach(fn func(key K, value T)) {
	r.data.Range(func(key, value any) bool {
		fn(key.(K), value.(T))
		return true
	})
}

// ForEachKey calls the given function for each key in the container.
func (r *SMapKeyValue[K, T]) ForEachKey(fn func(key K)) {
	r.data.Range(func(key, value any) bool {
		fn(key.(K))
		return true
	})
}

// ForEachValue calls the given function for each value in the container.
func (r *SMapKeyValue[K, T]) ForEachValue(fn func(value T)) {
	r.data.Range(func(key, value any) bool {
		fn(value.(T))
		return true
	})
}

// Clone returns a new SMapKeyValue with a copy of the underlying data.
func (r *SMapKeyValue[K, T]) Clone() *SMapKeyValue[K, T] {
	clone := NewSMapKeyValue[K, T]()

	r.data.Range(func(key, value any) bool {
		clone.Set(key.(K), value.(T))
		return true
	})

	return clone
}

// CloneAndClear returns a new SMapKeyValue with a copy of the underlying data and clears the container.
func (r *SMapKeyValue[K, T]) CloneAndClear() *SMapKeyValue[K, T] {
	clone := NewSMapKeyValue[K, T]()
	r.data.Range(func(key, value any) bool {
		clone.Set(key.(K), value.(T))
		return true
	})
	r.Clear()
	return clone
}

// DeepEqual returns true if the given kv is deep equal to the SMapKeyValue container
func (r *SMapKeyValue[K, T]) DeepEqual(kv *SMapKeyValue[K, T]) bool {
	if r.Size() != kv.Size() {
		return false
	}
	if (r.Size() == kv.Size()) && r.Size() == 0 {
		return true
	}

	var ret bool
	r.data.Range(func(key, value any) bool {
		kk, ok := kv.GetAndCheck(key.(K))
		if !ok {
			ret = false
			return false
		} else {
			ret = reflect.DeepEqual(kk, value.(T))
			return true
		}
	})

	return ret
}

// Map returns a new SMapKeyValue after applying the given function fn to each key-value pair.
func (r *SMapKeyValue[K, T]) Map(fn func(key K, value T) (newKey K, newValue T)) *SMapKeyValue[K, T] {
	m := NewSMapKeyValue[K, T]()
	r.data.Range(func(key, value any) bool {
		newKey, newValue := fn(key.(K), value.(T))
		m.Set(newKey, newValue)
		return true
	})

	return m
}

// MapKey returns a new SMapKeyValue after applying the given function fn to each key.
func (r *SMapKeyValue[K, T]) MapKey(fn func(key K) K) *SMapKeyValue[K, T] {
	m := NewSMapKeyValue[K, T]()
	r.data.Range(func(key, value any) bool {
		newKey := fn(key.(K))
		m.Set(newKey, value.(T))
		return true
	})
	return m
}

// MapValue returns a new SMapKeyValue after applying the given function fn to each value.
func (r *SMapKeyValue[K, T]) MapValue(fn func(value T) T) *SMapKeyValue[K, T] {
	m := NewSMapKeyValue[K, T]()
	r.data.Range(func(key, value any) bool {
		newValue := fn(value.(T))
		m.Set(key.(K), newValue)
		return true
	})
	return m
}

// Filter returns a new SMapKeyValue after applying the given function fn to each key-value pair.
func (r *SMapKeyValue[K, T]) Filter(fn func(key K, value T) bool) *SMapKeyValue[K, T] {
	m := NewSMapKeyValue[K, T]()
	r.data.Range(func(key, value any) bool {
		if fn(key.(K), value.(T)) {
			m.Set(key.(K), value.(T))
		}
		return true
	})
	return m
}

// FilterKey returns a new SMapKeyValue after applying the given function fn to each key.
func (r *SMapKeyValue[K, T]) FilterKey(fn func(key K) bool) *SMapKeyValue[K, T] {
	m := NewSMapKeyValue[K, T]()
	r.data.Range(func(key, value any) bool {
		if fn(key.(K)) {
			m.Set(key.(K), value.(T))
		}
		return true
	})
	return m
}

// FilterValue returns a new SMapKeyValue after applying the given function fn to each value.
func (r *SMapKeyValue[K, T]) FilterValue(fn func(value T) bool) *SMapKeyValue[K, T] {
	m := NewSMapKeyValue[K, T]()
	r.data.Range(func(key, value any) bool {
		if fn(value.(T)) {
			m.Set(key.(K), value.(T))
		}
		return true
	})
	return m
}

// Partition returns two new SMapKeyValue. One with all the elements that satisfy the predicate and
// another with the rest. The predicate is applied to each element.
func (r *SMapKeyValue[K, T]) Partition(fn func(key K, value T) bool) (match, others *SMapKeyValue[K, T]) {
	match = NewSMapKeyValue[K, T]()
	others = NewSMapKeyValue[K, T]()
	r.data.Range(func(key, value any) bool {
		if fn(key.(K), value.(T)) {
			match.Set(key.(K), value.(T))
		} else {
			others.Set(key.(K), value.(T))
		}
		return true
	})

	return
}

// PartitionKey returns two new SMapKeyValue. One with all the elements that satisfy the predicate and
// another with the rest. The predicate is applied to each key.
func (r *SMapKeyValue[K, T]) PartitionKey(fn func(key K) bool) (match, others *SMapKeyValue[K, T]) {
	match = NewSMapKeyValue[K, T]()
	others = NewSMapKeyValue[K, T]()
	r.data.Range(func(key, value any) bool {
		if fn(key.(K)) {
			match.Set(key.(K), value.(T))
		} else {
			others.Set(key.(K), value.(T))
		}
		return true
	})

	return
}

// PartitionValue returns two new SMapKeyValue. One with all the elements that satisfy the predicate and
// another with the rest. The predicate is applied to each value.
func (r *SMapKeyValue[K, T]) PartitionValue(fn func(value T) bool) (match, others *SMapKeyValue[K, T]) {
	match = NewSMapKeyValue[K, T]()
	others = NewSMapKeyValue[K, T]()
	r.data.Range(func(key, value any) bool {
		if fn(value.(T)) {
			match.Set(key.(K), value.(T))
		} else {
			others.Set(key.(K), value.(T))
		}
		return true
	})
	return
}

// SortKeys returns a []*K (keys) after sorting the keys using the given sortFn function.
func (r *SMapKeyValue[K, T]) SortKeys(sortFn func(key1, key2 K) bool) []*K {
	keys := r.Keys()

	sort.Slice(keys, func(i, j int) bool {
		return sortFn(keys[i], keys[j])
	})

	m := make([]*K, len(keys))
	for i, key := range keys {
		k := key
		m[i] = &k
	}
	return m
}

// SortValues returns a []*T (values) after sorting the values using given function sortFn.
func (r *SMapKeyValue[K, T]) SortValues(sortFn func(value1, value2 T) bool) []*T {
	kvs := make([]*skv[K, T], 0, r.Size())
	r.data.Range(func(key, value any) bool {
		kvs = append(kvs, &skv[K, T]{key.(K), value.(T)})
		return true
	})

	sort.Slice(kvs, func(i, j int) bool {
		return sortFn(kvs[i].value, kvs[j].value)
	})

	m := make([]*T, len(kvs))
	for i, pair := range kvs {
		m[i] = &pair.value
	}

	return m
}
