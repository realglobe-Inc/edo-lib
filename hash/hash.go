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

// ハッシュ値計算用関数。
package hash

import (
	"crypto"
	"hash"
	"strings"
)

var hashToAlg = map[crypto.Hash]string{
	crypto.MD4:       "MD4",
	crypto.MD5:       "MD5",
	crypto.SHA1:      "SHA1",
	crypto.SHA224:    "SHA224",
	crypto.SHA256:    "SHA256",
	crypto.SHA384:    "SHA384",
	crypto.SHA512:    "SHA512",
	crypto.MD5SHA1:   "MD5SHA1",
	crypto.RIPEMD160: "RIPEMD160",
	crypto.SHA3_224:  "SHA3_224",
	crypto.SHA3_256:  "SHA3_256",
	crypto.SHA3_384:  "SHA3_384",
	crypto.SHA3_512:  "SHA3_512",
}

var algToHash = map[string]crypto.Hash{}

func init() {
	for hGen, alg := range hashToAlg {
		algToHash[alg] = hGen
	}
}

// 入力の大文字、小文字は区別しない。
// 該当するものが無ければ nil を返す。
func Generator(alg string) crypto.Hash {
	return algToHash[strings.ToUpper(alg)]
}

// ハッシュアルゴリズム名を返す。
// 該当するものが無ければ空文字列。
func Algorithm(hGen crypto.Hash) string {
	return hashToAlg[hGen]
}

// ハッシュ値を計算して返す。
func Hashing(h hash.Hash, data ...[]byte) []byte {
	for _, d := range data {
		h.Write(d)
	}
	return h.Sum(nil)
}
