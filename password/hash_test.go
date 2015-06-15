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

package password

import (
	"crypto"
	"testing"
)

func TestHashFunction(t *testing.T) {
	if hGen, err := hashFunction("sha256"); err != nil {
		t.Fatal(err)
	} else if hGen != crypto.SHA256 {
		t.Error(hGen)
		t.Fatal(crypto.SHA256)
	} else if hGen, err := hashFunction("sha384"); err != nil {
		t.Fatal(err)
	} else if hGen != crypto.SHA384 {
		t.Error(hGen)
		t.Fatal(crypto.SHA384)
	} else if hGen, err := hashFunction("sha512"); err != nil {
		t.Fatal(err)
	} else if hGen != crypto.SHA512 {
		t.Error(hGen)
		t.Fatal(crypto.SHA512)
	}
}

func TestUnknownHashFunction(t *testing.T) {
	if _, err := hashFunction("unknown"); err == nil {
		t.Fatal("no error")
	}
}
