package r9e

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

const kvSize = 8192

type TestStruct struct {
	a string
	b int
	c int64
	d float64
}

var (
	kv_int_int       = NewMapKeyValue[int, int](WithCapacity(kvSize))
	kv_string_string = NewMapKeyValue[string, string](WithCapacity(kvSize))
	kv_string_struct = NewMapKeyValue[string, TestStruct](WithCapacity(kvSize))
)

func init() {
	rand.Seed(time.Now().UnixNano())

	// fill the kv_int_int
	for i := 0; i < kvSize; i++ {
		kv_int_int.Set(rand.Intn(kvSize), rand.Intn(kvSize))
	}

	// fill the kv_string_string
	for i := 0; i < kvSize; i++ {
		keyval := fmt.Sprintf("%x", md5.Sum([]byte(strconv.Itoa(rand.Intn(kvSize)))))
		kv_string_string.Set(keyval, keyval)
	}

	// fill the kv_string_struct
	for i := 0; i < kvSize; i++ {
		keyval := fmt.Sprintf("%x", md5.Sum([]byte(strconv.Itoa(rand.Intn(kvSize)))))
		s := TestStruct{
			a: keyval,
			b: rand.Intn(kvSize),
			c: int64(rand.Intn(kvSize)),
			d: rand.Float64(),
		}
		kv_string_struct.Set(keyval, s)
	}
}

// **************************************************
// ******************** Tests ***********************
func TestNewMapKeyValue(t *testing.T) {
	t.Run("test NewMapKeyValue[int, int] with capacity", func(t *testing.T) {
		kv := NewMapKeyValue[int, int](WithCapacity(kvSize))

		if kv.Size() != 0 {
			t.Errorf("Expected size to be %v, got %v", kvSize, kv.Size())
		}

		kv.Set(1, 8096)

		if kv.Size() != 1 {
			t.Errorf("Expected size to be %v, got %v", 1, kv.Size())
		}

		value := kv.Get(1)
		VKind := reflect.TypeOf(value).Kind().String()

		if VKind != "int" {
			t.Errorf("Expected type to be %s, got %s", "int", VKind)
		}

		key := kv.Keys()[0]
		kKind := reflect.TypeOf(key).Kind().String()

		if kKind != "int" {
			t.Errorf("Expected type to be %s, got %s", "int", kKind)
		}
	})

	t.Run("test NewMapKeyValue[int, int] without capacity", func(t *testing.T) {
		kv := NewMapKeyValue[int, int]()

		if kv.Size() != 0 {
			t.Errorf("Expected size to be %v, got %v", kvSize, kv.Size())
		}

		kv.Set(1, 8096)

		if kv.Size() != 1 {
			t.Errorf("Expected size to be %v, got %v", 1, kv.Size())
		}

		value := kv.Get(1)
		VKind := reflect.TypeOf(value).Kind().String()

		if VKind != "int" {
			t.Errorf("Expected type to be %s, got %s", "int", VKind)
		}

		key := kv.Keys()[0]
		kKind := reflect.TypeOf(key).Kind().String()

		if kKind != "int" {
			t.Errorf("Expected type to be %s, got %s", "int", kKind)
		}
	})

	t.Run("test NewMapKeyValue[float64, string] without capacity", func(t *testing.T) {
		kv := NewMapKeyValue[float64, string]()

		if kv.Size() != 0 {
			t.Errorf("Expected size to be %v, got %v", kvSize, kv.Size())
		}

		kv.Set(1, "test string")

		if kv.Size() != 1 {
			t.Errorf("Expected size to be %v, got %v", 1, kv.Size())
		}

		value := kv.Get(1)
		VKind := reflect.TypeOf(value).Kind().String()

		if VKind != "string" {
			t.Errorf("Expected type to be %s, got %s", "string", VKind)
		}

		key := kv.Keys()[0]
		kKind := reflect.TypeOf(key).Kind().String()

		if kKind != "float64" {
			t.Errorf("Expected type to be %s, got %s", "float64", kKind)
		}
	})

	t.Run("test NewMapKeyValue[int, custom struct] without capacity", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[int, testStruct]()

		if kv.Size() != 0 {
			t.Errorf("Expected size to be %v, got %v", kvSize, kv.Size())
		}

		kv.Set(1, testStruct{"This is Archimedes' Constant (Pi)", 3.1415})

		if kv.Size() != 1 {
			t.Errorf("Expected size to be %v, got %v", 1, kv.Size())
		}

		value := kv.Get(1)
		typeOf := reflect.TypeOf(value)
		kind := typeOf.Kind().String()

		if kind != "struct" {
			t.Errorf("Expected type to be %s, got %s", "struct", kind)
		}

		if typeOf.Name() != "testStruct" {
			t.Errorf("Expected type to be %s, got %s", "testStruct", kind)
		}

		key := kv.Keys()[0]
		kKind := reflect.TypeOf(key).Kind().String()

		if kKind != "int" {
			t.Errorf("Expected type to be %s, got %s", "int", kKind)
		}
	})
}

func TestSet(t *testing.T) {
	t.Run("test Set for NewMapKeyValue[string, struct] key exist", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		kv.Set("Archimedes", testStruct{"This is Archimedes' Constant (Pi)", 3.1415})

		if kv.Size() != 1 {
			t.Errorf("Expected size to be %v, got %v", 1, kv.Size())
		}

		value := kv.Get("Archimedes")

		if value.Name != "This is Archimedes' Constant (Pi)" {
			t.Errorf("Expected value to be %s, got %s", "This is Archimedes' Constant (Pi)", value.Name)
		}
		if value.value != 3.1415 {
			t.Errorf("Expected value to be %v, got %v", 3.1415, value.value)
		}
	})
}

func TestGetCheck(t *testing.T) {
	t.Run("test GetCheck for NewMapKeyValue[string, struct] key exist", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		kv.Set("Archimedes", testStruct{"This is Archimedes' Constant (Pi)", 3.1415})

		if kv.Size() != 1 {
			t.Errorf("Expected size to be %v, got %v", 1, kv.Size())
		}

		value, ok := kv.GetCheck("Archimedes")
		if !ok {
			t.Errorf("Expected GetCheck to return true, got %v", ok)
		}

		if value.Name != "This is Archimedes' Constant (Pi)" {
			t.Errorf("Expected value to be %s, got %s", "This is Archimedes' Constant (Pi)", value.Name)
		}
		if value.value != 3.1415 {
			t.Errorf("Expected value to be %v, got %v", 3.1415, value.value)
		}
	})

	t.Run("test GetCheck for NewMapKeyValue[string, struct] key doesn't exist", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		kv.Set("Archimedes", testStruct{"This is Archimedes' Constant (Pi)", 3.1415})

		if kv.Size() != 1 {
			t.Errorf("Expected size to be %v, got %v", 1, kv.Size())
		}

		_, ok := kv.GetCheck("Euler")
		if ok {
			t.Errorf("Expected GetCheck to return true, got %v", ok)
		}
	})
}

func TestGet(t *testing.T) {
	t.Run("test Get for NewMapKeyValue[string, struct] key exist", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		kv.Set("Archimedes", testStruct{"This is Archimedes' Constant (Pi)", 3.1415})

		if kv.Size() != 1 {
			t.Errorf("Expected size to be %v, got %v", 1, kv.Size())
		}

		value := kv.Get("Archimedes")

		if value.Name != "This is Archimedes' Constant (Pi)" {
			t.Errorf("Expected value to be %s, got %s", "This is Archimedes' Constant (Pi)", value.Name)
		}
		if value.value != 3.1415 {
			t.Errorf("Expected value to be %v, got %v", 3.1415, value.value)
		}
	})

	t.Run("test Get for NewMapKeyValue[string, struct] key doesn't exist", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		kv.Set("Archimedes", testStruct{"This is Archimedes' Constant (Pi)", 3.1415})

		if kv.Size() != 1 {
			t.Errorf("Expected size to be %v, got %v", 1, kv.Size())
		}

		value := kv.Get("Euler")
		if value.value != 0 {
			t.Errorf("Expected value to be %v, got %v", nil, value)
		}
	})
}

