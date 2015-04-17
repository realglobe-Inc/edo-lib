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

package list

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestJson(t *testing.T) {
	for a := New(); a.Len() < 10; a.PushFront(float64(a.Len())) {

		buff, err := json.Marshal(a)
		if err != nil {
			t.Error(err)
			return
		} else if buff[0] != '[' {
			// JSON 配列じゃない。
			t.Error(string(buff))
		}

		var b List
		if err := json.Unmarshal(buff, &b); err != nil {
			t.Fatal(err, string(buff))
		}

		if !reflect.DeepEqual(&b, a) {
			t.Error(string(buff))
			t.Error(&b)
			t.Error(a)
		}
	}
}
