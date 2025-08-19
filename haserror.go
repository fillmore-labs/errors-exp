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

// HasError finds the first error in `err`'s tree that is of type `T`.
// If a matching error is found, it is returned along with `true`.
// Otherwise, the zero value for `T` (`nil` in case of a pointer type)
// and `false` are returned.
//
// This function provides a generic, type-safe alternative to the standard library's `errors.As`.
//
// The error tree is traversed depth-first, starting with `err` itself.
// The tree is explored by repeatedly calling `Unwrap() error` or `Unwrap() []error`.
//
// An error is considered to be of type `T` if:
//   - The error's concrete value is assignable to `T`.
//   - The error has a method `As(any) bool`, and `As(&T)` returns `true`. In this case,
//     the `As` method is responsible for the result.
func HasError[T error](err error) (T, bool) {
	for err := range DepthFirstErrorTree(err) {
		if target, ok := err.(T); ok {
			return target, true
		}

		if x, ok := err.(interface{ As(any) bool }); ok {
			// Try the standard errors.As contract.
			var target T
			if x.As(&target) {
				return target, true
			}
		}
	}

	var zero T

	return zero, false
}