func TestDelete(t *testing.T) {
	t.Run("test Delete for NewMapKeyValue[string, struct] with keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		kv.Set("Archimedes", testStruct{"This is Archimedes' Constant (Pi)", 3.1415})
		kv.Set("Euler", testStruct{"This is Euler's Number (e)", 2.7182})
		kv.Set("Golden Ratio", testStruct{"This is The Golden Ratio", 1.6180})

		if kv.Size() != 3 {
			t.Errorf("Expected size to be %v, got %v", 3, kv.Size())
		}

		value := kv.Get("Archimedes")

		if value.Name != "This is Archimedes' Constant (Pi)" {
			t.Errorf("Expected value to be %s, got %s", "This is Archimedes' Constant (Pi)", value.Name)
		}
		if value.value != 3.1415 {
			t.Errorf("Expected value to be %v, got %v", 3.1415, value.value)
		}

		kv.Delete("Archimedes")

		if kv.Size() != 2 {
			t.Errorf("Expected size to be %v, got %v", 2, kv.Size())
		}

		value = kv.Get("Golden Ratio")

		if value.Name != "This is The Golden Ratio" {
			t.Errorf("Expected value to be %s, got %s", "This is The Golden Ratio", value.Name)
		}
		if value.value != 1.6180 {
			t.Errorf("Expected value to be %v, got %v", 1.6180, value.value)
		}
	})

	t.Run("test Delete for NewMapKeyValue[string, struct] without keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		if kv.Size() != 0 {
			t.Errorf("Expected size to be %v, got %v", 0, kv.Size())
		}

		kv.Delete("Archimedes")

		if kv.Size() != 0 {
			t.Errorf("Expected size to be %v, got %v", 0, kv.Size())
		}
	})
}

func TestClear(t *testing.T) {
	t.Run("test Clear for NewMapKeyValue[string, struct] with keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		kv.Set("Archimedes", testStruct{"This is Archimedes' Constant (Pi)", 3.1415})
		kv.Set("Euler", testStruct{"This is Euler's Number (e)", 2.7182})
		kv.Set("Golden Ratio", testStruct{"This is The Golden Ratio", 1.6180})

		if kv.Size() != 3 {
			t.Errorf("Expected size to be %v, got %v", 1, kv.Size())
		}

		value := kv.Get("Archimedes")

		if value.Name != "This is Archimedes' Constant (Pi)" {
			t.Errorf("Expected value to be %s, got %s", "This is Archimedes' Constant (Pi)", value.Name)
		}
		if value.value != 3.1415 {
			t.Errorf("Expected value to be %v, got %v", 3.1415, value.value)
		}

		kv.Clear()

		if kv.Size() != 0 {
			t.Errorf("Expected size to be %v, got %v", 0, kv.Size())
		}
	})

	t.Run("test Clear for NewMapKeyValue[string, struct] without keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		if kv.Size() != 0 {
			t.Errorf("Expected size to be %v, got %v", 0, kv.Size())
		}

		kv.Clear()

		if kv.Size() != 0 {
			t.Errorf("Expected size to be %v, got %v", 0, kv.Size())
		}
	})
}

func TestSize(t *testing.T) {
	t.Run("test Size for NewMapKeyValue[string, struct] with keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		kv.Set("Archimedes", testStruct{"This is Archimedes' Constant (Pi)", 3.1415})
		kv.Set("Euler", testStruct{"This is Euler's Number (e)", 2.7182})
		kv.Set("Golden Ratio", testStruct{"This is The Golden Ratio", 1.6180})

		if kv.Size() != 3 {
			t.Errorf("Expected size to be %v, got %v", 3, kv.Size())
		}

		kv.Delete("Archimedes")

		if kv.Size() != 2 {
			t.Errorf("Expected size to be %v, got %v", 2, kv.Size())
		}

		kv.Delete("Euler")

		if kv.Size() != 1 {
			t.Errorf("Expected size to be %v, got %v", 1, kv.Size())
		}

		kv.Clear()

		if kv.Size() != 0 {
			t.Errorf("Expected size to be %v, got %v", 0, kv.Size())
		}
	})

	t.Run("test Size for NewMapKeyValue[string, struct] without keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		if kv.Size() != 0 {
			t.Errorf("Expected size to be %v, got %v", 0, kv.Size())
		}

		kv.Delete("Euler")

		if kv.Size() != 0 {
			t.Errorf("Expected size to be %v, got %v", 0, kv.Size())
		}
	})
}

func TestIsEmpty(t *testing.T) {
	t.Run("test IsEmpty for NewMapKeyValue[string, struct] with keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		kv.Set("Archimedes", testStruct{"This is Archimedes' Constant (Pi)", 3.1415})
		kv.Set("Euler", testStruct{"This is Euler's Number (e)", 2.7182})
		kv.Set("Golden Ratio", testStruct{"This is The Golden Ratio", 1.6180})

		if kv.Size() != 3 {
			t.Errorf("Expected size to be %v, got %v", 3, kv.Size())
		}

		if kv.IsEmpty() != false {
			t.Errorf("Expected IsEmpty to be %v, got %v", false, kv.IsEmpty())
		}

		kv.Clear()

		if kv.IsEmpty() != true {
			t.Errorf("Expected IsEmpty to be %v, got %v", true, kv.IsEmpty())
		}
	})
}

func TestIsFull(t *testing.T) {
	t.Run("test IsFull for NewMapKeyValue[string, struct] with keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		kv.Set("Archimedes", testStruct{"This is Archimedes' Constant (Pi)", 3.1415})
		kv.Set("Euler", testStruct{"This is Euler's Number (e)", 2.7182})
		kv.Set("Golden Ratio", testStruct{"This is The Golden Ratio", 1.6180})

		if kv.Size() != 3 {
			t.Errorf("Expected size to be %v, got %v", 3, kv.Size())
		}

		if kv.IsFull() != true {
			t.Errorf("Expected IsFull to be %v, got %v", false, kv.IsFull())
		}

		kv.Clear()

		if kv.IsFull() != false {
			t.Errorf("Expected IsFull to be %v, got %v", true, kv.IsFull())
		}
	})
}

func TestContainsKey(t *testing.T) {
	t.Run("test ContainsKey for NewMapKeyValue[string, struct] with keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		kv.Set("Archimedes", testStruct{"This is Archimedes' Constant (Pi)", 3.1415})

		if kv.Size() != 1 {
			t.Errorf("Expected size to be %v, got %v", 1, kv.Size())
		}

		if kv.ContainsKey("Archimedes") != true {
			t.Errorf("Expected key to be %v, got %v", true, kv.ContainsKey("Archimedes"))
		}

		if kv.ContainsKey("Do not exist") != false {
			t.Errorf("Expected key to be %v, got %v", false, kv.ContainsKey("Do not exist"))
		}
	})

	t.Run("test ContainsKey for NewMapKeyValue[string, struct] without keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		if kv.ContainsKey("Do not exist") != false {
			t.Errorf("Expected key to be %v, got %v", false, kv.ContainsKey("Do not exist"))
		}
	})
}

