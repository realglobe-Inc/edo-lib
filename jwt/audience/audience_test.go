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

package audience

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	if aud := New("abcde"); !reflect.DeepEqual(aud, Audience{"abcde": true}) {
		t.Error(aud)
		t.Fatal(Audience{"abcde": true})
	} else if aud := New("abcde", "fghij"); !reflect.DeepEqual(aud, Audience{"abcde": true, "fghij": true}) {
		t.Error(aud)
		t.Fatal(Audience{"abcde": true, "fghij": true})
	}
}

func TestMarshalSingle(t *testing.T) {
	aud := Audience{"abcde": true}
	if data, err := json.Marshal(aud); err != nil {
		t.Fatal(err)
	} else if string(data) != `"abcde"` {
		t.Error(aud)
		t.Fatal(string(data))
	}
}

func TestMarshal(t *testing.T) {
	aud := Audience{"abcde": true, "fghij": true}
	if data, err := json.Marshal(aud); err != nil {
		t.Fatal(err)
	} else if string(data) != `["abcde","fghij"]` &&
		string(data) != `["fghij","abcde"]` {
		t.Error(aud)
		t.Fatal(string(data))
	}
}

func TestUnmarshalSingle(t *testing.T) {
	var aud Audience
	if err := json.Unmarshal([]byte(`"abcde"`), &aud); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(aud, Audience{"abcde": true}) {
		t.Error(`"abcde"`)
		t.Fatal(aud)
	}
}

func TestUnmarshal(t *testing.T) {
	var aud Audience
	if err := json.Unmarshal([]byte(`["abcde","fghij"]`), &aud); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(aud, Audience{"abcde": true, "fghij": true}) {
		t.Error(`["abcde","fghij"]`)
		t.Fatal(aud)
	}
}
