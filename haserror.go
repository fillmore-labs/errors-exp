// Copyright 2025 Oliver Eikemeier. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package errors

import goerrors "errors"

// HasError finds the first error in `err`'s tree that has type `T`, and if one is found,
// returns that error and true. Otherwise, it returns nil and false.
//
// The tree consists of `err` itself, followed by the errors obtained by repeatedly
// calling its `Unwrap() error` or `Unwrap() []error` method. When `err` wraps multiple
// errors, `HasError` examines `err` followed by a depth-first traversal of its children.
//
// An error has the type `T` if the error's concrete value is assignable to `T`, or if
// the error has a method As(any) bool such that As(*T) returns true. In the latter case,
// the `As` method is responsible for the result.
//
// An error type might provide an `As` method, so it can be treated as if it were a
// different error type.
func HasError[T error](err error) (T, bool) {
	var target T
	ok := goerrors.As(err, &target)

	return target, ok
}
