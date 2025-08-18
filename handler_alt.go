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

import (
	"reflect"
)

type altHandler[T error] interface {
	handleAssert(err error) (T, bool)
	handleAs(x interface{ As(any) bool }) (T, bool)
}

func newAltHandler[T error]() altHandler[T] {
	targetType := reflect.TypeFor[T]()

	isPointerType := targetType.Kind() == reflect.Pointer

	var altType reflect.Type
	if isPointerType {
		// T is a Pointer (e.g., `*MyError`), so maybe value errors `MyError` match
		altType = targetType.Elem()
	} else {
		// T is a value (e.g., `MyError`), so maybe pointer errors `*MyError` match
		altType = reflect.PointerTo(targetType)
	}

	if !altType.Implements(errorType) {
		// Do not look for non-error alternatives
		return noneHandler[T]{}
	}

	if isPointerType {
		// altType is a value type.
		// handle value alternatives for the queried pointer type
		return valueHandler[T]{altType: altType}
	}

	// altType is a pointer type.
	// handle pointer alternatives for the queried value type
	return pointerHandler[T]{altType: altType}
}

var errorType = reflect.TypeOf((*error)(nil)).Elem()
