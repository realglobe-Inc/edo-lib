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

package prand

import (
	"encoding/base64"
	"math/rand"
	"sync"
	"time"
)

// 安全な乱数が使えない場合の代替。
type Random struct {
	interval time.Duration

	lock sync.Mutex
	exp  time.Time
	rand *rand.Rand
}

// interval: 乱数種を切り替える間隔。
func New(interval time.Duration) *Random {
	now := time.Now()
	return &Random{
		interval: interval,
		exp:      now.Add(interval),
		rand:     rand.New(rand.NewSource(now.UnixNano())),
	}
}

func (this *Random) Bytes(length int) []byte {
	now := time.Now()
	func() {
		this.lock.Lock()
		defer this.lock.Unlock()
		if now.After(this.exp) {
			this.exp = now.Add(this.interval)
			this.rand = rand.New(rand.NewSource(now.UnixNano()))
		}
	}()

	buff := make([]byte, length)
	for i := 0; i < len(buff); i++ {
		buff[i] = byte(this.rand.Intn(256))
	}
	return buff
}

func (this *Random) String(length int) string {
	return base64.URLEncoding.EncodeToString(this.Bytes((length*6 + 7) / 8))[:length]
}
