package r9e

import (
	"encoding/json"
	"maps"
	"testing"
)

func TestGetOrSet_MapKeyValue(t *testing.T) {
	kv := NewMapKeyValue[string, int]()

	if v, loaded := kv.GetOrSet("a", 1); v != 1 || loaded {
		t.Fatalf("first GetOrSet = (%d, %v), want (1, false)", v, loaded)
	}
	if v, loaded := kv.GetOrSet("a", 99); v != 1 || !loaded {
		t.Fatalf("second GetOrSet = (%d, %v), want (1, true)", v, loaded)
	}
	if kv.Size() != 1 {
		t.Fatalf("Size = %d, want 1", kv.Size())
	}
}

func TestGetOrSet_SMapKeyValue(t *testing.T) {
	kv := NewSMapKeyValue[string, int]()

	if v, loaded := kv.GetOrSet("a", 1); v != 1 || loaded {
		t.Fatalf("first GetOrSet = (%d, %v), want (1, false)", v, loaded)
	}
	if v, loaded := kv.GetOrSet("a", 99); v != 1 || !loaded {
		t.Fatalf("second GetOrSet = (%d, %v), want (1, true)", v, loaded)
	}
	if kv.Size() != 1 {
		t.Fatalf("Size = %d, want 1", kv.Size())
	}
}

func TestSetOverwriteSize_SMapKeyValue(t *testing.T) {
	kv := NewSMapKeyValue[string, int]()

	kv.Set("a", 1)
	kv.Set("a", 2)
	kv.Set("a", 3)

	if kv.Size() != 1 {
		t.Fatalf("Size after overwriting same key = %d, want 1", kv.Size())
	}
	if got := kv.Get("a"); got != 3 {
		t.Fatalf("Get = %d, want 3", got)
	}

	kv.Delete("a")
	if kv.Size() != 0 {
		t.Fatalf("Size after delete = %d, want 0", kv.Size())
	}
}

func TestMerge_MapKeyValue(t *testing.T) {
	a := NewMapKeyValue[string, int]()
	a.Set("x", 1)
	a.Set("y", 2)

	b := NewMapKeyValue[string, int]()
	b.Set("y", 20)
	b.Set("z", 30)

	a.Merge(b)

	want := map[string]int{"x": 1, "y": 20, "z": 30}
	if got := maps.Collect(a.All()); !maps.Equal(got, want) {
		t.Fatalf("after Merge = %v, want %v", got, want)
	}

	// nil and self merges are no-ops.
	a.Merge(nil)
	a.Merge(a)
	if a.Size() != 3 {
		t.Fatalf("Size after nil/self merge = %d, want 3", a.Size())
	}
}

func TestMerge_SMapKeyValue(t *testing.T) {
	a := NewSMapKeyValue[string, int]()
	a.Set("x", 1)

	b := NewSMapKeyValue[string, int]()
	b.Set("y", 2)

	a.Merge(b)
	if a.Size() != 2 {
		t.Fatalf("Size = %d, want 2", a.Size())
	}
	a.Merge(nil)
	a.Merge(a)
	if a.Size() != 2 {
		t.Fatalf("Size after nil/self merge = %d, want 2", a.Size())
	}
}

func TestAll_MapKeyValue(t *testing.T) {
	kv := NewMapKeyValue[string, int]()
	kv.Set("a", 1)
	kv.Set("b", 2)
	kv.Set("c", 3)

	got := maps.Collect(kv.All())
	want := map[string]int{"a": 1, "b": 2, "c": 3}
	if !maps.Equal(got, want) {
		t.Fatalf("All collected = %v, want %v", got, want)
	}

	// Early break must stop iteration without panicking.
	count := 0
	for range kv.All() {
		count++
		break
	}
	if count != 1 {
		t.Fatalf("early-break visited %d entries, want 1", count)
	}
}

func TestAll_SMapKeyValue(t *testing.T) {
	kv := NewSMapKeyValue[string, int]()
	kv.Set("a", 1)
	kv.Set("b", 2)

	got := maps.Collect(kv.All())
	want := map[string]int{"a": 1, "b": 2}
	if !maps.Equal(got, want) {
		t.Fatalf("All collected = %v, want %v", got, want)
	}
}

func TestJSON_MapKeyValue(t *testing.T) {
	kv := NewMapKeyValue[string, int]()
	kv.Set("a", 1)
	kv.Set("b", 2)

	encoded, err := json.Marshal(kv)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	if string(encoded) != `{"a":1,"b":2}` {
		t.Fatalf("Marshal = %s, want {\"a\":1,\"b\":2}", encoded)
	}

	decoded := NewMapKeyValue[string, int]()
	if err := json.Unmarshal([]byte(`{"x":10,"y":20}`), decoded); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if decoded.Get("x") != 10 || decoded.Get("y") != 20 {
		t.Fatalf("decoded = %v, want x=10 y=20", maps.Collect(decoded.All()))
	}

	// Empty container marshals to an empty object, not null.
	empty := NewMapKeyValue[string, int]()
	encoded, _ = json.Marshal(empty)
	if string(encoded) != `{}` {
		t.Fatalf("empty Marshal = %s, want {}", encoded)
	}
}

func TestJSON_SMapKeyValue(t *testing.T) {
	kv := NewSMapKeyValue[string, int]()
	kv.Set("a", 1)

	encoded, err := json.Marshal(kv)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	if string(encoded) != `{"a":1}` {
		t.Fatalf("Marshal = %s, want {\"a\":1}", encoded)
	}

	decoded := NewSMapKeyValue[string, int]()
	if err := json.Unmarshal([]byte(`{"x":10}`), decoded); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if decoded.Get("x") != 10 || decoded.Size() != 1 {
		t.Fatalf("decoded x=%d size=%d, want x=10 size=1", decoded.Get("x"), decoded.Size())
	}
}

func TestDeepEqualNil_MapKeyValue(t *testing.T) {
	empty := NewMapKeyValue[string, int]()
	if !empty.DeepEqual(nil) {
		t.Fatal("empty.DeepEqual(nil) = false, want true")
	}

	nonEmpty := NewMapKeyValue[string, int]()
	nonEmpty.Set("a", 1)
	if nonEmpty.DeepEqual(nil) {
		t.Fatal("nonEmpty.DeepEqual(nil) = true, want false")
	}
}

func TestDeepEqualNil_SMapKeyValue(t *testing.T) {
	empty := NewSMapKeyValue[string, int]()
	if !empty.DeepEqual(nil) {
		t.Fatal("empty.DeepEqual(nil) = false, want true")
	}

	nonEmpty := NewSMapKeyValue[string, int]()
	nonEmpty.Set("a", 1)
	if nonEmpty.DeepEqual(nil) {
		t.Fatal("nonEmpty.DeepEqual(nil) = true, want false")
	}
}