func TestContainsValue(t *testing.T) {
	t.Run("test ContainsValue for NewMapKeyValue[string, struct] with keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		kv.Set("Archimedes", testStruct{"This is Archimedes' Constant (Pi)", 3.1415})

		if kv.Size() != 1 {
			t.Errorf("Expected size to be %v, got %v", 1, kv.Size())
		}

		if kv.ContainsValue(testStruct{"This is Archimedes' Constant (Pi)", 3.1415}) != true {
			t.Errorf("Expected key to be %v, got %v", true, kv.ContainsValue(testStruct{"This is Archimedes' Constant (Pi)", 3.1415}))
		}

		if kv.ContainsValue(testStruct{"This is other constant", 0.00000}) != false {
			t.Errorf("Expected key to be %v, got %v", false, kv.ContainsValue(testStruct{"This is other constant", 0.00000}))
		}
	})

	t.Run("test ContainsValue for NewMapKeyValue[string, struct] without keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		if kv.ContainsValue(testStruct{"This is other constant", 0.00000}) != false {
			t.Errorf("Expected key to be %v, got %v", false, kv.ContainsValue(testStruct{"This is other constant", 0.00000}))
		}
	})
}

func TestKey(t *testing.T) {
	t.Run("test Key for NewMapKeyValue[string, struct] with keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		kv.Set("Archimedes", testStruct{"This is Archimedes' Constant (Pi)", 3.1415})
		kv.Set("Euler", testStruct{"This is Euler's Number (e)", 2.7182})
		kv.Set("Golden Ratio", testStruct{"This is The Golden Ratio", 1.6180})

		if kv.Size() != 3 {
			t.Errorf("Expected size to be %v, got %v", 3, kv.Size())
		}

		if kv.Key("Archimedes") != "Archimedes" {
			t.Errorf("Expected key to be %v, got %v", "Archimedes", kv.Key("Archimedes"))
		}

		if kv.Key("Do Not Exist") != "" {
			t.Errorf("Expected key to be %v, got %v", "Archimedes", kv.Key("Do Not Exist"))
		}
	})
}

func TestKeys(t *testing.T) {
	t.Run("test Keys for NewMapKeyValue[string, struct] with keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		kv.Set("Archimedes", testStruct{"This is Archimedes' Constant (Pi)", 3.1415})
		kv.Set("Euler", testStruct{"This is Euler's Number (e)", 2.7182})
		kv.Set("Golden Ratio", testStruct{"This is The Golden Ratio", 1.6180})

		if kv.Size() != 3 {
			t.Errorf("Expected size to be %v, got %v", 3, kv.Size())
		}

		keys := kv.Keys()

		if len(keys) != 3 {
			t.Errorf("Expected size to be %v, got %v", 3, len(keys))
		}

		for _, key := range keys {
			if key != "Archimedes" && key != "Euler" && key != "Golden Ratio" {
				t.Errorf("Expected key to be %v, got %v", true, key)
			}
		}
	})

	t.Run("test Keys for NewMapKeyValue[string, struct] without keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		keys := kv.Keys()

		if len(keys) != 0 {
			t.Errorf("Expected size to be %v, got %v", 0, len(keys))
		}
	})
}

func TestValues(t *testing.T) {
	t.Run("test Values for NewMapKeyValue[string, struct] with keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		kv.Set("Archimedes", testStruct{"This is Archimedes' Constant (Pi)", 3.1415})
		kv.Set("Euler", testStruct{"This is Euler's Number (e)", 2.7182})
		kv.Set("Golden Ratio", testStruct{"This is The Golden Ratio", 1.6180})

		if kv.Size() != 3 {
			t.Errorf("Expected size to be %v, got %v", 3, kv.Size())
		}

		values := kv.Values()

		if len(values) != 3 {
			t.Errorf("Expected size to be %v, got %v", 3, len(values))
		}

		for _, val := range values {
			if val.Name != "This is Archimedes' Constant (Pi)" && val.Name != "This is Euler's Number (e)" && val.Name != "This is The Golden Ratio" {
				t.Errorf("Expected value to be %v, got %v", true, val)
			}

			if val.value != 3.1415 && val.value != 2.7182 && val.value != 1.6180 {
				t.Errorf("Expected value to be %v, got %v", true, val)
			}
		}
	})

	t.Run("test Values for NewMapKeyValue[string, struct] without keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		values := kv.Values()

		if len(values) != 0 {
			t.Errorf("Expected size to be %v, got %v", 0, len(values))
		}
	})
}

func TestEach(t *testing.T) {
	t.Run("test Each for NewMapKeyValue[string, struct] with keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		kv.Set("Archimedes", testStruct{"This is Archimedes' Constant (Pi)", 3.1415})
		kv.Set("Euler", testStruct{"This is Euler's Number (e)", 2.7182})
		kv.Set("Golden Ratio", testStruct{"This is The Golden Ratio", 1.6180})

		if kv.Size() != 3 {
			t.Errorf("Expected size to be %v, got %v", 3, kv.Size())
		}

		kv.ForEach(func(key string, value testStruct) {
			if key != "Archimedes" && key != "Euler" && key != "Golden Ratio" {
				t.Errorf("Expected key to be %v, got %v", true, key)
			}

			if value.Name != "This is Archimedes' Constant (Pi)" && value.Name != "This is Euler's Number (e)" && value.Name != "This is The Golden Ratio" {
				t.Errorf("Expected value to be %v, got %v", true, value)
			}

			if value.value != 3.1415 && value.value != 2.7182 && value.value != 1.6180 {
				t.Errorf("Expected value to be %v, got %v", true, value)
			}
		})
	})

	t.Run("test Each for NewMapKeyValue[string, struct] without keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		kv.ForEach(func(key string, value testStruct) {
			t.Errorf("Expected Each to not be called, got %v", true)
		})
	})
}

func TestEachKey(t *testing.T) {
	t.Run("test EachKey for NewMapKeyValue[string, struct] with keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		kv.Set("Archimedes", testStruct{"This is Archimedes' Constant (Pi)", 3.1415})
		kv.Set("Euler", testStruct{"This is Euler's Number (e)", 2.7182})
		kv.Set("Golden Ratio", testStruct{"This is The Golden Ratio", 1.6180})

		if kv.Size() != 3 {
			t.Errorf("Expected size to be %v, got %v", 3, kv.Size())
		}

		kv.ForEachKey(func(key string) {
			if key != "Archimedes" && key != "Euler" && key != "Golden Ratio" {
				t.Errorf("Expected key to be %v, got %v", true, key)
			}
		})
	})

	t.Run("test EachKey for NewMapKeyValue[string, struct] without keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		kv.ForEachKey(func(key string) {
			t.Errorf("Expected EachKey to not be called, got %v", true)
		})
	})
}

func TestEachValue(t *testing.T) {
	t.Run("test EachValue for NewMapKeyValue[string, struct] with keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		kv.Set("Archimedes", testStruct{"This is Archimedes' Constant (Pi)", 3.1415})
		kv.Set("Euler", testStruct{"This is Euler's Number (e)", 2.7182})
		kv.Set("Golden Ratio", testStruct{"This is The Golden Ratio", 1.6180})

		if kv.Size() != 3 {
			t.Errorf("Expected size to be %v, got %v", 3, kv.Size())
		}

		kv.ForEachValue(func(value testStruct) {
			if value.Name != "This is Archimedes' Constant (Pi)" && value.Name != "This is Euler's Number (e)" && value.Name != "This is The Golden Ratio" {
				t.Errorf("Expected value to be %v, got %v", true, value)
			}

			if value.value != 3.1415 && value.value != 2.7182 && value.value != 1.6180 {
				t.Errorf("Expected value to be %v, got %v", true, value)
			}
		})
	})

	t.Run("test EachValue for NewMapKeyValue[string, struct] without keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		kv.ForEachValue(func(value testStruct) {
			t.Errorf("Expected EachValue to not be called, got %v", true)
		})
	})
}

