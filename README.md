# RamStorage (r9e)

[![CodeQL Analysis](https://github.com/slashdevops/r9e/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/slashdevops/r9e/actions/workflows/codeql-analysis.yml)
[![Gosec](https://github.com/slashdevops/r9e/actions/workflows/gosec.yml/badge.svg)](https://github.com/slashdevops/r9e/actions/workflows/gosec.yml)
[![Unit Test](https://github.com/slashdevops/r9e/actions/workflows/main.yml/badge.svg)](https://github.com/slashdevops/r9e/actions/workflows/main.yml)
[![Release](https://github.com/slashdevops/r9e/actions/workflows/release.yml/badge.svg)](https://github.com/slashdevops/r9e/actions/workflows/release.yml)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/slashdevops/r9e?style=plastic)
[![license](https://img.shields.io/github/license/slashdevops/r9e.svg)](https://github.com/slashdevops/r9e/blob/main/LICENSE)
[![codecov](https://codecov.io/gh/slashdevops/r9e/branch/main/graph/badge.svg?token=UNTP5C1P6C)](https://codecov.io/gh/slashdevops/r9e)
[![release](https://img.shields.io/github/release/slashdevops/r9e/all.svg)](https://github.com/slashdevops/r9e/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/slashdevops/r9e.svg)](https://pkg.go.dev/github.com/slashdevops/r9e)

**RamStorage (r9e)** is a `Thread-Safe` [Golang](https://go.dev/) library used for memory storage with convenient methods to store and retrieve data.

This is focused on `usability and simplicity` rather than performance, but it doesn't mean that it's not fast. [Discover it for yourself](#how-fast)

## Overview

Taking advantage of the [Golang Generics](https://go.dev/blog/intro-generics) and internal golang data structures like [sync.Map](https://golang.org/pkg/sync/#Map), `RamStorage (r9e)` provides a simple way to store and retrieve data.

The goal is to provide a simple and easy way to use a library to store and retrieve data from memory using a simple API and data structures.

### Available Containers

* [MapKeyValue[K comparable, T any]](https://pkg.go.dev/github.com/slashdevops/r9e#MapKeyValue)

### Documentation

Official documentation is available on [pkg.go.dev -> slashdevops/r9e](https://pkg.go.dev/github.com/slashdevops/r9e)

## Installing

Latest release:

```bash
go get -u github.com/slashdevops/r9e@latest
```

Specific release:

```bash
go get -u github.com/slashdevops/r9e@vx.y.z
```

Adding it to your project:

```go
import "github.com/slashdevops/r9e"
```

## Example

```go
package main

import (
  "fmt"

  "github.com/slashdevops/r9e"
)

func main() {
  type MathematicalConstants struct {
    Name  string
    Value float64
  }

  // With Capacity allocated
  //kv := r9e.NewMapKeyValue[string, MathematicalConstants](r9e.WithCapacity(5))
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
  kvFilteredValues.Each(func(key string, value MathematicalConstants) {
    fmt.Printf("Key: %v, Name: %v, Value: %v\n", key, value.Name, value.Value)
  })

  fmt.Printf("\n")
  fmt.Printf("The most famous mathematical constant:\n")
  fmt.Printf("Name: %v, Value: %v\n", kv.Get("pi").Name, kv.Get("pi").Value)
}
```

Output:

```bash
Mathematical Constants:
Key: ρ, Name: Plastic number ρ (or silver constant), Value: 2.414213
Key: e, Name: Euler number, Napier's constant, Value: 2.718281
Key: pi, Name: Archimedes' constant, Value: 3.141592

The most famous mathematical constant:
Name: Archimedes' constant, Value: 3.141592
```

## How Fast?

Discover it for yourself:

```bash
git clone git@github.com:slashdevops/r9e.git
cd r9e/
make bench
```

## License

RamStorage (r9e)  is released under the Apache License Version 2.0:

* [http://www.apache.org/licenses/LICENSE-2.0.html](http://www.apache.org/licenses/LICENSE-2.0.html)
