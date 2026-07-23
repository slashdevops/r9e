package r9e

import (
	"encoding/json"
	"iter"
	"maps"
	"reflect"
	"slices"
	"sync"
)

type mapKeyValueOptions struct {
	size int
}

// MapKeyValueOptions configures a [MapKeyValue] container at construction time.
type MapKeyValueOptions func(*mapKeyValueOptions)

// WithCapacity presizes the underlying map to hold at least size entries
// without reallocating. It is a performance hint: when the approximate number
// of entries is known up front, presizing avoids incremental map growth.
func WithCapacity(size int) MapKeyValueOptions {
	return func(o *mapKeyValueOptions) {
		o.size = size
	}
}

// MapKeyValue is a thread-safe, generic key-value container backed by a native
// Go map guarded by a [sync.RWMutex]. Reads are served concurrently under a
// read lock; writes take the exclusive lock.
//
// Prefer MapKeyValue when the workload is read-heavy or mixed and callers want
// predictable, snapshot-consistent bulk operations (Keys, Values, Clone, Map,
// Filter, Partition). For workloads dominated by disjoint keys written from
// many goroutines, see [SMapKeyValue].
//
// The zero value is not ready for use; construct one with [NewMapKeyValue].
type MapKeyValue[K comparable, T any] struct {
	mu   sync.RWMutex
	data map[K]T
}

// NewMapKeyValue returns a ready-to-use MapKeyValue. Pass [WithCapacity] to
// presize the container.
func NewMapKeyValue[K comparable, T any](options ...MapKeyValueOptions) *MapKeyValue[K, T] {
	var o mapKeyValueOptions
	for _, opt := range options {
		opt(&o)
	}

	return &MapKeyValue[K, T]{
		data: make(map[K]T, o.size),
	}
}

// snapshot returns a shallow copy of the underlying data taken under the read
// lock. It is the building block for operations that must not hold the lock
// while touching another container (avoiding lock-ordering deadlocks).
func (r *MapKeyValue[K, T]) snapshot() map[K]T {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make(map[K]T, len(r.data))
	maps.Copy(out, r.data)
	return out
}

// Set stores value under key, replacing any existing value.
func (r *MapKeyValue[K, T]) Set(key K, value T) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[key] = value
}

// Get returns the value stored under key, or the zero value of T if the key is
// absent. Use [MapKeyValue.GetAndCheck] to distinguish an absent key from a
// stored zero value.
func (r *MapKeyValue[K, T]) Get(key K) T {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.data[key]
}

// GetAndCheck returns the value stored under key and a boolean reporting whether
// the key was present.
func (r *MapKeyValue[K, T]) GetAndCheck(key K) (T, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	value, ok := r.data[key]
	return value, ok
}

// GetOrSet returns the existing value for key if present. Otherwise it stores
// and returns value. The loaded result is true if the value was already
// present. The lookup and store are performed atomically under a single lock.
func (r *MapKeyValue[K, T]) GetOrSet(key K, value T) (actual T, loaded bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if existing, ok := r.data[key]; ok {
		return existing, true
	}
	r.data[key] = value
	return value, false
}

// GetAndDelete returns the value stored under key and deletes it. The loaded
// result reports whether the key was present.
func (r *MapKeyValue[K, T]) GetAndDelete(key K) (value T, loaded bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	value, loaded = r.data[key]
	if loaded {
		delete(r.data, key)
	}
	return value, loaded
}

// Delete removes key from the container. Deleting an absent key is a no-op.
func (r *MapKeyValue[K, T]) Delete(key K) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.data, key)
}

// Clear removes all entries, retaining the allocated capacity for reuse.
func (r *MapKeyValue[K, T]) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	clear(r.data)
}

// Size returns the number of entries stored.
func (r *MapKeyValue[K, T]) Size() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.data)
}

// IsEmpty reports whether the container has no entries.
func (r *MapKeyValue[K, T]) IsEmpty() bool {
	return r.Size() == 0
}

// ContainsKey reports whether key is present.
func (r *MapKeyValue[K, T]) ContainsKey(key K) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, ok := r.data[key]
	return ok
}

// ContainsValue reports whether any stored value is deeply equal to value,
// using [reflect.DeepEqual]. This is O(n) in the number of entries.
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