func TestClone(t *testing.T) {
	t.Run("test Clone for NewMapKeyValue[string, struct] with keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		kv.Set("Archimedes", testStruct{"This is Archimedes' Constant (Pi)", 3.1415})
		kv.Set("Euler", testStruct{"This is Euler's Number (e)", 2.7182})
		kv.Set("Golden Ratio", testStruct{"This is The Golden Ratio", 1.6180})

		if kv.Size() != 3 {
			t.Errorf("Expected size to be %v, got %v", 3, kv.Size())
		}

		kvClone := kv.Clone()

		if kvClone.Size() != kv.Size() {
			t.Errorf("Expected size to be %v, got %v", 3, kvClone.Size())
		}

		if reflect.DeepEqual(kv, kvClone) == false {
			t.Errorf("Expected Clone to be equal to original, got %v", true)
		}

		kvKeys := kv.Keys()
		kvValues := kv.Values()

		kvCloneKeys := kvClone.Keys()

		sort.Strings(kvKeys)
		sort.Strings(kvCloneKeys)

		if reflect.DeepEqual(kvKeys, kvCloneKeys) == false {
			t.Errorf("Expected keys to be equal, got %v", true)
		}

		for _, kvValue := range kvValues {
			if kvClone.ContainsValue(kvValue) == false {
				t.Errorf("Expected Clone to contain value, got %v", true)
			}
		}
	})

	t.Run("test Clone for NewMapKeyValue[string, struct] without keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		kvClone := kv.Clone()

		if kvClone.Size() != kv.Size() {
			t.Errorf("Expected size to be %v, got %v", 3, kvClone.Size())
		}

		if reflect.DeepEqual(kv, kvClone) == false {
			t.Errorf("Expected Clone to be equal to original, got %v", true)
		}
	})
}

func TestCloneAndClear(t *testing.T) {
	t.Run("test CloneAndClear for NewMapKeyValue[string, struct] with keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		kv.Set("Archimedes", testStruct{"This is Archimedes' Constant (Pi)", 3.1415})
		kv.Set("Euler", testStruct{"This is Euler's Number (e)", 2.7182})
		kv.Set("Golden Ratio", testStruct{"This is The Golden Ratio", 1.6180})

		if kv.Size() != 3 {
			t.Errorf("Expected size to be %v, got %v", 3, kv.Size())
		}

		kvClone := kv.CloneAndClear()

		if kv.Size() != 0 {
			t.Errorf("Expected size to be %v, got %v", 0, kv.Size())
		}

		if kvClone.Size() == kv.Size() {
			t.Errorf("Expected size to be %v, got %v", 3, kvClone.Size())
		}

		if reflect.DeepEqual(kv, kvClone) == true {
			t.Errorf("Expected Clone to be not equal to original, got %v", true)
		}

		kvKeys := kv.Keys()
		kvCloneKeys := kvClone.Keys()

		sort.Strings(kvKeys)
		sort.Strings(kvCloneKeys)

		if reflect.DeepEqual(kvKeys, kvCloneKeys) == true {
			t.Errorf("Expected keys to be not equal, got %v", true)
		}
	})

	t.Run("test CloneAndClear for NewMapKeyValue[string, struct] without keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		kvClone := kv.CloneAndClear()

		if kv.Size() != 0 {
			t.Errorf("Expected size to be %v, got %v", 0, kv.Size())
		}

		if kvClone.Size() != 0 {
			t.Errorf("Expected size to be %v, got %v", 3, kvClone.Size())
		}

		if reflect.DeepEqual(kv, kvClone) == false {
			t.Errorf("Expected Clone to be equal to original, got %v", true)
		}
	})
}

func TestDeepEqual(t *testing.T) {
	t.Run("test DeepEqual for NewMapKeyValue[string, struct] with keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv1 := NewMapKeyValue[string, testStruct]()
		kv2 := NewMapKeyValue[string, testStruct]()

		kv1.Set("Archimedes", testStruct{"This is Archimedes' Constant (Pi)", 3.1415})
		kv1.Set("Euler", testStruct{"This is Euler's Number (e)", 2.7182})
		kv1.Set("Golden Ratio", testStruct{"This is The Golden Ratio", 1.6180})

		kv2.Set("Archimedes", testStruct{"This is Archimedes' Constant (Pi)", 3.1415})
		kv2.Set("Euler", testStruct{"This is Euler's Number (e)", 2.7182})
		kv2.Set("Golden Ratio", testStruct{"This is The Golden Ratio", 1.6180})

		if kv1.Size() != 3 {
			t.Errorf("Expected size to be %v, got %v", 3, kv1.Size())
		}
		if kv2.Size() != 3 {
			t.Errorf("Expected size to be %v, got %v", 3, kv2.Size())
		}

		if kv1.DeepEqual(kv2) == false {
			t.Errorf("Expected DeepEqual to be equal, got %v", true)
		}

		if kv2.DeepEqual(kv1) == false {
			t.Errorf("Expected DeepEqual to be equal, got %v", true)
		}
	})

	t.Run("test DeepEqual for NewMapKeyValue[string, struct] with keys and it is not equal same size", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv1 := NewMapKeyValue[string, testStruct]()
		kv2 := NewMapKeyValue[string, testStruct]()

		kv1.Set("Archimedes", testStruct{"This is Archimedes' Constant (Pi)", 3.1415})
		kv1.Set("Euler", testStruct{"This is Euler's Number (e)", 2.7182})
		kv1.Set("Golden Ratio", testStruct{"This is The Golden Ratio", 1.6180})

		kv2.Set("key 1", testStruct{"This is Archimedes' Constant (Pi)", 3.1415})
		kv2.Set("Euler", testStruct{"This is Euler's Number (e)", 2.7182})
		kv2.Set("Golden Ratio", testStruct{"This is The Golden Ratio", 1.6180})

		if kv1.Size() != 3 {
			t.Errorf("Expected size to be %v, got %v", 3, kv1.Size())
		}
		if kv2.Size() != 3 {
			t.Errorf("Expected size to be %v, got %v", 3, kv2.Size())
		}

		if kv1.DeepEqual(kv2) == true {
			t.Errorf("Expected DeepEqual to be equal, got %v", true)
		}

		if kv2.DeepEqual(kv1) == true {
			t.Errorf("Expected DeepEqual to be equal, got %v", true)
		}
	})

	t.Run("test DeepEqual for NewMapKeyValue[string, struct] with keys and it is not equal different size", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv1 := NewMapKeyValue[string, testStruct]()
		kv2 := NewMapKeyValue[string, testStruct]()

		kv1.Set("Archimedes", testStruct{"This is Archimedes' Constant (Pi)", 3.1415})
		kv1.Set("Euler", testStruct{"This is Euler's Number (e)", 2.7182})
		kv1.Set("Golden Ratio", testStruct{"This is The Golden Ratio", 1.6180})

		kv2.Set("key 1", testStruct{"This is Archimedes' Constant (Pi)", 3.1415})

		if kv1.Size() != 3 {
			t.Errorf("Expected size to be %v, got %v", 3, kv1.Size())
		}
		if kv2.Size() != 1 {
			t.Errorf("Expected size to be %v, got %v", 1, kv2.Size())
		}

		if kv1.DeepEqual(kv2) == true {
			t.Errorf("Expected DeepEqual to be equal, got %v", true)
		}

		if kv2.DeepEqual(kv1) == true {
			t.Errorf("Expected DeepEqual to be equal, got %v", true)
		}
	})

	t.Run("test DeepEqual for NewMapKeyValue[string, struct] without keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv1 := NewMapKeyValue[string, testStruct]()
		kv2 := NewMapKeyValue[string, testStruct]()

		if kv1.Size() != 0 {
			t.Errorf("Expected size to be %v, got %v", 0, kv1.Size())
		}
		if kv2.Size() != 0 {
			t.Errorf("Expected size to be %v, got %v", 0, kv2.Size())
		}

		if kv1.DeepEqual(kv2) == false {
			t.Errorf("Expected DeepEqual to be equal, got %v", true)
		}
		if kv2.DeepEqual(kv1) == false {
			t.Errorf("Expected DeepEqual to be equal, got %v", true)
		}
	})
}

