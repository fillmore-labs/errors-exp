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
	"errors"
	"fmt"

	errorx "fillmore-labs.com/exp/errors"
)

func ExampleHas() {
	key := []byte("My kung fu is better than yours")
	_, err := aes.NewCipher(key)

	// With errors.As - this check fails silently.
	var target *aes.KeySizeError
	if errors.As(err, &target) {
		fmt.Printf("Wrong AES key size: %d bytes.\n", *target)
	}

	// With Has - the check succeeds.
	if kse, ok := errorx.Has[*aes.KeySizeError](err); ok {
		fmt.Printf("AES keys must be 16, 24, or 32 bytes long, got %d bytes.\n", *kse)
	}
	// Output: AES keys must be 16, 24, or 32 bytes long, got 31 bytes.
}
