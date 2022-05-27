package r9e

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
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

// ******************** Test ********************
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
	grades.Each(func(key string, value float64) {
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

	filterValues.Each(func(key string, value float64) {
		fmt.Printf("name: %v, grade: %v\n", key, value)
	})
}

// ******************** Benchmarks ********************
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