func TestMap(t *testing.T) {
	t.Run("test Map for NewMapKeyValue[string, struct] with keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		kv.Set("Archimedes", testStruct{"This is Archimedes' Constant (Pi)", 3.1415})
		kv.Set("Euler", testStruct{"This is Euler's Number (e)", 2.7182})
		kv.Set("Golden Ratio", testStruct{"This is The Golden Ratio", 1.6180})

		if kv.Size() != 3 {
			t.Errorf("Expected size to be %v, got %v", 3, kv.Size())
		}

		newKv := kv.Map(func(key string, value testStruct) (newKey string, newValue testStruct) {
			newKey = key
			newValue.Name = strings.ToUpper(value.Name)
			newValue.value = value.value * 2
			return
		})

		newKv.ForEach(func(key string, value testStruct) {
			if kv.Key(key) != key {
				t.Errorf("Expected key to be uppercase, want: %v, got %v", strings.ToUpper(key), key)
			}
			if strings.ToUpper(kv.Get(key).Name) != value.Name {
				t.Errorf("Expected value.Name to be uppercase, want: %v, got %v", strings.ToUpper(kv.Get(key).Name), value.Name)
			}
			if kv.Get(key).value*2 != value.value {
				t.Errorf("Expected value.value to be doubled, want: %v, got %v", value.value*2, value.value)
			}
		})
	})

	t.Run("test Map for NewMapKeyValue[string, struct] without keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		newKv := kv.Map(func(key string, value testStruct) (newKey string, newValue testStruct) {
			newKey = strings.ToUpper(key)
			newValue = value
			return
		})

		if kv.Size() != 0 {
			t.Errorf("Expected size to be %v, got %v", 0, kv.Size())
		}
		if newKv.Size() != 0 {
			t.Errorf("Expected size to be %v, got %v", 0, newKv.Size())
		}
	})
}

func TestMapKey(t *testing.T) {
	t.Run("test MapKey for NewMapKeyValue[string, struct] with keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		kv.Set("Archimedes", testStruct{"This is Archimedes' Constant (Pi)", 3.1415})
		kv.Set("Euler", testStruct{"This is Euler's Number (e)", 2.7182})
		kv.Set("Golden Ratio", testStruct{"This is The Golden Ratio", 1.6180})

		if kv.Size() != 3 {
			t.Errorf("Expected size to be %v, got %v", 3, kv.Size())
		}

		newKv := kv.MapKey(func(key string) string {
			return strings.ToUpper(key)
		})

		newKv.ForEach(func(key string, value testStruct) {
			if strings.ToUpper(kv.Key(strings.Title(strings.ToLower(key)))) != key {
				t.Errorf("Expected key to be uppercase, want: %v, got %v", kv.Key(strings.Title(strings.ToLower(key))), key)
			}
			if kv.Get(strings.Title(strings.ToLower(key))).Name != value.Name {
				t.Errorf("Expected value.Name to be uppercase, want: %v, got %v", kv.Get(strings.Title(strings.ToLower(key))).Name, value.Name)
			}
			if kv.Get(strings.Title(strings.ToLower(key))).value != value.value {
				t.Errorf("Expected value.value to be doubled, want: %v, got %v", kv.Get(strings.Title(strings.ToLower(key))).value, value.value)
			}
		})
	})

	t.Run("test MapKey for NewMapKeyValue[string, struct] without keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		newKv := kv.MapKey(func(key string) string {
			return strings.ToUpper(key)
		})

		if kv.Size() != 0 {
			t.Errorf("Expected size to be %v, got %v", 0, kv.Size())
		}
		if newKv.Size() != 0 {
			t.Errorf("Expected size to be %v, got %v", 0, newKv.Size())
		}
	})
}

func TestMapValue(t *testing.T) {
	t.Run("test MapValue for NewMapKeyValue[string, struct] with keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		kv.Set("Archimedes", testStruct{"This is Archimedes' Constant (Pi)", 3.1415})
		kv.Set("Euler", testStruct{"This is Euler's Number (e)", 2.7182})
		kv.Set("Golden Ratio", testStruct{"This is The Golden Ratio", 1.6180})

		if kv.Size() != 3 {
			t.Errorf("Expected size to be %v, got %v", 3, kv.Size())
		}

		newKv := kv.MapValue(func(value testStruct) testStruct {
			value.Name = strings.ToUpper(value.Name)
			value.value = value.value * 2
			return value
		})

		newKv.ForEach(func(key string, value testStruct) {
			if kv.Key(key) != key {
				t.Errorf("Expected key to be uppercase, want: %v, got %v", kv.Key(key), key)
			}
			if strings.ToUpper(kv.Get(key).Name) != value.Name {
				t.Errorf("Expected value.Name to be uppercase, want: %v, got %v", kv.Get(key).Name, value.Name)
			}
			if kv.Get(key).value*2 != value.value {
				t.Errorf("Expected value.value to be doubled, want: %v, got %v", kv.Get(key).value, value.value)
			}
		})
	})

	t.Run("test MapValue for NewMapKeyValue[string, struct] without keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		newKv := kv.MapValue(func(value testStruct) testStruct {
			value.Name = strings.ToUpper(value.Name)
			value.value = value.value * 2
			return value
		})

		if kv.Size() != 0 {
			t.Errorf("Expected size to be %v, got %v", 0, kv.Size())
		}
		if newKv.Size() != 0 {
			t.Errorf("Expected size to be %v, got %v", 0, newKv.Size())
		}
	})
}

func TestFilter(t *testing.T) {
	t.Run("test Filter for NewMapKeyValue[string, struct] with keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		kv.Set("Archimedes", testStruct{"This is Archimedes' Constant (Pi)", 3.1415})
		kv.Set("Euler", testStruct{"This is Euler's Number (e)", 2.7182})
		kv.Set("Golden Ratio", testStruct{"This is The Golden Ratio", 1.6180})

		if kv.Size() != 3 {
			t.Errorf("Expected size to be %v, got %v", 3, kv.Size())
		}

		newKv := kv.Filter(func(key string, value testStruct) bool {
			return strings.Contains(value.Name, "Constant")
		})

		newKv.ForEach(func(key string, value testStruct) {
			if key != "Archimedes" {
				t.Errorf("Expected key to be uppercase, want: %v, got %v", kv.Key(key), key)
			}
			if value.Name != "This is Archimedes' Constant (Pi)" {
				t.Errorf("Expected value.Name to be uppercase, want: %v, got %v", "This is Archimedes' Constant (Pi)", value.Name)
			}
			if value.value != 3.1415 {
				t.Errorf("Expected value.value to be doubled, want: %v, got %v", kv.Get(key).value, value.value)
			}
		})
	})

	t.Run("test Filter for NewMapKeyValue[string, struct] without keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		newKv := kv.Filter(func(key string, value testStruct) bool {
			return strings.Contains(value.Name, "Constant")
		})

		if kv.Size() != 0 {
			t.Errorf("Expected size to be %v, got %v", 0, kv.Size())
		}
		if newKv.Size() != 0 {
			t.Errorf("Expected size to be %v, got %v", 0, newKv.Size())
		}
	})
}

