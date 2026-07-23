# 🗄️ JSON

Both containers implement `json.Marshaler` and `json.Unmarshaler`, so they
encode to and decode from a plain JSON **object** — exactly like the underlying
`map[K]T` would.

## 📤 Marshalling

```go
kv := r9e.NewMapKeyValue[string, int]()
kv.Set("a", 1)
kv.Set("b", 2)

encoded, _ := json.Marshal(kv)
fmt.Println(string(encoded)) // {"a":1,"b":2}
```

`MarshalJSON` takes a consistent snapshot (under the read lock for
`MapKeyValue`) and encodes it, so it is safe to call while other goroutines
read the container.

```mermaid
flowchart LR
    KV["MapKeyValue / SMapKeyValue"] -->|"MarshalJSON"| MAP["map[K]T snapshot"]
    MAP -->|"json.Marshal"| JSON["{\"a\":1,\"b\":2}"]
    JSON -->|"UnmarshalJSON"| KV2["container<br/>(entries merged in)"]
```

## 📥 Unmarshalling

`UnmarshalJSON` decodes the JSON object and **merges** the decoded entries over
any existing ones (it does not clear the container first). Construct the
container before unmarshalling into it.

```go
decoded := r9e.NewMapKeyValue[string, int]()
_ = json.Unmarshal([]byte(`{"x":10,"y":20}`), decoded)

fmt.Println(decoded.Get("x"), decoded.Get("y")) // 10 20
```

## 🧩 As a struct field

Because the container marshals like a map, you can embed it directly in a struct
and it serializes as a nested object:

```go
type Server struct {
    Name   string                          `json:"name"`
    Labels *r9e.MapKeyValue[string, string] `json:"labels"`
}

s := Server{
    Name:   "web-1",
    Labels: r9e.NewMapKeyValue[string, string](),
}
s.Labels.Set("env", "prod")
s.Labels.Set("tier", "frontend")

out, _ := json.Marshal(s)
// {"name":"web-1","labels":{"env":"prod","tier":"frontend"}}
```

When decoding into such a struct, make sure the field is non-nil first (e.g.
initialize it, or the surrounding decode allocates it) so `UnmarshalJSON` has a
container to merge into.

## 🔑 Key-type constraints

JSON object keys are always strings, so `encoding/json` only accepts key types
it can turn into (and back from) a string. Marshalling **succeeds only** when
`K` is one of:

| Key type `K` | JSON key example | Works? |
| ------------ | ---------------- | ------ |
| `string` | `"env"` | ✅ |
| integer kinds (`int`, `int64`, `uint`, ...) | `"42"` | ✅ (quoted) |
| type implementing `encoding.TextMarshaler` | its `MarshalText` output | ✅ |
| `bool`, `float`, struct, pointer, ... | — | ❌ marshalling errors |

This is a constraint of `encoding/json`, not of r9e — the same rules apply to a
raw `map[K]T`. Value type `T` may be anything JSON can encode.

```go
// int keys marshal as quoted strings, per encoding/json
ports := r9e.NewMapKeyValue[int, string]()
ports.Set(80, "http")
ports.Set(443, "https")

b, _ := json.Marshal(ports)
fmt.Println(string(b)) // {"443":"https","80":"http"}
```

## ⚠️ Error handling

Do not ignore the error from `json.Marshal`/`json.Unmarshal` in production code
— an unsupported key type, or malformed input, is reported there:

```go
b, err := json.Marshal(kv)
if err != nil {
    return fmt.Errorf("encode store: %w", err)
}
```

## ➡️ Next

- Tune throughput and memory in [Performance](performance.md).
- Browse common questions in the [FAQ](faq.md).
