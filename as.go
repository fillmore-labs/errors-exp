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

// As finds the first error in `err`'s tree that has type `T`, and if one is found,
// sets target to that error value and returns true. Otherwise, it returns false.
//
// The tree consists of `err` itself, followed by the errors obtained by repeatedly
// calling its `Unwrap() error` or `Unwrap() []error` method. When `err` wraps multiple
// errors, `As` examines `err` followed by a depth-first traversal of its children.
//
// An error has the type `T` if the error's value is assignable to `T`, including
// cases where there are pointer-value mismatches (e.g., if `T` is `*MyError` but
// the found error is `MyError`, or vice versa), or if the error has a method
// `As(any) bool` that returns `true`. To accommodate pointer-value mismatches
// in `As` implementations, `As` tries different variations of the target type.
// In the latter case, the `As` method is responsible for setting `target`.
//
// An error type might provide an `As` method, so it can be treated as if it were a
// different error type.
//
// As panics if `target` is a nil pointer.
func As[T error](err error, target *T) bool {
	if target == nil {
		panic("errors: target cannot be nil")
	}

	var handler altHandler[T]

	for err := range DepthFirstErrorTree(err) {
		if result, ok := err.(T); ok {
			*target = result

			return true
		}

		if handler == nil {
			// Lazily initialize the handler only when a direct type assertion fails.
			handler = newAltHandler(target)
		}

		if result, ok := handler.handleAssert(err); ok {
			*target = result

			return true
		}

		if x, ok := err.(interface{ As(any) bool }); ok {
			// First, try the standard errors.As contract. This works when T matches
			// the type expected by the As method.
			if x.As(target) {
				return true
			}

			// If the standard call fails, it might be due to a pointer-vs-value mismatch
			// between T and the type the As method is designed to handle.
			if result, ok := handler.handleAs(x); ok {
				*target = result

				return true
			}
		}
	}

	return false
}
