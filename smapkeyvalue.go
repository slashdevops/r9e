package r9e

import (
	"encoding/json"
	"iter"
	"reflect"
	"slices"
	"sync"
	"sync/atomic"
)

// SMapKeyValue is a thread-safe, generic key-value container backed by a
// [sync.Map]. It keeps an atomic entry counter so [SMapKeyValue.Size] is O(1).
//
// Prefer SMapKeyValue for workloads where many goroutines write disjoint keys,
// or where a key is written once and read many times, matching the cases
// sync.Map is optimized for. For read-heavy or mixed workloads that also need
// consistent bulk snapshots, [MapKeyValue] is usually simpler and faster.
//
// The zero value is not ready for use; construct one with [NewSMapKeyValue].
type SMapKeyValue[K comparable, T any] struct {
	count atomic.Int64
	data  sync.Map
}

// NewSMapKeyValue returns a ready-to-use SMapKeyValue.
func NewSMapKeyValue[K comparable, T any]() *SMapKeyValue[K, T] {
	return &SMapKeyValue[K, T]{}
}

// cast converts a value loaded from the underlying sync.Map back to T,
// returning the zero value of T when the stored value is absent or of an
// unexpected type.
func cast[T any](v any) T {
	if t, ok := v.(T); ok {
		return t
	}
	var zero T
	return zero
}

// Set stores value under key, replacing any existing value.
func (r *SMapKeyValue[K, T]) Set(key K, value T) {
	if _, loaded := r.data.Swap(key, value); !loaded {
		r.count.Add(1)
	}
}

// Get returns the value stored under key, or the zero value of T if the key is
// absent.
func (r *SMapKeyValue[K, T]) Get(key K) T {
	value, _ := r.data.Load(key)
	return cast[T](value)
}

// GetAndCheck returns the value stored under key and a boolean reporting whether
// the key was present.
func (r *SMapKeyValue[K, T]) GetAndCheck(key K) (T, bool) {
	value, ok := r.data.Load(key)
	return cast[T](value), ok
}

// GetOrSet returns the existing value for key if present. Otherwise it stores
// and returns value. The loaded result reports whether the value was already
// present. The operation is atomic.
func (r *SMapKeyValue[K, T]) GetOrSet(key K, value T) (actual T, loaded bool) {
	v, loaded := r.data.LoadOrStore(key, value)
	if !loaded {
		r.count.Add(1)
	}
	return cast[T](v), loaded
}

// GetAndDelete returns the value stored under key and deletes it. The loaded
// result reports whether the key was present.
func (r *SMapKeyValue[K, T]) GetAndDelete(key K) (value T, loaded bool) {
	v, loaded := r.data.LoadAndDelete(key)
	if loaded {
		r.count.Add(-1)
	}
	return cast[T](v), loaded
}

// Delete removes key from the container. Deleting an absent key is a no-op.
func (r *SMapKeyValue[K, T]) Delete(key K) {
	if _, loaded := r.data.LoadAndDelete(key); loaded {
		r.count.Add(-1)
	}
}

// Clear removes all entries.
func (r *SMapKeyValue[K, T]) Clear() {
	r.data.Clear()
	r.count.Store(0)
}

// Size returns the number of entries stored. It is O(1).
func (r *SMapKeyValue[K, T]) Size() int {
	return int(r.count.Load())
}

// IsEmpty reports whether the container has no entries.
func (r *SMapKeyValue[K, T]) IsEmpty() bool {
	return r.Size() == 0
}

// ContainsKey reports whether key is present.
func (r *SMapKeyValue[K, T]) ContainsKey(key K) bool {
	_, ok := r.data.Load(key)
	return ok
}

// ContainsValue reports whether any stored value is deeply equal to value,
// using [reflect.DeepEqual]. This is O(n) in the number of entries.
func (r *SMapKeyValue[K, T]) ContainsValue(value T) bool {
	found := false
	r.data.Range(func(_, v any) bool {
		if reflect.DeepEqual(v, value) {
			found = true
			return false
		}
		return true
	})
	return found
}

// Keys returns a snapshot slice of all keys. The order is unspecified.
func (r *SMapKeyValue[K, T]) Keys() []K {
	keys := make([]K, 0, r.Size())
	r.data.Range(func(key, _ any) bool {
		keys = append(keys, key.(K))
		return true
	})
	return keys
}

// Values returns a snapshot slice of all values. The order is unspecified.
func (r *SMapKeyValue[K, T]) Values() []T {
	values := make([]T, 0, r.Size())
	r.data.Range(func(_, value any) bool {
		values = append(values, value.(T))
		return true
	})
	return values
}

// All returns an iterator over all key-value pairs, suitable for use with a
// range-over-func loop:
//
//	for k, v := range sm.All() {
//		// ...
//	}
//
// Iteration reflects a moment-in-time view of the map; concurrent writes may or
// may not be observed. Break out of the loop early to stop iterating.
func (r *SMapKeyValue[K, T]) All() iter.Seq2[K, T] {
	return func(yield func(K, T) bool) {
		r.data.Range(func(key, value any) bool {
			return yield(key.(K), value.(T))
		})
	}
}

// ForEach calls fn for every key-value pair. Prefer [SMapKeyValue.All] with a
// range-over-func loop in new code.
func (r *SMapKeyValue[K, T]) ForEach(fn func(key K, value T)) {
	r.data.Range(func(key, value any) bool {
		fn(key.(K), value.(T))
		return true
	})
}

// ForEachKey calls fn for every key.
func (r *SMapKeyValue[K, T]) ForEachKey(fn func(key K)) {
	r.data.Range(func(key, _ any) bool {
		fn(key.(K))
		return true
	})
}

