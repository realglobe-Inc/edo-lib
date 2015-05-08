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

package json

import (
	"encoding/json"
	"testing"
)

func TestStringEscape(t *testing.T) {
	s := ` !"#$%'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_\` +
		"`" + `abcdefghijklmnopqrstuvwxyz{|}~` +
		`いろは` + "\n\r\t\b\f"

	data1, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	data2 := `"` + StringEscape(s) + `"`

	var s1, s2 string
	if err := json.Unmarshal(data1, &s1); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal([]byte(data2), &s2); err != nil {
		t.Fatal(err)
	}

	if s2 != s1 {
		t.Error(s)
		t.Error(s2)
		t.Fatal(s1)
	}
}
