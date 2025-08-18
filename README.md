# Errors

[![Go Reference](https://pkg.go.dev/badge/fillmore-labs.com/exp/errors.svg)](https://pkg.go.dev/fillmore-labs.com/exp/errors)
[![Test](https://github.com/fillmore-labs/errors-exp/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/fillmore-labs/errors-exp/actions/workflows/test.yml)
[![CodeQL](https://github.com/fillmore-labs/errors-exp/actions/workflows/github-code-scanning/codeql/badge.svg?branch=main)](https://github.com/fillmore-labs/errors-exp/actions/workflows/github-code-scanning/codeql)
[![Coverage](https://codecov.io/gh/fillmore-labs/errors-exp/branch/main/graph/badge.svg?token=LEOQNYK9KB)](https://codecov.io/gh/fillmore-labs/errors-exp)
[![Go Report Card](https://goreportcard.com/badge/fillmore-labs.com/exp/errors)](https://goreportcard.com/report/fillmore-labs.com/exp/errors)
[![License](https://img.shields.io/github/license/fillmore-labs/errors-exp)](https://www.apache.org/licenses/LICENSE-2.0)

`fillmore-labs.com/exp/errors` is an experimental Go library providing two enhanced, generic alternatives to `errors.As`
for inspecting error chains with improved ergonomics and type safety.

## Motivation

In Go, error types can be designed for use as either values or pointers. When inspecting an error chain, a mismatch
between the type you're looking for (e.g., `MyError`) and the type that was actually wrapped (e.g., `*MyError`) can lead
to subtle, hard-to-find bugs where errors are missed.

See [this blog post](https://blog.fillmore-labs.com/posts/errors-1/) to read more about the issue and the `errortype`
linter.

## Function Overview

This library provides two complementary functions:

- **`HasError`** - A drop-in, type-safe replacement for `errors.As` with better ergonomics
- **`Has`** - An enhanced version that also handles pointer-value type mismatches automatically

| Feature                         | `errors.As` | `HasError` | `Has` |
| ------------------------------- | ----------- | ---------- | ----- |
| Generic return type             | ❌          | ✅         | ✅    |
| No target variable needed       | ❌          | ✅         | ✅    |
| Interface support               | ✅          | ✅         | ✅    |
| Pointer-value mismatch handling | ❌          | ❌         | ✅    |

## `HasError` - Enhanced Ergonomics

### `HasError` Overview

`HasError` is a direct, type-safe replacement for `errors.As` that uses Go generics for improved ergonomics and
readability.

```go
func HasError[T error](err error) (T, bool)
```

### `HasError` Key Benefits

#### Eliminates Target Variables

Instead of this with `errors.As`:

```go
  var myErr *MyError
  if errors.As(err, &myErr) { /* ... use myErr */ }
```

Write this with `HasError`:

```go
  if myErr, ok := HasError[*MyError](err); ok { /* ... use myErr */ }
```

#### Interface Support

`HasError` works seamlessly with interfaces. To check if an error implements a specific interface:

```go
  // Single-line variant
  if e, ok := HasError[interface { error; Temporary() bool }](err); ok && e.Temporary() { /* handle temporary error */ }

  // Or, using a named interface type:
  type temporary interface { error; Temporary() bool }
	if e, ok := HasError[temporary](err); ok && e.Temporary() { /* handle temporary error */}
```

### When to Use `HasError`

- You want better ergonomics than `errors.As`
- You need exact type matching behavior
- You're working with interfaces
- You don't have pointer-value mismatch concerns

## `Has` - Pointer-Value Flexibility

### `Has` Overview

`Has` provides all the benefits of `HasError` plus automatic handling of pointer-value type mismatches, preventing
subtle bugs.

```go
func Has[T error](err error) (T, bool)
```

### `Has` Key Benefits

#### Automatic Pointer-Value Matching

`Has` finds matching errors regardless of pointer-value mismatches:

```go
	// Scenario 1: Looking for a value, but a pointer was wrapped
	err := &MyError{msg: "oops"}
	if myErr, ok := Has[MyError](err); ok { /* This matches! Has automatically handles the mismatch */ }

	// Scenario 2: Looking for a pointer, but a value was wrapped
	err2 := MyError{msg: "oops"}
	if myErr, ok := Has[*MyError](err2); ok { /* This also matches! */ }
```

#### Prevents Common Bugs

Without `Has`, these mismatches silently fail:

```go
	err := &MyError{msg: "oops"}

	// With errors.As - this fails silently
	var myErr MyError
	if errors.As(err, &myErr) { /* false - type mismatch */ }

	// With Has - this succeeds
	if myErr, ok := Has[MyError](err); ok { /* true - automatic handling */ }
```

### When to Use `Has`

- You want to prevent pointer-value mismatch bugs
- You're unsure whether errors in your chain are wrapped as pointers or values
- You want the most robust error detection
- You need interface support with error embedding

### Limitations

Unlike `errors.As`, the interface type provided to `Has` must embed the `error` interface.

```go
  // OK with errors.As
  var temp interface{ Temporary() bool }
  if errors.As(err, &temp) && temp.Temporary() { /* handle temporary error */ }

  // With Has, the interface must embed `error`:
  if temp, ok := Has[interface { error; Temporary() bool }](err); ok && temp.Temporary() { /* handle temporary error */ }
```

## Migration Guide

### From `errors.As` to `HasError`

```go
	// Before
	var myErr *MyError
	if errors.As(err, &myErr) { return fmt.Errorf("unexpected MyError: %w", myErr) }

	// After
	if myErr, ok := HasError[*MyError](err); ok { return fmt.Errorf("unexpected MyError: %w", myErr) }
```

### From `HasError` to `Has`

If you're experiencing pointer-value mismatch issues:

```go
// If this sometimes fails unexpectedly
if myErr, ok := HasError[MyError](err); ok { /* ... */ }

// Try this instead
if myErr, ok := Has[MyError](err); ok { /* ... */ }
```

## License

This project is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.