func TestFilterKey(t *testing.T) {
	t.Run("test FilterKey for NewMapKeyValue[string, struct] with keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		kv.Set("Archimedes", testStruct{"This is Archimedes' Constant (Pi)", 3.1415})
		kv.Set("Euler", testStruct{"This is Euler's Number (e)", 2.7182})
		kv.Set("Golden Ratio", testStruct{"This is The Golden Ratio", 1.6180})

		if kv.Size() != 3 {
			t.Errorf("Expected size to be %v, got %v", 3, kv.Size())
		}

		newKv := kv.FilterKey(func(key string) bool {
			return strings.Contains(key, "chime")
		})

		newKv.ForEach(func(key string, value testStruct) {
			if key != "Archimedes" {
				t.Errorf("Expected key to be uppercase, want: %v, got %v", kv.Key(key), key)
			}
			if value.Name != "This is Archimedes' Constant (Pi)" {
				t.Errorf("Expected value.Name to be uppercase, want: %v, got %v", "This is Archimedes' Constant (Pi)", value.Name)
			}
			if value.value != 3.1415 {
				t.Errorf("Expected value.value to be doubled, want: %v, got %v", kv.Get(key).value, value.value)
			}
		})
	})

	t.Run("test FilterKey for NewMapKeyValue[string, struct] without keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		newKv := kv.FilterKey(func(key string) bool {
			return strings.Contains(key, "chime")
		})

		if kv.Size() != 0 {
			t.Errorf("Expected size to be %v, got %v", 0, kv.Size())
		}
		if newKv.Size() != 0 {
			t.Errorf("Expected size to be %v, got %v", 0, newKv.Size())
		}
	})
}

func TestFilterValue(t *testing.T) {
	t.Run("test FilterValue for NewMapKeyValue[string, struct] with keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		kv.Set("Archimedes", testStruct{"This is Archimedes' Constant (Pi)", 3.1415})
		kv.Set("Euler", testStruct{"This is Euler's Number (e)", 2.7182})
		kv.Set("Golden Ratio", testStruct{"This is The Golden Ratio", 1.6180})

		if kv.Size() != 3 {
			t.Errorf("Expected size to be %v, got %v", 3, kv.Size())
		}

		newKv := kv.FilterValue(func(value testStruct) bool {
			return value.value > 3
		})

		if newKv.Size() != 1 {
			t.Errorf("Expected size to be %v, got %v", 1, newKv.Size())
		}

		newKv.ForEach(func(key string, value testStruct) {
			if key != "Archimedes" {
				t.Errorf("Expected key to be uppercase, want: %v, got %v", kv.Key(key), key)
			}
			if value.Name != "This is Archimedes' Constant (Pi)" {
				t.Errorf("Expected value.Name to be uppercase, want: %v, got %v", "This is Archimedes' Constant (Pi)", value.Name)
			}
			if value.value != 3.1415 {
				t.Errorf("Expected value.value to be doubled, want: %v, got %v", kv.Get(key).value, value.value)
			}
		})
	})

	t.Run("test FilterValue for NewMapKeyValue[string, struct] without keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		newKv := kv.FilterValue(func(value testStruct) bool {
			return value.value > 3
		})

		if kv.Size() != 0 {
			t.Errorf("Expected size to be %v, got %v", 0, kv.Size())
		}
		if newKv.Size() != 0 {
			t.Errorf("Expected size to be %v, got %v", 0, newKv.Size())
		}
	})
}

func TestPartition(t *testing.T) {
	t.Run("test Partition for NewMapKeyValue[string, struct] with keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		kv.Set("Archimedes", testStruct{"This is Archimedes' Constant (Pi)", 3.1415})
		kv.Set("Euler", testStruct{"This is Euler's Number (e)", 2.7182})
		kv.Set("Golden Ratio", testStruct{"This is The Golden Ratio", 1.6180})

		if kv.Size() != 3 {
			t.Errorf("Expected size to be %v, got %v", 3, kv.Size())
		}

		grp1Kv, grp2Kv := kv.Partition(func(key string, value testStruct) bool {
			return value.value > 3
		})

		if grp1Kv.Size() != 1 {
			t.Errorf("Expected size to be %v, got %v", 1, grp1Kv.Size())
		}
		if grp2Kv.Size() != 2 {
			t.Errorf("Expected size to be %v, got %v", 2, grp1Kv.Size())
		}

		grp1Kv.ForEach(func(key string, value testStruct) {
			if key != "Archimedes" {
				t.Errorf("Expected key to be uppercase, want: %v, got %v", kv.Key(key), key)
			}
			if value.Name != "This is Archimedes' Constant (Pi)" {
				t.Errorf("Expected value.Name to be uppercase, want: %v, got %v", "This is Archimedes' Constant (Pi)", value.Name)
			}
			if value.value != 3.1415 {
				t.Errorf("Expected value.value to be doubled, want: %v, got %v", kv.Get(key).value, value.value)
			}
		})

		grp2Kv.ForEach(func(key string, value testStruct) {
			if key != "Euler" && key != "Golden Ratio" {
				t.Errorf("Expected key to be uppercase, want: %v, got %v", kv.Key(key), key)
			}
			if value.Name != "This is Euler's Number (e)" && value.Name != "This is The Golden Ratio" {
				t.Errorf("Expected value.Name to be uppercase, want: %v, got %v", "This is Euler's Number (e)", value.Name)
			}
			if value.value != 2.7182 && value.value != 1.6180 {
				t.Errorf("Expected value.value to be doubled, want: %v, got %v", kv.Get(key).value, value.value)
			}
		})
	})

	t.Run("test Partition for NewMapKeyValue[string, struct] without keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		grp1Kv, grp2Kv := kv.Partition(func(key string, value testStruct) bool {
			return value.value > 3
		})

		if kv.Size() != 0 {
			t.Errorf("Expected size to be %v, got %v", 0, kv.Size())
		}
		if grp1Kv.Size() != 0 {
			t.Errorf("Expected size to be %v, got %v", 0, grp1Kv.Size())
		}

		if grp2Kv.Size() != 0 {
			t.Errorf("Expected size to be %v, got %v", 0, grp2Kv.Size())
		}
	})
}

