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
	"container/list"
	"encoding/json"
	"reflect"
	"testing"
)

func TestJson(t *testing.T) {
	for a := list.New(); a.Len() < 10; a.PushFront(float64(a.Len())) {

		l := (*List)(a)
		buff, err := json.Marshal(l)
		if err != nil {
			t.Fatal(err)
		} else if buff[0] != '[' {
			// JSON 配列じゃない。
			t.Fatal(string(buff))
		}

		var l2 List
		if err := json.Unmarshal(buff, &l2); err != nil {
			t.Fatal(err, string(buff))
		}

		if !reflect.DeepEqual(&l2, l) {
			t.Error(a.Len())
			t.Error(string(buff))
			t.Error(&l2)
			t.Fatal(l)
		}
	}
}
