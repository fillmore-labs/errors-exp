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

package errors_test

import (
	"testing"

	. "fillmore-labs.com/exp/errors"
)

func TestAsAs(t *testing.T) {
	t.Parallel()

	errValue := MyAsError(8)
	errPointer := &errValue

	testAsAs(t, "value receiver", errValue)
	testAsAs(t, "pointer receiver", errPointer)
}

func testAsAs(t *testing.T, name string, err error) {
	t.Helper()

	t.Run(name, func(t *testing.T) {
		t.Parallel()

		t.Run("find MyPointerError", func(t *testing.T) {
			t.Parallel()

			var e MyPointerError
			if !As(err, &e) {
				t.Errorf("Expected to find MyPointerError, but didn't.")
			} else if e != MyPointerError(8) {
				t.Errorf("Expected MyPointerError(8), but got %d", int(e))
			}
		})

		t.Run("find *MyPointerError", func(t *testing.T) {
			t.Parallel()

			var e *MyPointerError
			if !As(err, &e) {
				t.Errorf("Expected to find *MyPointerError, but didn't.")
			} else if *e != MyPointerError(8) {
				t.Errorf("Expected *MyPointerError(8), but got %d", int(*e))
			}
		})

		t.Run("find MyValueError", func(t *testing.T) {
			t.Parallel()

			var e MyValueError
			if !As(err, &e) {
				t.Errorf("Expected to find MyValueError, but didn't.")
			} else if e != MyValueError(8) {
				t.Errorf("Expected MyValueError(8), but got %d", int(e))
			}
		})

		t.Run("find *MyValueError", func(t *testing.T) {
			t.Parallel()

			var e *MyValueError
			if !As(err, &e) {
				t.Errorf("Expected to find *MyValueError, but didn't.")
			} else if *e != MyValueError(8) {
				t.Errorf("Expected *MyValueError(8), but got %d", int(*e))
			}
		})
	})
}
