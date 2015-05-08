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
	"bytes"
	"testing"
)

func TestDef(t *testing.T) {
	for plain := []byte{}; len(plain) < 100; plain = append(plain, byte(len(plain))) {
		if d, err := defCompress(plain); err != nil {
			t.Fatal(err)
		} else if b2, err := defDecompress(d); err != nil {
			t.Fatal(err)
		} else if !bytes.Equal(b2, plain) {
			t.Error(b2)
			t.Fatal(plain)
		}
	}
}
