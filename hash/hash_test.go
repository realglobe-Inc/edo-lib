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
	"crypto/sha256"
	"testing"
)

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
