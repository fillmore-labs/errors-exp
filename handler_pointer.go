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

import "reflect"

// pointerHandler gets chosen when the queried error type T is a value.
// It handles cases where a found error is a pointer to that value type.
type pointerHandler[T error] struct {
	noneHandler[T]
	altType reflect.Type // alternative pointer type to a value type T, *T = altType
}

func (h pointerHandler[T]) handleAssert(err error) (T, bool) {
	if !reflect.TypeOf(err).AssignableTo(h.altType) {
		return h.zero()
	}

	// pointer-value mismatch
	errv := reflect.ValueOf(err)

	// Handle the case where T is a value type, but err is a pointer type.
	var ptr *T
	reflect.ValueOf(&ptr).Elem().Set(errv) // This is equivalent to `ptr = err.(*T)` but works with the generic type `T`.

	if ptr == nil { // Found a (*T)(nil) error
		return h.zero()
	}

	return *ptr, true // Dereference the pointer to return the value.
}

func (h pointerHandler[T]) handleAs(x interface{ As(any) bool }) (T, bool) {
	// When T is a value type (e.g., MyError), some `As` implementations might
	// expect to populate a pointer to *T, requiring a pointer-to-pointer argument
	// (e.g., target **MyError).
	var ptr *T
	if x.As(&ptr) && ptr != nil {
		// If `As` succeeds, ptr is now a valid *T.
		// We dereference it to return the value type T.
		return *ptr, true
	}

	return h.zero()
}
