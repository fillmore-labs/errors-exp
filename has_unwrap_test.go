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
	"strconv"
	"testing"

	. "fillmore-labs.com/exp/errors"
)

func TestHasUnwrap(t *testing.T) {
	t.Parallel()

	errValue := MyUnwrapError(8)
	errPointer := &errValue

	testHasUnwrap(t, "value receiver", errValue)
	testHasUnwrap(t, "pointer receiver", errPointer)
}

func testHasUnwrap(t *testing.T, name string, err error) {
	t.Helper()

	t.Run(name, func(t *testing.T) {
		t.Parallel()

		t.Run("find MyPointerError", func(t *testing.T) {
			t.Parallel()

			if e, ok := Has[MyPointerError](err); !ok {
				t.Errorf("Expected to find MyPointerError, but didn't.")
			} else if e != MyPointerError(8) {
				t.Errorf("Expected MyPointerError(8), but got %d", int(e))
			}
		})

		t.Run("find *MyPointerError", func(t *testing.T) {
			t.Parallel()

			if e, ok := Has[*MyPointerError](err); !ok {
				t.Errorf("Expected to find *MyPointerError, but didn't.")
			} else if *e != MyPointerError(8) {
				t.Errorf("Expected *MyPointerError(8), but got %d", int(*e))
			}
		})

		t.Run("find MyValueError", func(t *testing.T) {
			t.Parallel()

			if e, ok := Has[MyValueError](err); !ok {
				t.Errorf("Expected to find MyValueError, but didn't.")
			} else if e != MyValueError(8) {
				t.Errorf("Expected MyValueError(8), but got %d", int(e))
			}
		})

		t.Run("find *MyValueError", func(t *testing.T) {
			t.Parallel()

			if e, ok := Has[*MyValueError](err); !ok {
				t.Errorf("Expected to find *MyValueError, but didn't.")
			} else if *e != MyValueError(8) {
				t.Errorf("Expected *MyValueError(8), but got %d", int(*e))
			}
		})
	})
}

func TestHasSimpleUnwrap(t *testing.T) {
	t.Parallel()

	err := func() error {
		err := MySimpleUnwrapError(8)
		return &err
	}()

	e1, ok := Has[MyValueError](err)

	if !ok {
		t.Errorf("Expected to find MyValueError, but didn't.")
	} else if e1 != MyValueError(8) {
		t.Errorf("Expected MyValueError(8), but got %d", int(e1))
	}

	e2, ok := Has[*MyValueError](err)

	if !ok {
		t.Errorf("Expected to find *MyValueError, but didn't.")
	} else if *e2 != MyValueError(8) {
		t.Errorf("Expected *MyValueError(8), but got %d", int(*e2))
	}
}

type (
	MyUnwrapError       int
	MySimpleUnwrapError int
)

func (e MyUnwrapError) Error() string {
	return "MyError" + strconv.Itoa(int(e))
}

func (e MyUnwrapError) Unwrap() []error {
	pe, ve := MyPointerError(e), MyValueError(e)
	return []error{nil, &pe, nil, ve, nil}
}

func (e MySimpleUnwrapError) Error() string {
	return "MyError" + strconv.Itoa(int(e))
}

func (e MySimpleUnwrapError) Unwrap() error {
	return MyValueError(e)
}

var _, _ error = MyUnwrapError(0), MySimpleUnwrapError(0)
