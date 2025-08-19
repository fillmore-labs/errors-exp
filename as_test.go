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

func TestAsMatch(t *testing.T) {
	t.Parallel()

	err := func() error {
		err := MyValueError(8)
		return err
	}()

	var e1 MyValueError
	ok := As(err, &e1)

	if !ok {
		t.Errorf("Expected to find MyValueError, but didn't.")
	} else if e1 != MyValueError(8) {
		t.Errorf("Expected MyValueError(8), but got %d", int(e1))
	}

	err2 := func() error {
		err := MyPointerError(8)
		return &err
	}()

	var e2 *MyPointerError
	ok = As(err2, &e2)

	if !ok {
		t.Errorf("Expected to find *MyPointerError, but didn't.")
	} else if *e2 != MyPointerError(8) {
		t.Errorf("Expected *MyPointerError(8), but got %d", int(*e2))
	}
}

func TestAsMismatch(t *testing.T) {
	t.Parallel()

	err := func() error {
		err := MyPointerError(8)
		return &err
	}()

	var e1 MyPointerError
	ok := As(err, &e1)

	if !ok {
		t.Errorf("Expected to find MyPointerError, but didn't.")
	} else if e1 != MyPointerError(8) {
		t.Errorf("Expected MyPointerError(8), but got %d", int(e1))
	}

	err2 := func() error {
		err := MyValueError(8)
		return err
	}()

	var e2 *MyValueError
	ok = As(err2, &e2)

	if !ok {
		t.Errorf("Expected to find *MyValueError, but didn't.")
	} else if *e2 != MyValueError(8) {
		t.Errorf("Expected *MyValueError(8), but got %d", int(*e2))
	}
}

func TestAsNil(t *testing.T) {
	t.Parallel()

	var e MyValueError
	ok := As((*MyValueError)(nil), &e)

	if ok {
		t.Errorf("Expected to not find MyValueError in nil, but did.")
	}
}

func TestAsNotFound(t *testing.T) {
	t.Parallel()

	err := func() error {
		err := MyPointerError(8)
		return &err
	}()

	var e *MyPointerOnlyError
	ok := As(err, &e)

	if ok {
		t.Errorf("Expected to not find *MyPointerOnlyError in *MyPointerError, but did.")
	}
}

func TestAs_PanicsWithNilTarget(t *testing.T) {
	t.Parallel()

	defer func() {
		// We expect a panic, so we recover to prevent the test runner from
		// marking the test as failed.
		_ = recover()
	}()

	// Call As with a nil target to trigger the panic.
	As(MyValueError(1), (*MyValueError)(nil))

	// If this line is reached, the function did not panic, which is an error.
	t.Error("Expected As to panic with a nil target, but it did not")
}
