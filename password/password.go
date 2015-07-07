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

// パスワード関係。
package password

import (
	"strconv"
	"strings"

	"github.com/realglobe-Inc/edo-lib/hash"
	"github.com/realglobe-Inc/go-lib/erro"
	"golang.org/x/crypto/pbkdf2"
)

const algSep = ":"

// 保存用ハッシュ値を計算する。
// alg が pbkdf2:... なら、params は salt ([]byte), passwd (string)。
func Calculate(alg string, params ...interface{}) ([]byte, error) {
	parts := strings.Split(alg, algSep)
	if len(parts) == 0 {
		return nil, erro.New("no algorithm")
	}

	switch parts[0] {
	case tagPbkdf2:
		if len(parts) < 2 {
			return nil, erro.New("no " + parts[0] + " hash algorithm")
		} else if len(parts) < 3 {
			return nil, erro.New("no " + parts[0] + " iteration number")
		} else if len(params) < 1 {
			return nil, erro.New("no " + parts[0] + " salt")
		} else if len(params) < 2 {
			return nil, erro.New("no " + parts[0] + " password")
		}

		hGen := hash.Generator(parts[1])
		if hGen == 0 {
			return nil, erro.New("unsupported hash algorithm " + parts[1])
		}
		iter, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, erro.Wrap(err)
		}
		salt, _ := params[0].([]byte)
		if salt == nil {
			return nil, erro.New("invalid salt")
		}
		passwd, _ := params[1].(string)
		if passwd == "" {
			return nil, erro.New("invalid password")
		}
		return pbkdf2.Key([]byte(passwd), salt, iter, hGen.Size(), hGen.New), nil
	default:
		return nil, erro.New("unsupported algorithm " + parts[0])
	}
}
