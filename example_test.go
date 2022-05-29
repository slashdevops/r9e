package r9e_test

import (
	"fmt"

	"github.com/slashdevops/r9e"
)

func ExampleMapKeyValue_basic() {
	type MathematicalConstants struct {
		Name  string
		Value float64
	}

	// With Capacity allocated
	// kv := r9e.NewMapKeyValue[string, MathematicalConstants](r9e.WithCapacity(5))
	kv := r9e.NewMapKeyValue[string, MathematicalConstants]()

	kv.Set("pi", MathematicalConstants{"Archimedes' constant", 3.141592})
	kv.Set("e", MathematicalConstants{"Euler number, Napier's constant", 2.718281})
	kv.Set("γ", MathematicalConstants{"Euler number, Napier's constant", 0.577215})
	kv.Set("Φ", MathematicalConstants{"Golden ratio constant", 1.618033})
	kv.Set("ρ", MathematicalConstants{"Plastic number ρ (or silver constant)", 2.414213})

	kvFilteredValues := kv.FilterValue(func(value MathematicalConstants) bool {
		return value.Value > 2.0
	})

	fmt.Println("Mathematical Constants:")
	kvFilteredValues.ForEach(func(key string, value MathematicalConstants) {
		fmt.Printf("Key: %v, Name: %v, Value: %v\n", key, value.Name, value.Value)
	})

	fmt.Printf("\n")
	fmt.Printf("The most famous mathematical constant:\n")
	fmt.Printf("Name: %v, Value: %v\n", kv.Get("pi").Name, kv.Get("pi").Value)

	lst := kv.SortValues(func(value1, value2 MathematicalConstants) bool {
		return value1.Value > value2.Value
	})

	fmt.Printf("\n")
	fmt.Printf("The most famous mathematical constant sorted by value:\n")
	for i, value := range lst {
		fmt.Printf("i: %v, Name: %v, Value: %v\n", i, value.Name, value.Value)
	}

	kvHigh, kvLow := kv.Partition(func(key string, value MathematicalConstants) bool {
		return value.Value > 2.5
	})

	fmt.Printf("\n")
	fmt.Printf("Mathematical constants which value is greater than 2.5:\n")
	kvHigh.ForEach(func(key string, value MathematicalConstants) {
		fmt.Printf("Key: %v, Name: %v, Value: %v\n", key, value.Name, value.Value)
	})

	fmt.Printf("\n")
	fmt.Printf("Mathematical constants which value is less than 2.5:\n")
	kvLow.ForEach(func(key string, value MathematicalConstants) {
		fmt.Printf("Key: %v, Name: %v, Value: %v\n", key, value.Name, value.Value)
	})
}

func ExampleSMapKeyValue_basic() {
	type MathematicalConstants struct {
		Name  string
		Value float64
	}

	kv := r9e.NewSMapKeyValue[string, MathematicalConstants]()

	kv.Set("pi", MathematicalConstants{"Archimedes' constant", 3.141592})
	kv.Set("e", MathematicalConstants{"Euler number, Napier's constant", 2.718281})
	kv.Set("γ", MathematicalConstants{"Euler number, Napier's constant", 0.577215})
	kv.Set("Φ", MathematicalConstants{"Golden ratio constant", 1.618033})
	kv.Set("ρ", MathematicalConstants{"Plastic number ρ (or silver constant)", 2.414213})

	kvFilteredValues := kv.FilterValue(func(value MathematicalConstants) bool {
		return value.Value > 2.0
	})

	fmt.Println("Mathematical Constants:")
	kvFilteredValues.ForEach(func(key string, value MathematicalConstants) {
		fmt.Printf("Key: %v, Name: %v, Value: %v\n", key, value.Name, value.Value)
	})

	fmt.Printf("\n")
	fmt.Printf("The most famous mathematical constant:\n")
	fmt.Printf("Name: %v, Value: %v\n", kv.Get("pi").Name, kv.Get("pi").Value)

	lst := kv.SortValues(func(value1, value2 MathematicalConstants) bool {
		return value1.Value > value2.Value
	})

	fmt.Printf("\n")
	fmt.Printf("The most famous mathematical constant sorted by value:\n")
	for i, value := range lst {
		fmt.Printf("i: %v, Name: %v, Value: %v\n", i, value.Name, value.Value)
	}

	kvHigh, kvLow := kv.Partition(func(key string, value MathematicalConstants) bool {
		return value.Value > 2.5
	})

	fmt.Printf("\n")
	fmt.Printf("Mathematical constants which value is greater than 2.5:\n")
	kvHigh.ForEach(func(key string, value MathematicalConstants) {
		fmt.Printf("Key: %v, Name: %v, Value: %v\n", key, value.Name, value.Value)
	})

	fmt.Printf("\n")
	fmt.Printf("Mathematical constants which value is less than 2.5:\n")
	kvLow.ForEach(func(key string, value MathematicalConstants) {
		fmt.Printf("Key: %v, Name: %v, Value: %v\n", key, value.Name, value.Value)
	})
}
