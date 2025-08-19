# Errors

[![Go Reference](https://pkg.go.dev/badge/fillmore-labs.com/exp/errors.svg)](https://pkg.go.dev/fillmore-labs.com/exp/errors)
[![Test](https://github.com/fillmore-labs/errors-exp/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/fillmore-labs/errors-exp/actions/workflows/test.yml)
[![CodeQL](https://github.com/fillmore-labs/errors-exp/actions/workflows/github-code-scanning/codeql/badge.svg?branch=main)](https://github.com/fillmore-labs/errors-exp/actions/workflows/github-code-scanning/codeql)
[![Coverage](https://codecov.io/gh/fillmore-labs/errors-exp/branch/main/graph/badge.svg?token=LEOQNYK9KB)](https://codecov.io/gh/fillmore-labs/errors-exp)
[![Go Report Card](https://goreportcard.com/badge/fillmore-labs.com/exp/errors)](https://goreportcard.com/report/fillmore-labs.com/exp/errors)
[![License](https://img.shields.io/github/license/fillmore-labs/errors-exp)](https://www.apache.org/licenses/LICENSE-2.0)

`fillmore-labs.com/exp/errors` is an experimental Go library that provides two enhanced, generic alternatives to
`errors.As` for inspecting error chains with improved ergonomics and type safety.

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

`HasError` is a generic, type-safe replacement for `errors.As` that offers improved ergonomics and readability.

```go
func HasError[T error](err error) (T, bool)
```

### `HasError` Key Benefits

#### No Target Variable Needed

With `errors.As`, you must declare a variable beforehand:

```go
  var myErr *MyError
  if errors.As(err, &myErr) { /* ... use myErr */ }
```

`HasError` allows you to declare and check in one line:

```go
  if myErr, ok := HasError[*MyError](err); ok { /* ... use myErr */ }
```

#### Improved Readability

The syntax for `errors.As` can sometimes obscure intent, especially when you only need to check for an error's presence
without using its value. Note how `errors.As` requires a **pointer** to a struct literal to check for a
[**value** type](https://pkg.go.dev/crypto/x509#UnknownAuthorityError):

```go
  // This check is valid, but not immediately clear.
  if errors.As(err, &x509.UnknownAuthorityError{}) { /* ... */ }
```

`HasError` makes the intent explicit and is easier to read:

```go
  // The type is clearly specified as a generic parameter.
  if _, ok := HasError[x509.UnknownAuthorityError](err); ok { /* ... */ }
```

#### Interface Support

`HasError` works seamlessly with interfaces. To check if an error implements a specific interface:

```go
  // Single-line variant
  if e, ok := HasError[interface { error; Temporary() bool }](err); ok && e.Temporary() { /* handle temporary error */ }

  // Or, using a named interface:
  type temporary interface { error; Temporary() bool }
  if e, ok := HasError[temporary](err); ok && e.Temporary() { /* handle temporary error */ }
```

### When to Use `HasError`

Choose `HasError` when you need:

- **Improved ergonomics** compared to `errors.As`.
- **Strict type matching**, where pointers and values are treated as distinct types.
- To check if an error in the chain **implements a specific interface**.

## `Has` - Pointer-Value Flexibility

### `Has` Overview

`Has` provides all the benefits of `HasError` plus automatic handling of pointer-value type mismatches, preventing
subtle bugs.

```go
func Has[T error](err error) (T, bool)
```

### `Has` Key Benefits

#### Automatic Pointer-Value Matching

`Has` automatically resolves pointer-value mismatches, finding errors that `errors.As` would miss:

```go
  // Scenario 1: Looking for a value (MyError), but a pointer (*MyError) was wrapped.
  err := &MyError{msg: "oops"}
  if myErr, ok := Has[MyError](err); ok { /* This matches! */ }

  // Scenario 2: Looking for a pointer (*MyError), but a value (MyError) was wrapped.
  err2 := MyError{msg: "oops"}
  if myErr, ok := Has[*MyError](err2); ok { /* This also matches! */ }
```

#### Prevents Common Bugs

This mismatch would silently fail with `errors.As`:

```go
  key := []byte("My kung fu is better than yours")
  _, err := aes.NewCipher(key)

  // With errors.As - this check fails silently.
  var kse *aes.KeySizeError
  if errors.As(err, &kse) {
    fmt.Printf("Wrong AES key size: %d bytes.\n", *kse)
  }

  // With Has - the check succeeds.
  if kse, ok := Has[*aes.KeySizeError](err); ok {
    fmt.Printf("AES keys must be 16, 24, or 32 bytes long, got %d bytes.\n", *kse)
  }
```

### When to Use `Has`

Choose `Has` when you want:

- **The most robust error detection**, with all the ergonomic benefits of `HasError`.
- **Automatic handling of pointer-value mismatches** to prevent subtle bugs caused by inconsistent error wrapping.

### Limitations

Unlike `errors.As`, interface types provided to `Has` or `HasError` must embed the `error` interface.

```go
  // This is valid with errors.As:
  var temp interface{ Temporary() bool }
  if errors.As(err, &temp) && temp.Temporary() { /* handle temporary error */ }

  // With Has or HasError, the interface must embed `error`:
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

If you suspect pointer-value mismatches are causing issues, replace `HasError` with `Has`.

```go
  // If this check sometimes fails unexpectedly:
  if myErr, ok := HasError[MyError](err); ok { /* ... */ }

  // Switch to Has for more robust matching:
  if myErr, ok := Has[MyError](err); ok { /* ... */ }
```

## License

This project is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.