func TestPartitionKey(t *testing.T) {
	t.Run("test PartitionKey for NewMapKeyValue[string, struct] with keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		kv.Set("Archimedes", testStruct{"This is Archimedes' Constant (Pi)", 3.1415})
		kv.Set("Euler", testStruct{"This is Euler's Number (e)", 2.7182})
		kv.Set("Golden Ratio", testStruct{"This is The Golden Ratio", 1.6180})

		if kv.Size() != 3 {
			t.Errorf("Expected size to be %v, got %v", 3, kv.Size())
		}

		grp1Kv, grp2Kv := kv.PartitionKey(func(key string) bool {
			return key == "Archimedes"
		})

		if grp1Kv.Size() != 1 {
			t.Errorf("Expected size to be %v, got %v", 1, grp1Kv.Size())
		}
		if grp2Kv.Size() != 2 {
			t.Errorf("Expected size to be %v, got %v", 2, grp1Kv.Size())
		}

		grp1Kv.ForEach(func(key string, value testStruct) {
			if key != "Archimedes" {
				t.Errorf("Expected key to be uppercase, want: %v, got %v", kv.Key(key), key)
			}
			if value.Name != "This is Archimedes' Constant (Pi)" {
				t.Errorf("Expected value.Name to be uppercase, want: %v, got %v", "This is Archimedes' Constant (Pi)", value.Name)
			}
			if value.value != 3.1415 {
				t.Errorf("Expected value.value to be doubled, want: %v, got %v", kv.Get(key).value, value.value)
			}
		})

		grp2Kv.ForEach(func(key string, value testStruct) {
			if key != "Euler" && key != "Golden Ratio" {
				t.Errorf("Expected key to be uppercase, want: %v, got %v", kv.Key(key), key)
			}
			if value.Name != "This is Euler's Number (e)" && value.Name != "This is The Golden Ratio" {
				t.Errorf("Expected value.Name to be uppercase, want: %v, got %v", "This is Euler's Number (e)", value.Name)
			}
			if value.value != 2.7182 && value.value != 1.6180 {
				t.Errorf("Expected value.value to be doubled, want: %v, got %v", kv.Get(key).value, value.value)
			}
		})
	})

	t.Run("test PartitionKey for NewMapKeyValue[string, struct] without keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		grp1Kv, grp2Kv := kv.PartitionKey(func(key string) bool {
			return key == "Archimedes"
		})

		if kv.Size() != 0 {
			t.Errorf("Expected size to be %v, got %v", 0, kv.Size())
		}
		if grp1Kv.Size() != 0 {
			t.Errorf("Expected size to be %v, got %v", 0, grp1Kv.Size())
		}

		if grp2Kv.Size() != 0 {
			t.Errorf("Expected size to be %v, got %v", 0, grp2Kv.Size())
		}
	})
}

func TestPartitionValue(t *testing.T) {
	t.Run("test PartitionValue for NewMapKeyValue[string, struct] with keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		kv.Set("Archimedes", testStruct{"This is Archimedes' Constant (Pi)", 3.1415})
		kv.Set("Euler", testStruct{"This is Euler's Number (e)", 2.7182})
		kv.Set("Golden Ratio", testStruct{"This is The Golden Ratio", 1.6180})

		if kv.Size() != 3 {
			t.Errorf("Expected size to be %v, got %v", 3, kv.Size())
		}

		grp1Kv, grp2Kv := kv.PartitionValue(func(value testStruct) bool {
			return value.value > 3
		})

		if grp1Kv.Size() != 1 {
			t.Errorf("Expected size to be %v, got %v", 1, grp1Kv.Size())
		}
		if grp2Kv.Size() != 2 {
			t.Errorf("Expected size to be %v, got %v", 2, grp1Kv.Size())
		}

		grp1Kv.ForEach(func(key string, value testStruct) {
			if key != "Archimedes" {
				t.Errorf("Expected key to be uppercase, want: %v, got %v", kv.Key(key), key)
			}
			if value.Name != "This is Archimedes' Constant (Pi)" {
				t.Errorf("Expected value.Name to be uppercase, want: %v, got %v", "This is Archimedes' Constant (Pi)", value.Name)
			}
			if value.value != 3.1415 {
				t.Errorf("Expected value.value to be doubled, want: %v, got %v", kv.Get(key).value, value.value)
			}
		})

		grp2Kv.ForEach(func(key string, value testStruct) {
			if key != "Euler" && key != "Golden Ratio" {
				t.Errorf("Expected key to be uppercase, want: %v, got %v", kv.Key(key), key)
			}
			if value.Name != "This is Euler's Number (e)" && value.Name != "This is The Golden Ratio" {
				t.Errorf("Expected value.Name to be uppercase, want: %v, got %v", "This is Euler's Number (e)", value.Name)
			}
			if value.value != 2.7182 && value.value != 1.6180 {
				t.Errorf("Expected value.value to be doubled, want: %v, got %v", kv.Get(key).value, value.value)
			}
		})
	})

	t.Run("test PartitionValue for NewMapKeyValue[string, struct] without keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		grp1Kv, grp2Kv := kv.PartitionValue(func(value testStruct) bool {
			return value.value > 3
		})

		if kv.Size() != 0 {
			t.Errorf("Expected size to be %v, got %v", 0, kv.Size())
		}
		if grp1Kv.Size() != 0 {
			t.Errorf("Expected size to be %v, got %v", 0, grp1Kv.Size())
		}

		if grp2Kv.Size() != 0 {
			t.Errorf("Expected size to be %v, got %v", 0, grp2Kv.Size())
		}
	})
}

func TestSortKeys(t *testing.T) {
	t.Run("test SortKeys for NewMapKeyValue[string, struct] with keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		kv.Set("Golden Ratio", testStruct{"This is The Golden Ratio", 1.6180})
		kv.Set("Archimedes", testStruct{"This is Archimedes' Constant (Pi)", 3.1415})
		kv.Set("Euler", testStruct{"This is Euler's Number (e)", 2.7182})

		if kv.Size() != 3 {
			t.Errorf("Expected size to be %v, got %v", 3, kv.Size())
		}

		kSorted := kv.SortKeys(func(key1 string, key2 string) bool {
			return key1 < key2
		})

		if len(kSorted) != 3 {
			t.Errorf("Expected size to be %v, got %v", 3, len(kSorted))
		}

		if *kSorted[0] != "Archimedes" {
			t.Errorf("Expected key to be uppercase, want: %v, got %v", "Archimedes", *kSorted[0])
		}
	})

	t.Run("test SortKeys for NewMapKeyValue[string, struct] without keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		kSorted := kv.SortKeys(func(key1 string, key2 string) bool {
			return key1 < key2
		})

		if kv.Size() != 0 {
			t.Errorf("Expected size to be %v, got %v", 0, kv.Size())
		}
		if len(kSorted) != 0 {
			t.Errorf("Expected size to be %v, got %v", 0, len(kSorted))
		}
	})
}

func TestSortValues(t *testing.T) {
	t.Run("test SortValues for NewMapKeyValue[string, struct] with keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		kv.Set("Archimedes", testStruct{"This is Archimedes' Constant (Pi)", 3.1415})
		kv.Set("Golden Ratio", testStruct{"This is The Golden Ratio", 1.6180})
		kv.Set("Euler", testStruct{"This is Euler's Number (e)", 2.7182})

		if kv.Size() != 3 {
			t.Errorf("Expected size to be %v, got %v", 3, kv.Size())
		}

		vSorted := kv.SortValues(func(value1 testStruct, value2 testStruct) bool {
			return value1.value < value2.value
		})

		if len(vSorted) != 3 {
			t.Errorf("Expected size to be %v, got %v", 3, len(vSorted))
		}

		if vSorted[0].Name != "This is The Golden Ratio" {
			t.Errorf("Expected key to be uppercase, want: %v, got %v", "This is The Golden Ratio", vSorted[0].Name)
		}
		if vSorted[1].Name != "This is Euler's Number (e)" {
			t.Errorf("Expected key to be uppercase, want: %v, got %v", "This is Euler's Number (e)", vSorted[1].Name)
		}
		if vSorted[2].Name != "This is Archimedes' Constant (Pi)" {
			t.Errorf("Expected key to be uppercase, want: %v, got %v", "This is Archimedes' Constant (Pi)", vSorted[2].Name)
		}
	})

	t.Run("test SortValues for NewMapKeyValue[string, struct] without keys", func(t *testing.T) {
		type testStruct struct {
			Name  string
			value float64
		}
		kv := NewMapKeyValue[string, testStruct]()

		vSorted := kv.SortValues(func(value1 testStruct, value2 testStruct) bool {
			return value1.value < value2.value
		})

		if kv.Size() != 0 {
			t.Errorf("Expected size to be %v, got %v", 0, kv.Size())
		}
		if len(vSorted) != 0 {
			t.Errorf("Expected size to be %v, got %v", 0, len(vSorted))
		}
	})
}

