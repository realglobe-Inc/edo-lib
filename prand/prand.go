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
type Generator struct {
	intrv time.Duration

	lock sync.Mutex
	exp  time.Time
	rand *rand.Rand
}

// intrv: 乱数種を切り替える間隔。
func New(intrv time.Duration) *Generator {
	now := time.Now()
	return &Generator{
		intrv: intrv,
		exp:   now.Add(intrv),
		rand:  rand.New(rand.NewSource(now.UnixNano())),
	}
}

func (this *Generator) Bytes(n int) []byte {
	func() {
		now := time.Now()
		this.lock.Lock()
		defer this.lock.Unlock()
		if now.After(this.exp) {
			this.exp = now.Add(this.intrv)
			this.rand = rand.New(rand.NewSource(now.UnixNano()))
		}
	}()

	buff := make([]byte, n)
	for i := 0; i < n; i++ {
		buff[i] = byte(this.rand.Intn(256))
	}
	return buff
}

func (this *Generator) String(n int) string {
	return base64.URLEncoding.EncodeToString(this.Bytes((n*6 + 7) / 8))[:n]
}
