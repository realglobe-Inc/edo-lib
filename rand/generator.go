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

// ランダムな文字列とバイト列の生成器。
package rand

import (
	"github.com/realglobe-Inc/edo-lib/prand"
	"github.com/realglobe-Inc/edo-lib/secrand"
	"github.com/realglobe-Inc/go-lib/erro"
	"time"
)

type Generator interface {
	// n バイトのランダムバイト列を返す。
	Bytes(n int) []byte
	// 長さ n のランダム文字列を返す。
	String(n int) string
}

func New(intrv time.Duration) Generator {
	return &generator{prand.New(intrv)}
}

type generator struct {
	// 安全な乱数が使えないときの代替。
	p *prand.Generator
}

func (this *generator) String(n int) string {
	id, err := secrand.String(n)
	if err != nil {
		log.Err(erro.Wrap(err))
		id = this.p.String(n)
	}
	return id
}

func (this *generator) Bytes(n int) []byte {
	id, err := secrand.Bytes(n)
	if err != nil {
		log.Err(erro.Wrap(err))
		id = this.p.Bytes(n)
	}
	return id
}
