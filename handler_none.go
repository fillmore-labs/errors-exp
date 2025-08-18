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

// noneHandler gets chosen when the queried error type has no alternate form.
type noneHandler[T error] struct{}

func (noneHandler[T]) zero() (T, bool) {
	var zero T

	return zero, false
}

var _ altHandler[error] = noneHandler[error]{}

func (q noneHandler[T]) handleAssert(_ error) (T, bool) {
	return q.zero()
}

func (q noneHandler[T]) handleAs(_ interface{ As(any) bool }) (T, bool) {
	return q.zero()
}
