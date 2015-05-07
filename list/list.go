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
)

// JSON にしたときに配列になる list。
type List list.List

func (this *List) MarshalJSON() ([]byte, error) {
	if this == nil {
		return json.Marshal(nil)
	}

	a := []interface{}{}
	l := (*list.List)(this)
	for elem := l.Front(); elem != nil; elem = elem.Next() {
		a = append(a, elem.Value)
	}
	return json.Marshal(a)
}

func (this *List) UnmarshalJSON(data []byte) error {
	var a []interface{}
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}
	if a == nil {
		return nil
	}

	l := list.New()
	for _, val := range a {
		l.PushBack(val)
	}
	*this = List(*l)
	return nil
}
