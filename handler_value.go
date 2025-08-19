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

// valueHandler gets chosen when the queried error type T is a pointer.
// It handles finding an error value when querying for an error pointer type
// and allocates a new pointer and copies the value.
type valueHandler[T error] struct {
	noneHandler[T]
	altType reflect.Type // alternative value type to a pointer type T = *altType
	ptr     any
}

func (h *valueHandler[T]) handleAssert(err error) (T, bool) {
	if !reflect.TypeOf(err).AssignableTo(h.altType) {
		return h.zero()
	}

	// Handle the case where T is a pointer type, but err is a value type.
	ptr := reflect.New(h.altType)        // Create a new pointer to a zero value of the error's type.
	ptr.Elem().Set(reflect.ValueOf(err)) // Copy the error value into the pointed-to value.

	// Type asserts the pointer to T. This is safe because T is a pointer to h.altType,
	// which is what `reflect.New` creates.
	return typeAssert[T](ptr)
}

func (h *valueHandler[T]) handleAs(x interface{ As(any) bool }) (T, bool) {
	// Here, T is a pointer type (*altType). Some `As` implementations might
	// be designed to populate a value (altType), so they expect a pointer to
	// that value (*altType).
	if h.ptr == nil { // Create a new (non-nil) pointer to a zero value of the error's type, *altType = T.
		h.ptr = reflect.New(h.altType).Interface()
	}

	if x.As(h.ptr) { // And pass that as a target.
		return h.ptr.(T), true // We can then assert the (non-nil) pointer to T.
	}

	return h.zero()
}
