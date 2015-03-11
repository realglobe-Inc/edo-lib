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
	"time"
)

// 並列使用時に便利なメソッドを持つ制限時間付きデータ用コンテナ。
type ConcurrentVolatileKeyValueStore interface {
	VolatileKeyValueStore

	// 以下、key と eKey の範囲が被っていた場合の挙動は保証しない。
	// エントリを返す。
	Entry(eKey string) (eVal string, err error)
	// エントリを設定する。
	SetEntry(eKey, eVal string, eExpDate time.Time) error
	// エントリを設定しつつ、値を返す。
	GetAndSetEntry(key string, caStmp *Stamp, eKey, eVal string, eExpiDate time.Time) (val interface{}, newCaStmp *Stamp, err error)
	// eVal がエントリに設定されていれば、値を設定する。
	PutIfEntered(key string, val interface{}, expiDate time.Time, eKey, eVal string) (entered bool, newCaStmp *Stamp, err error)
}
