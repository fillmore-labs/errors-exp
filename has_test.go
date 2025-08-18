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

func TestHasMatch(t *testing.T) {
	t.Parallel()

	err := func() error {
		err := MyValueError(8)
		return err
	}()

	e1, ok := Has[MyValueError](err)

	if !ok {
		t.Errorf("Expected to find MyValueError, but didn't.")
	} else if e1 != MyValueError(8) {
		t.Errorf("Expected MyValueError(8), but got %d", int(e1))
	}

	err2 := func() error {
		err := MyPointerError(8)
		return &err
	}()

	e2, ok := Has[*MyPointerError](err2)

	if !ok {
		t.Errorf("Expected to find *MyPointerError, but didn't.")
	} else if *e2 != MyPointerError(8) {
		t.Errorf("Expected *MyPointerError(8), but got %d", int(*e2))
	}
}

func TestHasMismatch(t *testing.T) {
	t.Parallel()

	err := func() error {
		err := MyPointerError(8)
		return &err
	}()

	e1, ok := Has[MyPointerError](err)

	if !ok {
		t.Errorf("Expected to find MyPointerError, but didn't.")
	} else if e1 != MyPointerError(8) {
		t.Errorf("Expected MyPointerError(8), but got %d", int(e1))
	}

	err2 := func() error {
		err := MyValueError(8)
		return err
	}()

	e2, ok := Has[*MyValueError](err2)

	if !ok {
		t.Errorf("Expected to find *MyValueError, but didn't.")
	} else if *e2 != MyValueError(8) {
		t.Errorf("Expected *MyValueError(8), but got %d", int(*e2))
	}
}

func TestHasNil(t *testing.T) {
	t.Parallel()

	_, ok := Has[MyValueError]((*MyValueError)(nil))

	if ok {
		t.Errorf("Expected to not find MyValueError in nil, but did.")
	}
}

func TestHasNotFound(t *testing.T) {
	t.Parallel()

	err := func() error {
		err := MyPointerError(8)
		return &err
	}()

	_, ok := Has[*MyPointerOnlyError](err)

	if ok {
		t.Errorf("Expected to not find *MyPointerOnlyError in *MyPointerError, but did.")
	}
}

type (
	MyPointerError     int
	MyValueError       int
	MyPointerOnlyError struct{ v int }
)

func (e MyPointerError) Error() string {
	return "MyPointerError" + strconv.Itoa(int(e))
}

func (e MyPointerError) As(other any) bool {
	if o, ok := other.(*MyPointerError); ok {
		return int(e) == int(*o)
	}

	return false
}

func (e MyValueError) Error() string {
	return "MyValueError" + strconv.Itoa(int(e))
}

func (e *MyPointerOnlyError) Error() string {
	return "MyPointerOnlyError" + strconv.Itoa(e.v)
}

var _, _, _ error = (*MyPointerError)(nil), MyValueError(0), (*MyPointerOnlyError)(nil)
