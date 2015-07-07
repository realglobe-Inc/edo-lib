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

// JSON にしたときに 72h3m0.5s みたいな文字列になる time.Duration のラッパー。
package duration

import (
	"encoding/json"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Duration time.Duration

func (this Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(this).String())
}

func (this *Duration) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*this = Duration(d)
	return nil
}

func (this Duration) GetBSON() (interface{}, error) {
	return time.Duration(this).String(), nil
}

func (this *Duration) SetBSON(raw bson.Raw) error {
	var s string
	if err := raw.Unmarshal(&s); err != nil {
		return err
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*this = Duration(d)
	return nil
}
