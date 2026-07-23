// Package r9e (RamStorage) provides thread-safe, generic in-memory key-value
// containers built only on the Go standard library.
//
// r9e takes advantage of Go generics and the standard library's own concurrency
// primitives to offer a small, dependency-free API for storing and retrieving
// data from memory. The focus is usability and simplicity without giving up
// performance.
//
// # Containers
//
// The package exposes two containers with the same method set, so code can be
// written against either and switched based on the workload:
//
//   - [MapKeyValue] is backed by a native Go map guarded by a [sync.RWMutex].
//     Reads run concurrently under a read lock; writes take the exclusive lock.
//     It is the best default for read-heavy or mixed workloads and for
//     operations that need a consistent snapshot (Keys, Values, Clone, Map,
//     Filter, Partition).
//
//   - [SMapKeyValue] is backed by a [sync.Map] with an atomic size counter, so
//     Size is O(1). It suits workloads where many goroutines write disjoint
//     keys, or where each key is written once and read many times.
//
// # Quick Start
//
//	kv := r9e.NewMapKeyValue[string, int]()
//	kv.Set("answer", 42)
//
//	if v, ok := kv.GetAndCheck("answer"); ok {
//		fmt.Println(v) // 42
//	}
//
//	for k, v := range kv.All() {
//		fmt.Printf("%s=%d\n", k, v)
//	}
//
// # Iteration
//
// Both containers implement range-over-func iteration via [MapKeyValue.All] and
// [SMapKeyValue.All], which return an [iter.Seq2] over key-value pairs:
//
//	for k, v := range kv.All() {
//		// ...
//	}
//
// For MapKeyValue the read lock is held for the duration of the loop, so the
// body must not call methods that mutate the same container. Break out of the
// loop early to stop iterating.
//
// # Functional Helpers
//
// Map, MapKey, MapValue, Filter, FilterKey, FilterValue, Partition,
// PartitionKey, and PartitionValue each return new, independent containers and
// never mutate the receiver. Clone, CloneAndClear, and Merge provide snapshot
// and combination semantics. SortKeys and SortValues return sorted slices.
//
// # Choosing a Container
//
// Prefer [MapKeyValue] unless profiling shows that [sync.Map]'s access pattern
// (write-once/read-many, or disjoint keys across many goroutines) fits your
// workload better. When in doubt, benchmark both with your own key and value
// types.
//
// # Dependencies
//
// r9e has zero third-party dependencies. It uses only the Go standard library
// (sync, sync/atomic, iter, maps, slices, reflect, encoding/json).
package r9e
