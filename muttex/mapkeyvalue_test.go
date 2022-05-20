package muttex

import (
	"crypto/md5"
	"fmt"
	"math/rand"
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
	kv_int_int       = NewMapKeyValue[int, int](WithSize(kvSize))
	kv_string_string = NewMapKeyValue[string, string](WithSize(kvSize))
	kv_string_struct = NewMapKeyValue[string, TestStruct](WithSize(kvSize))
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

// ******************** Benchmarks ********************
func BenchmarkMapKeyValue_Set_int_int(b *testing.B) {
	kv := NewMapKeyValue[int, int](WithSize(kvSize))

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
	kv := NewMapKeyValue[int, int](WithSize(kvSize))

	for i := 0; i < b.N; i++ {
		kv.Set(rand.Intn(kvSize), rand.Intn(kvSize))
	}

	for i := 0; i < b.N; i++ {
		kv.Get(rand.Intn(kvSize))
	}
}

func BenchmarkMapKeyValue_Set_string_string(b *testing.B) {
	kv := NewMapKeyValue[string, string](WithSize(kvSize))

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
	kv := NewMapKeyValue[string, string](WithSize(kvSize))

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
	kv := NewMapKeyValue[string, TestStruct](WithSize(kvSize))

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
	kv := NewMapKeyValue[string, TestStruct](WithSize(kvSize))

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
