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
	ptr     *T
}

func (h *pointerHandler[T]) handleAssert(err error) (T, bool) {
	if !reflect.TypeOf(err).AssignableTo(h.altType) {
		return h.zero()
	}

	// Handle the case where T is a value type, but err is a pointer type.
	val := reflect.ValueOf(err)
	if val.IsNil() { // Found a (*T)(nil) error
		return h.zero()
	}

	// Dereference the pointer, and assert to T.
	return typeAssert[T](val.Elem())
}

func (h *pointerHandler[T]) handleAs(x interface{ As(any) bool }) (T, bool) {
	if h.ptr == nil {
		h.ptr = new(T)
	}

	// When T is a value type (e.g., MyError), some `As` implementations might
	// expect to populate a pointer to *T, requiring a pointer-to-pointer argument
	// (e.g., target **MyError).
	if x.As(&h.ptr) && h.ptr != nil {
		// If `As` succeeds, ptr is now a valid *T.
		// We dereference it to return the value type T.
		return *h.ptr, true
	}

	return h.zero()
}
