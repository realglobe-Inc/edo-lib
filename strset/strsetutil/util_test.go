// Copyright 2015 realglobe, Inc.
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

package strsetutil

import (
	"testing"
)

func newSet(elems ...string) map[string]bool {
	s := map[string]bool{}
	for _, elem := range elems {
		s[elem] = true
	}
	return s
}

func TestContains(t *testing.T) {
	if s1, s2 := newSet(), newSet(); !Contains(s1, s2) {
		t.Error(s1)
		t.Fatal(s2)
	} else if s1, s2 := newSet("a"), newSet(); !Contains(s1, s2) {
		t.Error(s1)
		t.Fatal(s2)
	} else if s1, s2 := newSet("a"), newSet("a"); !Contains(s1, s2) {
		t.Error(s1)
		t.Fatal(s2)
	} else if s1, s2 := newSet("a", "b"), newSet("a"); !Contains(s1, s2) {
		t.Error(s1)
		t.Fatal(s2)
	} else if s1, s2 := newSet("a"), newSet("a", "b"); Contains(s1, s2) {
		t.Error(s1)
		t.Fatal(s2)
	}
}
