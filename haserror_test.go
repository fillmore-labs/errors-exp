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

func TestHasError(t *testing.T) {
	t.Parallel()

	errVal := aes.KeySizeError(8)
	errPtr := &errVal

	t.Run("Match", func(t *testing.T) {
		t.Parallel()
		t.Run("Value", func(t *testing.T) {
			t.Parallel()

			e, ok := HasError[aes.KeySizeError](errVal)
			if !ok {
				t.Fatal("Expected to find aes.KeySizeError, but didn't")
			}

			if e != errVal {
				t.Errorf("Expected %v, but got %v", errVal, e)
			}
		})

		t.Run("Pointer", func(t *testing.T) {
			t.Parallel()

			e, ok := HasError[*aes.KeySizeError](errPtr)
			if !ok {
				t.Fatal("Expected to find *aes.KeySizeError, but didn't")
			}

			if e != errPtr {
				t.Errorf("Expected %v, but got %v", errPtr, e)
			}
		})
	})

	t.Run("Mismatch", func(t *testing.T) {
		t.Parallel()
		t.Run("Value", func(t *testing.T) {
			t.Parallel()

			if _, ok := HasError[aes.KeySizeError](errPtr); ok {
				t.Error("Expected to not find aes.KeySizeError, but did")
			}
		})

		t.Run("Pointer", func(t *testing.T) {
			t.Parallel()

			if _, ok := HasError[*aes.KeySizeError](errVal); ok {
				t.Error("Expected to not find *aes.KeySizeError, but did")
			}
		})
	})
}

func TestHasErrorAs(t *testing.T) {
	t.Parallel()

	err := MyAsError(8)
	want := MyValueError(8)

	if e, ok := HasError[MyValueError](err); !ok {
		t.Errorf("Expected to find MyValueError, but didn't.")
	} else if e != want {
		t.Errorf("Expected %v, but got %v", want, e)
	}
}