// Keys returns a snapshot slice of all keys. The order is unspecified.
func (r *MapKeyValue[K, T]) Keys() []K {
	r.mu.RLock()
	defer r.mu.RUnlock()

	keys := make([]K, 0, len(r.data))
	for key := range r.data {
		keys = append(keys, key)
	}
	return keys
}

// Values returns a snapshot slice of all values. The order is unspecified.
func (r *MapKeyValue[K, T]) Values() []T {
	r.mu.RLock()
	defer r.mu.RUnlock()

	values := make([]T, 0, len(r.data))
	for _, value := range r.data {
		values = append(values, value)
	}
	return values
}

// All returns an iterator over all key-value pairs, suitable for use with a
// range-over-func loop:
//
//	for k, v := range kv.All() {
//		// ...
//	}
//
// The read lock is held for the duration of the iteration, so the callback must
// not call methods that mutate the same container (Set, Delete, Clear, ...);
// doing so deadlocks. Break out of the loop early to stop iterating.
func (r *MapKeyValue[K, T]) All() iter.Seq2[K, T] {
	return func(yield func(K, T) bool) {
		r.mu.RLock()
		defer r.mu.RUnlock()

		for key, value := range r.data {
			if !yield(key, value) {
				return
			}
		}
	}
}

// ForEach calls fn for every key-value pair. The read lock is held for the
// duration; fn must not mutate the same container. Prefer [MapKeyValue.All]
// with a range-over-func loop in new code.
func (r *MapKeyValue[K, T]) ForEach(fn func(key K, value T)) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for key, value := range r.data {
		fn(key, value)
	}
}

// ForEachKey calls fn for every key. See [MapKeyValue.ForEach] for locking
// semantics.
func (r *MapKeyValue[K, T]) ForEachKey(fn func(key K)) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for key := range r.data {
		fn(key)
	}
}

// ForEachValue calls fn for every value. See [MapKeyValue.ForEach] for locking
// semantics.
func (r *MapKeyValue[K, T]) ForEachValue(fn func(value T)) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, value := range r.data {
		fn(value)
	}
}

// Clone returns a new independent container holding a shallow copy of the data.
func (r *MapKeyValue[K, T]) Clone() *MapKeyValue[K, T] {
	return &MapKeyValue[K, T]{data: r.snapshot()}
}

// CloneAndClear atomically copies the data into a new container and clears the
// receiver.
func (r *MapKeyValue[K, T]) CloneAndClear() *MapKeyValue[K, T] {
	r.mu.Lock()
	defer r.mu.Unlock()

	out := make(map[K]T, len(r.data))
	maps.Copy(out, r.data)
	clear(r.data)
	return &MapKeyValue[K, T]{data: out}
}

// Merge copies every entry from other into the receiver, overwriting existing
// keys. other is read via a consistent snapshot, so the two containers are
// never locked simultaneously. Merging a container into itself is a no-op-safe
// operation. A nil other is ignored.
func (r *MapKeyValue[K, T]) Merge(other *MapKeyValue[K, T]) {
	if other == nil || other == r {
		return
	}
	src := other.snapshot()

	r.mu.Lock()
	defer r.mu.Unlock()
	maps.Copy(r.data, src)
}

// DeepEqual reports whether the receiver and other hold the same keys mapped to
// deeply equal values ([reflect.DeepEqual]). A nil other equals an empty
// receiver only when the receiver is also empty. The two containers are never
// locked at the same time.
func (r *MapKeyValue[K, T]) DeepEqual(other *MapKeyValue[K, T]) bool {
	var otherData map[K]T
	if other != nil {
		otherData = other.snapshot()
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.data) != len(otherData) {
		return false
	}
	for key, value := range r.data {
		ov, ok := otherData[key]
		if !ok || !reflect.DeepEqual(value, ov) {
			return false
		}
	}
	return true
}

// Map returns a new container produced by applying fn to every pair.
func (r *MapKeyValue[K, T]) Map(fn func(key K, value T) (K, T)) *MapKeyValue[K, T] {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make(map[K]T, len(r.data))
	for key, value := range r.data {
		nk, nv := fn(key, value)
		out[nk] = nv
	}
	return &MapKeyValue[K, T]{data: out}
}

// MapKey returns a new container with each key transformed by fn.
func (r *MapKeyValue[K, T]) MapKey(fn func(key K) K) *MapKeyValue[K, T] {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make(map[K]T, len(r.data))
	for key, value := range r.data {
		out[fn(key)] = value
	}
	return &MapKeyValue[K, T]{data: out}
}

