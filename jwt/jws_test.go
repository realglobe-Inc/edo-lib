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

package jwt

import (
	"crypto"
	"github.com/realglobe-Inc/edo-lib/jwk"
	"testing"
)

func TestHashFunction(t *testing.T) {
	for hGen, algs := range map[crypto.Hash][]string{
		crypto.SHA256: {"HS256", "ES256", "RS256", "PS256"},
		crypto.SHA384: {"HS384", "ES384", "RS384", "PS384"},
		crypto.SHA512: {"HS512", "ES512", "RS512", "PS512"},
	} {
		for _, alg := range algs {
			if hGen2, err := HashFunction(alg); err != nil {
				t.Error(alg)
				t.Fatal(err)
			} else if hGen2 != hGen {
				t.Error(hGen2)
				t.Fatal(hGen)
			}
		}
	}
}

func TestHs(t *testing.T) {
	buff := []byte{}
	for ; len(buff) < 50; buff = append(buff, byte(len(buff))) {
	}
	type param struct {
		key jwk.Key
		crypto.Hash
	}
	for _, p := range []param{{test_256Key, crypto.SHA256}, {test_384Key, crypto.SHA384}, {test_512Key, crypto.SHA512}} {
		for i := 0; i < 50; i += 10 {
			plain := make([]byte, i)
			copy(plain, buff)
			if sig, err := hsSign(p.key, p.Hash, plain); err != nil {
				t.Fatal(err)
			} else if err := hsVerify(p.key, p.Hash, sig, plain); err != nil {
				t.Error(sig)
				t.Fatal(err)
			}
		}
	}
}

func TestRs(t *testing.T) {
	buff := []byte{}
	for ; len(buff) < 50; buff = append(buff, byte(len(buff))) {
	}
	for _, hGen := range []crypto.Hash{crypto.SHA256, crypto.SHA384, crypto.SHA512} {
		for i := 0; i < 50; i += 10 {
			plain := make([]byte, i)
			copy(plain, buff)
			if sig, err := rsSign(test_rsaKey, hGen, plain); err != nil {
				t.Fatal(err)
			} else if err := rsVerify(test_rsaKey, hGen, sig, plain); err != nil {
				t.Error(sig)
				t.Fatal(err)
			}
		}
	}
}

func TestEs(t *testing.T) {
	buff := []byte{}
	for ; len(buff) < 50; buff = append(buff, byte(len(buff))) {
	}
	type param struct {
		key jwk.Key
		crypto.Hash
	}
	for _, p := range []param{{test_ec256Key, crypto.SHA256}, {test_ec384Key, crypto.SHA384}, {test_ec521Key, crypto.SHA512}} {
		for i := 0; i < 50; i += 10 {
			plain := make([]byte, i)
			copy(plain, buff)
			if sig, err := esSign(p.key, p.Hash, plain); err != nil {
				t.Fatal(err)
			} else if err := esVerify(p.key, p.Hash, sig, plain); err != nil {
				t.Error(sig)
				t.Fatal(err)
			}
		}
	}
}

func TestPs(t *testing.T) {
	buff := []byte{}
	for ; len(buff) < 50; buff = append(buff, byte(len(buff))) {
	}
	for _, hGen := range []crypto.Hash{crypto.SHA256, crypto.SHA384, crypto.SHA512} {
		for i := 0; i < 50; i += 10 {
			plain := make([]byte, i)
			copy(plain, buff)
			if sig, err := psSign(test_rsaKey, hGen, plain); err != nil {
				t.Fatal(err)
			} else if err := psVerify(test_rsaKey, hGen, sig, plain); err != nil {
				t.Error(sig)
				t.Fatal(err)
			}
		}
	}
}
