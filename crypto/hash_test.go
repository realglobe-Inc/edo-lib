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

package crypto

import (
	"crypto"
	"testing"
)

func TestHashFunction(t *testing.T) {
	for _, h := range []crypto.Hash{
		crypto.MD4,
		crypto.MD5,
		crypto.SHA1,
		crypto.SHA224,
		crypto.SHA256,
		crypto.SHA384,
		crypto.SHA512,
		crypto.MD5SHA1,
		crypto.RIPEMD160,
		crypto.SHA3_224,
		crypto.SHA3_256,
		crypto.SHA3_384,
		crypto.SHA3_512,
	} {
		s := HashFunctionString(h)
		h2, err := ParseHashFunction(s)
		if err != nil {
			t.Fatal(err)
		} else if h2 != h {
			t.Fatal(h2, h)
		}
	}
}

func TestParseUnknownHashFunction(t *testing.T) {
	_, err := ParseHashFunction("unknown")
	if err == nil {
		t.Fatal("no error")
	}
}

func TestUnknownHashFunctionString(t *testing.T) {
	str := HashFunctionString(crypto.Hash(1000000))
	if _, ok := strToHash[str]; ok {
		t.Fatal(str)
	}
}