// MapValue returns a new container with each value transformed by fn.
func (r *MapKeyValue[K, T]) MapValue(fn func(value T) T) *MapKeyValue[K, T] {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make(map[K]T, len(r.data))
	for key, value := range r.data {
		out[key] = fn(value)
	}
	return &MapKeyValue[K, T]{data: out}
}

// Filter returns a new container with the pairs for which fn reports true.
func (r *MapKeyValue[K, T]) Filter(fn func(key K, value T) bool) *MapKeyValue[K, T] {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make(map[K]T)
	for key, value := range r.data {
		if fn(key, value) {
			out[key] = value
		}
	}
	return &MapKeyValue[K, T]{data: out}
}

// FilterKey returns a new container with the pairs whose key satisfies fn.
func (r *MapKeyValue[K, T]) FilterKey(fn func(key K) bool) *MapKeyValue[K, T] {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make(map[K]T)
	for key, value := range r.data {
		if fn(key) {
			out[key] = value
		}
	}
	return &MapKeyValue[K, T]{data: out}
}

// FilterValue returns a new container with the pairs whose value satisfies fn.
func (r *MapKeyValue[K, T]) FilterValue(fn func(value T) bool) *MapKeyValue[K, T] {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make(map[K]T)
	for key, value := range r.data {
		if fn(value) {
			out[key] = value
		}
	}
	return &MapKeyValue[K, T]{data: out}
}

// Partition splits the container into match (pairs for which fn is true) and
// others (the rest), returning two new containers.
func (r *MapKeyValue[K, T]) Partition(fn func(key K, value T) bool) (match, others *MapKeyValue[K, T]) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	m := make(map[K]T)
	o := make(map[K]T)
	for key, value := range r.data {
		if fn(key, value) {
			m[key] = value
		} else {
			o[key] = value
		}
	}
	return &MapKeyValue[K, T]{data: m}, &MapKeyValue[K, T]{data: o}
}

// PartitionKey splits the container by applying fn to each key.
func (r *MapKeyValue[K, T]) PartitionKey(fn func(key K) bool) (match, others *MapKeyValue[K, T]) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	m := make(map[K]T)
	o := make(map[K]T)
	for key, value := range r.data {
		if fn(key) {
			m[key] = value
		} else {
			o[key] = value
		}
	}
	return &MapKeyValue[K, T]{data: m}, &MapKeyValue[K, T]{data: o}
}

// PartitionValue splits the container by applying fn to each value.
func (r *MapKeyValue[K, T]) PartitionValue(fn func(value T) bool) (match, others *MapKeyValue[K, T]) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	m := make(map[K]T)
	o := make(map[K]T)
	for key, value := range r.data {
		if fn(value) {
			m[key] = value
		} else {
			o[key] = value
		}
	}
	return &MapKeyValue[K, T]{data: m}, &MapKeyValue[K, T]{data: o}
}

// SortKeys returns all keys sorted by the less function, which must report
// whether a should sort before b.
func (r *MapKeyValue[K, T]) SortKeys(less func(a, b K) bool) []K {
	keys := r.Keys()
	slices.SortFunc(keys, func(a, b K) int {
		switch {
		case less(a, b):
			return -1
		case less(b, a):
			return 1
		default:
			return 0
		}
	})
	return keys
}

// SortValues returns all values sorted by the less function, which must report
// whether a should sort before b.
func (r *MapKeyValue[K, T]) SortValues(less func(a, b T) bool) []T {
	values := r.Values()
	slices.SortFunc(values, func(a, b T) int {
		switch {
		case less(a, b):
			return -1
		case less(b, a):
			return 1
		default:
			return 0
		}
	})
	return values
}

// MarshalJSON encodes the container as a JSON object, so a MapKeyValue can be
// used directly as a struct field. Encoding succeeds only for key types that
// encoding/json accepts as object keys (strings, integers, and
// encoding.TextMarshaler implementations).
func (r *MapKeyValue[K, T]) MarshalJSON() ([]byte, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make(map[K]T, len(r.data))
	maps.Copy(out, r.data)
	return json.Marshal(out)
}

// UnmarshalJSON decodes a JSON object into the container, merging the decoded
// entries over any existing ones.
func (r *MapKeyValue[K, T]) UnmarshalJSON(data []byte) error {
	var m map[K]T
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if r.data == nil {
		r.data = make(map[K]T, len(m))
	}
	maps.Copy(r.data, m)
	return nil
}
