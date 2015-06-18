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
	"bytes"
	"crypto"
	"crypto/sha256"
	"testing"
)

func TestGeneratorAndAlgorithm(t *testing.T) {
	for _, hGen := range []crypto.Hash{
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
		alg := Algorithm(hGen)
		hGen2 := Generator(alg)
		if hGen2 != hGen {
			t.Error(hGen2)
			t.Error(alg)
			t.Fatal(hGen)
		}
	}
}

func TestUnknownGenerator(t *testing.T) {
	if hGen := Generator("unknown"); hGen != 0 {
		t.Fatal(hGen)
	}
}

func TestUnknownAlgorithm(t *testing.T) {
	if alg := Algorithm(crypto.Hash(1000000)); alg != "" {
		t.Fatal(alg)
	}
}

func TestHashing(t *testing.T) {
	hFun := sha256.New()
	data := [][]byte{}
	for i := 0; i < 100; i++ {
		hFun.Write([]byte{byte(i)})
		data = append(data, []byte{byte(i)})
	}
	h := hFun.Sum(nil)

	hFun.Reset()
	h2 := Hashing(hFun, data...)
	if !bytes.Equal(h2, h) {
		t.Error(h)
		t.Fatal(h2)
	}
}
