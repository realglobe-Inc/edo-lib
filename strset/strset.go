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

// コピーするだけ。
func New(m map[string]bool) StringSet {
	s := map[string]bool{}
	for elem, ok := range m {
		if ok {
			s[elem] = true
		}
	}
	return StringSet(s)
}

func FromSlice(l []string) StringSet {
	s := map[string]bool{}
	for _, elem := range l {
		s[elem] = true
	}
	return StringSet(s)
}

func (this StringSet) MarshalJSON() ([]byte, error) {
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
	}
	s := map[string]bool{}
	for _, elem := range a {
		s[elem] = true
	}
	*this = StringSet(s)
	return nil
}

func (this StringSet) GetBSON() (interface{}, error) {
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
	}
	s := map[string]bool{}
	for _, elem := range a {
		s[elem] = true
	}
	*this = StringSet(s)
	return nil
}

func (this StringSet) Copy() StringSet {
	c := StringSet{}
	for elem, ok := range this {
		if ok {
			c[elem] = true
		}
	}
	return c
}

func (this StringSet) OneOf() string {
	return OneOf(this)
}

func OneOf(s StringSet) string {
	for elem := range s {
		return elem
	}
	return ""
}
