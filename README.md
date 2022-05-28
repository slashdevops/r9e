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

This is focused on `usability and simplicity` rather than performance, but it doesn't mean that it's not fast.

## Overview

Taking advantage of the [Golang Generics](https://go.dev/blog/intro-generics) and internal golang data structures like [sync.Map](https://golang.org/pkg/sync/#Map), `RamStorage (r9e)` provides a simple way to store and retrieve data.

The goal is to provide a simple and easy way to use a library to store and retrieve data from memory using a simple API and data structures.

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
