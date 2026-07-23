package r9e_test

import (
	"encoding/json"
	"fmt"
	"maps"

	"github.com/slashdevops/r9e"
)

// Store and read back a value.
func ExampleMapKeyValue() {
	kv := r9e.NewMapKeyValue[string, int]()

	kv.Set("answer", 42)

	if v, ok := kv.GetAndCheck("answer"); ok {
		fmt.Println(v)
	}
	// Output: 42
}

// GetOrSet returns the existing value or stores and returns a new one.
func ExampleMapKeyValue_GetOrSet() {
	kv := r9e.NewMapKeyValue[string, int]()

	v1, loaded1 := kv.GetOrSet("hits", 1)
	fmt.Printf("first:  value=%d loaded=%v\n", v1, loaded1)

	v2, loaded2 := kv.GetOrSet("hits", 999)
	fmt.Printf("second: value=%d loaded=%v\n", v2, loaded2)
	// Output:
	// first:  value=1 loaded=false
	// second: value=1 loaded=true
}

// GetAndDelete removes a key and returns its former value.
func ExampleMapKeyValue_GetAndDelete() {
	kv := r9e.NewMapKeyValue[string, string]()
	kv.Set("token", "abc123")

	value, loaded := kv.GetAndDelete("token")
	fmt.Printf("value=%q loaded=%v size=%d\n", value, loaded, kv.Size())
	// Output: value="abc123" loaded=true size=0
}

// Iterate deterministically by sorting the keys first.
func ExampleMapKeyValue_All() {
	kv := r9e.NewMapKeyValue[string, int]()
	kv.Set("c", 3)
	kv.Set("a", 1)
	kv.Set("b", 2)

	// All() yields in map order; maps.Collect drains the iterator into a plain
	// map that we can look up while iterating the sorted keys for stable output.
	snapshot := maps.Collect(kv.All())

	for _, k := range kv.SortKeys(func(a, b string) bool { return a < b }) {
		fmt.Printf("%s=%d\n", k, snapshot[k])
	}
	// Output:
	// a=1
	// b=2
	// c=3
}

// FilterValue returns a new container with the matching entries.
func ExampleMapKeyValue_FilterValue() {
	kv := r9e.NewMapKeyValue[string, int]()
	kv.Set("one", 1)
	kv.Set("two", 2)
	kv.Set("three", 3)
	kv.Set("four", 4)

	even := kv.FilterValue(func(v int) bool { return v%2 == 0 })

	fmt.Println(even.SortValues(func(a, b int) bool { return a < b }))
	// Output: [2 4]
}

// Partition splits a container into matching and non-matching halves.
func ExampleMapKeyValue_Partition() {
	kv := r9e.NewMapKeyValue[string, int]()
	kv.Set("a", 10)
	kv.Set("b", 20)
	kv.Set("c", 30)

	big, small := kv.Partition(func(_ string, v int) bool { return v >= 20 })

	fmt.Println("big:  ", big.SortValues(func(a, b int) bool { return a < b }))
	fmt.Println("small:", small.SortValues(func(a, b int) bool { return a < b }))
	// Output:
	// big:   [20 30]
	// small: [10]
}

// SortValues returns values ordered by a custom comparison.
func ExampleMapKeyValue_SortValues() {
	kv := r9e.NewMapKeyValue[string, float64]()
	kv.Set("pi", 3.14)
	kv.Set("e", 2.71)
	kv.Set("phi", 1.61)

	fmt.Println(kv.SortValues(func(a, b float64) bool { return a < b }))
	// Output: [1.61 2.71 3.14]
}

// Merge copies entries from another container, overwriting existing keys.
func ExampleMapKeyValue_Merge() {
	a := r9e.NewMapKeyValue[string, int]()
	a.Set("x", 1)
	a.Set("y", 2)

	b := r9e.NewMapKeyValue[string, int]()
	b.Set("y", 20)
	b.Set("z", 30)

	a.Merge(b)

	for _, k := range a.SortKeys(func(a, b string) bool { return a < b }) {
		fmt.Printf("%s=%d\n", k, a.Get(k))
	}
	// Output:
	// x=1
	// y=20
	// z=30
}

// A MapKeyValue marshals to and from a JSON object.
func ExampleMapKeyValue_json() {
	kv := r9e.NewMapKeyValue[string, int]()
	kv.Set("a", 1)
	kv.Set("b", 2)

	encoded, _ := json.Marshal(kv)
	fmt.Println(string(encoded))

	decoded := r9e.NewMapKeyValue[string, int]()
	_ = json.Unmarshal([]byte(`{"x":10,"y":20}`), decoded)
	fmt.Println(decoded.Get("x"), decoded.Get("y"))
	// Output:
	// {"a":1,"b":2}
	// 10 20
}

// SMapKeyValue has the same API, backed by sync.Map.
func ExampleSMapKeyValue() {
	sm := r9e.NewSMapKeyValue[string, int]()

	sm.Set("a", 1)
	sm.Set("a", 2) // overwrite: Size stays 1
	sm.Set("b", 3)

	fmt.Println("size:", sm.Size())
	fmt.Println("a:", sm.Get("a"))
	// Output:
	// size: 2
	// a: 2
}

// Using a struct value type.
func ExampleMapKeyValue_struct() {
	type Constant struct {
		Name  string
		Value float64
	}

	kv := r9e.NewMapKeyValue[string, Constant]()
	kv.Set("pi", Constant{"Archimedes' constant", 3.141592})

	c := kv.Get("pi")
	fmt.Printf("%s = %v\n", c.Name, c.Value)
	// Output: Archimedes' constant = 3.141592
}