// ForEachValue calls fn for every value.
func (r *SMapKeyValue[K, T]) ForEachValue(fn func(value T)) {
	r.data.Range(func(_, value any) bool {
		fn(value.(T))
		return true
	})
}

// Clone returns a new independent container holding a copy of the data.
func (r *SMapKeyValue[K, T]) Clone() *SMapKeyValue[K, T] {
	clone := NewSMapKeyValue[K, T]()
	r.data.Range(func(key, value any) bool {
		clone.Set(key.(K), value.(T))
		return true
	})
	return clone
}

// CloneAndClear copies the data into a new container and clears the receiver.
// The two operations are not a single atomic step: concurrent writes that land
// between the copy and the clear are observed by neither container reliably.
func (r *SMapKeyValue[K, T]) CloneAndClear() *SMapKeyValue[K, T] {
	clone := r.Clone()
	r.Clear()
	return clone
}

// Merge copies every entry from other into the receiver, overwriting existing
// keys. A nil other, or merging a container into itself, is ignored.
func (r *SMapKeyValue[K, T]) Merge(other *SMapKeyValue[K, T]) {
	if other == nil || other == r {
		return
	}
	other.data.Range(func(key, value any) bool {
		r.Set(key.(K), value.(T))
		return true
	})
}

// DeepEqual reports whether the receiver and other hold the same keys mapped to
// deeply equal values ([reflect.DeepEqual]). A nil other equals the receiver
// only when the receiver is empty.
func (r *SMapKeyValue[K, T]) DeepEqual(other *SMapKeyValue[K, T]) bool {
	otherSize := 0
	if other != nil {
		otherSize = other.Size()
	}
	if r.Size() != otherSize {
		return false
	}

	equal := true
	r.data.Range(func(key, value any) bool {
		ov, ok := other.GetAndCheck(key.(K))
		if !ok || !reflect.DeepEqual(value.(T), ov) {
			equal = false
			return false
		}
		return true
	})
	return equal
}

// Map returns a new container produced by applying fn to every pair.
func (r *SMapKeyValue[K, T]) Map(fn func(key K, value T) (K, T)) *SMapKeyValue[K, T] {
	out := NewSMapKeyValue[K, T]()
	r.data.Range(func(key, value any) bool {
		nk, nv := fn(key.(K), value.(T))
		out.Set(nk, nv)
		return true
	})
	return out
}

// MapKey returns a new container with each key transformed by fn.
func (r *SMapKeyValue[K, T]) MapKey(fn func(key K) K) *SMapKeyValue[K, T] {
	out := NewSMapKeyValue[K, T]()
	r.data.Range(func(key, value any) bool {
		out.Set(fn(key.(K)), value.(T))
		return true
	})
	return out
}

// MapValue returns a new container with each value transformed by fn.
func (r *SMapKeyValue[K, T]) MapValue(fn func(value T) T) *SMapKeyValue[K, T] {
	out := NewSMapKeyValue[K, T]()
	r.data.Range(func(key, value any) bool {
		out.Set(key.(K), fn(value.(T)))
		return true
	})
	return out
}

// Filter returns a new container with the pairs for which fn reports true.
func (r *SMapKeyValue[K, T]) Filter(fn func(key K, value T) bool) *SMapKeyValue[K, T] {
	out := NewSMapKeyValue[K, T]()
	r.data.Range(func(key, value any) bool {
		if fn(key.(K), value.(T)) {
			out.Set(key.(K), value.(T))
		}
		return true
	})
	return out
}

// FilterKey returns a new container with the pairs whose key satisfies fn.
func (r *SMapKeyValue[K, T]) FilterKey(fn func(key K) bool) *SMapKeyValue[K, T] {
	out := NewSMapKeyValue[K, T]()
	r.data.Range(func(key, value any) bool {
		if fn(key.(K)) {
			out.Set(key.(K), value.(T))
		}
		return true
	})
	return out
}

// FilterValue returns a new container with the pairs whose value satisfies fn.
func (r *SMapKeyValue[K, T]) FilterValue(fn func(value T) bool) *SMapKeyValue[K, T] {
	out := NewSMapKeyValue[K, T]()
	r.data.Range(func(key, value any) bool {
		if fn(value.(T)) {
			out.Set(key.(K), value.(T))
		}
		return true
	})
	return out
}

// Partition splits the container into match (pairs for which fn is true) and
// others (the rest), returning two new containers.
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
	return match, others
}

// PartitionKey splits the container by applying fn to each key.
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
	return match, others
}

// PartitionValue splits the container by applying fn to each value.
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
	return match, others
}

// SortKeys returns all keys sorted by the less function, which must report
// whether a should sort before b.
func (r *SMapKeyValue[K, T]) SortKeys(less func(a, b K) bool) []K {
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
func (r *SMapKeyValue[K, T]) SortValues(less func(a, b T) bool) []T {
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

// MarshalJSON encodes the container as a JSON object. Encoding succeeds only for
// key types that encoding/json accepts as object keys (strings, integers, and
// encoding.TextMarshaler implementations).
func (r *SMapKeyValue[K, T]) MarshalJSON() ([]byte, error) {
	out := make(map[K]T, r.Size())
	r.data.Range(func(key, value any) bool {
		out[key.(K)] = value.(T)
		return true
	})
	return json.Marshal(out)
}

// UnmarshalJSON decodes a JSON object into the container, merging the decoded
// entries over any existing ones.
func (r *SMapKeyValue[K, T]) UnmarshalJSON(data []byte) error {
	var m map[K]T
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	for key, value := range m {
		r.Set(key, value)
	}
	return nil
}
