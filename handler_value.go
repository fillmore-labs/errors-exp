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

//go:build go1.25

package errors

import "reflect"

// valueHandler gets chosen when the queried error type T is a pointer.
// It handles cases where a found error is a value of the pointer's element type.
type valueHandler[T error] struct {
	noneHandler[T]
	altType reflect.Type // alternative value type to a pointer type T = *altType
}

func (h valueHandler[T]) handleAssert(err error) (T, bool) {
	if !reflect.TypeOf(err).AssignableTo(h.altType) {
		return h.zero()
	}

	// pointer-value mismatch
	errv := reflect.ValueOf(err)

	// Handle the case where T is a pointer type, but err is a value type.
	ptr := reflect.New(h.altType) // Create a new pointer to a zero value of the error's type.
	ptr.Elem().Set(errv)          // Copy the error value into the new pointer.

	// Type asserts the pointer to T. This is safe because T is a pointer to h.altType,
	// which is what `reflect.New` creates.
	t, _ := reflect.TypeAssert[T](ptr)

	return t, true
}

func (h valueHandler[T]) handleAs(x interface{ As(any) bool }) (T, bool) {
	// Here, T is a pointer type (e.g., *MyError). Some `As` implementations might
	// be designed to populate a value (MyError), so they expect a pointer to
	// that value (*MyError). This handler creates a pointer to the underlying
	// value type (h.altType) and passes it to As.
	ptr := reflect.New(h.altType) // We create a non-nil pointer to the underlying value type and pass that.
	if x.As(ptr.Interface()) {    // If `As` succeeds, we can then assert the pointer to T.
		t, _ := reflect.TypeAssert[T](ptr)

		return t, true
	}

	return h.zero()
}
