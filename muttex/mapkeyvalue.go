package muttex

import (
	"sync"
)

type mapKeyValueOptions struct {
	size int
}

type MapKeyValueOptions func(*mapKeyValueOptions)

func WithSize(size int) MapKeyValueOptions {
	return func(kv *mapKeyValueOptions) {
		kv.size = size
	}
}

type MapKeyValue[K comparable, T any] struct {
	mu   sync.RWMutex
	data map[K]T
}

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

func (r *MapKeyValue[K, T]) Get(key K) (T, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	value, ok := r.data[key]
	return value, ok
}

func (r *MapKeyValue[K, T]) Set(key K, value T) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[key] = value
}

func (r *MapKeyValue[K, T]) Delete(key K) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.data, key)
}

func (r *MapKeyValue[K, T]) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data = make(map[K]T)
}

func (r *MapKeyValue[K, T]) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.data)
}

func (r *MapKeyValue[K, T]) Keys() []K {
	r.mu.RLock()
	defer r.mu.RUnlock()

	keys := make([]K, 0, len(r.data))
	for key := range r.data {
		keys = append(keys, key)
	}
	return keys
}

func (r *MapKeyValue[K, T]) Values() []T {
	r.mu.RLock()
	defer r.mu.RUnlock()

	values := make([]T, 0, len(r.data))
	for _, value := range r.data {
		values = append(values, value)
	}
	return values
}

func (r *MapKeyValue[K, T]) Each(fn func(key K, value T)) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for key, value := range r.data {
		fn(key, value)
	}
}

func (r *MapKeyValue[K, T]) EachKey(fn func(key K)) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for key := range r.data {
		fn(key)
	}
}

func (r *MapKeyValue[K, T]) EachValue(fn func(value T)) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, value := range r.data {
		fn(value)
	}
}

func (r *MapKeyValue[K, T]) Clone() *MapKeyValue[K, T] {
	r.mu.RLock()
	defer r.mu.RUnlock()

	clone := NewMapKeyValue[K, T](WithSize(r.Len()))
	for key, value := range r.data {
		clone.Set(key, value)
	}
	return clone
}

func (r *MapKeyValue[K, T]) CloneAndClear() *MapKeyValue[K, T] {
	r.mu.RLock()
	defer r.mu.RUnlock()

	clone := NewMapKeyValue[K, T](WithSize(r.Len()))
	for key, value := range r.data {
		clone.Set(key, value)
	}
	r.Clear()
	return clone
}

func (r *MapKeyValue[K, T]) Map(fn func(key K, value T) (newKey K, newValue T)) *MapKeyValue[K, T] {
	r.mu.Lock()
	defer r.mu.Unlock()

	m := NewMapKeyValue[K, T](WithSize(r.Len()))
	for key, value := range r.data {
		newKey, newValue := fn(key, value)
		m.Set(newKey, newValue)
	}
	return m
}

func (r *MapKeyValue[K, T]) MapKey(fn func(key K) K) *MapKeyValue[K, T] {
	r.mu.Lock()
	defer r.mu.Unlock()

	m := NewMapKeyValue[K, T](WithSize(r.Len()))
	for key := range r.data {
		newKey := fn(key)
		m.Set(newKey, r.data[key])
	}
	return m
}

func (r *MapKeyValue[K, T]) MapValue(fn func(value T) T) *MapKeyValue[K, T] {
	r.mu.Lock()
	defer r.mu.Unlock()

	m := NewMapKeyValue[K, T](WithSize(r.Len()))
	for key, value := range r.data {
		newValue := fn(value)
		m.Set(key, newValue)
	}
	return m
}

func (r *MapKeyValue[K, T]) Filter(fn func(key K, value T) bool) *MapKeyValue[K, T] {
	r.mu.Lock()
	defer r.mu.Unlock()

	m := NewMapKeyValue[K, T](WithSize(r.Len()))
	for key, value := range r.data {
		if fn(key, value) {
			m.Set(key, value)
		}
	}
	return m
}

func (r *MapKeyValue[K, T]) FilterKey(fn func(key K) bool) *MapKeyValue[K, T] {
	r.mu.Lock()
	defer r.mu.Unlock()

	m := NewMapKeyValue[K, T](WithSize(r.Len()))
	for key := range r.data {
		if fn(key) {
			m.Set(key, r.data[key])
		}
	}
	return m
}

func (r *MapKeyValue[K, T]) FilterValue(fn func(value T) bool) *MapKeyValue[K, T] {
	r.mu.Lock()
	defer r.mu.Unlock()

	m := NewMapKeyValue[K, T](WithSize(r.Len()))
	for key, value := range r.data {
		if fn(value) {
			m.Set(key, value)
		}
	}
	return m
}

func (r *MapKeyValue[K, T]) Partition(fn func(key K, value T) bool) (match, others *MapKeyValue[K, T]) {
	r.mu.Lock()
	defer r.mu.Unlock()

	match = NewMapKeyValue[K, T](WithSize(r.Len()))
	others = NewMapKeyValue[K, T](WithSize(r.Len()))
	for key, value := range r.data {
		if fn(key, value) {
			match.Set(key, value)
		} else {
			others.Set(key, value)
		}
	}
	return
}

func (r *MapKeyValue[K, T]) PartitionKey(fn func(key K) bool) (match, others *MapKeyValue[K, T]) {
	r.mu.Lock()
	defer r.mu.Unlock()

	match = NewMapKeyValue[K, T](WithSize(r.Len()))
	others = NewMapKeyValue[K, T](WithSize(r.Len()))
	for key := range r.data {
		if fn(key) {
			match.Set(key, r.data[key])
		} else {
			others.Set(key, r.data[key])
		}
	}
	return
}

func (r *MapKeyValue[K, T]) PartitionValue(fn func(value T) bool) (match, others *MapKeyValue[K, T]) {
	r.mu.Lock()
	defer r.mu.Unlock()

	match = NewMapKeyValue[K, T](WithSize(r.Len()))
	others = NewMapKeyValue[K, T](WithSize(r.Len()))
	for key, value := range r.data {
		if fn(value) {
			match.Set(key, value)
		} else {
			others.Set(key, value)
		}
	}
	return
}
