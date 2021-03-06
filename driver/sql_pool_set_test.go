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
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

// テストするなら、mysql を立てる必要あり。
// 立ってなかったらテストはスキップ。
func TestSqlPoolSet(t *testing.T) {
	addr := "root@tcp(localhost:3306)"
	var sqlPool, _ = sql.Open("mysql", addr)
	if sqlPool == nil {
		t.SkipNow()
	}

	poolSet := NewSqlPoolSet("mysql")
	defer poolSet.Close()

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
