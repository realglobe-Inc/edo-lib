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
	"github.com/realglobe-Inc/edo-lib/jwk"
	"testing"
)

func TestFindKeySample(t *testing.T) {
	// JWS Appendix A.1 より。
	key, err := jwk.FromMap(map[string]interface{}{
		"kty": "oct",
		"k":   "AyM1SysPpbyDfgZld3umj1qzKObwVMkoqQ-EstJQLr_T-1qS0gZH75aKtMN3Yj0iPS4hcgUuTwjAzZr1Z9CAow",
	})
	if err != nil {
		t.Fatal(err)
	}

	if k := findKey([]jwk.Key{key}, "", "oct", "sig", "verify", "HS256"); k != key {
		t.Error(k)
		t.Fatal(key)
	}
}

func TestFindKey(t *testing.T) {
	keys := []jwk.Key{
		jwk.New(test_rsaKey.Private(), map[string]interface{}{
			"use": "sig",
		}),
		jwk.New(test_ec256Key.Private(), map[string]interface{}{
			"kid":     "2",
			"key_ops": []interface{}{"sign", "verify"},
		}),
		jwk.New(test_128Key.Common(), map[string]interface{}{
			"kid": "3",
			"alg": "A128KW",
		}),
		jwk.New(test_ec521Key.Private(), map[string]interface{}{
			"kid": "4",
		}),
	}

	if key := findKey(keys, ""); key == nil || key.Id() != "" {
		t.Fatal(key)
	} else if key := findKey(keys, "3"); key == nil || key.Id() != "3" {
		t.Fatal(key)
	} else if key := findKey(keys, "99"); key != nil {
		t.Fatal(key)
	} else if key := findKey(keys, "", "EC"); key == nil || key.Id() != "2" {
		t.Fatal(key)
	} else if key := findKey(keys, "", "RSA"); key == nil || key.Id() != "" {
		t.Fatal(key)
	} else if key := findKey(keys, "", "oct"); key == nil || key.Id() != "3" {
		t.Fatal(key)
	} else if key := findKey(keys, "", "", "sig"); key == nil || key.Id() != "" {
		t.Fatal(key)
	} else if key := findKey(keys, "", "", "enc"); key == nil || key.Id() != "3" {
		t.Fatal(key)
	} else if key := findKey(keys, "", "", "", "sign"); key == nil || key.Id() != "" {
		t.Fatal(key)
	} else if key := findKey(keys, "", "", "", "verify"); key == nil || key.Id() != "" {
		t.Fatal(key)
	} else if key := findKey(keys, "", "", "", "encrypt"); key == nil || key.Id() != "4" {
		t.Fatal(key)
	} else if key := findKey(keys, "", "", "", "decrypt"); key == nil || key.Id() != "4" {
		t.Fatal(key)
	} else if key := findKey(keys, "", "", "", "wrapKey"); key == nil || key.Id() != "3" {
		t.Fatal(key)
	} else if key := findKey(keys, "", "", "", "", "HS256"); key != nil {
		t.Fatal(key)
	} else if key := findKey(keys, "", "", "", "", "RS256"); key == nil || key.Id() != "" {
		t.Error(key)
		t.Fatal()
	} else if key := findKey(keys, "", "", "", "", "ES256"); key == nil || key.Id() != "2" {
		t.Fatal(key)
	} else if key := findKey(keys, "", "", "", "", "A128KW"); key == nil || key.Id() != "3" {
		t.Fatal(key)
	}
}
