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

package driver

import (
	"encoding/json"
	"github.com/realglobe-Inc/go-lib/erro"
	"reflect"
	"time"
)

const (
	testStaleDur  = 0
	testCaExpiDur = 0
)

const (
	testLabel = "edo-test"
	testKey   = "test-key"
)

var testData = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

var testVal = map[string]interface{}{
	"array":  []interface{}{"elem-1", "elem-2"},
	"date":   time.Now(),
	"digest": "xyz",
}

// JSON を通して等しいかどうか調べる。
func jsonEqual(v1 interface{}, v2 interface{}) (equal bool) {
	b1, err := json.Marshal(v1)
	if err != nil {
		log.Err(erro.Wrap(err))
		return false
	}
	var w1 interface{}
	if err := json.Unmarshal(b1, &w1); err != nil {
		log.Err(erro.Wrap(err))
		return false
	}

	b2, err := json.Marshal(v2)
	if err != nil {
		log.Err(erro.Wrap(err))
		return false
	}
	var w2 interface{}
	if err := json.Unmarshal(b2, &w2); err != nil {
		log.Err(erro.Wrap(err))
		return false
	}

	return reflect.DeepEqual(w1, w2)
}

func jsonUnmarshal(data []byte) (interface{}, error) {
	var res interface{}
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, erro.Wrap(err)
	}
	return res, nil
}

// 剰余切り捨て。
func cutOff(val, thres int64) int64 {
	return val - val%thres
}
