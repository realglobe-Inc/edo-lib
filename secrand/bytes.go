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

// 安全なランダム文字列とバイト列の生成器。
package secrand

import (
	"crypto/rand"
	"github.com/realglobe-Inc/go-lib/erro"
	"io"
)

func Bytes(length int) ([]byte, error) {
	buff := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, buff); err != nil {
		return nil, erro.Wrap(err)
	}
	return buff, nil
}
