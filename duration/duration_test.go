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

package duration

import (
	"encoding/json"
	"reflect"
	"strconv"
	"testing"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func TestJson(t *testing.T) {
	a := Duration(time.Duration(time.Now().UnixNano()) % (24 * time.Hour))

	buff, err := json.Marshal(a)
	if err != nil {
		t.Fatal(err)
	} else if buff[0] != '"' {
		// JSON 文字列じゃない。
		t.Fatal(string(buff))
	}

	var b Duration
	if err := json.Unmarshal(buff, &b); err != nil {
		t.Fatal(err, string(buff))
	}

	if b != a {
		t.Fatal(b, a, string(buff))
	}
}

// 何かの中に入ってても大丈夫か。
func TestNestedJson(t *testing.T) {
	type testType struct {
		D Duration
	}

	var a testType
	a.D = Duration(time.Duration(time.Now().UnixNano()) % (24 * time.Hour))

	buff, err := json.Marshal(&a)
	if err != nil {
		t.Fatal(err)
	}

	var b testType
	if err := json.Unmarshal(buff, &b); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(b, a) {
		t.Fatal(b, a, string(buff))
	}
}

// テストするなら、mongodb を立てる必要あり。
// 立ってなかったらテストはスキップ。
var monPool, _ = mgo.DialWithTimeout("localhost", time.Minute)

func init() {
	if monPool != nil {
		monPool.SetSyncTimeout(time.Minute)
	}
}

const (
	test_coll = "test-collection"
)

func TestBson(t *testing.T) {
	if monPool == nil {
		t.SkipNow()
	}

	test_db := "test-db-" + strconv.FormatInt(time.Now().UnixNano(), 16)
	conn := monPool.New()
	defer conn.Close()

	type testType struct {
		K string   `bson:"key"`
		D Duration `bson:"duration"`
	}

	var a testType
	a.K = strconv.FormatInt(time.Now().UnixNano(), 16)
	a.D = Duration(time.Duration(time.Now().UnixNano()) % (24 * time.Hour))

	if err := conn.DB(test_db).C(test_coll).Insert(&a); err != nil {
		t.Fatal(err)
	}
	defer conn.DB(test_db).DropDatabase()

	var b testType
	if err := conn.DB(test_db).C(test_coll).Find(bson.M{"key": a.K}).One(&b); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(b, a) {
		t.Fatal(b, a)
	}
}
