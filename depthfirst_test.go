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
	"errors"
	"slices"
	"testing"
)

// singleWrapError wraps a single error, implementing `Unwrap() error`.
type singleWrapError struct {
	msg string
	err error
}

func (e *singleWrapError) Error() string { return e.msg }
func (e *singleWrapError) Unwrap() error { return e.err }

// multiWrapError wraps multiple errors, implementing `Unwrap() []error`.
type multiWrapError struct {
	msg  string
	errs []error
}

func (e *multiWrapError) Error() string   { return e.msg }
func (e *multiWrapError) Unwrap() []error { return e.errs }

var (
	err1 = errors.New("err1")
	err2 = errors.New("err2")

	errGrand1 = errors.New("grandchild 1")
	errChild2 = errors.New("child 2")
)

func TestDepthFirstErrorTree(t *testing.T) {
	t.Parallel()

	t.Run("CorrectTraversalOrder", func(t *testing.T) {
		t.Parallel()

		// Build an error tree to verify depth-first traversal order.
		// A breadth-first search would visit errChild2 before errGrand1.
		//       root (multi)
		//       /      \
		// child1(single) child2
		//    |
		// grand1
		child1 := &singleWrapError{msg: "child 1", err: errGrand1}
		root := &multiWrapError{msg: "root", errs: []error{child1, errChild2}}
		expected := []error{root, child1, errGrand1, errChild2}

		if got := slices.Collect(DepthFirstErrorTree(root)); !slices.Equal(got, expected) {
			t.Errorf("DepthFirstErrorTree() traversal order incorrect\n got: %v\nwant: %v", got, expected)
		}
	})

	t.Run("NilRoot", func(t *testing.T) {
		t.Parallel()

		if got := slices.Collect(DepthFirstErrorTree(nil)); len(got) != 0 {
			t.Errorf("DepthFirstErrorTree(nil) should yield no errors, got %d", len(got))
		}
	})

	t.Run("NilInMultiError", func(t *testing.T) {
		t.Parallel()

		rootWithNil := &multiWrapError{msg: "root with nil", errs: []error{err1, nil, err2}}
		expected := []error{rootWithNil, err1, err2}

		if got := slices.Collect(DepthFirstErrorTree(rootWithNil)); !slices.Equal(got, expected) {
			t.Errorf("DepthFirstErrorTree() did not skip nil in multi-error\n got: %v\nwant: %v", got, expected)
		}
	})
}
