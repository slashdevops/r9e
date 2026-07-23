# ❓ FAQ

Common questions, gotchas, and thread-safety notes. If something here is
unclear, the [pkg.go.dev reference](https://pkg.go.dev/github.com/slashdevops/r9e)
has the authoritative per-method documentation.

## Which container should I use?

Default to **`MapKeyValue`**. It is a native map behind a `sync.RWMutex`:
simple, fast for read-heavy and mixed workloads, and it gives consistent
snapshots for bulk operations. Reach for **`SMapKeyValue`** only when your
access pattern matches what `sync.Map` is optimized for — many goroutines
writing **disjoint** keys, or keys written **once and read many** times. When
unsure, [benchmark both](performance.md). See [Containers](containers.md) for
the full decision tree.

## Is it safe for concurrent use?

Yes. Both containers are safe for use by multiple goroutines with **no external
locking**. Every method acquires whatever synchronization it needs internally.
See [Concurrency](concurrency.md).

## Can I iterate and modify at the same time?

**Not on `MapKeyValue`.** `All()` and `ForEach*` hold the read lock for the
whole loop, so calling a mutating method (`Set`, `Delete`, `Clear`, `Merge`,
`GetOrSet`, `GetAndDelete`, `CloneAndClear`) on the *same* container from inside
the loop **deadlocks**. Collect the changes and apply them after the loop, or
build a new container with `Filter`/`Map`.

**On `SMapKeyValue`** mutation during iteration does not deadlock, but the scan
reflects a *moment-in-time* view (`sync.Map.Range` semantics): concurrent writes
may or may not be visited. Don't rely on either outcome. See
[Iteration](iteration.md).

## What key and value types are allowed?

- **Keys** — any `comparable` type: strings, all integer/float kinds, `bool`,
  pointers, and comparable structs/arrays.
- **Values** — literally `any` type, including structs, slices, maps, pointers,
  and interfaces.

For **JSON** there is one extra restriction: `encoding/json` only supports map
keys that are strings, integers, or `encoding.TextMarshaler` implementations. A
container with, say, a struct key stores and retrieves fine but cannot be
JSON-encoded. See [JSON](json.md).

## Does it persist to disk?

**No.** r9e is purely in-memory — the name stands for *RamStorage*. Data lives
for the lifetime of the process and disappears when it exits. To persist, encode
the container (e.g. with [JSON](json.md)) and write the bytes wherever you like;
r9e itself never touches the filesystem, the network, or any external store.

## Are the slices and containers returned by methods safe to use?

Yes — they are independent copies you fully own:

- `Keys()`, `Values()`, `SortKeys`, `SortValues` return **fresh slices**.
  Mutating them never affects the source container.
- `Clone`, `CloneAndClear`, and the `Map`/`Filter`/`Partition` family return
  **new independent containers**. The receiver is unchanged (except
  `CloneAndClear`, which clears it).

One caveat: copies are **shallow**. If `T` is a pointer or contains reference
types (slices, maps, pointers), the copy shares the pointed-to data with the
original. Mutating that shared data is visible through both.

```go
type Bag struct{ Items []string }

kv := r9e.NewMapKeyValue[string, *Bag]()
kv.Set("a", &Bag{Items: []string{"x"}})

clone := kv.Clone()
clone.Get("a").Items[0] = "MUTATED" // both see it — same *Bag
```

Store values instead of pointers, or deep-copy inside a `MapValue`, if you need
full isolation.

## Why is `Get` returning zero for a key I never set?

`Get` returns the **zero value** of `T` for absent keys, which is
indistinguishable from a stored zero. Use `GetAndCheck` to tell them apart:

```go
v, ok := kv.GetAndCheck("maybe")
if !ok {
    // truly absent
}
```

## Why is iteration order different every run?

Both containers are unordered, and Go randomizes map iteration order on purpose.
For stable output, sort first with `SortKeys` / `SortValues`. See
[Iteration](iteration.md#-deterministic-iteration).

## How do I count entries cheaply?

`Size()` is **O(1)** on both containers — a `len()` under the read lock for
`MapKeyValue`, and an atomic load for `SMapKeyValue`. `IsEmpty()` is just
`Size() == 0`. Avoid `ContainsValue`, which is O(n).

## Does r9e have any third-party dependencies?

No. It uses only the Go standard library (`sync`, `sync/atomic`, `iter`, `maps`,
`slices`, `reflect`, `encoding/json`) and requires **Go 1.26+**.

## ➡️ Next

- Back to the [documentation index](README.md).
- Deep-dive the internals in [Containers](containers.md) and
  [Concurrency](concurrency.md).
