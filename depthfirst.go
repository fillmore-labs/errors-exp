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

import "iter"

// DepthFirstErrorTree traverses an error tree depth-first and returns a sequence of errors starting from the root error.
// It supports both single error unwrapping (`Unwrap() error`) and multi-error unwrapping (`Unwrap() []error`) mechanisms.
// Nil errors or nil results from unwrapping are skipped during traversal.
func DepthFirstErrorTree(root error) iter.Seq[error] {
	stack := []error{root}

	return func(yield func(error) bool) {
		for len(stack) > 0 {
			top := len(stack) - 1
			err := stack[top]
			stack = stack[:top]

			if err == nil {
				continue
			}

			if !yield(err) {
				return
			}

			switch x := err.(type) {
			case interface{ Unwrap() []error }:
				unwrap := x.Unwrap()
				// Push children in reverse order to visit them in their original order (depth-first).
				for i := len(unwrap) - 1; i >= 0; i-- {
					stack = append(stack, unwrap[i])
				}

			case interface{ Unwrap() error }:
				stack = append(stack, x.Unwrap())
			}
		}
	}
}
