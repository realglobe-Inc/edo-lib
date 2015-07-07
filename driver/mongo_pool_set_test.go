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
	"testing"
	"time"

	"gopkg.in/mgo.v2"
)

// テストするなら、mongodb を立てる必要あり。
// 立ってなかったらテストはスキップ。
var monPool, _ = mgo.DialWithTimeout("localhost", time.Minute)

func init() {
	if monPool != nil {
		monPool.SetSyncTimeout(time.Minute)
	}
}

func TestMongoPoolSet(t *testing.T) {
	if monPool == nil {
		t.SkipNow()
	}

	poolSet := NewMongoPoolSet(time.Minute)
	defer poolSet.Close()

	addr := monPool.LiveServers()[0]

	if pool, err := poolSet.Get(addr); err != nil {
		t.Fatal(err)
	} else if pool2, err := poolSet.Get(addr); err != nil {
		t.Fatal(err)
	} else if pool2 != pool {
		t.Error("cannot reuse")
		t.Error(pool2)
		t.Fatal(pool)
	}
}
