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

// JWT の aud クレーム。
package audience

import (
	"encoding/json"

	"github.com/realglobe-Inc/edo-lib/strset"
	"github.com/realglobe-Inc/go-lib/erro"
)

// JSON にしたときに、適切に単文字列か文字列配列になる。
type Audience map[string]bool

func New(aud ...string) Audience {
	s := map[string]bool{}
	for _, a := range aud {
		s[a] = true
	}
	return s
}

func (this Audience) MarshalJSON() ([]byte, error) {
	switch len(this) {
	case 1:
		for k := range this {
			return json.Marshal(k)
		}
	default:
		return json.Marshal(strset.Set(this))
	}
	panic("logic error")
}

func (this *Audience) UnmarshalJSON(data []byte) error {
	if data[0] == '[' {
		var buff strset.Set
		if err := json.Unmarshal(data, &buff); err != nil {
			return erro.Wrap(err)
		}
		*this = Audience(buff)
	} else {
		var buff string
		if err := json.Unmarshal(data, &buff); err != nil {
			return erro.Wrap(err)
		}
		*this = map[string]bool{buff: true}
	}
	return nil
}
