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

package hash

import (
	"crypto"
	"github.com/realglobe-Inc/go-lib/erro"
)

// related_users に入れるハッシュ値の文字列としての長さを返す。
func StringSize(alg string) (int, error) {
	switch alg {
	case "SHA256":
		return (128 + 5) / 6, nil
	case "SHA384":
		return (192 + 5) / 6, nil
	case "SHA512":
		return (256 + 5) / 6, nil
	default:
		return 0, erro.New("unsupported algorithm " + alg)
	}
}

func HashFunction(alg string) (crypto.Hash, error) {
	switch alg {
	case "SHA256":
		return crypto.SHA256, nil
	case "SHA384":
		return crypto.SHA384, nil
	case "SHA512":
		return crypto.SHA512, nil
	default:
		return 0, erro.New("unsupported algorithm " + alg)
	}
}