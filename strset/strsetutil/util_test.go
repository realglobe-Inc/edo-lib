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
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	if s1, s2 := map[string]bool{}, New(); !reflect.DeepEqual(s2, s1) {
		t.Error(s2)
		t.Fatal(s1)
	} else if s1, s2 := map[string]bool{"a": true}, New("a"); !reflect.DeepEqual(s2, s1) {
		t.Error(s2)
		t.Fatal(s1)
	} else if s1, s2 := map[string]bool{"a": true, "b": true}, New("a", "b"); !reflect.DeepEqual(s2, s1) {
		t.Error(s2)
		t.Fatal(s1)
	}
}

func TestContains(t *testing.T) {
	if s1, s2 := New(), New(); !Contains(s1, s2) {
		t.Error(s1)
		t.Fatal(s2)
	} else if s1, s2 := New("a"), New(); !Contains(s1, s2) {
		t.Error(s1)
		t.Fatal(s2)
	} else if s1, s2 := New("a"), New("a"); !Contains(s1, s2) {
		t.Error(s1)
		t.Fatal(s2)
	} else if s1, s2 := New("a", "b"), New("a"); !Contains(s1, s2) {
		t.Error(s1)
		t.Fatal(s2)
	} else if s1, s2 := New("a"), New("a", "b"); Contains(s1, s2) {
		t.Error(s1)
		t.Fatal(s2)
	}
}
