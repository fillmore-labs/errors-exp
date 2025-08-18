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
	"crypto/aes"
	"testing"

	. "fillmore-labs.com/exp/errors"
)

func TestHasErrorMatch(t *testing.T) {
	t.Parallel()

	err := func() error {
		err := aes.KeySizeError(8)
		return err
	}()

	e1, ok := HasError[aes.KeySizeError](err)

	if !ok {
		t.Errorf("Expected to find aes.KeySizeError, but didn't.")
	} else if e1 != aes.KeySizeError(8) {
		t.Errorf("Expected aes.KeySizeError(8), but got %d", int(e1))
	}

	err2 := func() error {
		err := aes.KeySizeError(8)
		return &err
	}()

	e2, ok := HasError[*aes.KeySizeError](err2)

	if !ok {
		t.Errorf("Expected to find *aes.KeySizeError, but didn't.")
	} else if *e2 != aes.KeySizeError(8) {
		t.Errorf("Expected *aes.KeySizeError(8), but got %d", int(*e2))
	}
}

func TestHasErrorMismatch(t *testing.T) {
	t.Parallel()

	err := func() error {
		err := aes.KeySizeError(8)
		return &err
	}()

	_, ok := HasError[aes.KeySizeError](err)

	if ok {
		t.Errorf("Expected to not find aes.KeySizeError, but did.")
	}

	err2 := func() error {
		err := aes.KeySizeError(8)
		return err
	}()

	_, ok = HasError[*aes.KeySizeError](err2)

	if ok {
		t.Errorf("Expected to not find *aes.KeySizeError, but did.")
	}
}
