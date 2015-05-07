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

package strset

import (
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
)

// JSON にしたときに要素の配列になる文字列集合型。
type StringSet map[string]bool

func (this StringSet) MarshalJSON() ([]byte, error) {
	if this == nil {
		return json.Marshal(nil)
	}

	a := []string{}
	for elem, ok := range this {
		if ok {
			a = append(a, elem)
		}
	}
	return json.Marshal(a)
}

func (this *StringSet) UnmarshalJSON(data []byte) error {
	var a []string
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	} else if a == nil {
		return nil
	}

	s := map[string]bool{}
	for _, elem := range a {
		s[elem] = true
	}
	*this = StringSet(s)
	return nil
}

func (this StringSet) GetBSON() (interface{}, error) {
	if this == nil {
		return nil, nil
	}

	a := []string{}
	for elem, ok := range this {
		if ok {
			a = append(a, elem)
		}
	}
	return a, nil
}

func (this *StringSet) SetBSON(raw bson.Raw) error {
	var a []string
	if err := raw.Unmarshal(&a); err != nil {
		return err
	} else if a == nil {
		return nil
	}

	s := map[string]bool{}
	for _, elem := range a {
		s[elem] = true
	}
	*this = StringSet(s)
	return nil
}
