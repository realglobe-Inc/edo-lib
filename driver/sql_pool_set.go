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

	"github.com/realglobe-Inc/go-lib/erro"
)

// sql のコネクションプール集。
type SqlPoolSet struct {
	driverName string
	pools      map[string]*sql.DB
}

func NewSqlPoolSet(driverName string) *SqlPoolSet {
	return &SqlPoolSet{
		driverName,
		map[string]*sql.DB{},
	}
}

func (this *SqlPoolSet) Get(addr string) (*sql.DB, error) {
	pool := this.pools[addr]
	if pool != nil {
		return pool, nil
	}

	pool, err := sql.Open(this.driverName, addr)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	this.pools[addr] = pool
	return pool, nil
}

func (this *SqlPoolSet) Close() {
	for _, pool := range this.pools {
		pool.Close()
	}
}