// **************************************************
// ******************** Examples ********************

// Using int data types
func ExampleNewMapKeyValue_int() {
	kv := NewMapKeyValue[int, int]()

	kv.Set(1, 8096)
	kv.Set(25, 4096)

	fmt.Printf("key 1: %v, value 1: %v\nkey 2: %v, value 2: %v", 1, kv.Get(1), 25, kv.Get(25))
	// Output:
	// key 1: 1, value 1: 8096
	// key 2: 25, value 2: 4096
}

// Using string as key and struct as value data types.
func ExampleNewMapKeyValue_struct() {
	type testStruct struct {
		Name  string
		value float64
	}

	MathConstants := NewMapKeyValue[string, testStruct]()

	MathConstants.Set("Archimedes", testStruct{"This is Archimedes' Constant (Pi)", 3.1415})
	MathConstants.Set("Euler", testStruct{"This is Euler's Number (e)", 2.7182})
	MathConstants.Set("Golden Ratio", testStruct{"This is The Golden Ratio", 1.6180})

	fmt.Printf("name: %v, value: %v\n", MathConstants.Get("Archimedes").Name, MathConstants.Get("Archimedes").value)

	// Output:
	// name: This is Archimedes' Constant (Pi), value: 3.1415
}

// Using string as key and int as value data types.
func Example() {
	grades := NewMapKeyValue[string, float64]()

	grades.Set("John Doe", 7.456)
	grades.Set("Jane Doe", 9.876)
	grades.Set("Donato Ricupero", 9.123)
	grades.Set("Joe Blow", 9.123)
	grades.Set("Joe Doakes", 9.123)
	grades.Set("Joe Sixpack", 9.123)

	// show elements
	grades.ForEach(func(key string, value float64) {
		fmt.Printf("name: %v, grade: %v\n", key, value)
	})

	// show values
	values := grades.Values()
	for _, value := range values {
		fmt.Printf("grade: %v\n", value)
	}

	// show keys
	keys := grades.Values()
	for _, key := range keys {
		fmt.Printf("student: %v\n", key)
	}

	filterValues := grades.Filter(func(key string, value float64) bool {
		return value > 8
	})

	filterValues.ForEach(func(key string, value float64) {
		fmt.Printf("name: %v, grade: %v\n", key, value)
	})
}

// **************************************************
// ******************** Benchmarks ******************
func BenchmarkMapKeyValue_Set_int_int(b *testing.B) {
	kv := NewMapKeyValue[int, int](WithCapacity(kvSize))

	for i := 0; i < b.N; i++ {
		kv.Set(rand.Intn(kvSize), rand.Intn(kvSize))
	}
}

func BenchmarkMapKeyValue_Get_int_int(b *testing.B) {
	for i := 0; i < b.N; i++ {
		kv_int_int.Get(rand.Intn(kvSize))
	}
}

func BenchmarkMapKeyValue_Set_Get_int_int(b *testing.B) {
	kv := NewMapKeyValue[int, int](WithCapacity(kvSize))

	for i := 0; i < b.N; i++ {
		kv.Set(rand.Intn(kvSize), rand.Intn(kvSize))
	}

	for i := 0; i < b.N; i++ {
		kv.Get(rand.Intn(kvSize))
	}
}

func BenchmarkMapKeyValue_Set_string_string(b *testing.B) {
	kv := NewMapKeyValue[string, string](WithCapacity(kvSize))

	for i := 0; i < b.N; i++ {
		keyval := fmt.Sprintf("%x", md5.Sum([]byte(strconv.Itoa(rand.Intn(kvSize)))))
		kv.Set(keyval, keyval)
	}
}

func BenchmarkMapKeyValue_Get_string_string(b *testing.B) {
	for i := 0; i < b.N; i++ {
		keyval := fmt.Sprintf("%x", md5.Sum([]byte(strconv.Itoa(rand.Intn(kvSize)))))
		kv_string_string.Get(keyval)
	}
}

func BenchmarkMapKeyValue_Set_Get_string_string(b *testing.B) {
	kv := NewMapKeyValue[string, string](WithCapacity(kvSize))

	for i := 0; i < b.N; i++ {
		keyval := fmt.Sprintf("%x", md5.Sum([]byte(strconv.Itoa(rand.Intn(kvSize)))))
		kv.Set(keyval, keyval)
	}

	for i := 0; i < b.N; i++ {
		keyval := fmt.Sprintf("%x", md5.Sum([]byte(strconv.Itoa(rand.Intn(kvSize)))))
		kv.Get(keyval)
	}
}

func BenchmarkMapKeyValue_Set_string_struct(b *testing.B) {
	kv := NewMapKeyValue[string, TestStruct](WithCapacity(kvSize))

	for i := 0; i < b.N; i++ {
		keyval := fmt.Sprintf("%x", md5.Sum([]byte(strconv.Itoa(rand.Intn(kvSize)))))
		s := TestStruct{
			a: keyval,
			b: rand.Intn(kvSize),
			c: int64(rand.Intn(kvSize)),
			d: rand.Float64(),
		}
		kv.Set(keyval, s)
	}
}

func BenchmarkMapKeyValue_Get_string_struct(b *testing.B) {
	for i := 0; i < b.N; i++ {
		keyval := fmt.Sprintf("%x", md5.Sum([]byte(strconv.Itoa(rand.Intn(kvSize)))))
		kv_string_struct.Get(keyval)
	}
}

func BenchmarkMapKeyValue_Set_Get_string_struct(b *testing.B) {
	kv := NewMapKeyValue[string, TestStruct](WithCapacity(kvSize))

	for i := 0; i < b.N; i++ {
		keyval := fmt.Sprintf("%x", md5.Sum([]byte(strconv.Itoa(rand.Intn(kvSize)))))
		s := TestStruct{
			a: keyval,
			b: rand.Intn(kvSize),
			c: int64(rand.Intn(kvSize)),
			d: rand.Float64(),
		}
		kv.Set(keyval, s)
	}

	for i := 0; i < b.N; i++ {
		keyval := fmt.Sprintf("%x", md5.Sum([]byte(strconv.Itoa(rand.Intn(kvSize)))))
		kv.Get(keyval)
	}
}

func BenchmarkMapKeyValue_Set_Get_string_struct_concurrent(b *testing.B) {
	kv := NewMapKeyValue[string, TestStruct](WithCapacity(kvSize))

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		for i := 0; i < b.N; i++ {
			keyval := fmt.Sprintf("%x", md5.Sum([]byte(strconv.Itoa(rand.Intn(kvSize)))))
			s := TestStruct{
				a: keyval,
				b: rand.Intn(kvSize),
				c: int64(rand.Intn(kvSize)),
				d: rand.Float64(),
			}
			kv.Set(keyval, s)
		}
	}()

	wg.Add(1)
	go func() {
		for i := 0; i < b.N; i++ {
			keyval := fmt.Sprintf("%x", md5.Sum([]byte(strconv.Itoa(rand.Intn(kvSize)))))
			kv.Get(keyval)
		}
	}()

	wg.Done()
	wg.Done()

	wg.Wait()
}
