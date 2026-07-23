# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) and other AI
assistants when working with code in this repository. Follow these guidelines
precisely to ensure consistency and maintainability.

## Stack

- Language: Go (Go 1.26+)
- Framework: Go standard library **only** — this library is intentionally
  dependency-free
- Testing: Go's built-in testing package
- Dependency Management: Go modules
- Version Control: Git
- Documentation: `go doc` / pkg.go.dev, and Markdown files under `docs/`
- Code Review: Pull requests on GitHub
- CI/CD: GitHub Actions

## What this library is

`RamStorage (r9e)` is a small, dependency-free Go library of **thread-safe,
generic in-memory key-value containers**. It exposes two containers built on Go
generics:

- `MapKeyValue[K comparable, T any]` — backed by a native map guarded by a
  `sync.RWMutex`; best for read-heavy or mixed workloads and consistent bulk
  snapshots.
- `SMapKeyValue[K comparable, T any]` — backed by `sync.Map` with an atomic
  size counter; best for disjoint-key writes and write-once/read-many workloads.

The focus is usability and simplicity without sacrificing performance.

## Key Conventions

- **Style:** Follow the Google Go Style Guide and Effective Go:
  - <https://google.github.io/styleguide/go/guide>
  - <https://google.github.io/styleguide/go/decisions>
  - <https://google.github.io/styleguide/go/best-practices>
  - <https://go.dev/doc/effective_go>
- Keep functions small and focused on a single task.
- Use meaningful names for variables, functions, and packages.
- Use comments to explain complex logic or decisions.
- Use `any`, never `interface{}`.
- Prefer `for b.Loop()` over `for i := 0; i < b.N; i++` in benchmarks (Go 1.24+).
- Never hold a lock while calling a method that takes the same lock (no
  recursive `RLock`); build results directly under a single lock, or take a
  snapshot first.
- **Tests:** Table-driven tests where practical. Test files are co-located with
  source (`*_test.go`). Executable examples live in `example_test.go` and are
  verified by `go test`.
- **No external dependencies:** Do not add third-party modules without explicit
  approval. `go.sum` should not exist unless a dependency is intentionally
  introduced.

## Project Structure

This is a single-package Go library; source files live in the repository root.

- `*.go` — library source code.
- `*_test.go` — unit tests and benchmarks.
- `example_test.go` — executable examples rendered by pkg.go.dev.
- `doc.go` — package-level documentation.
- `docs/` — extensive usage documentation in Markdown.
- `.github/` — GitHub Actions workflows, CodeQL, Dependabot, release metadata.
- `.golangci.yaml` — optional local golangci-lint configuration.
- `.vscode/` — editor settings.
- `LICENSE` — Apache License 2.0.
- `README.md` — project overview and usage guide.
- `SECURITY.md` — vulnerability reporting policy.
- `go.mod` — module definition with no external requirements.

## Post-Change Checklist

Use these standard Go commands after making changes, before committing:

```bash
go fix ./...
go fmt ./...
go vet ./...
betteralign -apply ./...
go test -race -coverprofile=/tmp/r9e-coverage.txt -covermode=atomic ./...
go build ./...
```

Do not add external module dependencies without explicit approval; this project
is intentionally standard-library-only.

## Commit Message & Pull Request Guidelines

- Always work on a new branch unless explicitly told to use the current one.
- Use semantic (Conventional Commits) messages, e.g. `feat:`, `fix:`, `docs:`,
  `test:`, `chore:`, `refactor:`.
- Keep the commit subject under 72 characters and the whole message concise.
- Group related changes into focused commits rather than one large commit.
- Open pull requests against `main` with a semantic title and a description that
  summarizes the changes and references any related issues.
- Keep changes small, idiomatic, tested, documented, and dependency-free unless
  there is a clear reason to expand the project scope.
